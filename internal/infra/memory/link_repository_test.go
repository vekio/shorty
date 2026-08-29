package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

func TestSaveAndFindByCode(t *testing.T) {
	repository := NewLinkRepository()
	link, err := domain.New("https://example.com/docs")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}

	if err := repository.Save(t.Context(), link); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, err := repository.FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("FindByCode() error = %v", err)
	}
	gotURL := got.OriginURL()
	wantURL := link.OriginURL()
	if got.Code() != link.Code() || gotURL.String() != wantURL.String() {
		t.Errorf("FindByCode() = %#v, want %#v", got, link)
	}
}

func TestFindByCodeReturnsNotFound(t *testing.T) {
	_, err := NewLinkRepository().FindByCode(t.Context(), "missing")
	if !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("FindByCode() error = %v, want %v", err, ports.ErrLinkNotFound)
	}
}

func TestFindPageReturnsRequestedCreationOrderedPage(t *testing.T) {
	repository := NewLinkRepository()
	first, err := domain.New("https://example.com/first")
	if err != nil {
		t.Fatalf("create first link: %v", err)
	}
	second, err := domain.New("https://example.com/second")
	if err != nil {
		t.Fatalf("create second link: %v", err)
	}
	if err := repository.Save(t.Context(), first); err != nil {
		t.Fatalf("save first link: %v", err)
	}
	if err := repository.Save(t.Context(), second); err != nil {
		t.Fatalf("save second link: %v", err)
	}

	page, err := repository.FindPage(t.Context(), 1, 1)
	if err != nil {
		t.Fatalf("FindPage() error = %v", err)
	}
	if page.Total != 2 {
		t.Errorf("FindPage() total = %d, want 2", page.Total)
	}
	if len(page.Links) != 1 || page.Links[0].Code() != second.Code() {
		t.Errorf("FindPage() links = %#v, want second link", page.Links)
	}

	empty, err := repository.FindPage(t.Context(), 1, 10)
	if err != nil {
		t.Fatalf("FindPage() beyond total error = %v", err)
	}
	if empty.Total != 2 || empty.Links == nil || len(empty.Links) != 0 {
		t.Errorf("FindPage() beyond total = %#v, want non-nil empty page with total 2", empty)
	}
}

func TestUpdateLinkVisitsReturnsNotFound(t *testing.T) {
	link, err := domain.New("https://example.com")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	if err := NewLinkRepository().UpdateLinkVisits(t.Context(), link); !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("UpdateLinkVisits() error = %v, want %v", err, ports.ErrLinkNotFound)
	}
}

func TestDeleteRemovesLink(t *testing.T) {
	repository := NewLinkRepository()
	link, err := domain.New("https://example.com")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	if err := repository.Save(t.Context(), link); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if err := repository.Delete(t.Context(), link.Code()); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repository.FindByCode(t.Context(), link.Code()); !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("FindByCode() error = %v after deletion, want %v", err, ports.ErrLinkNotFound)
	}
}

func TestDeleteReturnsNotFound(t *testing.T) {
	if err := NewLinkRepository().Delete(t.Context(), "missing"); !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("Delete() error = %v, want %v", err, ports.ErrLinkNotFound)
	}
}

func TestOperationsRespectCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	repository := NewLinkRepository()
	link, err := domain.New("https://example.com")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}

	if err := repository.Save(ctx, link); !errors.Is(err, context.Canceled) {
		t.Errorf("Save() error = %v, want context canceled", err)
	}
	if _, err := repository.FindByCode(ctx, link.Code()); !errors.Is(err, context.Canceled) {
		t.Errorf("FindByCode() error = %v, want context canceled", err)
	}
	if _, err := repository.FindPage(ctx, 20, 0); !errors.Is(err, context.Canceled) {
		t.Errorf("FindPage() error = %v, want context canceled", err)
	}
	if err := repository.UpdateLinkVisits(ctx, link); !errors.Is(err, context.Canceled) {
		t.Errorf("UpdateLinkVisits() error = %v, want context canceled", err)
	}
	if err := repository.Delete(ctx, link.Code()); !errors.Is(err, context.Canceled) {
		t.Errorf("Delete() error = %v, want context canceled", err)
	}
}

func TestLinkValuesMustBeSavedAfterMutation(t *testing.T) {
	repository := NewLinkRepository()
	link, err := domain.New("https://example.com/docs")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	if err := repository.Save(t.Context(), link); err != nil {
		t.Fatalf("save link: %v", err)
	}

	stored, err := repository.FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("find link: %v", err)
	}
	stored.RegisterVisit()

	unchanged, err := repository.FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("find unchanged link: %v", err)
	}
	if unchanged.Visits() != 0 {
		t.Fatalf("stored visits = %d before saving the changed value, want 0", unchanged.Visits())
	}

	if err := repository.UpdateLinkVisits(t.Context(), stored); err != nil {
		t.Fatalf("save changed link: %v", err)
	}
	updated, err := repository.FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("find updated link: %v", err)
	}
	if updated.Visits() != 1 {
		t.Errorf("stored visits = %d, want 1", updated.Visits())
	}
}
