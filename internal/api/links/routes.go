// Package links exposes the HTTP routes for Shorty's link resource.
package links

import (
	"net/http"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

type handler struct {
	createLink app.CreateLinkHandler
	getLink    app.GetLinkHandler
	listLinks  app.ListLinksHandler
	visitLink  app.VisitLinkHandler
}

func newHandler(application app.Application) *handler {
	return &handler{
		createLink: application.Commands.CreateLink,
		getLink:    application.Queries.GetLink,
		listLinks:  application.Queries.ListLinks,
		visitLink:  application.Commands.VisitLink,
	}
}

// Register adds the link resource and its endpoints to api.
func Register(api *amigo.API, application app.Application, middlewares ...amigo.Middleware) {
	handler := newHandler(application)
	router := api.Group("/links", middlewares...)

	router.POST("", handler.CreateLink,
		amigo.WithStatus(http.StatusCreated),
		amigo.WithErrorMapping(domain.ErrOriginURLRequired, http.StatusBadRequest, "origin URL is required"),
		amigo.WithErrorMapping(domain.ErrOriginURLInvalid, http.StatusBadRequest, "origin URL must be an absolute HTTP or HTTPS URL"),
	)
	router.GET("", handler.ListLinks)
	router.GET("/{id}", handler.GetLink,
		amigo.WithErrorMapping(ports.ErrLinkNotFound, http.StatusNotFound, "link not found"),
	)
	router.POST("/{id}/visit", handler.VisitLink,
		amigo.WithStatus(http.StatusNoContent),
		amigo.WithErrorMapping(ports.ErrLinkNotFound, http.StatusNotFound, "link not found"),
	)
}
