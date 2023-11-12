package server

import (
	"context"
	"myproject/internal/config"
	"net/http"
)

type server struct {
	s *http.Server
}

func New(cfg *config.Config) *server {

	s := &http.Server{
		Addr:         cfg.HttpServer.HostPort(),
		ReadTimeout:  cfg.HttpServer.ReadTimeout(),
		WriteTimeout: cfg.HttpServer.WriteTimeout(),
		Handler:      nil,
	}
	return &server{s}
}
func (s *server) SetHandler(h http.Handler) {
	s.s.Handler = h
}
func (s *server) ListenAndServe() error {
	return s.s.ListenAndServe()
}
func (s *server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
