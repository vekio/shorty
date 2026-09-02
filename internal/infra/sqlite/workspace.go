package sqlite

import (
	"context"
	"database/sql"

	"github.com/vekio/shorty/internal/infra/sqlite/sqlitedb"
)

// FindWorkspaceByName returns the persisted identifier and name of a workspace.
func FindWorkspaceByName(ctx context.Context, database *sql.DB, name string) (string, string, error) {
	row, err := sqlitedb.New(database).GetWorkspaceByName(ctx, name)
	return row.ID, row.Name, err
}
