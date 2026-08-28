package api

import (
	"encoding/json/v2"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vekio/shorty/internal/app"
)

type routingProblem struct {
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func newTestAPI() http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(app.Application{}, logger)
}

func TestServerRejectsUnregisteredMethod(t *testing.T) {
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/links", nil))
	assertRoutingProblem(t, response, http.StatusMethodNotAllowed, "method not allowed", "/links")
	if got := response.Header().Get("Allow"); got != "GET, HEAD, POST" {
		t.Errorf("Allow = %q, want %q", got, "GET, HEAD, POST")
	}
}

func TestServerReturnsNotFoundForUnknownRoute(t *testing.T) {
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/missing", nil))
	assertRoutingProblem(t, response, http.StatusNotFound, "resource not found", "/missing")
	if got := response.Header().Get("Allow"); got != "" {
		t.Errorf("Allow = %q, want empty header", got)
	}
}

func TestNewRegistersLinkRoutesAndSharedValidators(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(`{"origin_url":"invalid"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, request)

	assertRoutingProblem(t, response, http.StatusUnprocessableEntity, "request validation failed", "/links")
}

func assertRoutingProblem(
	t *testing.T,
	response *httptest.ResponseRecorder,
	status int,
	detail string,
	instance string,
) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, status, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Errorf("Content-Type = %q, want application/problem+json", got)
	}

	var problem routingProblem
	if err := json.Unmarshal(response.Body.Bytes(), &problem); err != nil {
		t.Fatalf("decode problem: %v; body = %s", err, response.Body.String())
	}
	if problem.Title != http.StatusText(status) || problem.Status != status ||
		problem.Detail != detail || problem.Instance != instance {
		t.Errorf("problem = %#v, want status %d, detail %q, instance %q", problem, status, detail, instance)
	}
}
