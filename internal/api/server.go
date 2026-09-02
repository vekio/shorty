// Package api exposes Shorty's JSON HTTP API.
package api

import (
	"log/slog"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/api/links"
	"github.com/vekio/shorty/internal/api/validator"
	"github.com/vekio/shorty/internal/app"
)

// New builds the JSON API and registers its resource routes.
func New(application app.Application, logger *slog.Logger) *amigo.API {
	httpAPI := amigo.New(amigo.WithLogger(logger))
	validator.Register(httpAPI)

	links.Register(httpAPI.Group("/links"), application)

	return httpAPI
}
