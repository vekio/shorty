// Package resolve exposes the public redirect operation for shortened codes.
package resolve

import (
	"net/http"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/ports"
)

// Register adds the shortened-code redirect endpoint below router's prefix.
func Register(router *amigo.Router, application app.Application) {
	handler := &handler{resolveLink: application.Commands.ResolveLink}
	router.GET("/{code}", handler.ResolveLink,
		amigo.WithStatus(http.StatusFound),
		amigo.WithErrorMapping(ports.ErrLinkNotFound, http.StatusNotFound, "link not found"),
	)
}
