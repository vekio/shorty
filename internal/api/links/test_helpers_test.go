package links

import (
	"bytes"
	"encoding/json/v2"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/amigo"
	apivalidator "github.com/vekio/shorty/internal/api/validator"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/deletelink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/infra/memory"
)

type fieldErrorResponse struct {
	Location string `json:"location"`
	Message  string `json:"message"`
}

type problemResponse struct {
	Title    string               `json:"title"`
	Status   int                  `json:"status"`
	Detail   string               `json:"detail"`
	Instance string               `json:"instance"`
	Errors   []fieldErrorResponse `json:"errors"`
}

func newTestAPI() http.Handler {
	return newTestAPIWithApplication(newTestApplication())
}

func newTestApplication() app.Application {
	repository := memory.NewLinkRepository()
	return app.Application{
		Commands: app.Commands{
			CreateLink: createlink.NewCreateLinkHandler(repository),
			DeleteLink: deletelink.NewDeleteLinkHandler(repository),
		},
		Queries: app.Queries{
			GetLink:   getlink.NewGetLinkHandler(repository),
			ListLinks: listlinks.NewListLinksHandler(repository),
		},
	}
}

func newTestAPIWithApplication(application app.Application) http.Handler {
	httpAPI := amigo.New(amigo.WithLogger(testLogger()))
	apivalidator.Register(httpAPI)
	Register(httpAPI.Group("/links"), application)
	return httpAPI
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func createLink(t *testing.T, httpAPI http.Handler, originURL string) CreateLinkOutput {
	t.Helper()
	body, err := json.Marshal(CreateLinkInput{OriginURL: originURL})
	if err != nil {
		t.Fatalf("encode create request: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/links", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body = %s", response.Code, http.StatusCreated, response.Body.String())
	}
	var output CreateLinkOutput
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	return output
}

func decodeResponse[T any](t *testing.T, response *httptest.ResponseRecorder) T {
	t.Helper()
	var body T
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v; body = %s", err, response.Body.String())
	}
	return body
}

func assertProblem(
	t *testing.T,
	response *httptest.ResponseRecorder,
	status int,
	detail string,
	instance string,
) problemResponse {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, status, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Errorf("Content-Type = %q, want application/problem+json", got)
	}

	problem := decodeResponse[problemResponse](t, response)
	if problem.Status != status || problem.Title != http.StatusText(status) ||
		problem.Detail != detail || problem.Instance != instance {
		t.Errorf("problem = %#v, want status %d, detail %q, instance %q", problem, status, detail, instance)
	}
	return problem
}
