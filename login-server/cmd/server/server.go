package main

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

func NewServer(s *http.Server) *Server {
	return &Server{Server: s}
}

func Run(ctx context.Context, s *Server) error {
	e := make(chan error, 1)
	go func() { e <- s.ListenAndServe() }()

	shutdown := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(ctx)
		return <-e
	}

	select {
	case err := <-e:
		return err
	case <-ctx.Done():
		return shutdown()
	}
}
