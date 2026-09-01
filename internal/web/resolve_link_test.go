package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

func TestResolveLinkMapsAPINotFound(t *testing.T) {
	api := &linkAPIStub{resolveErr: &shortysdk.ProblemError{
		Status: http.StatusNotFound,
		Detail: "link not found",
	}}
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/r/missing", nil),
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusNotFound, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want HTML", got)
	}
	if body := response.Body.String(); !strings.Contains(body, "Short code not found") ||
		!strings.Contains(body, "missing") {
		t.Errorf("body = %s, want missing-code explanation", body)
	}
}

func TestResolveLinkKeepsUnknownAPIErrorsPrivate(t *testing.T) {
	api := &linkAPIStub{resolveErr: errors.New("private upstream failure")}
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/r/abc123", nil),
	)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusInternalServerError, response.Body.String())
	}
	if strings.Contains(response.Body.String(), "private upstream failure") {
		t.Errorf("response exposes private error: %s", response.Body.String())
	}
}

func TestResolveLinkRedirectsBrowserToFinalDestination(t *testing.T) {
	api := &linkAPIStub{resolveDestination: "https://example.com/docs"}
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/r/abc123", nil),
	)

	if response.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusFound, response.Body.String())
	}
	if got := response.Header().Get("Location"); got != "https://example.com/docs" {
		t.Errorf("Location = %q, want final destination", got)
	}
	if api.resolvedCode != "abc123" {
		t.Errorf("resolved code = %q, want abc123", api.resolvedCode)
	}
	if response.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", response.Body.String())
	}
}

func TestSingleSegmentPathIsNotTreatedAsShortCode(t *testing.T) {
	api := &linkAPIStub{resolveDestination: "https://example.com/docs"}
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/abc123", nil),
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusNotFound, response.Body.String())
	}
	if api.resolvedCode != "" {
		t.Errorf("unexpectedly resolved code %q", api.resolvedCode)
	}
	if body := response.Body.String(); !strings.Contains(body, "Page not found") ||
		!strings.Contains(body, "/abc123") {
		t.Errorf("body = %s, want missing-page explanation", body)
	}
}
