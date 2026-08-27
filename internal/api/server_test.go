package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerRejectsUnregisteredMethod(t *testing.T) {
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/links", nil))
	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
}

func TestServerReturnsNotFoundForUnknownRoute(t *testing.T) {
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/missing", nil))
	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
