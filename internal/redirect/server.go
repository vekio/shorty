// Package redirect exposes public shortened-link navigation.
package redirect

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

type input struct {
	Code string `path:"code" json:"-"`
}

type output struct {
	Location string `header:"Location" json:"-"`
}

type handler struct {
	resolveLink app.ResolveLinkHandler
}

// New builds the public redirect handler independently from the management API.
func New(application app.Application, logger *slog.Logger) http.Handler {
	httpAPI := amigo.New(amigo.WithLogger(logger))
	handler := &handler{resolveLink: application.Commands.ResolveLink}
	httpAPI.GET("/r/{code}", handler.ResolveLink,
		amigo.WithStatus(http.StatusFound),
		amigo.WithErrorMapping(ports.ErrLinkNotFound, http.StatusNotFound, "link not found"),
	)
	return httpAPI
}

func (h *handler) ResolveLink(ctx context.Context, input input) (output, error) {
	result, err := h.resolveLink.Handle(ctx, resolvelink.ResolveLinkCommand{Code: input.Code})
	if err != nil {
		return output{}, err
	}
	return output{Location: result.OriginURL}, nil
}
