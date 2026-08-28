// Package api exposes Shorty's JSON HTTP API.
package api

import (
	"log/slog"
	"net/http"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

// New builds the HTTP API and registers all routes.
func New(application app.Application, logger *slog.Logger) *amigo.API {
	endpoints := newHandlers(application)
	api := amigo.New(amigo.WithLogger(logger))
	registerValidators(api)
	links := api.Group("/links", logRequest(logger))
	links.POST("", endpoints.CreateLink,
		amigo.WithStatus(http.StatusCreated),
		amigo.WithErrorMapping(domain.ErrOriginURLRequired, http.StatusBadRequest, "origin URL is required"),
		amigo.WithErrorMapping(domain.ErrOriginURLInvalid, http.StatusBadRequest, "origin URL must be an absolute HTTP or HTTPS URL"),
	)
	links.GET("", endpoints.ListLinks)
	links.GET("/{id}", endpoints.GetLink,
		amigo.WithErrorMapping(ports.ErrLinkNotFound, http.StatusNotFound, "link not found"),
	)
	links.POST("/{id}/visit", endpoints.VisitLink,
		amigo.WithStatus(http.StatusNoContent),
		amigo.WithErrorMapping(ports.ErrLinkNotFound, http.StatusNotFound, "link not found"),
	)
	return api
}
