package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServerMountsHandlersAndStripsPrefixes(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantBody string
	}{
		{name: "redirect", path: "/r/abc123", wantBody: "redirect:/r/abc123"},
		{name: "API", path: "/api/v1/links", wantBody: "api:/links"},
		{name: "admin", path: "/_/settings", wantBody: "admin:/settings"},
	}

	handler := NewServer(pathHandler("redirect"), pathHandler("api"), pathHandler("admin"))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, test.path, nil))
			if response.Body.String() != test.wantBody {
				t.Errorf("body = %q, want %q", response.Body.String(), test.wantBody)
			}
		})
	}
}

func TestNewServerReservesAppPrefix(t *testing.T) {
	handler := NewServer(pathHandler("redirect"), pathHandler("api"), pathHandler("admin"))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/app/", nil))
	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func pathHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		_, _ = w.Write([]byte(name + ":" + request.URL.Path))
	})
}
