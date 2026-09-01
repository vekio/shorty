package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

func TestDeleteLinkReturnsUpdatedHTMXFragment(t *testing.T) {
	api := &linkAPIStub{page: shortysdk.LinkPage{
		Links: []shortysdk.Link{{Code: "abc123", OriginURL: "https://example.com"}},
		Total: 1,
	}}
	request := httptest.NewRequest(http.MethodDelete, "/links/abc123", nil)
	request.Header.Set("HX-Request", "true")
	response := httptest.NewRecorder()

	newTestWeb(api).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if api.deletedCode != "abc123" {
		t.Errorf("deleted code = %q, want abc123", api.deletedCode)
	}
	if !strings.Contains(response.Body.String(), `id="workspace"`) ||
		!strings.Contains(response.Body.String(), "No links yet.") {
		t.Errorf("body = %s, want empty updated workspace", response.Body.String())
	}
}
