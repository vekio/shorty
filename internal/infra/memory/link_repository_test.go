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
	link, err := domain.NewLink("https://example.com/docs")
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

func TestFindPageReturnsRequestedNewestFirstPage(t *testing.T) {
	repository := NewLinkRepository()
	first, err := domain.NewLink("https://example.com/first")
	if err != nil {
		t.Fatalf("create first link: %v", err)
	}
	second, err := domain.NewLink("https://example.com/second")
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
	if len(page.Links) != 1 || page.Links[0].Code() != first.Code() {
		t.Errorf("FindPage() links = %#v, want oldest link on second page", page.Links)
	}

	empty, err := repository.FindPage(t.Context(), 1, 10)
	if err != nil {
		t.Fatalf("FindPage() beyond total error = %v", err)
	}
	if empty.Total != 2 || empty.Links == nil || len(empty.Links) != 0 {
		t.Errorf("FindPage() beyond total = %#v, want non-nil empty page with total 2", empty)
	}
}

func TestResolveByCodeReturnsNotFound(t *testing.T) {
	if _, err := NewLinkRepository().ResolveByCode(t.Context(), "missing"); !errors.Is(err, ports.ErrLinkNotFound) {
		t.Errorf("ResolveByCode() error = %v, want %v", err, ports.ErrLinkNotFound)
	}
}

func TestDeleteRemovesLink(t *testing.T) {
	repository := NewLinkRepository()
	link, err := domain.NewLink("https://example.com")
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
	link, err := domain.NewLink("https://example.com")
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
	if _, err := repository.ResolveByCode(ctx, link.Code()); !errors.Is(err, context.Canceled) {
		t.Errorf("ResolveByCode() error = %v, want context canceled", err)
	}
	if err := repository.Delete(ctx, link.Code()); !errors.Is(err, context.Canceled) {
		t.Errorf("Delete() error = %v, want context canceled", err)
	}
}

func TestResolveByCodeAtomicallyRegistersVisit(t *testing.T) {
	repository := NewLinkRepository()
	link, err := domain.NewLink("https://example.com/docs")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	if err := repository.Save(t.Context(), link); err != nil {
		t.Fatalf("save link: %v", err)
	}

	resolved, err := repository.ResolveByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("ResolveByCode() error = %v", err)
	}
	updated, err := repository.FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("find updated link: %v", err)
	}
	if resolved.Visits() != 1 || updated.Visits() != 1 {
		t.Errorf("resolved visits = %d, stored visits = %d, want 1", resolved.Visits(), updated.Visits())
	}
}

func TestUpdateLinkOriginPreservesVisitsRegisteredAfterRead(t *testing.T) {
	repository := NewLinkRepository()
	link, err := domain.NewLink("https://example.com/old")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	if err := repository.Save(t.Context(), link); err != nil {
		t.Fatalf("save link: %v", err)
	}

	stale, err := repository.FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("find link before visit: %v", err)
	}
	if _, err := repository.ResolveByCode(t.Context(), link.Code()); err != nil {
		t.Fatalf("resolve link: %v", err)
	}
	if err := stale.ChangeOriginURL("https://example.com/new"); err != nil {
		t.Fatalf("change origin URL: %v", err)
	}
	if err := repository.UpdateLinkOrigin(t.Context(), stale); err != nil {
		t.Fatalf("update link origin: %v", err)
	}

	stored, err := repository.FindByCode(t.Context(), link.Code())
	if err != nil {
		t.Fatalf("find updated link: %v", err)
	}
	storedURL := stored.OriginURL()
	if storedURL.String() != "https://example.com/new" || stored.Visits() != 1 {
		t.Errorf("stored link = %#v, want updated URL and one visit", stored)
	}
}
