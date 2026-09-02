package deletelink

import (
	"context"
	"errors"
	"testing"
)

type repositoryStub struct {
	deletedCode string
	err         error
}

func (repository *repositoryStub) Delete(_ context.Context, code string) error {
	repository.deletedCode = code
	return repository.err
}

func TestHandleDeletesLinkByCode(t *testing.T) {
	repository := &repositoryStub{}
	_, err := NewDeleteLinkHandler(repository).Handle(t.Context(), DeleteLinkCommand{
		Code: "abc123",
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if repository.deletedCode != "abc123" {
		t.Errorf("Delete() code = %q, want abc123", repository.deletedCode)
	}
}

func TestHandleReturnsDeleteError(t *testing.T) {
	wantErr := errors.New("delete failed")
	_, err := NewDeleteLinkHandler(&repositoryStub{err: wantErr}).Handle(t.Context(), DeleteLinkCommand{Code: "abc123"})
	if !errors.Is(err, wantErr) {
		t.Errorf("Handle() error = %v, want %v", err, wantErr)
	}
}
