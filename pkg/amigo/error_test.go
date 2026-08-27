package amigo

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteErrorReturnsSafeInternalServerError(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/things", nil)
	response := httptest.NewRecorder()

	writeError(response, request, route{}, errors.New("private failure"))

	if response.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	if response.Header().Get("Content-Type") != "application/problem+json" {
		t.Errorf("Content-Type = %q", response.Header().Get("Content-Type"))
	}
	var output problem
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if output.Detail != "internal server error" || output.Instance != "/things" {
		t.Errorf("problem = %#v", output)
	}
}

func TestWriteErrorUsesRouteMapping(t *testing.T) {
	target := errors.New("thing not found")
	request := httptest.NewRequest(http.MethodGet, "/things/42", nil)
	response := httptest.NewRecorder()
	route := newRoute(http.MethodGet, "/things/{id}", MapError(target, http.StatusNotFound))

	writeError(response, request, route, errors.Join(errors.New("repository"), target))

	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
	var output problem
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if output.Detail != target.Error() || output.Instance != "/things/42" {
		t.Errorf("problem = %#v", output)
	}
}
