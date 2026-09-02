package resolvelink

import (
	"context"
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/domain"
)

type repositoryStub struct {
	link        domain.Link
	resolveCode string
	resolveErr  error
}

func (repository *repositoryStub) ResolveByCode(_ context.Context, code string) (domain.Link, error) {
	repository.resolveCode = code
	return repository.link, repository.resolveErr
}

func TestHandleResolvesDestinationAndRegistersVisit(t *testing.T) {
	link, err := domain.NewLink("https://example.com/docs")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	repository := &repositoryStub{link: link}

	result, err := NewResolveLinkHandler(repository).Handle(t.Context(), ResolveLinkCommand{Code: link.Code()})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if repository.resolveCode != link.Code() {
		t.Errorf("ResolveByCode() code = %q, want %q", repository.resolveCode, link.Code())
	}
	if result.OriginURL != "https://example.com/docs" {
		t.Errorf("OriginURL = %q, want https://example.com/docs", result.OriginURL)
	}
}

func TestHandleReturnsResolveError(t *testing.T) {
	wantErr := errors.New("resolve failed")
	repository := &repositoryStub{resolveErr: wantErr}

	_, err := NewResolveLinkHandler(repository).Handle(t.Context(), ResolveLinkCommand{Code: "missing"})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
