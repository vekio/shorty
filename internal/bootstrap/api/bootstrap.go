// Package api wires the application handlers required by the JSON API.
package api

import (
	"github.com/vekio/shorty/internal/api"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/deletelink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/app/resolvelink"
	"github.com/vekio/shorty/internal/bootstrap"
	apiconfig "github.com/vekio/shorty/internal/config/api"
	"github.com/vekio/shorty/internal/domain"
	"github.com/vekio/shorty/internal/httpmiddleware"
)

// New composes the complete JSON API process.
func New(config apiconfig.Config) (bootstrap.Runtime, error) {
	originURLPolicy, err := domain.DisallowOriginHost(config.ShortURL)
	if err != nil {
		return bootstrap.Runtime{}, err
	}
	repository := bootstrap.NewLinkRepository()
	application := app.Application{
		Commands: app.Commands{
			CreateLink:  createlink.NewCreateLinkHandler(repository, originURLPolicy),
			DeleteLink:  deletelink.NewDeleteLinkHandler(repository),
			ResolveLink: resolvelink.NewResolveLinkHandler(repository),
		},
		Queries: app.Queries{
			GetLink:   getlink.NewGetLinkHandler(repository),
			ListLinks: listlinks.NewListLinksHandler(repository),
		},
	}
	processLogger, err := bootstrap.NewLogger(config.Logger)
	if err != nil {
		return bootstrap.Runtime{}, err
	}
	handler := api.New(application, processLogger)

	return bootstrap.Runtime{
		Handler: httpmiddleware.LogRequests(processLogger)(handler),
		Logger:  processLogger,
	}, nil
}
