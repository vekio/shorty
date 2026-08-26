package amigo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewUsesSafeDefaults(t *testing.T) {
	app := New()
	if app.mux == nil {
		t.Fatal("mux is nil")
	}
	if app.maxBodyBytes != defaultMaxBodyBytes {
		t.Errorf("maxBodyBytes = %d, want %d", app.maxBodyBytes, defaultMaxBodyBytes)
	}
	if app.Handler() != app {
		t.Error("Handler() does not return the application")
	}
}

func TestServeHTTPDelegatesToMux(t *testing.T) {
	response := httptest.NewRecorder()
	New().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/missing", nil))
	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func TestSetMaxBodyBytes(t *testing.T) {
	app := New()
	app.SetMaxBodyBytes(0)
	if app.maxBodyBytes != 0 {
		t.Errorf("maxBodyBytes = %d, want 0", app.maxBodyBytes)
	}
}

func TestSetMaxBodyBytesRejectsNegativeLimit(t *testing.T) {
	assertPanics(t, func() { New().SetMaxBodyBytes(-1) })
}

func TestMapErrorsRejectsNilMapper(t *testing.T) {
	assertPanics(t, func() { New().MapErrors(nil) })
}
