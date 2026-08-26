package amigo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type bodyInput struct {
	Body struct {
		Name string `json:"name"`
	}
}

func bodyEcho(_ context.Context, input bodyInput) (string, error) {
	return input.Body.Name, nil
}

func TestBindBody(t *testing.T) {
	app := New()
	app.POST("/things", bodyEcho)
	request := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(`{"name":"shorty"}`))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response := httptest.NewRecorder()

	app.ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Body.String() != `"shorty"` {
		t.Errorf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestBindBodyRequiresBody(t *testing.T) {
	app := New()
	app.POST("/things", bodyEcho)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/things", nil))
	assertProblem(t, response, http.StatusBadRequest, "request body is required")
}

func TestBindBodyRequiresJSONContentType(t *testing.T) {
	app := New()
	app.POST("/things", bodyEcho)
	request := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(`{"name":"shorty"}`))
	request.Header.Set("Content-Type", "text/plain")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	assertProblem(t, response, http.StatusUnsupportedMediaType, "content type must be application/json")
}

func TestBindBodyRejectsInvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed", body: `{"name":`},
		{name: "unknown member", body: `{"name":"shorty","unknown":true}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := New()
			app.POST("/things", bodyEcho)
			request := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			app.ServeHTTP(response, request)
			assertProblem(t, response, http.StatusBadRequest, "invalid JSON request body")
		})
	}
}

func TestBindBodyEnforcesSizeLimit(t *testing.T) {
	app := New()
	app.SetMaxBodyBytes(5)
	app.POST("/things", bodyEcho)
	request := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(`{"name":"shorty"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	assertProblem(t, response, http.StatusRequestEntityTooLarge, "request body exceeds the maximum allowed size")
}

func TestBindBodyCanDisableSizeLimit(t *testing.T) {
	app := New()
	app.SetMaxBodyBytes(0)
	app.POST("/things", bodyEcho)
	request := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(`{"name":"shorty"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", response.Code, http.StatusOK)
	}
}

func TestBindBodyIgnoresBodyWhenInputDoesNotDeclareOne(t *testing.T) {
	app := New()
	app.GET("/things", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, nil })
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/things", strings.NewReader("ignored")))
	if response.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", response.Code, http.StatusOK)
	}
}
