package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/vekio/shorty/internal/bootstrap"
	"github.com/vekio/shorty/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("shorty stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	if len(os.Args) > 1 && os.Args[1] != "serve" {
		return fmt.Errorf("unknown command %q: expected serve", os.Args[1])
	}
	runtime, err := bootstrap.New(cfg)
	if err != nil {
		return fmt.Errorf("bootstrap shorty: %w", err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			runtime.Logger.Error("close shorty", "error", err)
		}
	}()

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           runtime.Handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	runtime.Logger.Info("shorty listening", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	return nil
}
