package amigo

import (
	"errors"
	"net/http"
	"testing"
)

func TestMapErrorAddsRouteMapping(t *testing.T) {
	target := errors.New("conflict")
	route := newRoute(http.MethodPost, "/things", MapError(target, http.StatusConflict))

	if len(route.errors) != 1 || !errors.Is(route.errors[0].target, target) || route.errors[0].status != http.StatusConflict {
		t.Errorf("errors = %#v", route.errors)
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
