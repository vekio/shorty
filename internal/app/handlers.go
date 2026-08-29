package app

import (
	"context"

	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/deletelink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

// CreateLinkHandler is the create-link command capability exposed to input
// adapters.
type CreateLinkHandler interface {
	Handle(context.Context, createlink.CreateLinkCommand) (createlink.CreateLinkResult, error)
}

// ResolveLinkHandler is the resolve-link command capability exposed to input
// adapters.
type ResolveLinkHandler interface {
	Handle(context.Context, resolvelink.ResolveLinkCommand) (resolvelink.ResolveLinkResult, error)
}

// DeleteLinkHandler is the delete-link command capability exposed to input
// adapters.
type DeleteLinkHandler interface {
	Handle(context.Context, deletelink.DeleteLinkCommand) (deletelink.DeleteLinkResult, error)
}

// GetLinkHandler is the get-link query capability exposed to input
// adapters.
type GetLinkHandler interface {
	Handle(context.Context, getlink.GetLinkQuery) (getlink.GetLinkResult, error)
}

// ListLinksHandler is the list-links query capability exposed to input adapters.
type ListLinksHandler interface {
	Handle(context.Context, listlinks.ListLinksQuery) (listlinks.ListLinksResult, error)
}
