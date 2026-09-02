package links

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateLinkChangesOriginURL(t *testing.T) {
	httpAPI := newTestAPI()
	created := createLink(t, httpAPI, "https://example.com/old")
	request := httptest.NewRequest(
		http.MethodPatch,
		"/links/"+created.Code,
		strings.NewReader(`{"origin_url":"https://example.com/new"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", response.Code, response.Body.String())
	}
	output := decodeResponse[LinkOutput](t, response)
	if output.Code != created.Code || output.OriginURL != "https://example.com/new" {
		t.Errorf("output = %#v", output)
	}
}
