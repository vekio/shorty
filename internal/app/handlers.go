package app

import (
	"context"

	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/visitlink"
)

// CreateLinkHandler is the create-link command capability exposed to input
// adapters.
type CreateLinkHandler interface {
	Handle(context.Context, createlink.CreateLinkCommand) (createlink.CreateLinkResult, error)
}

// VisitLinkHandler is the visit-link command capability exposed to input adapters.
type VisitLinkHandler interface {
	Handle(context.Context, visitlink.VisitLinkCommand) (visitlink.VisitLinkResult, error)
}

// GetLinkHandler is the get-link query capability exposed to input
// adapters.
type GetLinkHandler interface {
	Handle(context.Context, getlink.GetLinkQuery) (getlink.GetLinkResult, error)
}
