package listlinks

import (
	"context"
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/domain"
)

type repositoryStub struct {
	links   []domain.Link
	findErr error
}

func (*repositoryStub) Save(context.Context, domain.Link) error {
	panic("Save must not be called by ListLinks")
}
func (*repositoryStub) FindByCode(context.Context, string) (domain.Link, error) {
	panic("FindByCode must not be called by ListLinks")
}
func (repository *repositoryStub) FindAll(context.Context) ([]domain.Link, error) {
	return repository.links, repository.findErr
}
func (*repositoryStub) UpdateLinkVisits(context.Context, domain.Link) error {
	panic("UpdateLinkVisits must not be called by ListLinks")
}

func TestHandleReturnsLinksOrderedByCreation(t *testing.T) {
	first, err := domain.New("https://example.com/first")
	if err != nil {
		t.Fatalf("create first link: %v", err)
	}
	second, err := domain.New("https://example.com/second")
	if err != nil {
		t.Fatalf("create second link: %v", err)
	}
	second.RegisterVisit()

	result, err := NewListLinksHandler(&repositoryStub{links: []domain.Link{second, first}}).Handle(t.Context(), ListLinksQuery{})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if len(result.Links) != 2 {
		t.Fatalf("links length = %d, want 2", len(result.Links))
	}
	if result.Links[0].Code != first.Code() || result.Links[0].OriginURL != "https://example.com/first" {
		t.Errorf("first result = %#v", result.Links[0])
	}
	if result.Links[1].Code != second.Code() || result.Links[1].Visits != 1 {
		t.Errorf("second result = %#v", result.Links[1])
	}
}

func TestHandleReturnsEmptyCollection(t *testing.T) {
	result, err := NewListLinksHandler(&repositoryStub{}).Handle(t.Context(), ListLinksQuery{})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.Links == nil || len(result.Links) != 0 {
		t.Errorf("links = %#v, want non-nil empty collection", result.Links)
	}
}

func TestHandleReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("find all failed")
	_, err := NewListLinksHandler(&repositoryStub{findErr: wantErr}).Handle(t.Context(), ListLinksQuery{})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
