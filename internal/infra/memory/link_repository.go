// Package memory provides in-memory implementations of application ports.
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

type LinkRepository struct {
	mu    sync.RWMutex
	links map[string]ownedLink
}

type ownedLink struct {
	ownerID string
	link    domain.Link
}

var _ ports.LinkRepository = (*LinkRepository)(nil)

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{links: make(map[string]ownedLink)}
}

func (repository *LinkRepository) Save(ctx context.Context, ownerID string, link domain.Link) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.links[link.Code()] = ownedLink{ownerID: ownerID, link: link}
	return nil
}

func (repository *LinkRepository) FindByCode(ctx context.Context, code string) (domain.Link, error) {
	if err := ctx.Err(); err != nil {
		return domain.Link{}, err
	}

	repository.mu.RLock()
	defer repository.mu.RUnlock()
	stored, exists := repository.links[code]
	if !exists {
		return domain.Link{}, ports.ErrLinkNotFound
	}
	return stored.link, nil
}

func (repository *LinkRepository) FindOwnedByCode(
	ctx context.Context,
	ownerID string,
	code string,
) (domain.Link, error) {
	if err := ctx.Err(); err != nil {
		return domain.Link{}, err
	}

	repository.mu.RLock()
	defer repository.mu.RUnlock()
	stored, exists := repository.links[code]
	if !exists || stored.ownerID != ownerID {
		return domain.Link{}, ports.ErrLinkNotFound
	}
	return stored.link, nil
}

func (repository *LinkRepository) FindPage(
	ctx context.Context,
	ownerID string,
	limit int,
	offset int,
) (ports.LinkPage, error) {
	if err := ctx.Err(); err != nil {
		return ports.LinkPage{}, err
	}

	repository.mu.RLock()
	defer repository.mu.RUnlock()
	links := make([]domain.Link, 0, len(repository.links))
	for _, stored := range repository.links {
		if stored.ownerID == ownerID {
			links = append(links, stored.link)
		}
	}
	sort.Slice(links, func(i, j int) bool {
		if links[i].CreatedAt().Equal(links[j].CreatedAt()) {
			return links[i].Code() > links[j].Code()
		}
		return links[i].CreatedAt().After(links[j].CreatedAt())
	})

	total := len(links)
	start := min(offset, total)
	pageSize := min(limit, total-start)
	return ports.LinkPage{
		Links: links[start : start+pageSize],
		Total: total,
	}, nil
}

func (repository *LinkRepository) UpdateLinkVisits(ctx context.Context, link domain.Link) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()
	stored, exists := repository.links[link.Code()]
	if !exists {
		return ports.ErrLinkNotFound
	}
	stored.link = link
	repository.links[link.Code()] = stored
	return nil
}

func (repository *LinkRepository) Delete(ctx context.Context, ownerID string, code string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()
	stored, exists := repository.links[code]
	if !exists || stored.ownerID != ownerID {
		return ports.ErrLinkNotFound
	}
	delete(repository.links, code)
	return nil
}
