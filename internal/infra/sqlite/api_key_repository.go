package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vekio/shorty/internal/auth"
	"github.com/vekio/shorty/internal/infra/sqlite/sqlitedb"
)

// APIKeyRepository persists management credentials in SQLite.
type APIKeyRepository struct {
	queries *sqlitedb.Queries
}

var _ auth.Repository = (*APIKeyRepository)(nil)

func NewAPIKeyRepository(database *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{queries: sqlitedb.New(database)}
}

func (repository *APIKeyRepository) Save(ctx context.Context, key auth.APIKey) error {
	return repository.queries.CreateAPIKey(ctx, sqlitedb.CreateAPIKeyParams{
		ID:          key.ID,
		WorkspaceID: key.WorkspaceID,
		Name:        key.Name,
		TokenHash:   key.TokenHash,
		CreatedAt:   formatTime(key.CreatedAt),
	})
}

func (repository *APIKeyRepository) Authenticate(
	ctx context.Context,
	tokenHash string,
	usedAt time.Time,
) (auth.APIKey, error) {
	row, err := repository.queries.AuthenticateAPIKey(ctx, sqlitedb.AuthenticateAPIKeyParams{
		LastUsedAt: sql.NullString{String: formatTime(usedAt), Valid: true},
		TokenHash:  tokenHash,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return auth.APIKey{}, auth.ErrInvalidAPIKey
	}
	if err != nil {
		return auth.APIKey{}, err
	}
	return restoreAPIKey(row)
}

func (repository *APIKeyRepository) List(ctx context.Context, workspaceID string) ([]auth.APIKey, error) {
	rows, err := repository.queries.ListAPIKeys(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	keys := make([]auth.APIKey, 0, len(rows))
	for _, row := range rows {
		key, err := restoreAPIKey(row)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (repository *APIKeyRepository) Revoke(
	ctx context.Context,
	workspaceID string,
	id string,
	revokedAt time.Time,
) error {
	updated, err := repository.queries.RevokeAPIKey(ctx, sqlitedb.RevokeAPIKeyParams{
		RevokedAt:   sql.NullString{String: formatTime(revokedAt), Valid: true},
		WorkspaceID: workspaceID,
		ID:          id,
	})
	if err != nil {
		return err
	}
	if updated == 0 {
		return auth.ErrAPIKeyNotFound
	}
	return nil
}

func restoreAPIKey(row sqlitedb.ApiKey) (auth.APIKey, error) {
	createdAt, err := parseTime(row.CreatedAt)
	if err != nil {
		return auth.APIKey{}, fmt.Errorf("parse API key creation time: %w", err)
	}
	lastUsedAt, err := parseOptionalTime(row.LastUsedAt)
	if err != nil {
		return auth.APIKey{}, fmt.Errorf("parse API key last-used time: %w", err)
	}
	revokedAt, err := parseOptionalTime(row.RevokedAt)
	if err != nil {
		return auth.APIKey{}, fmt.Errorf("parse API key revocation time: %w", err)
	}
	return auth.APIKey{
		ID:          row.ID,
		WorkspaceID: row.WorkspaceID,
		Name:        row.Name,
		TokenHash:   row.TokenHash,
		CreatedAt:   createdAt,
		LastUsedAt:  lastUsedAt,
		RevokedAt:   revokedAt,
	}, nil
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, value)
}

func parseOptionalTime(value sql.NullString) (*time.Time, error) {
	if !value.Valid {
		return nil, nil
	}
	parsed, err := parseTime(value.String)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
