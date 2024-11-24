package cmd

import (
	"context"
	"fmt"

	execute "github.com/alexellis/go-execute/v2"
)

// TestArgs represents arguments for test commands
type TestArgs struct {
	Verbose bool   `flag:"verbose" default:"false" usage:"Enable verbose output"`
	Pattern string `flag:"pattern" default:"./..." usage:"Test pattern to run"`
	Race    bool   `flag:"race" default:"false" usage:"Enable race detection"`
}

// LintArgs represents arguments for lint commands
type LintArgs struct {
	Fix     bool   `flag:"fix" default:"false" usage:"Auto-fix issues when possible"`
	Config  string `flag:"config" default:".golangci.yml" usage:"Path to config file"`
	Verbose bool   `flag:"verbose" default:"false" usage:"Enable verbose output"`
}

func TestCommand(ctx context.Context, args TestArgs) {
	cmdArgs := []string{"test"}
	if args.Verbose {
		cmdArgs = append(cmdArgs, "-v")
	}
	if args.Race {
		cmdArgs = append(cmdArgs, "-race")
	}
	cmdArgs = append(cmdArgs, args.Pattern)

	cmd := execute.ExecTask{
		Command:     "go",
		Args:        cmdArgs,
		StreamStdio: true,
	}

	res, err := cmd.Execute(ctx)
	if err != nil {
		panic(err)
	}

	if res.ExitCode != 0 {
		panic("Non-zero exit code: " + res.Stderr)
	}

	fmt.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
}

func LintCommand(ctx context.Context, args LintArgs) {
	cmdArgs := []string{"run"}
	if args.Fix {
		cmdArgs = append(cmdArgs, "--fix")
	}
	if args.Verbose {
		cmdArgs = append(cmdArgs, "-v")
	}


	cmd := execute.ExecTask{
		Command:     "golangci-lint",
		Args:        cmdArgs,
		StreamStdio: true,
	}

	res, err := cmd.Execute(ctx)
	if err != nil {
		panic(err)
	}

	if res.ExitCode != 0 {
		panic("Non-zero exit code: " + res.Stderr)
	}

	fmt.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
}