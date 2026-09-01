package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

func TestHomeRendersLinkFormAndCollection(t *testing.T) {
	api := &linkAPIStub{page: shortysdk.LinkPage{
		Links: []shortysdk.Link{{
			Code:      "abc123",
			OriginURL: `https://example.com/docs?q=<script>`,
			CreatedAt: time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC),
			Visits:    3,
		}},
		Total: 1,
	}}
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/", nil),
	)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html; charset=utf-8", got)
	}
	for _, expected := range []string{
		`hx-post="/links?limit=20&amp;offset=0"`,
		`href="https://sho.rt/r/abc123"`,
		`data-copy="https://sho.rt/r/abc123"`,
		`3`,
		`29 Aug 2026`,
		`&lt;script&gt;`,
	} {
		if !strings.Contains(response.Body.String(), expected) {
			t.Errorf("body does not contain %q", expected)
		}
	}
	if strings.Contains(response.Body.String(), "<script>") {
		t.Error("body contains unescaped link data")
	}
}

func TestHomePreservesPaginationInHTMXControls(t *testing.T) {
	api := &linkAPIStub{listPage: func(options shortysdk.ListOptions) shortysdk.LinkPage {
		return shortysdk.LinkPage{
			Links: []shortysdk.Link{
				{Code: "third", OriginURL: "https://example.com/third"},
				{Code: "fourth", OriginURL: "https://example.com/fourth"},
			},
			Total:  5,
			Limit:  options.Limit,
			Offset: options.Offset,
		}
	}}
	request := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=2", nil)
	request.Header.Set("HX-Request", "true")
	response := httptest.NewRecorder()

	newTestWeb(api).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if api.listOptions != (shortysdk.ListOptions{Limit: 2, Offset: 2}) {
		t.Errorf("ListLinks() options = %#v, want limit 2 and offset 2", api.listOptions)
	}
	for _, expected := range []string{
		`hx-get="/?limit=2&amp;offset=0"`,
		`hx-get="/?limit=2&amp;offset=4"`,
		`3–4 of 5`,
		`hx-delete="/links/third?limit=2&amp;offset=2"`,
	} {
		if !strings.Contains(response.Body.String(), expected) {
			t.Errorf("body does not contain %q", expected)
		}
	}
}

func TestHomeRendersEmptyState(t *testing.T) {
	response := httptest.NewRecorder()
	newTestWeb(&linkAPIStub{}).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/", nil),
	)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), "No links yet.") {
		t.Errorf("body = %s, want empty state", response.Body.String())
	}
}

func TestHomeReturnsSafeInternalServerError(t *testing.T) {
	api := &linkAPIStub{listErr: errors.New("private API failure")}
	response := httptest.NewRecorder()
	newTestWeb(api).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/", nil),
	)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusInternalServerError, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want net/http text response", got)
	}
	if strings.Contains(response.Body.String(), "private API failure") {
		t.Errorf("body exposes private error: %s", response.Body.String())
	}
}
