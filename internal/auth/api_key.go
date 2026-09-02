// Package auth manages credentials for Shorty's management API.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

const tokenPrefix = "shorty_"

var (
	// ErrInvalidAPIKey indicates that a presented credential is unknown or revoked.
	ErrInvalidAPIKey = errors.New("invalid API key")
	// ErrAPIKeyNotFound indicates that an API key identifier does not exist.
	ErrAPIKeyNotFound = errors.New("API key not found")
)

// APIKey is the persisted metadata for one management credential.
// The clear-text token is deliberately never kept in this model.
type APIKey struct {
	ID          string     `json:"id"`
	WorkspaceID string     `json:"workspace_id"`
	Name        string     `json:"name"`
	TokenHash   string     `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}

// Repository persists and authenticates API keys.
type Repository interface {
	Save(context.Context, APIKey) error
	Authenticate(context.Context, string, time.Time) (APIKey, error)
	List(context.Context, string) ([]APIKey, error)
	Revoke(context.Context, string, string, time.Time) error
}

// Service creates and manages opaque API keys.
type Service struct {
	repository Repository
	now        func() time.Time
}

// NewService creates an API key service backed by repository.
func NewService(repository Repository) *Service {
	return &Service{repository: repository, now: func() time.Time { return time.Now().UTC() }}
}

// Create persists a new key and returns its clear-text token exactly once.
func (service *Service) Create(ctx context.Context, workspaceID string, name string) (APIKey, string, error) {
	if strings.TrimSpace(workspaceID) == "" {
		return APIKey{}, "", errors.New("workspace ID is required")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return APIKey{}, "", errors.New("API key name is required")
	}
	id, err := randomString(16)
	if err != nil {
		return APIKey{}, "", fmt.Errorf("generate API key ID: %w", err)
	}
	secret, err := randomString(32)
	if err != nil {
		return APIKey{}, "", fmt.Errorf("generate API key token: %w", err)
	}
	token := tokenPrefix + secret
	key := APIKey{
		ID:          "key_" + id,
		WorkspaceID: workspaceID,
		Name:        name,
		TokenHash:   HashToken(token),
		CreatedAt:   service.now(),
	}
	if err := service.repository.Save(ctx, key); err != nil {
		return APIKey{}, "", err
	}
	return key, token, nil
}

// AuthenticateToken verifies an opaque API token and registers its use.
func (service *Service) AuthenticateToken(ctx context.Context, token string) error {
	if !strings.HasPrefix(token, tokenPrefix) {
		return ErrInvalidAPIKey
	}
	_, err := service.repository.Authenticate(ctx, HashToken(token), service.now())
	return err
}

// List returns all keys without their hashes, newest first.
func (service *Service) List(ctx context.Context, workspaceID string) ([]APIKey, error) {
	keys, err := service.repository.List(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	for index := range keys {
		keys[index].TokenHash = ""
	}
	return keys, nil
}

// Revoke permanently prevents an API key from authenticating.
func (service *Service) Revoke(ctx context.Context, workspaceID string, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("API key ID is required")
	}
	return service.repository.Revoke(ctx, workspaceID, id, service.now())
}

// HashToken returns the non-reversible representation stored by repositories.
func HashToken(token string) string {
	digest := sha256.Sum256([]byte(token))
	return hex.EncodeToString(digest[:])
}

func randomString(size int) (string, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
