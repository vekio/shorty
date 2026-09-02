package sqlite

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

func TestLinkRepositoryPersistsAcrossConnections(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shorty.db")
	database, err := Open(t.Context(), path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	link, err := domain.NewLink("https://example.com/persisted")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	workspaceID := testWorkspaceID(t, database)
	if err := NewLinkRepository(database, workspaceID).Save(t.Context(), link); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close first connection: %v", err)
	}

	database, err = Open(t.Context(), path)
	if err != nil {
		t.Fatalf("reopen database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	stored, err := NewLinkRepository(database, workspaceID).FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("FindByCode() error = %v", err)
	}
	storedURL := stored.OriginURL()
	if stored.Code() != link.Code() || storedURL.String() != "https://example.com/persisted" {
		t.Errorf("stored link = %#v", stored)
	}
}

func TestLinkRepositorySupportsPaginationVisitsAndDeletion(t *testing.T) {
	repository := openTestRepository(t)
	first := newTestLink(t, "https://example.com/first")
	second := newTestLink(t, "https://example.com/second")
	for _, link := range []domain.Link{first, second} {
		if err := repository.Save(t.Context(), link); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	page, err := repository.FindPage(t.Context(), 1, 1)
	if err != nil {
		t.Fatalf("FindPage() error = %v", err)
	}
	if page.Total != 2 || len(page.Links) != 1 || page.Links[0].Code() != first.Code() {
		t.Errorf("FindPage() = %#v", page)
	}

	updated, err := repository.ResolveByCode(t.Context(), first.Code())
	if err != nil {
		t.Fatalf("ResolveByCode() error = %v", err)
	}
	if updated.Visits() != 1 {
		t.Errorf("updated link = %#v, error = %v", updated, err)
	}
	if err := updated.ChangeOriginURL("https://example.com/updated"); err != nil {
		t.Fatalf("ChangeOriginURL() error = %v", err)
	}
	if err := repository.UpdateLinkOrigin(t.Context(), updated); err != nil {
		t.Fatalf("UpdateLinkOrigin() error = %v", err)
	}
	stored, err := repository.FindByCode(t.Context(), first.Code())
	if err != nil {
		t.Fatalf("FindByCode() after update error = %v", err)
	}
	storedURL := stored.OriginURL()
	if storedURL.String() != "https://example.com/updated" || stored.Visits() != 1 {
		t.Errorf("stored link after update = %#v", stored)
	}

	if err := repository.Delete(t.Context(), first.Code()); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repository.FindByCode(t.Context(), first.Code()); !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("FindByCode() after delete error = %v, want %v", err, ports.ErrLinkNotFound)
	}
}

func TestMissingUpdatesAndDeletesReturnNotFound(t *testing.T) {
	repository := openTestRepository(t)
	if _, err := repository.ResolveByCode(t.Context(), "missing"); !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("ResolveByCode() error = %v, want %v", err, ports.ErrLinkNotFound)
	}
	if err := repository.Delete(t.Context(), "missing"); !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("Delete() error = %v, want %v", err, ports.ErrLinkNotFound)
	}
	missing := newTestLink(t, "https://example.com")
	if err := repository.UpdateLinkOrigin(t.Context(), missing); !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("UpdateLinkOrigin() error = %v, want %v", err, ports.ErrLinkNotFound)
	}
}

func TestResolveByCodeDoesNotLoseConcurrentVisits(t *testing.T) {
	repository := openTestRepository(t)
	link := newTestLink(t, "https://example.com/concurrent")
	if err := repository.Save(t.Context(), link); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	const visits = 32
	errorsByVisit := make(chan error, visits)
	var waitGroup sync.WaitGroup
	for range visits {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			_, err := repository.ResolveByCode(t.Context(), link.Code())
			errorsByVisit <- err
		}()
	}
	waitGroup.Wait()
	close(errorsByVisit)
	for err := range errorsByVisit {
		if err != nil {
			t.Fatalf("ResolveByCode() error = %v", err)
		}
	}

	stored, err := repository.FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("FindByCode() error = %v", err)
	}
	if stored.Visits() != visits {
		t.Errorf("visits = %d, want %d", stored.Visits(), visits)
	}
}

func openTestRepository(t *testing.T) *LinkRepository {
	t.Helper()
	database, err := Open(t.Context(), filepath.Join(t.TempDir(), "shorty.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	return NewLinkRepository(database, testWorkspaceID(t, database))
}

func testWorkspaceID(t *testing.T, database *sql.DB) string {
	t.Helper()
	var workspaceID string
	if err := database.QueryRowContext(
		t.Context(), "SELECT id FROM workspaces WHERE name = 'default'",
	).Scan(&workspaceID); err != nil {
		t.Fatalf("load default workspace: %v", err)
	}
	if !strings.HasPrefix(workspaceID, "ws_") || len(workspaceID) != len("ws_")+32 {
		t.Fatalf("default workspace ID = %q, want random ws_ identifier", workspaceID)
	}
	return workspaceID
}

func newTestLink(t *testing.T, originURL string) domain.Link {
	t.Helper()
	link, err := domain.NewLink(originURL)
	if err != nil {
		t.Fatalf("NewLink() error = %v", err)
	}
	return link
}
