package getlink

import (
	"context"
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/domain"
)

type repositoryStub struct {
	link     domain.Link
	ownerID  string
	findCode string
	findErr  error
}

func (repository *repositoryStub) FindOwnedByCode(_ context.Context, ownerID string, code string) (domain.Link, error) {
	repository.ownerID = ownerID
	repository.findCode = code
	return repository.link, repository.findErr
}

func TestHandleReturnsLinkWithoutRegisteringVisit(t *testing.T) {
	link, err := domain.NewLink("https://example.com/docs")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	repository := &repositoryStub{link: link}
	result, err := NewGetLinkHandler(repository).Handle(t.Context(), GetLinkQuery{
		OwnerID: "browser-a",
		Code:    link.Code(),
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if repository.findCode != link.Code() {
		t.Errorf("FindByCode() code = %q, want %q", repository.findCode, link.Code())
	}
	if repository.ownerID != "browser-a" {
		t.Errorf("FindOwnedByCode() owner = %q, want browser-a", repository.ownerID)
	}
	originURL := link.OriginURL()
	if result.Code != link.Code() || result.OriginURL != originURL.String() || !result.CreatedAt.Equal(link.CreatedAt()) || result.Visits != 0 {
		t.Errorf("Handle() result = %#v", result)
	}
	if link.Visits() != 0 {
		t.Errorf("link visits = %d, want 0", link.Visits())
	}
}

func TestHandleReturnsFindError(t *testing.T) {
	wantErr := errors.New("find failed")
	_, err := NewGetLinkHandler(&repositoryStub{findErr: wantErr}).Handle(t.Context(), GetLinkQuery{Code: "missing"})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
