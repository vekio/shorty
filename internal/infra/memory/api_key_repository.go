package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/vekio/shorty/internal/auth"
)

// APIKeyRepository keeps API key metadata for the lifetime of the process.
type APIKeyRepository struct {
	mu       sync.RWMutex
	keysByID map[string]auth.APIKey
	idByHash map[string]string
}

var _ auth.Repository = (*APIKeyRepository)(nil)

func NewAPIKeyRepository() *APIKeyRepository {
	return &APIKeyRepository{
		keysByID: make(map[string]auth.APIKey),
		idByHash: make(map[string]string),
	}
}

func (repository *APIKeyRepository) Save(ctx context.Context, key auth.APIKey) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.keysByID[key.ID] = key
	repository.idByHash[key.TokenHash] = key.ID
	return nil
}

func (repository *APIKeyRepository) Authenticate(
	ctx context.Context,
	tokenHash string,
	usedAt time.Time,
) (auth.APIKey, error) {
	if err := ctx.Err(); err != nil {
		return auth.APIKey{}, err
	}
	repository.mu.Lock()
	defer repository.mu.Unlock()
	id, exists := repository.idByHash[tokenHash]
	if !exists {
		return auth.APIKey{}, auth.ErrInvalidAPIKey
	}
	key := repository.keysByID[id]
	if key.RevokedAt != nil {
		return auth.APIKey{}, auth.ErrInvalidAPIKey
	}
	usedAt = usedAt.UTC()
	key.LastUsedAt = &usedAt
	repository.keysByID[id] = key
	return key, nil
}

func (repository *APIKeyRepository) List(ctx context.Context, workspaceID string) ([]auth.APIKey, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	keys := make([]auth.APIKey, 0, len(repository.keysByID))
	for _, key := range repository.keysByID {
		if key.WorkspaceID == workspaceID {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].CreatedAt.After(keys[j].CreatedAt)
	})
	return keys, nil
}

func (repository *APIKeyRepository) Revoke(
	ctx context.Context,
	workspaceID string,
	id string,
	revokedAt time.Time,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	repository.mu.Lock()
	defer repository.mu.Unlock()
	key, exists := repository.keysByID[id]
	if !exists || key.WorkspaceID != workspaceID {
		return auth.ErrAPIKeyNotFound
	}
	if key.RevokedAt == nil {
		revokedAt = revokedAt.UTC()
		key.RevokedAt = &revokedAt
		repository.keysByID[id] = key
	}
	return nil
}
