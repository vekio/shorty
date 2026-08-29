// Package root exposes Shorty's routes that live at the API root.
package root

import (
	"net/http"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/ports"
)

type handler struct {
	resolveLink app.ResolveLinkHandler
}

func newHandler(application app.Application) *handler {
	return &handler{
		resolveLink: application.Commands.ResolveLink,
	}
}

// Register adds the public root-level endpoints to api.
func Register(
	api *amigo.API,
	application app.Application,
	middlewares ...amigo.Middleware,
) {
	handler := newHandler(application)

	api.GET("/{code}", handler.ResolveLink,
		amigo.WithStatus(http.StatusFound),
		amigo.WithMiddleware(middlewares...),
		amigo.WithErrorMapping(
			ports.ErrLinkNotFound,
			http.StatusNotFound, "link not found"),
	)
}
