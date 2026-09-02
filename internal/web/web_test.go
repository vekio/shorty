package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/deletelink"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/auth"
	"github.com/vekio/shorty/internal/infra/memory"
)

const testWorkspaceID = "ws_test"

func TestNewServesDashboard(t *testing.T) {
	handler, _ := newTestHandler()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "Workspace · default") {
		t.Errorf("response = %d %s", response.Code, response.Body.String())
	}
}

func TestDashboardCreatesLinkThroughUseCase(t *testing.T) {
	handler, _ := newTestHandler()
	response := serveForm(handler, http.MethodPost, "/links", url.Values{
		"origin_url": {"https://example.com/docs"},
	})
	if response.Code != http.StatusSeeOther || response.Header().Get("Location") != "/_/" {
		t.Fatalf("response = %d Location %q", response.Code, response.Header().Get("Location"))
	}

	dashboard := httptest.NewRecorder()
	handler.ServeHTTP(dashboard, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(dashboard.Body.String(), "https://example.com/docs") {
		t.Errorf("dashboard does not contain created link: %s", dashboard.Body.String())
	}
}

func TestDashboardCreatesAPIKeyAndShowsTokenOnce(t *testing.T) {
	handler, service := newTestHandler()
	response := serveForm(handler, http.MethodPost, "/api-keys", url.Values{"name": {"CI"}})
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

	dashboard := httptest.NewRecorder()
	handler.ServeHTTP(dashboard, httptest.NewRequest(http.MethodGet, "/", nil))
	if strings.Contains(dashboard.Body.String(), "shorty_") {
		t.Errorf("dashboard exposed the API key again: %s", dashboard.Body.String())
	}
}

func TestDashboardRevokesAPIKey(t *testing.T) {
	handler, service := newTestHandler()
	key, token, err := service.Create(t.Context(), testWorkspaceID, "automation")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	response := serveForm(handler, http.MethodPost, "/api-keys/"+key.ID+"/revoke", nil)
	if response.Code != http.StatusSeeOther || response.Header().Get("Location") != "/_/" {
		t.Fatalf("response = %d Location %q", response.Code, response.Header().Get("Location"))
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
}

func TestNewServesEmbeddedStaticFiles(t *testing.T) {
	handler, _ := newTestHandler()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	if response.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", response.Code, http.StatusOK)
	}
}

func newTestHandler() (http.Handler, *auth.Service) {
	links := memory.NewLinkRepository()
	application := app.Application{
		Commands: app.Commands{
			CreateLink: createlink.NewCreateLinkHandler(links),
			DeleteLink: deletelink.NewDeleteLinkHandler(links),
		},
		Queries: app.Queries{
			ListLinks: listlinks.NewListLinksHandler(links),
		},
	}
	apiKeys := auth.NewService(memory.NewAPIKeyRepository())
	return New(application, apiKeys, testWorkspaceID, "default"), apiKeys
}

func serveForm(handler http.Handler, method string, target string, values url.Values) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
