// Package ports defines the dependencies required by the application layer.
package ports

import (
	"context"
	"errors"

	"github.com/vekio/shorty/internal/domain"
)

var ErrLinkNotFound = errors.New("link not found")

type LinkRepository interface {
	Save(context.Context, domain.Link) error
	FindByCode(context.Context, string) (domain.Link, error)
	UpdateLinkVisits(context.Context, domain.Link) error
}
