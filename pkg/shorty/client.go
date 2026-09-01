// Package shorty provides a Go client for Shorty's HTTP API.
package shorty

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client sends requests to one Shorty API instance.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// Link is the public representation returned by Shorty's link endpoints.
type Link struct {
	Code      string    `json:"code"`
	OriginURL string    `json:"origin_url"`
	CreatedAt time.Time `json:"created_at"`
	Visits    int       `json:"visits"`
}

// LinkPage is one page returned by GET /links.
type LinkPage struct {
	Links  []Link `json:"links"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// FieldError describes one invalid request field returned by the API.
type FieldError struct {
	Location string `json:"location"`
	Message  string `json:"message"`
}

// ProblemError is an RFC 9457 error response returned by the API.
type ProblemError struct {
	Type     string       `json:"type,omitempty"`
	Title    string       `json:"title"`
	Status   int          `json:"status"`
	Detail   string       `json:"detail"`
	Instance string       `json:"instance,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
}

func (err *ProblemError) Error() string {
	if err.Detail != "" {
		return err.Detail
	}
	return err.Title
}

// NewClient creates a client for one Shorty API instance.
func NewClient(baseURL string, httpClient *http.Client) (*Client, error) {
	parsed, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("shorty: API URL must be absolute")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("shorty: API URL must use HTTP or HTTPS")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, fmt.Errorf("shorty: API URL cannot contain user information, query, or fragment")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{baseURL: parsed, httpClient: httpClient}, nil
}

// ListOptions selects a page. Zero values use the API defaults.
type ListOptions struct {
	// Limit is the maximum number of links returned by the API.
	Limit int
	// Offset is the zero-based position of the first returned link.
	Offset int
}
