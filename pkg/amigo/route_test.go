package amigo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPMethodsRegisterRoutes(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{name: "get", method: http.MethodGet},
		{name: "post", method: http.MethodPost},
		{name: "put", method: http.MethodPut},
		{name: "patch", method: http.MethodPatch},
		{name: "delete", method: http.MethodDelete},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := New()
			handler := func(context.Context, struct{}) (struct{}, error) { return struct{}{}, nil }
			switch test.method {
			case http.MethodGet:
				app.GET("/things", handler)
			case http.MethodPost:
				app.POST("/things", handler)
			case http.MethodPut:
				app.PUT("/things", handler)
			case http.MethodPatch:
				app.PATCH("/things", handler)
			case http.MethodDelete:
				app.DELETE("/things", handler)
			}

			response := httptest.NewRecorder()
			app.ServeHTTP(response, httptest.NewRequest(test.method, "/things", nil))
			if response.Code != http.StatusOK {
				t.Errorf("status = %d, want %d", response.Code, http.StatusOK)
			}
		})
	}
}

func TestRouteOnlyAcceptsRegisteredMethod(t *testing.T) {
	app := New()
	app.GET("/things", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, nil })
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/things", nil))
	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
}

func TestRouteRejectsNilHandler(t *testing.T) {
	var handler Handler[struct{}, struct{}]
	assertPanics(t, func() { New().GET("/things", handler) })
}

func TestRouteRejectsNonStructInput(t *testing.T) {
	assertPanics(t, func() {
		New().GET("/things", func(context.Context, string) (struct{}, error) { return struct{}{}, nil })
	})
}

func TestRouteRejectsNilOption(t *testing.T) {
	assertPanics(t, func() {
		New().GET("/things", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, nil }, nil)
	})
}

func TestRoutePattern(t *testing.T) {
	configuredRoute := route{method: http.MethodGet, path: "/things"}
	if configuredRoute.pattern() != "GET /things" {
		t.Errorf("pattern() = %q", configuredRoute.pattern())
	}
}
