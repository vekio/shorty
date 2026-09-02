package redirect

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

type resolveLinkStub struct {
	code string
}

func (stub *resolveLinkStub) Handle(_ context.Context, command resolvelink.ResolveLinkCommand) (resolvelink.ResolveLinkResult, error) {
	stub.code = command.Code
	return resolvelink.ResolveLinkResult{OriginURL: "https://example.com/docs"}, nil
}

func TestNewRedirectsPublicNavigation(t *testing.T) {
	stub := &resolveLinkStub{}
	handler := New(
		app.Application{Commands: app.Commands{ResolveLink: stub}},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/r/abc123", nil))
	if response.Code != http.StatusFound || response.Header().Get("Location") != "https://example.com/docs" {
		t.Fatalf("response = %d Location %q", response.Code, response.Header().Get("Location"))
	}
	if stub.code != "abc123" {
		t.Errorf("resolved code = %q, want abc123", stub.code)
	}
}
