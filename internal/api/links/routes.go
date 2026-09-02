// Package links exposes the HTTP routes for Shorty's link resource.
package links

import (
	"net/http"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
)

type handler struct {
	createLink app.CreateLinkHandler
	updateLink app.UpdateLinkHandler
	deleteLink app.DeleteLinkHandler
	getLink    app.GetLinkHandler
	listLinks  app.ListLinksHandler
}

func newHandler(application app.Application) *handler {
	return &handler{
		createLink: application.Commands.CreateLink,
		updateLink: application.Commands.UpdateLink,
		deleteLink: application.Commands.DeleteLink,
		getLink:    application.Queries.GetLink,
		listLinks:  application.Queries.ListLinks,
	}
}

// Register adds the link resource and its endpoints to api.
func Register(router *amigo.Router, application app.Application) {
	handler := newHandler(application)

	createOptions := append([]amigo.RouteOption{amigo.WithStatus(http.StatusCreated)}, originURLErrors...)
	router.POST("", handler.CreateLink, createOptions...)
	router.GET("", handler.ListLinks, paginationErrors...)
	router.GET("/{code}", handler.GetLink, linkNotFoundError)
	updateOptions := []amigo.RouteOption{linkNotFoundError}
	updateOptions = append(updateOptions, originURLErrors...)
	router.PATCH("/{code}", handler.UpdateLink, updateOptions...)
	router.DELETE("/{code}", handler.DeleteLink,
		amigo.WithStatus(http.StatusNoContent),
		linkNotFoundError,
	)
}
