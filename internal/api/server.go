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
func New(application app.Application) *amigo.API {
	endpoints := newHandlers(application)
	api := amigo.New()
	registerValidators(api)
	links := api.Group("/links", logRequest)
	links.POST("", endpoints.CreateLink,
		amigo.Status(http.StatusCreated),
		amigo.MapError(domain.ErrOriginURLRequired, http.StatusBadRequest),
		amigo.MapError(domain.ErrOriginURLInvalid, http.StatusBadRequest),
	)
	links.GET("", endpoints.ListLinks)
	links.GET("/{id}", endpoints.GetLink,
		amigo.MapError(ports.ErrLinkNotFound, http.StatusNotFound),
	)
	links.POST("/{id}/visit", endpoints.VisitLink,
		amigo.Status(http.StatusNoContent),
		amigo.MapError(ports.ErrLinkNotFound, http.StatusNotFound),
	)
	return api
}
