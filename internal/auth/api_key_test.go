package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

const testWorkspaceID = "ws_test"

type repositoryStub struct {
	key authKeySlot
}

type authKeySlot struct {
	value  APIKey
	exists bool
}

func (repository *repositoryStub) Save(_ context.Context, key APIKey) error {
	repository.key = authKeySlot{value: key, exists: true}
	return nil
}

func (repository *repositoryStub) Authenticate(_ context.Context, hash string, usedAt time.Time) (APIKey, error) {
	if !repository.key.exists || repository.key.value.TokenHash != hash || repository.key.value.RevokedAt != nil {
		return APIKey{}, ErrInvalidAPIKey
	}
	repository.key.value.LastUsedAt = &usedAt
	return repository.key.value, nil
}

func (repository *repositoryStub) List(_ context.Context, workspaceID string) ([]APIKey, error) {
	if !repository.key.exists {
		return []APIKey{}, nil
	}
	if repository.key.value.WorkspaceID != workspaceID {
		return []APIKey{}, nil
	}
	return []APIKey{repository.key.value}, nil
}

func (repository *repositoryStub) Revoke(_ context.Context, workspaceID string, id string, revokedAt time.Time) error {
	if !repository.key.exists || repository.key.value.WorkspaceID != workspaceID || repository.key.value.ID != id {
		return ErrAPIKeyNotFound
	}
	repository.key.value.RevokedAt = &revokedAt
	return nil
}

func TestServiceCreatesOpaqueKeyAndAuthenticatesToken(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository)
	key, token, err := service.Create(t.Context(), testWorkspaceID, "integration")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if token == "" || token == key.TokenHash || key.TokenHash != HashToken(token) {
		t.Fatalf("token = %q, key = %#v", token, key)
	}

	if err := service.AuthenticateToken(t.Context(), token); err != nil {
		t.Fatalf("AuthenticateToken() error = %v", err)
	}
	if repository.key.value.LastUsedAt == nil {
		t.Errorf("key = %#v, want last use registered", repository.key.value)
	}
}

func TestServiceRejectsMalformedAndRevokedTokens(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository)
	key, token, err := service.Create(t.Context(), testWorkspaceID, "integration")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := service.Revoke(t.Context(), testWorkspaceID, key.ID); err != nil {
		t.Fatalf("Revoke() error = %v", err)
	}

	for _, candidate := range []string{"", "invalid", token} {
		if err := service.AuthenticateToken(t.Context(), candidate); !errors.Is(err, ErrInvalidAPIKey) {
			t.Errorf("AuthenticateToken(%q) error = %v", candidate, err)
		}
	}
}

func TestListDoesNotExposeHashes(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository)
	key, _, err := service.Create(t.Context(), testWorkspaceID, "integration")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	keys, err := service.List(t.Context(), testWorkspaceID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(keys) != 1 || keys[0].ID != key.ID || keys[0].TokenHash != "" {
		t.Errorf("List() = %#v", keys)
	}
	if err := service.Revoke(t.Context(), testWorkspaceID, "missing"); !errors.Is(err, ErrAPIKeyNotFound) {
		t.Errorf("Revoke() error = %v", err)
	}
}
