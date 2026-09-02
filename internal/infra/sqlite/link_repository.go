package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
	"github.com/vekio/shorty/internal/infra/sqlite/sqlitedb"
)

type LinkRepository struct {
	workspaceID string
	database    *sql.DB
	queries     *sqlitedb.Queries
}

var _ ports.LinkRepository = (*LinkRepository)(nil)

func NewLinkRepository(database *sql.DB, workspaceID string) *LinkRepository {
	return &LinkRepository{workspaceID: workspaceID, database: database, queries: sqlitedb.New(database)}
}

func (repository *LinkRepository) Save(ctx context.Context, link domain.Link) error {
	originURL := link.OriginURL()
	return repository.queries.CreateLink(ctx, sqlitedb.CreateLinkParams{
		Code:        link.Code(),
		WorkspaceID: repository.workspaceID,
		OriginUrl:   originURL.String(),
		CreatedAt:   link.CreatedAt().UTC().Format(time.RFC3339Nano),
		Visits:      int64(link.Visits()),
	})
}

func (repository *LinkRepository) FindByCode(ctx context.Context, code string) (domain.Link, error) {
	row, err := repository.queries.GetLinkByCode(ctx, sqlitedb.GetLinkByCodeParams{
		WorkspaceID: repository.workspaceID,
		Code:        code,
	})
	if err != nil {
		return domain.Link{}, mapNotFound(err)
	}
	return restoreLink(row)
}

func (repository *LinkRepository) FindPage(ctx context.Context, limit int, offset int) (ports.LinkPage, error) {
	tx, err := repository.database.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return ports.LinkPage{}, err
	}
	defer func() { _ = tx.Rollback() }()

	queries := repository.queries.WithTx(tx)
	total, err := queries.CountLinks(ctx, repository.workspaceID)
	if err != nil {
		return ports.LinkPage{}, err
	}
	rows, err := queries.ListLinks(ctx, sqlitedb.ListLinksParams{
		WorkspaceID: repository.workspaceID,
		Limit:       int64(limit),
		Offset:      int64(offset),
	})
	if err != nil {
		return ports.LinkPage{}, err
	}
	links := make([]domain.Link, 0, len(rows))
	for _, row := range rows {
		link, err := restoreLink(row)
		if err != nil {
			return ports.LinkPage{}, err
		}
		links = append(links, link)
	}
	if err := tx.Commit(); err != nil {
		return ports.LinkPage{}, err
	}
	return ports.LinkPage{Links: links, Total: int(total)}, nil
}

func (repository *LinkRepository) UpdateLinkOrigin(ctx context.Context, link domain.Link) error {
	originURL := link.OriginURL()
	updated, err := repository.queries.UpdateLinkOrigin(ctx, sqlitedb.UpdateLinkOriginParams{
		OriginUrl:   originURL.String(),
		WorkspaceID: repository.workspaceID,
		Code:        link.Code(),
	})
	if err != nil {
		return err
	}
	if updated == 0 {
		return ports.ErrLinkNotFound
	}
	return nil
}

func (repository *LinkRepository) ResolveByCode(ctx context.Context, code string) (domain.Link, error) {
	row, err := repository.queries.ResolveLink(ctx, code)
	if err != nil {
		return domain.Link{}, mapNotFound(err)
	}
	return restoreLink(row)
}

func (repository *LinkRepository) Delete(ctx context.Context, code string) error {
	deleted, err := repository.queries.DeleteLink(ctx, sqlitedb.DeleteLinkParams{
		WorkspaceID: repository.workspaceID,
		Code:        code,
	})
	if err != nil {
		return err
	}
	if deleted == 0 {
		return ports.ErrLinkNotFound
	}
	return nil
}

func restoreLink(row sqlitedb.Link) (domain.Link, error) {
	createdAt, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
	if err != nil {
		return domain.Link{}, fmt.Errorf("parse persisted link creation time: %w", err)
	}
	link, err := domain.RestoreLink(row.Code, row.OriginUrl, createdAt, int(row.Visits))
	if err != nil {
		return domain.Link{}, fmt.Errorf("restore persisted link: %w", err)
	}
	return link, nil
}

func mapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ports.ErrLinkNotFound
	}
	return err
}
