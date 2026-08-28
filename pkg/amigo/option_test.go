package amigo

import (
	"errors"
	"net/http"
	"testing"
)

func TestMapErrorAddsRouteMapping(t *testing.T) {
	target := errors.New("conflict")
	route := newRoute(http.MethodPost, "/things", MapError(target, http.StatusConflict))

	if len(route.errorMappings) != 1 || !errors.Is(route.errorMappings[0].target, target) || route.errorMappings[0].status != http.StatusConflict {
		t.Errorf("errorMappings = %#v", route.errorMappings)
	}
}

func TestRouteDefaults(t *testing.T) {
	route := newRoute(http.MethodGet, "/things")

	if route.status != http.StatusOK {
		t.Errorf("status = %d, want %d", route.status, http.StatusOK)
	}
	if route.maxBodyBytes != defaultMaxBodyBytes {
		t.Errorf("maxBodyBytes = %d, want %d", route.maxBodyBytes, defaultMaxBodyBytes)
	}
}

func TestRouteOptionsOverrideDefaults(t *testing.T) {
	route := newRoute(
		http.MethodPost,
		"/things",
		Status(http.StatusCreated),
		MaxBodyBytes(512),
	)

	if route.status != http.StatusCreated {
		t.Errorf("status = %d, want %d", route.status, http.StatusCreated)
	}
	if route.maxBodyBytes != 512 {
		t.Errorf("maxBodyBytes = %d, want %d", route.maxBodyBytes, 512)
	}
}

func TestUseAddsRouteMiddleware(t *testing.T) {
	middleware := func(next http.Handler) http.Handler { return next }
	route := newRoute(http.MethodGet, "/things", Use(middleware))

	if len(route.middlewares) != 1 {
		t.Errorf("middlewares = %d, want %d", len(route.middlewares), 1)
	}
}

func TestUseRejectsNilMiddleware(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Use() did not panic")
		}
	}()

	_ = Use(nil)
}

func TestMapErrorRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		target error
		status int
	}{
		{name: "nil error", status: http.StatusBadRequest},
		{name: "success status", target: errors.New("failure"), status: http.StatusOK},
		{name: "invalid status", target: errors.New("failure"), status: 600},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("MapError() did not panic")
				}
			}()
			_ = MapError(test.target, test.status)
		})
	}
}

func TestStatusRejectsNonSuccessStatus(t *testing.T) {
	for _, status := range []int{http.StatusContinue, http.StatusMultipleChoices} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			assertPanics(t, func() { Status(status) })
		})
	}
}

func TestMaxBodyBytesRejectsNegativeLimit(t *testing.T) {
	assertPanics(t, func() { MaxBodyBytes(-1) })
}

func TestNewRouteRejectsNilOption(t *testing.T) {
	assertPanics(t, func() { newRoute(http.MethodGet, "/things", nil) })
}
