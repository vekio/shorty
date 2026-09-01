package shorty

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListLinksDecodesPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/links" {
			t.Errorf("request = %s %s, want GET /links", request.Method, request.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"links":[{"code":"abc123","origin_url":"https://example.com","created_at":"2026-08-29T10:00:00Z","visits":2}],"total":1,"limit":20,"offset":0}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	page, err := client.ListLinks(t.Context(), ListOptions{})
	if err != nil {
		t.Fatalf("ListLinks() error = %v", err)
	}
	if len(page.Links) != 1 || page.Links[0].Code != "abc123" ||
		page.Links[0].Visits != 2 || page.Total != 1 {
		t.Errorf("page = %#v, want one visited link", page)
	}
}

func TestListLinksEncodesPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Get("limit") != "10" || request.URL.Query().Get("offset") != "20" {
			t.Errorf("query = %q, want limit=10 and offset=20", request.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"links":[],"total":0,"limit":10,"offset":20}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if _, err := client.ListLinks(t.Context(), ListOptions{Limit: 10, Offset: 20}); err != nil {
		t.Fatalf("ListLinks() error = %v", err)
	}
}
