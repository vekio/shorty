package root

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

type problemResponse struct {
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func TestResolveLinkRedirectsToOriginURL(t *testing.T) {
	handler := &resolveLinkHandlerStub{
		result: resolvelink.ResolveLinkResult{OriginURL: "https://example.com/docs"},
	}
	httpAPI := newTestAPI(handler)
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/abc123", nil))

	if response.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusFound, response.Body.String())
	}
	if handler.code != "abc123" {
		t.Errorf("resolved code = %q, want abc123", handler.code)
	}
	if got := response.Header().Get("Location"); got != "https://example.com/docs" {
		t.Errorf("Location = %q, want https://example.com/docs", got)
	}
	if got := response.Header().Get("Content-Type"); got != "" {
		t.Errorf("Content-Type = %q, want no content type", got)
	}
	if response.Body.Len() != 0 {
		t.Errorf("body = %q, want empty response", response.Body.String())
	}
}

func TestResolveLinkMapsNotFoundError(t *testing.T) {
	httpAPI := newTestAPI(&resolveLinkHandlerStub{
		err: errors.Join(errors.New("repository detail"), ports.ErrLinkNotFound),
	})
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/missing", nil))

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusNotFound, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Errorf("Content-Type = %q, want application/problem+json", got)
	}
	var problem problemResponse
	if err := json.Unmarshal(response.Body.Bytes(), &problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Status != http.StatusNotFound || problem.Detail != "link not found" || problem.Instance != "/missing" {
		t.Errorf("problem = %#v", problem)
	}
}
