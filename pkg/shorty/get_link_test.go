package shorty

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLinkReturnsResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/links/abc123" {
			t.Errorf("request = %s %s, want GET /links/abc123", request.Method, request.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"abc123","origin_url":"https://example.com","created_at":"2026-08-29T10:00:00Z","visits":3}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	link, err := client.GetLink(t.Context(), "abc123")
	if err != nil {
		t.Fatalf("GetLink() error = %v", err)
	}
	if link.Code != "abc123" || link.Visits != 3 {
		t.Errorf("link = %#v, want abc123 with 3 visits", link)
	}
}
