package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNotFoundRendersUnknownNestedPath(t *testing.T) {
	response := httptest.NewRecorder()
	newTestWeb(&linkAPIStub{}).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/links/missing/page", nil),
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusNotFound, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want HTML", got)
	}
	if body := response.Body.String(); !strings.Contains(body, "Page not found") ||
		!strings.Contains(body, "/links/missing/page") {
		t.Errorf("body = %s, want missing path", body)
	}
}

func TestNotFoundEscapesRequestPath(t *testing.T) {
	response := httptest.NewRecorder()
	newTestWeb(&linkAPIStub{}).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/%3Cscript%3E", nil),
	)

	if strings.Contains(response.Body.String(), "<script>") {
		t.Errorf("body contains unescaped request path: %s", response.Body.String())
	}
}
