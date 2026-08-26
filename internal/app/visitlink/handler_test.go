package visitlink

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

func (*repositoryStub) Save(context.Context, domain.Link) error {
	panic("Save must not be called by VisitLink")
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

func TestHandleRegistersAndSavesVisit(t *testing.T) {
	link, err := domain.New("https://example.com")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	repository := &repositoryStub{link: link}
	_, err = NewVisitLinkHandler(repository).Handle(t.Context(), VisitLinkCommand{Code: link.Code()})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if repository.findCode != link.Code() {
		t.Errorf("FindByCode() code = %q, want %q", repository.findCode, link.Code())
	}
	if repository.updateCalls != 1 || repository.updated.Visits() != 1 {
		t.Errorf("UpdateLinkVisits() calls = %d, visits = %d", repository.updateCalls, repository.updated.Visits())
	}
	if link.Visits() != 0 {
		t.Errorf("source link visits = %d, want 0", link.Visits())
	}
}

func TestHandleReturnsFindErrorWithoutUpdating(t *testing.T) {
	wantErr := errors.New("find failed")
	repository := &repositoryStub{findErr: wantErr}
	_, err := NewVisitLinkHandler(repository).Handle(t.Context(), VisitLinkCommand{Code: "missing"})
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
	_, err = NewVisitLinkHandler(&repositoryStub{link: link, updateErr: wantErr}).Handle(t.Context(), VisitLinkCommand{Code: link.Code()})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
