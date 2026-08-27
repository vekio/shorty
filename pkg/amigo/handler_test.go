package amigo

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouteMapsHandlerError(t *testing.T) {
	target := errors.New("thing not found")
	api := New()
	api.GET("/things/{id}", func(context.Context, struct {
		ID string `path:"id" json:"-"`
	}) (struct{}, error) {
		return struct{}{}, errors.Join(errors.New("repository"), target)
	}, MapError(target, http.StatusNotFound))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/things/42", nil)
	api.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func TestRouteLimitsRequestBody(t *testing.T) {
	api := New()
	api.POST("/things", func(context.Context, struct {
		Name string `json:"name"`
	}) (struct{}, error) {
		return struct{}{}, nil
	}, MaxBodyBytes(4))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(`{"name":"shorty"}`))
	request.Header.Set("Content-Type", "application/json")
	api.ServeHTTP(response, request)

	if response.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", response.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestRawEndpointOwnsResponse(t *testing.T) {
	api := New()
	api.RAW(http.MethodGet, "/download", func(w http.ResponseWriter, _ *http.Request) error {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("raw"))
		return nil
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/download", nil)
	api.ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Body.String() != "raw" {
		t.Errorf("status = %d, body = %q", response.Code, response.Body.String())
	}
}
