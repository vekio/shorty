package listlinks

import (
	"context"
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

type repositoryStub struct {
	ownerID string
	page    ports.LinkPage
	limit   int
	offset  int
	calls   int
	err     error
}

func (repository *repositoryStub) FindPage(_ context.Context, ownerID string, limit int, offset int) (ports.LinkPage, error) {
	repository.ownerID = ownerID
	repository.limit = limit
	repository.offset = offset
	repository.calls++
	return repository.page, repository.err
}

func TestHandleReturnsRequestedPage(t *testing.T) {
	link, err := domain.NewLink("https://example.com/second")
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	link.RegisterVisit()
	repository := &repositoryStub{page: ports.LinkPage{
		Links: []domain.Link{link},
		Total: 2,
	}}

	result, err := NewListLinksHandler(repository).Handle(t.Context(), ListLinksQuery{
		OwnerID: "browser-a",
		Limit:   1,
		Offset:  1,
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if repository.limit != 1 || repository.offset != 1 {
		t.Errorf("FindPage() pagination = (%d, %d), want (1, 1)", repository.limit, repository.offset)
	}
	if repository.ownerID != "browser-a" {
		t.Errorf("FindPage() owner = %q, want browser-a", repository.ownerID)
	}
	if result.Total != 2 || result.Limit != 1 || result.Offset != 1 {
		t.Errorf("pagination = total %d, limit %d, offset %d", result.Total, result.Limit, result.Offset)
	}
	if len(result.Links) != 1 || result.Links[0].Code != link.Code() ||
		result.Links[0].OriginURL != "https://example.com/second" || result.Links[0].Visits != 1 {
		t.Errorf("links = %#v", result.Links)
	}
}

func TestHandleAppliesDefaultPagination(t *testing.T) {
	repository := &repositoryStub{page: ports.LinkPage{Links: []domain.Link{}, Total: 0}}
	result, err := NewListLinksHandler(repository).Handle(t.Context(), ListLinksQuery{})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if repository.limit != DefaultLimit || repository.offset != 0 {
		t.Errorf("FindPage() pagination = (%d, %d), want (%d, 0)", repository.limit, repository.offset, DefaultLimit)
	}
	if result.Links == nil || len(result.Links) != 0 {
		t.Errorf("links = %#v, want non-nil empty collection", result.Links)
	}
	if result.Total != 0 || result.Limit != DefaultLimit || result.Offset != 0 {
		t.Errorf("result pagination = %#v", result)
	}
}

func TestHandleRejectsInvalidPagination(t *testing.T) {
	tests := []struct {
		name    string
		query   ListLinksQuery
		wantErr error
	}{
		{name: "negative limit", query: ListLinksQuery{Limit: -1}, wantErr: ErrInvalidLimit},
		{name: "limit above maximum", query: ListLinksQuery{Limit: MaximumLimit + 1}, wantErr: ErrInvalidLimit},
		{name: "negative offset", query: ListLinksQuery{Limit: 1, Offset: -1}, wantErr: ErrInvalidOffset},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &repositoryStub{}
			_, err := NewListLinksHandler(repository).Handle(t.Context(), test.query)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("Handle() error = %v, want %v", err, test.wantErr)
			}
			if repository.calls != 0 {
				t.Errorf("FindPage() calls = %d, want 0", repository.calls)
			}
		})
	}
}

func TestHandleReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("find page failed")
	_, err := NewListLinksHandler(&repositoryStub{err: wantErr}).Handle(t.Context(), ListLinksQuery{})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
