package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

func TestMapApplicationError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		mapped bool
	}{
		{name: "required URL", err: domain.ErrOriginURLRequired, status: http.StatusBadRequest, mapped: true},
		{name: "invalid URL", err: domain.ErrOriginURLInvalid, status: http.StatusBadRequest, mapped: true},
		{name: "wrapped not found", err: errors.Join(errors.New("repository"), ports.ErrLinkNotFound), status: http.StatusNotFound, mapped: true},
		{name: "unknown", err: errors.New("unknown"), mapped: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			problem, mapped := mapApplicationError(test.err)
			if mapped != test.mapped {
				t.Fatalf("mapped = %v, want %v", mapped, test.mapped)
			}
			if mapped && (problem.Status != test.status || problem.Detail != test.err.Error()) {
				t.Errorf("problem = %#v", problem)
			}
		})
	}
}
