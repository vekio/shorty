package main

import (
	"log"
	"net/http"
	"time"

	"github.com/vekio/shorty/internal/api"
	"github.com/vekio/shorty/internal/bootstrap"
)

func main() {
	application := bootstrap.New()
	httpAPI := api.New(application)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           httpAPI,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("shorty listening on http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
