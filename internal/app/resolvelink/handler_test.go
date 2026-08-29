package resolvelink

import (
	"context"
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/domain"
)

type repositoryStub struct {
	link        domain.Link
	findCode    string
	findErr     error
	updated     domain.Link
	updateCalls int
	updateErr   error
}

func (repository *repositoryStub) FindByCode(_ context.Context, code string) (domain.Link, error) {
	repository.findCode = code
	return repository.link, repository.findErr
}

func (repository *repositoryStub) UpdateLinkVisits(_ context.Context, link domain.Link) error {
	repository.updated = link
	repository.updateCalls++
	return repository.updateErr
}

func TestHandleResolvesDestinationAndRegistersVisit(t *testing.T) {
	link, err := domain.New("https://example.com/docs")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	repository := &repositoryStub{link: link}

	result, err := NewResolveLinkHandler(repository).Handle(t.Context(), ResolveLinkCommand{Code: link.Code()})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if repository.findCode != link.Code() {
		t.Errorf("FindByCode() code = %q, want %q", repository.findCode, link.Code())
	}
	if result.OriginURL != "https://example.com/docs" {
		t.Errorf("OriginURL = %q, want https://example.com/docs", result.OriginURL)
	}
	if repository.updateCalls != 1 || repository.updated.Visits() != 1 {
		t.Errorf("UpdateLinkVisits() calls = %d, visits = %d, want one visit", repository.updateCalls, repository.updated.Visits())
	}
}

func TestHandleReturnsFindErrorWithoutUpdating(t *testing.T) {
	wantErr := errors.New("find failed")
	repository := &repositoryStub{findErr: wantErr}

	_, err := NewResolveLinkHandler(repository).Handle(t.Context(), ResolveLinkCommand{Code: "missing"})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
	if repository.updateCalls != 0 {
		t.Errorf("UpdateLinkVisits() calls = %d, want 0", repository.updateCalls)
	}
}

func TestHandleReturnsUpdateError(t *testing.T) {
	link, err := domain.New("https://example.com")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	wantErr := errors.New("update failed")

	_, err = NewResolveLinkHandler(&repositoryStub{link: link, updateErr: wantErr}).Handle(
		t.Context(),
		ResolveLinkCommand{Code: link.Code()},
	)
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
