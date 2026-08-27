package api

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateLink(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(`{"origin_url":"https://example.com/docs"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusCreated, response.Body.String())
	}
	var output CreateLinkOutput
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if output.Code == "" {
		t.Error("response code is empty")
	}
	if got := response.Header().Get("Location"); got != "/links/"+output.Code {
		t.Errorf("Location = %q, want %q", got, "/links/"+output.Code)
	}
}

func TestCreateLinkMapsDomainError(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(`{"origin_url":"not-a-url"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusBadRequest, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"detail":"origin URL must be an absolute HTTP or HTTPS URL"`) {
		t.Errorf("body = %s, want mapped domain error", response.Body.String())
	}
}

func TestCreateLinkRejectsInvalidJSONRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		status      int
	}{
		{name: "missing body", contentType: "application/json", status: http.StatusBadRequest},
		{name: "wrong content type", body: `{}`, contentType: "text/plain", status: http.StatusUnsupportedMediaType},
		{name: "invalid JSON", body: `{"origin_url":`, contentType: "application/json", status: http.StatusBadRequest},
		{name: "unknown field", body: `{"origin_url":"https://example.com","extra":true}`, contentType: "application/json", status: http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()
			newTestAPI().ServeHTTP(response, request)
			if response.Code != test.status {
				t.Errorf("status = %d, want %d; body = %s", response.Code, test.status, response.Body.String())
			}
		})
	}
}
