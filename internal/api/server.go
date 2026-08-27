// Package api exposes Shorty's JSON HTTP API.
package api

import (
	"net/http"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
	"github.com/vekio/shorty/pkg/amigo"
)

// New builds the HTTP API and registers all routes.
func New(application app.Application) *amigo.Api {
	endpoints := newHandlers(application)
	api := amigo.New()
	api.POST("/links", endpoints.CreateLink,
		amigo.Status(http.StatusCreated),
		amigo.MapError(domain.ErrOriginURLRequired, http.StatusBadRequest),
		amigo.MapError(domain.ErrOriginURLInvalid, http.StatusBadRequest),
	)
	api.GET("/links", endpoints.ListLinks)
	api.GET("/links/{id}", endpoints.GetLink,
		amigo.MapError(ports.ErrLinkNotFound, http.StatusNotFound),
	)
	api.POST("/links/{id}/visit", endpoints.VisitLink,
		amigo.Status(http.StatusNoContent),
		amigo.MapError(ports.ErrLinkNotFound, http.StatusNotFound),
	)
	return api
}
