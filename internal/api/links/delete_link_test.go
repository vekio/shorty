package links

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteLinkRemovesResource(t *testing.T) {
	httpAPI := newTestAPI()
	created := createLink(t, httpAPI, "https://example.com/docs")

	deleteResponse := httptest.NewRecorder()
	httpAPI.ServeHTTP(
		deleteResponse,
		httptest.NewRequest(http.MethodDelete, "/links/"+created.Code, nil),
	)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d; body = %s", deleteResponse.Code, http.StatusNoContent, deleteResponse.Body.String())
	}
	if deleteResponse.Body.Len() != 0 {
		t.Errorf("delete body = %q, want empty response", deleteResponse.Body.String())
	}

	getResponse := httptest.NewRecorder()
	httpAPI.ServeHTTP(
		getResponse,
		httptest.NewRequest(http.MethodGet, "/links/"+created.Code, nil),
	)
	assertProblem(t, getResponse, http.StatusNotFound, "link not found", "/links/"+created.Code)
}

func TestDeleteLinkMapsNotFoundError(t *testing.T) {
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(
		response,
		httptest.NewRequest(http.MethodDelete, "/links/missing", nil),
	)

	assertProblem(t, response, http.StatusNotFound, "link not found", "/links/missing")
}
