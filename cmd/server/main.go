package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/go-task-api/internal/handler"
	"github.com/example/go-task-api/internal/store"
)

func main() {
	// Structured logger (JSON in production, text for local)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Dependencies
	taskStore := store.NewMemoryStore()
	h := handler.New(taskStore, logger)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Middleware chain
	var finalHandler http.Handler = mux
	finalHandler = handler.LoggingMiddleware(logger, finalHandler)
	finalHandler = handler.RecoveryMiddleware(logger, finalHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      finalHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}
