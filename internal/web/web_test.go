package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/vekio/shorty/internal/auth"
	"github.com/vekio/shorty/internal/infra/memory"
)

const testWorkspaceID = "ws_test"

func TestNewServesDashboard(t *testing.T) {
	handler, _ := newTestHandler()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "default workspace") {
		t.Errorf("response = %d %s", response.Code, response.Body.String())
	}
}

func TestHTMXCreatesAPIKeyAndShowsTokenOnce(t *testing.T) {
	handler, service := newTestHandler()
	response := serveHTMXForm(handler, http.MethodPost, "/api-keys", url.Values{"name": {"CI"}})
	if response.Code != http.StatusCreated {
		t.Fatalf("response = %d %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Cache-Control") != "no-store" ||
		!strings.Contains(response.Body.String(), "shorty_") ||
		!strings.Contains(response.Body.String(), `history.replaceState(null, "", "/_/")`) {
		t.Errorf("created key response headers = %#v, body = %s", response.Header(), response.Body.String())
	}
	keys, err := service.List(t.Context(), testWorkspaceID)
	if err != nil || len(keys) != 1 || keys[0].Name != "CI" {
		t.Fatalf("List() = %#v, %v", keys, err)
	}
	if !strings.Contains(response.Body.String(), keys[0].ID) ||
		!strings.Contains(response.Body.String(), "Never") ||
		!strings.Contains(response.Body.String(), "Active") {
		t.Errorf("created key metadata is missing from dashboard: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `id="api-key-management"`) ||
		strings.Contains(response.Body.String(), "<!doctype html>") {
		t.Errorf("HTMX response is not an API key panel fragment: %s", response.Body.String())
	}

	dashboard := httptest.NewRecorder()
	handler.ServeHTTP(dashboard, httptest.NewRequest(http.MethodGet, "/", nil))
	if strings.Contains(dashboard.Body.String(), "shorty_") {
		t.Errorf("dashboard exposed the API key again: %s", dashboard.Body.String())
	}
}

func TestHTMXReturnsValidationInsideSwappablePanel(t *testing.T) {
	handler, _ := newTestHandler()
	response := serveHTMXForm(handler, http.MethodPost, "/api-keys", url.Values{"name": {""}})
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "API key name is required") {
		t.Errorf("response = %d %s", response.Code, response.Body.String())
	}
}

func TestHTMXRevokesAPIKey(t *testing.T) {
	handler, service := newTestHandler()
	key, token, err := service.Create(t.Context(), testWorkspaceID, "automation")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	response := serveHTMXForm(handler, http.MethodPost, "/api-keys/"+key.ID+"/revoke", nil)
	if response.Code != http.StatusOK || response.Header().Get("Location") != "" {
		t.Fatalf("response = %d Location %q", response.Code, response.Header().Get("Location"))
	}
	if !strings.Contains(response.Body.String(), "Revoked") || strings.Contains(response.Body.String(), "<!doctype html>") {
		t.Errorf("HTMX response does not contain the updated panel: %s", response.Body.String())
	}

	if err := service.AuthenticateToken(t.Context(), token); !errors.Is(err, auth.ErrInvalidAPIKey) {
		t.Errorf("AuthenticateToken() error = %v, want %v", err, auth.ErrInvalidAPIKey)
	}

	keys, err := service.List(t.Context(), testWorkspaceID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(keys) != 1 || keys[0].RevokedAt == nil {
		t.Errorf("List() = %#v, want one revoked key", keys)
	}

	dashboard := httptest.NewRecorder()
	handler.ServeHTTP(dashboard, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(dashboard.Body.String(), "Revoked") {
		t.Errorf("dashboard does not show revoked status: %s", dashboard.Body.String())
	}
}

func TestNewServesEmbeddedStaticFiles(t *testing.T) {
	handler, _ := newTestHandler()
	for _, path := range []string{"/static/css/style.css", "/static/js/htmx.min.js"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK {
			t.Errorf("%s status = %d, want %d", path, response.Code, http.StatusOK)
		}
	}
}

func newTestHandler() (http.Handler, *auth.Service) {
	apiKeys := auth.NewService(memory.NewAPIKeyRepository())
	return New(apiKeys, testWorkspaceID, "default"), apiKeys
}

func serveHTMXForm(handler http.Handler, method string, target string, values url.Values) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("HX-Request", "true")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
