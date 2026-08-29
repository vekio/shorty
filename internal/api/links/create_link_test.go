package links

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/domain"
)

type createLinkHandlerStub struct {
	err error
}

func (stub createLinkHandlerStub) Handle(context.Context, createlink.CreateLinkCommand) (createlink.CreateLinkResult, error) {
	return createlink.CreateLinkResult{}, stub.err
}

func TestCreateLink(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(`{"origin_url":"https://example.com/docs"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusCreated, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	output := decodeResponse[CreateLinkOutput](t, response)
	if output.Code == "" {
		t.Error("response code is empty")
	}
	if got := response.Header().Get("Location"); got != "/links/"+output.Code {
		t.Errorf("Location = %q, want %q", got, "/links/"+output.Code)
	}
}

func TestCreateLinkValidatesOriginURL(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(`{"origin_url":"not-a-url"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, request)

	problem := assertProblem(t, response, http.StatusUnprocessableEntity, "request validation failed", "/links")
	want := fieldErrorResponse{
		Location: "body.origin_url",
		Message:  "origin URL must be an absolute HTTP or HTTPS URL",
	}
	if len(problem.Errors) != 1 || problem.Errors[0] != want {
		t.Errorf("errors = %#v, want %#v", problem.Errors, []fieldErrorResponse{want})
	}
}

func TestCreateLinkRequiresOriginURL(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "missing body"},
		{name: "missing property", body: `{}`},
		{name: "null property", body: `{"origin_url":null}`},
		{name: "blank property", body: `{"origin_url":"  "}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			newTestAPI().ServeHTTP(response, request)

			problem := assertProblem(t, response, http.StatusUnprocessableEntity, "request validation failed", "/links")
			want := fieldErrorResponse{Location: "body.origin_url", Message: "is required"}
			if len(problem.Errors) != 1 || problem.Errors[0] != want {
				t.Errorf("errors = %#v, want %#v", problem.Errors, []fieldErrorResponse{want})
			}
		})
	}
}

func TestCreateLinkRejectsInvalidJSONRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		status      int
		detail      string
	}{
		{name: "wrong content type", body: `{}`, contentType: "text/plain", status: http.StatusUnsupportedMediaType, detail: "content type must be application/json"},
		{name: "invalid JSON", body: `{"origin_url":`, contentType: "application/json", status: http.StatusBadRequest, detail: "invalid JSON request body"},
		{name: "unknown field", body: `{"origin_url":"https://example.com","extra":true}`, contentType: "application/json", status: http.StatusBadRequest, detail: "invalid JSON request body"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()
			newTestAPI().ServeHTTP(response, request)
			assertProblem(t, response, test.status, test.detail, "/links")
		})
	}
}

func TestCreateLinkMapsApplicationErrorWithoutExposingItsCause(t *testing.T) {
	httpAPI := newTestAPIWithApplication(app.Application{
		Commands: app.Commands{
			CreateLink: createLinkHandlerStub{
				err: fmt.Errorf("database detail: %w", domain.ErrOriginURLInvalid),
			},
		},
	})
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(`{"origin_url":"https://example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, request)

	assertProblem(
		t,
		response,
		http.StatusUnprocessableEntity,
		"origin URL must be an absolute HTTP or HTTPS URL",
		"/links",
	)
	if strings.Contains(response.Body.String(), "database detail") {
		t.Errorf("response exposes private error cause: %s", response.Body.String())
	}
}
