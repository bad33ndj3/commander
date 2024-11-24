package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	server *http.Server
	wg     sync.WaitGroup
}

func New(host string, port int) *Server {
	addr := fmt.Sprintf("%s:%d", host, port)
	return &Server{
		server: &http.Server{
			Addr: addr,
		},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	s.server.Handler = mux
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	s.wg.Wait()
	return nil
} 