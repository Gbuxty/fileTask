package server

import (
	"context"
	"fmt"
	"workFileData/internal/config"
	"workFileData/pkg/logger"

	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
	cfg        config.HTTPServer
}

func NewServer(
	cfg config.HTTPServer,
	logger *logger.Logger,
	handler http.Handler,
) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run() error {
    s.logger.Info("Starting HTTP server", zap.Int("port", s.cfg.Port))
    return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
    s.logger.Info("Shutting down server...")
    return s.httpServer.Shutdown(ctx)
}
