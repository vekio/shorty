package sqlite

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/vekio/shorty/internal/auth"
)

func TestAPIKeyRepositoryPersistsUsageAndRevocation(t *testing.T) {
	repository, workspaceID := openTestAPIKeyRepository(t)
	service := auth.NewService(repository)
	created, token, err := service.Create(t.Context(), workspaceID, "deployment")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	authenticated, err := repository.Authenticate(t.Context(), auth.HashToken(token), created.CreatedAt.Add(1))
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	if authenticated.LastUsedAt == nil || authenticated.TokenHash == token {
		t.Errorf("authenticated key = %#v", authenticated)
	}
	if err := service.Revoke(t.Context(), workspaceID, created.ID); err != nil {
		t.Fatalf("Revoke() error = %v", err)
	}
	if _, err := repository.Authenticate(t.Context(), auth.HashToken(token), created.CreatedAt.Add(2)); !errors.Is(err, auth.ErrInvalidAPIKey) {
		t.Errorf("Authenticate() after revoke error = %v", err)
	}
	keys, err := service.List(t.Context(), workspaceID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(keys) != 1 || keys[0].RevokedAt == nil || keys[0].TokenHash != "" {
		t.Errorf("List() = %#v", keys)
	}
}

func openTestAPIKeyRepository(t *testing.T) (*APIKeyRepository, string) {
	t.Helper()
	database, err := Open(t.Context(), filepath.Join(t.TempDir(), "shorty.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	return NewAPIKeyRepository(database), testWorkspaceID(t, database)
}
