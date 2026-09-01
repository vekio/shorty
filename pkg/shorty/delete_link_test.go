package shorty

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteLinkSendsDeleteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodDelete || request.URL.Path != "/links/abc123" {
			t.Errorf("request = %s %s, want DELETE /links/abc123", request.Method, request.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.DeleteLink(t.Context(), "abc123"); err != nil {
		t.Fatalf("DeleteLink() error = %v", err)
	}
}
