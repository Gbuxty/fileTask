package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workFileData/internal/config"

	"workFileData/internal/service"
	"workFileData/internal/transport/http/handlers"
	"workFileData/internal/transport/http/handlers/router"
	"workFileData/internal/transport/http/server"
	"workFileData/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	logger := logger.NewLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("init cfg", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ServerShutdownTimeout)
	defer cancel()
	svc := service.NewFileService(logger)
	svc.Start(ctx)

	handlers := handlers.NewFileDataHandlers(logger, svc)
	router := router.NewRouter(logger, handlers)
	srv := server.NewServer(cfg.HTTP, logger, router)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed:", zap.Error(err))
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-interrupt

	logger.Info("Received signal, shutting down server", zap.String("signal", sig.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}

	svc.Stop()

}
