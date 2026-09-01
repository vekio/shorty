package shorty

import (
	"context"
	"net/http"
	"testing"
)

func TestNewClientValidatesAPIURL(t *testing.T) {
	for _, rawURL := range []string{
		"",
		"/api",
		"ftp://example.com",
		"https://user@example.com",
		"https://example.com?token=secret",
		"https://example.com#fragment",
	} {
		t.Run(rawURL, func(t *testing.T) {
			if _, err := NewClient(rawURL, nil); err == nil {
				t.Errorf("NewClient(%q) accepted invalid URL", rawURL)
			}
		})
	}
}

func TestWithOwnerAddsOwnershipHeader(t *testing.T) {
	client, err := NewClient("https://api.example.com", nil)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	request, err := client.request(
		WithOwner(context.Background(), "browser-a"),
		http.MethodGet,
		"/links",
		nil,
	)
	if err != nil {
		t.Fatalf("request() error = %v", err)
	}
	if got := request.Header.Get("X-Shorty-Owner"); got != "browser-a" {
		t.Errorf("X-Shorty-Owner = %q, want browser-a", got)
	}
}

func TestNewClientDoesNotMutateProvidedHTTPClient(t *testing.T) {
	originalRedirect := func(*http.Request, []*http.Request) error { return nil }
	httpClient := &http.Client{CheckRedirect: originalRedirect}
	if _, err := NewClient("https://api.example.com", httpClient); err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if httpClient.CheckRedirect == nil {
		t.Error("NewClient() mutated the provided HTTP client")
	}
}
