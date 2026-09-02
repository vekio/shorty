// Package ports defines the dependencies required by the application layer.
package ports

import (
	"context"
	"errors"

	"github.com/vekio/shorty/internal/domain"
)

// ErrLinkNotFound indicates that no link exists for the requested code.
var ErrLinkNotFound = errors.New("link not found")

// LinkSaver persists a newly created link.
type LinkSaver interface {
	Save(context.Context, domain.Link) error
}

// LinkFinder retrieves one link by its public code.
type LinkFinder interface {
	FindByCode(context.Context, string) (domain.Link, error)
}

// LinkLister retrieves newest-first pages of links.
type LinkLister interface {
	FindPage(ctx context.Context, limit int, offset int) (LinkPage, error)
}

// LinkPage contains one creation-ordered slice and the unfiltered total.
type LinkPage struct {
	Links []domain.Link
	Total int
}

// LinkOriginUpdater persists a link whose destination has changed.
type LinkOriginUpdater interface {
	UpdateLinkOrigin(context.Context, domain.Link) error
}

// LinkEditor groups the capabilities needed to change a destination.
type LinkEditor interface {
	LinkFinder
	LinkOriginUpdater
}

// LinkResolver atomically registers a visit and returns the resolved link.
type LinkResolver interface {
	ResolveByCode(context.Context, string) (domain.Link, error)
}

// LinkDeleter removes one link by its public code.
type LinkDeleter interface {
	Delete(context.Context, string) error
}

// LinkRepository is the complete persistence port used by bootstrap.
type LinkRepository interface {
	LinkSaver
	LinkFinder
	LinkLister
	LinkOriginUpdater
	LinkResolver
	LinkDeleter
}
