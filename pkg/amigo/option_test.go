package amigo

import (
	"net/http"
	"testing"
)

func TestStatusConfiguresSuccessCode(t *testing.T) {
	configuredRoute := route{status: http.StatusOK}
	Status(http.StatusCreated)(&configuredRoute)
	if configuredRoute.status != http.StatusCreated {
		t.Errorf("status = %d, want %d", configuredRoute.status, http.StatusCreated)
	}
}

func TestStatusRejectsNonSuccessCode(t *testing.T) {
	assertPanics(t, func() { Status(http.StatusBadRequest) })
}
