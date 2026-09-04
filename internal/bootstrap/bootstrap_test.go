package bootstrap

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	shortyauth "github.com/vekio/shorty/internal/auth"
	shortyconfig "github.com/vekio/shorty/internal/config"
)

func TestNewComposesSingleProcessWithSharedApplication(t *testing.T) {
	config := shortyconfig.Default()
	config.Database.Path = filepath.Join(t.TempDir(), "shorty.db")
	token := provisionAPIKey(t, config)
	runtime, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() {
		if err := runtime.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})
	if runtime.Handler == nil || runtime.Logger == nil {
		t.Fatalf("runtime = %#v, want complete process", runtime)
	}

	adminResponse := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(adminResponse, httptest.NewRequest(http.MethodGet, "/_/", nil))
	if adminResponse.Code != http.StatusOK || !strings.Contains(adminResponse.Body.String(), "Shorty") {
		t.Fatalf("admin response = %d %s", adminResponse.Code, adminResponse.Body.String())
	}

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/links",
		strings.NewReader(`{"origin_url":"https://example.com"}`),
	)
	createRequest.Header.Set("Content-Type", "application/json")
	authorizeRequest(createRequest, token)
	createResponse := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d; body = %s", createResponse.Code, createResponse.Body.String())
	}

	var created struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	redirectResponse := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(
		redirectResponse,
		httptest.NewRequest(http.MethodGet, "/r/"+created.Code, nil),
	)
	if redirectResponse.Code != http.StatusFound || redirectResponse.Header().Get("Location") != "https://example.com" {
		t.Errorf("public redirect = %d Location %q", redirectResponse.Code, redirectResponse.Header().Get("Location"))
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/v1/links/"+created.Code, nil)
	authorizeRequest(getRequest, token)
	getResponse := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusOK {
		t.Errorf("get status = %d; body = %s", getResponse.Code, getResponse.Body.String())
	}
}

func TestNewProtectsManagementAPI(t *testing.T) {
	config := shortyconfig.Default()
	config.Database = shortyconfig.DatabaseConfig{Driver: shortyconfig.DatabaseDriverMemory}
	runtime, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = runtime.Close() })
	response := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/links", nil))
	if response.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestAdminCreatedAPIKeyAuthenticatesManagementAPI(t *testing.T) {
	config := shortyconfig.Default()
	config.Database = shortyconfig.DatabaseConfig{Driver: shortyconfig.DatabaseDriverMemory}
	runtime, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = runtime.Close() })

	createKeyRequest := httptest.NewRequest(
		http.MethodPost,
		"/_/api-keys",
		strings.NewReader("name=Bruno"),
	)
	createKeyRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createKeyResponse := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(createKeyResponse, createKeyRequest)
	if createKeyResponse.Code != http.StatusCreated {
		t.Fatalf("create key status = %d; body = %s", createKeyResponse.Code, createKeyResponse.Body.String())
	}

	match := regexp.MustCompile(`<code id="new-token"[^>]*>([^<]+)</code>`).FindStringSubmatch(createKeyResponse.Body.String())
	if len(match) != 2 {
		t.Fatalf("created token not found in response: %s", createKeyResponse.Body.String())
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
	authorizeRequest(listRequest, match[1])
	listResponse := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(listResponse, listRequest)
	if listResponse.Code != http.StatusOK {
		t.Errorf("list status = %d; body = %s", listResponse.Code, listResponse.Body.String())
	}
}

func TestNewInjectsSQLiteRepository(t *testing.T) {
	config := shortyconfig.Default()
	config.Database.Path = filepath.Join(t.TempDir(), "shorty.db")
	token := provisionAPIKey(t, config)

	firstRuntime, err := New(config)
	if err != nil {
		t.Fatalf("first New() error = %v", err)
	}
	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/links",
		strings.NewReader(`{"origin_url":"https://example.com/persisted"}`),
	)
	createRequest.Header.Set("Content-Type", "application/json")
	authorizeRequest(createRequest, token)
	createResponse := httptest.NewRecorder()
	firstRuntime.Handler.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d; body = %s", createResponse.Code, createResponse.Body.String())
	}
	var created struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if err := firstRuntime.Close(); err != nil {
		t.Fatalf("close first runtime: %v", err)
	}

	secondRuntime, err := New(config)
	if err != nil {
		t.Fatalf("second New() error = %v", err)
	}
	t.Cleanup(func() { _ = secondRuntime.Close() })
	getRequest := httptest.NewRequest(http.MethodGet, "/api/v1/links/"+created.Code, nil)
	authorizeRequest(getRequest, token)
	getResponse := httptest.NewRecorder()
	secondRuntime.Handler.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusOK {
		t.Errorf("get persisted link status = %d; body = %s", getResponse.Code, getResponse.Body.String())
	}
}

func provisionAPIKey(t *testing.T, config shortyconfig.Config) string {
	t.Helper()
	persistence, err := OpenPersistence(t.Context(), config.Database)
	if err != nil {
		t.Fatalf("open persistence to provision API key: %v", err)
	}
	_, token, err := shortyauth.NewService(persistence.APIKeys).Create(
		t.Context(), persistence.WorkspaceID, "test",
	)
	if err != nil {
		t.Fatalf("provision API key: %v", err)
	}
	if err := persistence.Close(); err != nil {
		t.Fatalf("close API key persistence: %v", err)
	}
	return token
}

func authorizeRequest(request *http.Request, token string) {
	request.Header.Set("Authorization", "Bearer "+token)
}
