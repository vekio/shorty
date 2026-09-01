package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmbeddedStylesheetIsServed(t *testing.T) {
	httpWeb := newTestWeb(&linkAPIStub{})

	for _, method := range []string{http.MethodGet, http.MethodHead} {
		t.Run(method, func(t *testing.T) {
			response := httptest.NewRecorder()
			httpWeb.ServeHTTP(
				response,
				httptest.NewRequest(method, "/assets/css/style.css", nil),
			)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
			}
			if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/css") {
				t.Errorf("Content-Type = %q, want text/css", got)
			}
			if method == http.MethodGet && !strings.Contains(response.Body.String(), "--color-accent") {
				t.Error("response does not contain the embedded stylesheet")
			}
			if method == http.MethodHead && response.Body.Len() != 0 {
				t.Errorf("HEAD body = %q, want empty", response.Body.String())
			}
		})
	}
}

func TestEmbeddedCopyScriptIsServed(t *testing.T) {
	response := httptest.NewRecorder()
	newTestWeb(&linkAPIStub{}).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/assets/js/app.js", nil),
	)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); !strings.Contains(got, "javascript") {
		t.Errorf("Content-Type = %q, want JavaScript", got)
	}
	if !strings.Contains(response.Body.String(), "navigator.clipboard.writeText") {
		t.Error("response does not contain the copy behavior")
	}
}

func TestEmbeddedAssetsRedirectPrefixToDirectory(t *testing.T) {
	response := httptest.NewRecorder()
	newTestWeb(&linkAPIStub{}).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/assets", nil),
	)

	if response.Code != http.StatusTemporaryRedirect {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusTemporaryRedirect)
	}
	if got := response.Header().Get("Location"); got != "/assets/" {
		t.Errorf("Location = %q, want /assets/", got)
	}
}

func TestEmbeddedAssetsRejectUnsupportedMethod(t *testing.T) {
	response := httptest.NewRecorder()
	newTestWeb(&linkAPIStub{}).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodPost, "/assets/css/style.css", nil),
	)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusMethodNotAllowed, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want net/http text response", got)
	}
	if got := response.Header().Get("Allow"); got != "GET, HEAD" {
		t.Errorf("Allow = %q, want GET, HEAD", got)
	}
}
