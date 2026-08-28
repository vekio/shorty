package validator

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/domain"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{name: "empty is owned by required", value: "  "},
		{name: "HTTP", value: "http://example.com"},
		{name: "HTTPS", value: "https://example.com/docs?q=go"},
		{name: "relative", value: "/docs", wantErr: domain.ErrOriginURLInvalid},
		{name: "unsupported scheme", value: "ftp://example.com/file", wantErr: domain.ErrOriginURLInvalid},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateURL(test.value)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("validateURL(%q) error = %v, want %v", test.value, err, test.wantErr)
			}
		})
	}
}

func TestRegisterAddsURLValidator(t *testing.T) {
	type input struct {
		URL string `json:"url" validate:"url"`
	}

	httpAPI := amigo.New()
	Register(httpAPI)
	httpAPI.POST("/links", func(context.Context, input) (struct{}, error) {
		return struct{}{}, nil
	})

	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(`{"url":"not-a-url"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusUnprocessableEntity, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Errorf("Content-Type = %q, want application/problem+json", got)
	}
}
