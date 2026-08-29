// Package api exposes Shorty's HTTP API.
package api

import (
	"log/slog"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/api/links"
	"github.com/vekio/shorty/internal/api/root"
	"github.com/vekio/shorty/internal/api/validator"
	"github.com/vekio/shorty/internal/app"
)

// New builds the HTTP API and registers all routes.
func New(application app.Application, logger *slog.Logger) *amigo.API {
	httpAPI := amigo.New(amigo.WithLogger(logger))
	validator.Register(httpAPI)

	requestLogger := logRequest(logger)

	links.Register(httpAPI.Group("/links", requestLogger), application)
	root.Register(httpAPI, application, requestLogger)

	return httpAPI
}
