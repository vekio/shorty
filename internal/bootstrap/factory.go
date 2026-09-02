package bootstrap

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/auth"
	shortyconfig "github.com/vekio/shorty/internal/config"
	"github.com/vekio/shorty/internal/domain"
	"github.com/vekio/shorty/internal/infra/memory"
	"github.com/vekio/shorty/internal/infra/sqlite"
)

// NewLogger creates the process logger using its configured output format.
func NewLogger(config shortyconfig.LoggerConfig) (*slog.Logger, error) {
	return newLogger(config, os.Stdout)
}

func newLogger(config shortyconfig.LoggerConfig, output io.Writer) (*slog.Logger, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}
	options := &slog.HandlerOptions{Level: level, AddSource: config.AddSource}

	var handler slog.Handler
	switch config.Format {
	case "", shortyconfig.LoggerFormatJSON:
		handler = slog.NewJSONHandler(output, options)
	case shortyconfig.LoggerFormatText:
		handler = slog.NewTextHandler(output, options)
	default:
		return nil, fmt.Errorf("invalid logger format %q", config.Format)
	}
	return slog.New(handler), nil
}

// Persistence groups the repositories backed by the selected database driver.
// New repositories can be added here without changing application composition.
type Persistence struct {
	WorkspaceID   string
	WorkspaceName string
	Links         ports.LinkRepository
	APIKeys       auth.Repository
	close         func() error
}

// Close releases resources owned by the selected persistence driver.
func (persistence Persistence) Close() error {
	if persistence.close == nil {
		return nil
	}
	return persistence.close()
}

// OpenPersistence selects a driver, opens it, and constructs its repositories.
func OpenPersistence(ctx context.Context, config shortyconfig.DatabaseConfig) (Persistence, error) {
	switch config.Driver {
	case shortyconfig.DatabaseDriverMemory:
		workspace, err := domain.NewWorkspace("default")
		if err != nil {
			return Persistence{}, err
		}
		return Persistence{
			WorkspaceID:   workspace.ID(),
			WorkspaceName: workspace.Name(),
			Links:         memory.NewLinkRepository(),
			APIKeys:       memory.NewAPIKeyRepository(),
			close:         func() error { return nil },
		}, nil
	case shortyconfig.DatabaseDriverSQLite:
		database, err := sqlite.Open(ctx, config.Path)
		if err != nil {
			return Persistence{}, err
		}
		workspaceID, workspaceName, err := sqlite.FindWorkspaceByName(ctx, database, "default")
		if err != nil {
			_ = database.Close()
			return Persistence{}, err
		}
		return Persistence{
			WorkspaceID:   workspaceID,
			WorkspaceName: workspaceName,
			Links:         sqlite.NewLinkRepository(database, workspaceID),
			APIKeys:       sqlite.NewAPIKeyRepository(database),
			close:         database.Close,
		}, nil
	default:
		return Persistence{}, fmt.Errorf("unsupported database driver %q", config.Driver)
	}
}
