package web

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

func TestCreateLinkReturnsHTMXFragment(t *testing.T) {
	api := &linkAPIStub{createCode: "abc123"}
	form := url.Values{"origin_url": {"https://example.com/docs"}}
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("HX-Request", "true")
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if api.createdURL != "https://example.com/docs" {
		t.Errorf("origin URL = %q, want https://example.com/docs", api.createdURL)
	}
	if strings.Contains(response.Body.String(), "<!doctype html>") {
		t.Error("HTMX response contains the complete document")
	}
	if !strings.Contains(response.Body.String(), `id="workspace"`) ||
		!strings.Contains(response.Body.String(), "https://sho.rt/r/abc123") {
		t.Errorf("body = %s, want updated workspace", response.Body.String())
	}
}

func TestCreateLinkRendersKnownDomainError(t *testing.T) {
	api := &linkAPIStub{
		createErr: fmt.Errorf("request failed: %w", &shortysdk.ProblemError{
			Status: http.StatusUnprocessableEntity,
			Detail: "request validation failed",
			Errors: []shortysdk.FieldError{{
				Location: "body.origin_url",
				Message:  "origin URL must be an absolute HTTP or HTTPS URL",
			}},
		}),
	}
	form := url.Values{"origin_url": {"not-a-url"}}
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "origin URL must be an absolute HTTP or HTTPS URL") {
		t.Errorf("body = %s, want public validation detail", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `value="not-a-url"`) {
		t.Errorf("body = %s, want submitted URL", response.Body.String())
	}
}

func TestCreateLinkRendersRequiredDomainError(t *testing.T) {
	api := &linkAPIStub{
		createErr: &shortysdk.ProblemError{
			Status: http.StatusUnprocessableEntity,
			Detail: "request validation failed",
			Errors: []shortysdk.FieldError{{Location: "body.origin_url", Message: "is required"}},
		},
	}
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader("origin_url="))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "is required") {
		t.Errorf("body = %s, want public required detail", response.Body.String())
	}
}

func TestCreateLinkReturnsSafeInternalServerError(t *testing.T) {
	api := &linkAPIStub{createErr: errors.New("private API failure")}
	form := url.Values{"origin_url": {"https://example.com"}}
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusInternalServerError, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want net/http text response", got)
	}
	if strings.Contains(response.Body.String(), "private API failure") {
		t.Errorf("body exposes private error: %s", response.Body.String())
	}
}
