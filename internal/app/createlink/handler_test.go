package createlink

import (
	"context"
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/domain"
)

type repositoryStub struct {
	ownerID string
	saved   domain.Link
	saveErr error
}

func (repository *repositoryStub) Save(_ context.Context, ownerID string, link domain.Link) error {
	repository.ownerID = ownerID
	repository.saved = link
	return repository.saveErr
}

func TestHandleCreatesAndSavesLink(t *testing.T) {
	repository := &repositoryStub{}
	result, err := NewCreateLinkHandler(repository).Handle(t.Context(), CreateLinkCommand{
		OwnerID:   "browser-a",
		OriginURL: "https://example.com/docs",
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.Code == "" || result.Code != repository.saved.Code() {
		t.Errorf("result code = %q, saved code = %q", result.Code, repository.saved.Code())
	}
	if repository.ownerID != "browser-a" {
		t.Errorf("saved owner = %q, want browser-a", repository.ownerID)
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

func TestHandleAppliesOriginURLPoliciesBeforeSaving(t *testing.T) {
	repository := &repositoryStub{}
	policy, err := domain.DisallowOriginHost("https://sho.rt")
	if err != nil {
		t.Fatalf("DisallowOriginHost() error = %v", err)
	}

	_, err = NewCreateLinkHandler(repository, policy).Handle(t.Context(), CreateLinkCommand{
		OwnerID:   "browser-a",
		OriginURL: "https://sho.rt/r/existing",
	})

	if !errors.Is(err, domain.ErrOriginURLSelfReference) {
		t.Errorf("Handle() error = %v, want %v", err, domain.ErrOriginURLSelfReference)
	}
	if repository.saved.Code() != "" {
		t.Error("repository saved a self-referencing link")
	}
}

func TestHandleReturnsSaveError(t *testing.T) {
	wantErr := errors.New("save failed")
	_, err := NewCreateLinkHandler(&repositoryStub{saveErr: wantErr}).Handle(t.Context(), CreateLinkCommand{OriginURL: "https://example.com"})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
