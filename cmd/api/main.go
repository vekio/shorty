package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	apibootstrap "github.com/vekio/shorty/internal/bootstrap/api"
	apiconfig "github.com/vekio/shorty/internal/config/api"
)

func main() {
	cfg, err := apiconfig.Load()
	if err != nil {
		slog.Error("load API configuration", "error", err)
		os.Exit(1)
	}
	runtime, err := apibootstrap.New(cfg)
	if err != nil {
		slog.Error("bootstrap API process", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           runtime.Handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	runtime.Logger.Info("shorty API listening", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		runtime.Logger.Error("HTTP server stopped", "error", err)
		os.Exit(1)
	}
}
