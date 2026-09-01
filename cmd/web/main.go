package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	webbootstrap "github.com/vekio/shorty/internal/bootstrap/web"
	webconfig "github.com/vekio/shorty/internal/config/web"
)

func main() {
	cfg, err := webconfig.Load()
	if err != nil {
		slog.Error("load Web configuration", "error", err)
		os.Exit(1)
	}
	runtime, err := webbootstrap.New(cfg)
	if err != nil {
		slog.Error("bootstrap Web process", "error", err)
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

	runtime.Logger.Info("shorty web listening", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		runtime.Logger.Error("HTTP server stopped", "error", err)
		os.Exit(1)
	}
}
