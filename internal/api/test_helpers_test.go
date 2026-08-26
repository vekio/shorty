package api

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/visitlink"
	"github.com/vekio/shorty/internal/infra/memory"
	"github.com/vekio/shorty/pkg/amigo"
)

func newTestAPI() *amigo.App {
	repository := memory.NewLinkRepository()
	application := app.Application{
		Commands: app.Commands{
			CreateLink: createlink.NewCreateLinkHandler(repository),
			VisitLink:  visitlink.NewVisitLinkHandler(repository),
		},
		Queries: app.Queries{GetLink: getlink.NewGetLinkHandler(repository)},
	}
	return New(application)
}

func createLink(t *testing.T, httpAPI http.Handler, originURL string) CreateLinkOutput {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(`{"origin_url":"`+originURL+`"}`))
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
