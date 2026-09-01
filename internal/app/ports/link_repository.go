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
	Save(context.Context, string, domain.Link) error
}

// LinkFinder retrieves one link by its public code.
type LinkFinder interface {
	FindByCode(context.Context, string) (domain.Link, error)
}

// OwnedLinkFinder retrieves one link only when it belongs to an owner.
type OwnedLinkFinder interface {
	FindOwnedByCode(context.Context, string, string) (domain.Link, error)
}

// LinkLister retrieves newest-first pages of links owned by one caller.
type LinkLister interface {
	FindPage(ctx context.Context, ownerID string, limit int, offset int) (LinkPage, error)
}

// LinkPage contains one creation-ordered slice and the unfiltered total.
type LinkPage struct {
	Links []domain.Link
	Total int
}

// LinkVisitsUpdater persists a link whose visit count has changed.
type LinkVisitsUpdater interface {
	UpdateLinkVisits(context.Context, domain.Link) error
}

// LinkVisitor groups the capabilities needed to register a visit.
type LinkVisitor interface {
	LinkFinder
	LinkVisitsUpdater
}

// LinkDeleter removes one link by its public code.
type LinkDeleter interface {
	Delete(context.Context, string, string) error
}

// LinkRepository is the complete persistence port used by bootstrap.
type LinkRepository interface {
	LinkSaver
	LinkFinder
	OwnedLinkFinder
	LinkLister
	LinkVisitsUpdater
	LinkDeleter
}
