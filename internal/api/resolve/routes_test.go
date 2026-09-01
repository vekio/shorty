package resolve

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

type resolveLinkStub struct {
	code string
}

func (stub *resolveLinkStub) Handle(
	_ context.Context,
	command resolvelink.ResolveLinkCommand,
) (resolvelink.ResolveLinkResult, error) {
	stub.code = command.Code
	return resolvelink.ResolveLinkResult{OriginURL: "https://example.com/docs"}, nil
}

func TestResolveLinkWritesRedirectFromTypedOutput(t *testing.T) {
	stub := &resolveLinkStub{}
	httpAPI := amigo.New(amigo.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	Register(httpAPI.Group("/r"), app.Application{Commands: app.Commands{ResolveLink: stub}})
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/r/abc123", nil))

	if response.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusFound, response.Body.String())
	}
	if stub.code != "abc123" || response.Header().Get("Location") != "https://example.com/docs" {
		t.Errorf("code = %q, Location = %q", stub.code, response.Header().Get("Location"))
	}
	if response.Body.Len() != 0 || response.Header().Get("Content-Type") != "" {
		t.Errorf("headers-only redirect has body %q and Content-Type %q", response.Body.String(), response.Header().Get("Content-Type"))
	}
}

func TestRedirectIsNotExposedAtAPIRoot(t *testing.T) {
	stub := &resolveLinkStub{}
	httpAPI := amigo.New(amigo.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	Register(httpAPI.Group("/r"), app.Application{Commands: app.Commands{ResolveLink: stub}})
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/abc123", nil))

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusNotFound, response.Body.String())
	}
	if stub.code != "" {
		t.Errorf("root request unexpectedly resolved code %q", stub.code)
	}
}
