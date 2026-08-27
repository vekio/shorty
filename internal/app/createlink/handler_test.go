package createlink

import (
	"context"
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/domain"
)

type repositoryStub struct {
	saved   domain.Link
	saveErr error
}

func (repository *repositoryStub) Save(_ context.Context, link domain.Link) error {
	repository.saved = link
	return repository.saveErr
}
func (*repositoryStub) FindByCode(context.Context, string) (domain.Link, error) {
	panic("FindByCode must not be called by CreateLink")
}
func (*repositoryStub) FindAll(context.Context) ([]domain.Link, error) {
	panic("FindAll must not be called by CreateLink")
}
func (*repositoryStub) UpdateLinkVisits(context.Context, domain.Link) error {
	panic("UpdateLinkVisits must not be called by CreateLink")
}

func TestHandleCreatesAndSavesLink(t *testing.T) {
	repository := &repositoryStub{}
	result, err := NewCreateLinkHandler(repository).Handle(t.Context(), CreateLinkCommand{URL: "https://example.com/docs"})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.Code == "" || result.Code != repository.saved.Code() {
		t.Errorf("result code = %q, saved code = %q", result.Code, repository.saved.Code())
	}
	originURL := repository.saved.OriginURL()
	if got := originURL.String(); got != "https://example.com/docs" {
		t.Errorf("saved origin URL = %q", got)
	}
}

func TestHandleReturnsDomainErrorWithoutSaving(t *testing.T) {
	repository := &repositoryStub{}
	_, err := NewCreateLinkHandler(repository).Handle(t.Context(), CreateLinkCommand{})
	if !errors.Is(err, domain.ErrOriginURLRequired) {
		t.Errorf("Handle() error = %v, want %v", err, domain.ErrOriginURLRequired)
	}
	if repository.saved.Code() != "" {
		t.Error("repository saved an invalid link")
	}
}

func TestHandleReturnsSaveError(t *testing.T) {
	wantErr := errors.New("save failed")
	_, err := NewCreateLinkHandler(&repositoryStub{saveErr: wantErr}).Handle(t.Context(), CreateLinkCommand{URL: "https://example.com"})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
