package shorty

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveLinkReturnsJSONDestination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/links/abc123/resolve" {
			t.Errorf("request = %s %s, want POST /links/abc123/resolve", request.Method, request.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"origin_url":"https://example.com/docs"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	location, err := client.ResolveLink(t.Context(), "abc123")
	if err != nil {
		t.Fatalf("ResolveLink() error = %v", err)
	}
	if location != "https://example.com/docs" {
		t.Errorf("location = %q, want https://example.com/docs", location)
	}
}

func TestRedirectURLReturnsAPINavigationEndpoint(t *testing.T) {
	client, err := NewClient("https://api.sho.rt", nil)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if got := client.RedirectURL("abc123"); got != "https://api.sho.rt/r/abc123" {
		t.Errorf("RedirectURL() = %q", got)
	}
}
