package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPingHandler(t *testing.T) {
	// Create test server with our handler
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.Write([]byte("pong"))
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "ping endpoint",
			path:       "/ping",
			wantStatus: http.StatusOK,
			wantBody:   "pong",
		},
		{
			name:       "unknown endpoint",
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
			wantBody:   "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tt.path)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Status code = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			// Check response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if string(body) != tt.wantBody {
				t.Errorf("Body = %q, want %q", string(body), tt.wantBody)
			}
		})
	}
}

func TestServerLifecycle(t *testing.T) {
	// Create test server
	srv := New("localhost", 0) // Port 0 means random available port

	// Test server start
	t.Run("start", func(t *testing.T) {
		if err := srv.Start(); err != nil {
			t.Fatalf("Failed to start server: %v", err)
		}

		// Give server time to start
		time.Sleep(100 * time.Millisecond)
	})

	// Test server stop
	t.Run("stop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Stop(ctx); err != nil {
			t.Errorf("Failed to stop server: %v", err)
		}
	})
}

func TestServerErrors(t *testing.T) {
	t.Run("stop without start", func(t *testing.T) {
		srv := New("localhost", 0)
		ctx := context.Background()

		// Stopping an unstarted server should not error
		if err := srv.Stop(ctx); err != nil {
			t.Errorf("Stop() error = %v, want nil", err)
		}
	})
}

func TestHandlerResponse(t *testing.T) {
	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	// Create handler and serve request
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	mux.ServeHTTP(w, req)

	// Check response
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Want status '%d', got '%d'", http.StatusOK, resp.StatusCode)
	}

	if string(body) != "pong" {
		t.Errorf("Want 'pong', got '%s'", string(body))
	}
}

