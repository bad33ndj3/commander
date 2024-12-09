package commander

import (
	"context"
	"strings"
	"testing"
)

// Test structs
type TestArgs struct {
	Flag1 bool   `flag:"flag1" default:"false" usage:"Test flag 1"`
	Flag2 string `flag:"flag2" default:"test" usage:"Test flag 2"`
}

func TestCommanderBasics(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  error
		wantHelp bool
	}{
		{
			name:     "no arguments",
			args:     []string{"prog"},
			wantErr:  ErrNoSubcommand,
			wantHelp: true,
		},
		{
			name:     "unknown command",
			args:     []string{"prog", "unknown"},
			wantErr:  ErrUnknownCommand,
			wantHelp: true,
		},
		{
			name:     "help command",
			args:     []string{"prog", "help"},
			wantErr:  nil,
			wantHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdr := newWithArgs(tt.args)
			var builder strings.Builder
			cmdr.SetOutput(&builder)
			err := cmdr.Run()

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStructArguments(t *testing.T) {
	var capturedArgs TestArgs
	testHandler := func(ctx context.Context, args TestArgs) {
		capturedArgs = args
	}

	tests := []struct {
		name    string
		args    []string
		want    TestArgs
		wantErr bool
	}{
		{
			name: "default values",
			args: []string{"prog", "test"},
			want: TestArgs{
				Flag1: false,
				Flag2: "test",
			},
		},
		{
			name: "set bool flag",
			args: []string{"prog", "test", "--flag1"},
			want: TestArgs{
				Flag1: true,
				Flag2: "test",
			},
		},
		{
			name: "set string flag",
			args: []string{"prog", "test", "--flag2", "value"},
			want: TestArgs{
				Flag1: false,
				Flag2: "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdr := newWithArgs(tt.args)
			var output strings.Builder
			cmdr.output = &output
			cat := cmdr.AddCategory("Test")
			cat.Register(&Command{
				Name:        "test",
				Description: "Test command",
				Handler:     testHandler,
			})

			err := cmdr.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if capturedArgs != tt.want {
				t.Errorf("Args = %v, want %v", capturedArgs, tt.want)
			}
		})
	}
}

func TestHelpOutput(t *testing.T) {
	cmdr := newWithArgs([]string{"prog", "help"})
	// Capture output
	var builder strings.Builder
	cmdr.SetOutput(&builder)

	cat := cmdr.AddCategory("Test")
	cat.Register(&Command{
		Name:        "test",
		Description: "Test command",
		Handler: func(ctx context.Context, args TestArgs) {
			_ = args
		},
	})

	err := cmdr.Run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	expected := []string{
		"Available Commands",
		"Test",
		"test",
		"Test command",
		"--flag1",
		"--flag2",
	}

	for _, exp := range expected {
		if !strings.Contains(builder.String(), exp) {
			t.Logf("Got output:\n%s", builder.String())
			t.Errorf("Help output missing: %s", exp)
		}
	}
}
