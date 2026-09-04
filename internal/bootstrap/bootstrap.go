// Package bootstrap composes the Shorty application and its HTTP adapters.
package bootstrap

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/vekio/shorty/internal/api"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/deletelink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/app/resolvelink"
	"github.com/vekio/shorty/internal/app/updatelink"
	shortyauth "github.com/vekio/shorty/internal/auth"
	shortyconfig "github.com/vekio/shorty/internal/config"
	"github.com/vekio/shorty/internal/domain"
	"github.com/vekio/shorty/internal/httpmiddleware"
	"github.com/vekio/shorty/internal/redirect"
	"github.com/vekio/shorty/internal/web"
)

// dependencies carries initialized process dependencies through composition.
// It keeps New small as logging, persistence, and adapters grow independently.
type dependencies struct {
	logger          *slog.Logger
	persistence     Persistence
	originURLPolicy domain.OriginURLPolicy
	auth            *shortyauth.Service
}

// New composes the single Shorty process in three stages: initialize shared
// dependencies, build the application, and finally mount its HTTP adapters.
func New(settings shortyconfig.Config) (Runtime, error) {
	deps, err := newDependencies(context.Background(), settings)
	if err != nil {
		return Runtime{}, err
	}

	application := newApplication(deps)
	handler := newHTTPHandler(application, deps)
	return Runtime{
		Handler: handler,
		Logger:  deps.logger,
		close:   deps.persistence.Close,
	}, nil
}

// newDependencies validates settings and opens process-wide dependencies once.
func newDependencies(ctx context.Context, settings shortyconfig.Config) (dependencies, error) {
	if err := settings.Validate(); err != nil {
		return dependencies{}, err
	}
	originURLPolicy, err := domain.DisallowOriginHost(settings.ShortURL)
	if err != nil {
		return dependencies{}, err
	}
	logger, err := NewLogger(settings.Logger)
	if err != nil {
		return dependencies{}, err
	}
	persistence, err := OpenPersistence(ctx, settings.Database)
	if err != nil {
		return dependencies{}, err
	}
	return dependencies{
		logger:          logger,
		persistence:     persistence,
		originURLPolicy: originURLPolicy,
		auth:            shortyauth.NewService(persistence.APIKeys),
	}, nil
}

// newApplication injects persistence and domain policies into every use case.
func newApplication(deps dependencies) app.Application {
	links := deps.persistence.Links
	return app.Application{
		Commands: app.Commands{
			CreateLink:  createlink.NewCreateLinkHandler(links, deps.originURLPolicy),
			UpdateLink:  updatelink.NewUpdateLinkHandler(links, deps.originURLPolicy),
			DeleteLink:  deletelink.NewDeleteLinkHandler(links),
			ResolveLink: resolvelink.NewResolveLinkHandler(links),
		},
		Queries: app.Queries{
			GetLink:   getlink.NewGetLinkHandler(links),
			ListLinks: listlinks.NewListLinksHandler(links),
		},
	}
}

// newHTTPHandler exposes the application through API, redirect, and Web routes.
func newHTTPHandler(application app.Application, deps dependencies) http.Handler {
	managementAPI := httpmiddleware.RequireAPIKey(deps.auth)(api.New(application, deps.logger))
	handler := NewServer(
		redirect.New(application, deps.logger),
		managementAPI,
		web.New(
			deps.auth,
			deps.persistence.WorkspaceID,
			deps.persistence.WorkspaceName,
		),
	)
	return httpmiddleware.LogRequests(deps.logger)(handler)
}
