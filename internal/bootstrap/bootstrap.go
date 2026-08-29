// Package bootstrap wires application ports to their adapters.
package bootstrap

import (
	"log/slog"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/deletelink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

// Dependencies contains the application services required by Shorty's
// delivery layer.
type Dependencies struct {
	Application app.Application
	Logger      *slog.Logger

	linkRepo ports.LinkRepository
}

func New() Dependencies {
	deps := Dependencies{
		Logger:   newLogger(),
		linkRepo: newLinkRepository(),
	}
	deps.Application = newApplication(deps)
	return deps
}

func newApplication(deps Dependencies) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateLink:  createlink.NewCreateLinkHandler(deps.linkRepo),
			DeleteLink:  deletelink.NewDeleteLinkHandler(deps.linkRepo),
			ResolveLink: resolvelink.NewResolveLinkHandler(deps.linkRepo),
		},
		Queries: app.Queries{
			GetLink:   getlink.NewGetLinkHandler(deps.linkRepo),
			ListLinks: listlinks.NewListLinksHandler(deps.linkRepo),
		},
	}
}
