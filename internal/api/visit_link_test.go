package api

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVisitLinkRegistersVisit(t *testing.T) {
	httpAPI := newTestAPI()
	created := createLink(t, httpAPI, "https://example.com/docs")
	visitRequest := httptest.NewRequest(http.MethodPost, "/links/"+created.Code+"/visit", nil)
	visitResponse := httptest.NewRecorder()
	httpAPI.ServeHTTP(visitResponse, visitRequest)

	if visitResponse.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body = %s", visitResponse.Code, http.StatusNoContent, visitResponse.Body.String())
	}
	if visitResponse.Body.Len() != 0 {
		t.Errorf("body = %s, want empty response", visitResponse.Body.String())
	}

	getResponse := httptest.NewRecorder()
	httpAPI.ServeHTTP(getResponse, httptest.NewRequest(http.MethodGet, "/links/"+created.Code, nil))
	var fetched GetLinkOutput
	if err := json.Unmarshal(getResponse.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if fetched.Visits != 1 {
		t.Errorf("visits = %d, want 1", fetched.Visits)
	}
}

func TestVisitLinkMapsNotFoundError(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/links/missing/visit", nil)
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusNotFound, response.Body.String())
	}
}
