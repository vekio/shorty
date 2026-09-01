package shorty

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateLinkSendsJSONAndReturnsCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/links" {
			t.Errorf("request = %s %s, want POST /links", request.Method, request.URL.Path)
		}
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"code":"abc123"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	code, err := client.CreateLink(t.Context(), "https://example.com")
	if err != nil {
		t.Fatalf("CreateLink() error = %v", err)
	}
	if code != "abc123" {
		t.Errorf("code = %q, want abc123", code)
	}
}

func TestCreateLinkReturnsAPIVvalidationProblem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"title":"Unprocessable Entity","status":422,"detail":"request validation failed","errors":[{"location":"body.origin_url","message":"is required"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	_, err = client.CreateLink(t.Context(), "")
	problem, ok := err.(*ProblemError)
	if !ok || problem.Status != http.StatusUnprocessableEntity ||
		len(problem.Errors) != 1 || problem.Errors[0].Message != "is required" {
		t.Fatalf("error = %#v, want validation ProblemError", err)
	}
}
