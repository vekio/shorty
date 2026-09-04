package shorty

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (function roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func TestNewClientValidatesSettings(t *testing.T) {
	tests := []struct {
		name      string
		serverURL string
		apiKey    string
		option    ClientOption
	}{
		{name: "server URL", serverURL: "/shorty", apiKey: "shorty_token"},
		{name: "API key", serverURL: "https://shorty.example"},
		{name: "HTTP client", serverURL: "https://shorty.example", apiKey: "shorty_token", option: WithHTTPClient(nil)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options := []ClientOption{}
			if test.option != nil {
				options = append(options, test.option)
			}
			if _, err := NewClient(test.serverURL, test.apiKey, options...); err == nil {
				t.Fatal("NewClient() returned nil error")
			}
		})
	}
}

func TestCreateLink(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/v1/links" {
			t.Errorf("request = %s %s", request.Method, request.URL.Path)
		}
		if got := request.Header.Get("Authorization"); got != "Bearer shorty_token" {
			t.Errorf("Authorization = %q", got)
		}
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q", got)
		}
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Content-Type": {"application/json"},
				"Location":     {"/api/v1/links/abc123"},
			},
			Body: io.NopCloser(strings.NewReader(`{"code":"abc123"}`)),
		}, nil
	})}

	client, err := NewClient(
		"https://shorty.example",
		"shorty_token",
		WithHTTPClient(httpClient),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.CreateLink(t.Context(), CreateLinkRequest{
		OriginURL: "https://example.com/docs",
	})
	if err != nil {
		t.Fatalf("CreateLink() error = %v", err)
	}
	if result.Code != "abc123" || result.Location != "/api/v1/links/abc123" {
		t.Errorf("result = %#v", result)
	}
}

func TestCreateLinkReturnsAPIProblem(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripperFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Header:     http.Header{"Content-Type": {"application/problem+json"}},
			Body: io.NopCloser(strings.NewReader(`{
				"title":"Unprocessable Entity",
				"status":422,
				"detail":"request validation failed",
				"instance":"/api/v1/links",
				"errors":[{"location":"body.origin_url","message":"is required"}]
			}`)),
		}, nil
	})}

	client, err := NewClient(
		"https://shorty.example",
		"shorty_token",
		WithHTTPClient(httpClient),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	_, err = client.CreateLink(t.Context(), CreateLinkRequest{})
	var apiError *APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("CreateLink() error = %v, want *APIError", err)
	}
	if apiError.StatusCode != http.StatusUnprocessableEntity ||
		apiError.Problem.Detail != "request validation failed" ||
		len(apiError.Problem.Errors) != 1 {
		t.Errorf("APIError = %#v", apiError)
	}
}
