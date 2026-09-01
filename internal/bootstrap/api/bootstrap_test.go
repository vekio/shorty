package api

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	shortyconfig "github.com/vekio/shorty/internal/config"
	apiconfig "github.com/vekio/shorty/internal/config/api"
)

func TestNewComposesAPIProcessWithSharedRepository(t *testing.T) {
	runtime, err := New(apiconfig.Config{
		Address:  ":9080",
		ShortURL: "https://sho.rt",
		Logger:   shortyconfig.DefaultLoggerConfig(),
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if runtime.Handler == nil || runtime.Logger == nil {
		t.Fatalf("runtime = %#v, want complete API process", runtime)
	}

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/links",
		strings.NewReader(`{"origin_url":"https://example.com"}`),
	)
	createRequest.Header.Set("Content-Type", "application/json")
	createRequest.Header.Set("X-Shorty-Owner", "browser-a")
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

	getResponse := httptest.NewRecorder()
	getRequest := httptest.NewRequest(http.MethodGet, "/links/"+created.Code, nil)
	getRequest.Header.Set("X-Shorty-Owner", "browser-a")
	runtime.Handler.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusOK {
		t.Errorf("get status = %d; body = %s", getResponse.Code, getResponse.Body.String())
	}
}

func TestNewRejectsLinksPointingToTheConfiguredShortURL(t *testing.T) {
	runtime, err := New(apiconfig.Config{
		Address:  ":9080",
		ShortURL: "https://sho.rt",
		Logger:   shortyconfig.DefaultLoggerConfig(),
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	request := httptest.NewRequest(
		http.MethodPost,
		"/links",
		strings.NewReader(`{"origin_url":"https://sho.rt/r/existing"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Shorty-Owner", "browser-a")
	response := httptest.NewRecorder()

	runtime.Handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusUnprocessableEntity, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "origin URL cannot point to this Shorty instance") {
		t.Errorf("body = %s, want safe self-reference detail", response.Body.String())
	}
}
