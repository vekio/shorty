package links

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveLinkReturnsDestinationAndRegistersVisit(t *testing.T) {
	httpAPI := newTestAPI()
	created := createLink(t, httpAPI, "https://example.com/docs")

	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(
		response,
		httptest.NewRequest(http.MethodPost, "/links/"+created.Code+"/resolve", nil),
	)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	output := decodeResponse[ResolveLinkOutput](t, response)
	if output.OriginURL != "https://example.com/docs" {
		t.Errorf("origin URL = %q, want https://example.com/docs", output.OriginURL)
	}

	getResponse := httptest.NewRecorder()
	httpAPI.ServeHTTP(
		getResponse,
		httptest.NewRequest(http.MethodGet, "/links/"+created.Code, nil),
	)
	link := decodeResponse[LinkOutput](t, getResponse)
	if link.Visits != 1 {
		t.Errorf("visits = %d, want 1", link.Visits)
	}
}

func TestResolveLinkMapsNotFoundError(t *testing.T) {
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(
		response,
		httptest.NewRequest(http.MethodPost, "/links/missing/resolve", nil),
	)

	assertProblem(t, response, http.StatusNotFound, "link not found", "/links/missing/resolve")
}
