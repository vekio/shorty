package api

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLinkBindsIDFromPath(t *testing.T) {
	httpAPI := newTestAPI()
	created := createLink(t, httpAPI, "https://example.com/docs")
	request := httptest.NewRequest(http.MethodGet, "/links/"+created.Code, nil)
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	var output LinkOutput
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if output.Code != created.Code || output.OriginURL != "https://example.com/docs" || output.Visits != 0 || output.CreatedAt.IsZero() {
		t.Errorf("output = %#v, want created link without registering a visit", output)
	}
}

func TestGetLinkMapsNotFoundError(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/links/missing", nil)
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusNotFound, response.Body.String())
	}
}
