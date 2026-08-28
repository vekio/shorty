package main

import (
	"net/http"
	"os"
	"time"

	"github.com/vekio/shorty/internal/api"
	"github.com/vekio/shorty/internal/bootstrap"
)

func main() {
	deps := bootstrap.New()
	httpAPI := api.New(deps.Application, deps.Logger)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           httpAPI,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	deps.Logger.Info("shorty listening", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		deps.Logger.Error("HTTP server stopped", "error", err)
		os.Exit(1)
	}
}
