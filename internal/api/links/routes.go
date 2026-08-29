// Package links exposes the HTTP routes for Shorty's link resource.
package links

import (
	"net/http"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

type handler struct {
	createLink app.CreateLinkHandler
	deleteLink app.DeleteLinkHandler
	getLink    app.GetLinkHandler
	listLinks  app.ListLinksHandler
}

func newHandler(application app.Application) *handler {
	return &handler{
		createLink: application.Commands.CreateLink,
		deleteLink: application.Commands.DeleteLink,
		getLink:    application.Queries.GetLink,
		listLinks:  application.Queries.ListLinks,
	}
}

// Register adds the link resource and its endpoints to api.
func Register(router *amigo.Router, application app.Application) {
	handler := newHandler(application)

	router.POST("", handler.CreateLink,
		amigo.WithStatus(http.StatusCreated),
		amigo.WithErrorMapping(
			domain.ErrOriginURLRequired,
			http.StatusUnprocessableEntity, "origin URL is required"),
		amigo.WithErrorMapping(
			domain.ErrOriginURLInvalid,
			http.StatusUnprocessableEntity, "origin URL must be an absolute HTTP or HTTPS URL"),
	)
	router.GET("", handler.ListLinks,
		amigo.WithErrorMapping(listlinks.ErrInvalidLimit,
			http.StatusBadRequest, listlinks.ErrInvalidLimit.Error()),
		amigo.WithErrorMapping(listlinks.ErrInvalidOffset,
			http.StatusBadRequest, listlinks.ErrInvalidOffset.Error()),
	)
	router.GET("/{id}", handler.GetLink,
		amigo.WithErrorMapping(
			ports.ErrLinkNotFound,
			http.StatusNotFound, "link not found"),
	)
	router.DELETE("/{id}", handler.DeleteLink,
		amigo.WithStatus(http.StatusNoContent),
		amigo.WithErrorMapping(
			ports.ErrLinkNotFound,
			http.StatusNotFound, "link not found"),
	)
}
