package updatelink

import (
	"context"
	"testing"

	"github.com/vekio/shorty/internal/domain"
)

type repositoryStub struct {
	link    domain.Link
	updated domain.Link
}

func (repository *repositoryStub) FindByCode(context.Context, string) (domain.Link, error) {
	return repository.link, nil
}

func (repository *repositoryStub) UpdateLinkOrigin(_ context.Context, link domain.Link) error {
	repository.updated = link
	return nil
}

func TestHandleChangesAndPersistsOriginURL(t *testing.T) {
	link, err := domain.NewLink("https://example.com/old")
	if err != nil {
		t.Fatalf("NewLink() error = %v", err)
	}
	repository := &repositoryStub{link: link}
	result, err := NewUpdateLinkHandler(repository).Handle(t.Context(), UpdateLinkCommand{
		Code: link.Code(), OriginURL: "https://example.com/new",
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	updatedURL := repository.updated.OriginURL()
	if result.OriginURL != "https://example.com/new" || updatedURL.String() != result.OriginURL {
		t.Errorf("result = %#v, persisted URL = %q", result, updatedURL.String())
	}
}
