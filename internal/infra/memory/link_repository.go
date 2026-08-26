// Package memory provides in-memory implementations of application ports.
package memory

import (
	"context"
	"sync"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

type LinkRepository struct {
	mu    sync.RWMutex
	links map[string]domain.Link
}

var _ ports.LinkRepository = (*LinkRepository)(nil)

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{links: make(map[string]domain.Link)}
}

func (repository *LinkRepository) Save(ctx context.Context, link domain.Link) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.links[link.Code()] = link
	return nil
}

func (repository *LinkRepository) FindByCode(ctx context.Context, code string) (domain.Link, error) {
	if err := ctx.Err(); err != nil {
		return domain.Link{}, err
	}

	repository.mu.RLock()
	defer repository.mu.RUnlock()
	link, exists := repository.links[code]
	if !exists {
		return domain.Link{}, ports.ErrLinkNotFound
	}
	return link, nil
}

func (repository *LinkRepository) UpdateLinkVisits(ctx context.Context, link domain.Link) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()
	_, exists := repository.links[link.Code()]
	if !exists {
		return ports.ErrLinkNotFound
	}
	repository.links[link.Code()] = link
	return nil
}
