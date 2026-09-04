// Package shorty provides a Go client for the Shorty management API.
package shorty

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 30 * time.Second

// Client calls a Shorty management API.
type Client struct {
	serverURL  *url.URL
	apiKey     string
	httpClient *http.Client
}

// ClientOption customizes a Client.
type ClientOption func(*Client) error

// WithHTTPClient replaces the default HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) error {
		if httpClient == nil {
			return fmt.Errorf("HTTP client is required")
		}
		client.httpClient = httpClient
		return nil
	}
}

// NewClient creates a client authenticated with a Shorty API key.
func NewClient(serverURL string, apiKey string, options ...ClientOption) (*Client, error) {
	parsedURL, err := parseServerURL(serverURL)
	if err != nil {
		return nil, err
	}
	if apiKey == "" || strings.TrimSpace(apiKey) != apiKey {
		return nil, fmt.Errorf("API key must be non-empty without surrounding whitespace")
	}

	client := &Client{
		serverURL:  parsedURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func parseServerURL(rawURL string) (*url.URL, error) {
	if rawURL == "" || strings.TrimSpace(rawURL) != rawURL {
		return nil, fmt.Errorf("server URL must be non-empty without surrounding whitespace")
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Host == "" ||
		(parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, fmt.Errorf("server URL must be an absolute HTTP or HTTPS URL")
	}
	if parsedURL.User != nil || parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return nil, fmt.Errorf("server URL cannot contain user information, query, or fragment")
	}
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")
	return parsedURL, nil
}

func (client *Client) endpoint(path string) string {
	endpoint := *client.serverURL
	endpoint.Path = client.serverURL.Path + path
	return endpoint.String()
}
