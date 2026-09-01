package deletelink

import (
	"context"
	"errors"
	"testing"
)

type repositoryStub struct {
	ownerID     string
	deletedCode string
	err         error
}

func (repository *repositoryStub) Delete(_ context.Context, ownerID string, code string) error {
	repository.ownerID = ownerID
	repository.deletedCode = code
	return repository.err
}

func TestHandleDeletesLinkByCode(t *testing.T) {
	repository := &repositoryStub{}
	_, err := NewDeleteLinkHandler(repository).Handle(t.Context(), DeleteLinkCommand{
		OwnerID: "browser-a",
		Code:    "abc123",
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if repository.deletedCode != "abc123" {
		t.Errorf("Delete() code = %q, want abc123", repository.deletedCode)
	}
	if repository.ownerID != "browser-a" {
		t.Errorf("Delete() owner = %q, want browser-a", repository.ownerID)
	}
}

func TestHandleReturnsDeleteError(t *testing.T) {
	wantErr := errors.New("delete failed")
	_, err := NewDeleteLinkHandler(&repositoryStub{err: wantErr}).Handle(t.Context(), DeleteLinkCommand{Code: "abc123"})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
