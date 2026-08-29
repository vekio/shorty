package validator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vekio/amigo"
)

type validationInput struct {
	URL    string `json:"url" validate:"url"`
	Limit  int    `query:"limit" json:"-" validate:"page_limit"`
	Offset int    `query:"offset" json:"-" validate:"page_offset"`
}

func TestRegisterAddsSharedValidators(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		target   string
		body     string
		location string
	}{
		{name: "URL", method: http.MethodPost, target: "/links", body: `{"url":"not-a-url"}`, location: "body.url"},
		{name: "limit", method: http.MethodGet, target: "/links?limit=101", location: "query.limit"},
		{name: "offset", method: http.MethodGet, target: "/links?offset=-1", location: "query.offset"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpAPI := amigo.New()
			Register(httpAPI)
			httpAPI.GET("/links", acceptValidationInput)
			httpAPI.POST("/links", acceptValidationInput)

			request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.body))
			if test.body != "" {
				request.Header.Set("Content-Type", "application/json")
			}
			response := httptest.NewRecorder()
			httpAPI.ServeHTTP(response, request)

			if response.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusUnprocessableEntity, response.Body.String())
			}
			if got := response.Header().Get("Content-Type"); got != "application/problem+json" {
				t.Errorf("Content-Type = %q, want application/problem+json", got)
			}
			if !strings.Contains(response.Body.String(), `"location":"`+test.location+`"`) {
				t.Errorf("body = %s, want location %q", response.Body.String(), test.location)
			}
		})
	}
}

func acceptValidationInput(context.Context, validationInput) (struct{}, error) {
	return struct{}{}, nil
}
