// Package bootstrap wires application ports to their adapters.
package bootstrap

import (
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/app/visitlink"
)

type dependencies struct {
	linkRepo ports.LinkRepository
}

func New() app.Application {
	deps := dependencies{
		linkRepo: newLinkRepository(),
	}
	return newApplication(deps)
}

func newApplication(deps dependencies) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateLink: createlink.NewCreateLinkHandler(deps.linkRepo),
			VisitLink:  visitlink.NewVisitLinkHandler(deps.linkRepo),
		},
		Queries: app.Queries{
			GetLink: getlink.NewGetLinkHandler(deps.linkRepo),
		},
	}
}
