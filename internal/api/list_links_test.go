package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/listlinks"
)

type listLinksHandlerStub struct {
	err error
}

func (handler listLinksHandlerStub) Handle(context.Context, listlinks.ListLinksQuery) (listlinks.ListLinksResult, error) {
	return listlinks.ListLinksResult{}, handler.err
}

func TestListLinksReturnsEmptyCollection(t *testing.T) {
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	var output ListLinksOutput
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if output.Links == nil || len(output.Links) != 0 {
		t.Errorf("links = %#v, want non-nil empty collection", output.Links)
	}
}

func TestListLinksReturnsCreatedLinks(t *testing.T) {
	httpAPI := newTestAPI()
	first := createLink(t, httpAPI, "https://example.com/first")
	second := createLink(t, httpAPI, "https://example.com/second")

	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	var output ListLinksOutput
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(output.Links) != 2 {
		t.Fatalf("links length = %d, want 2", len(output.Links))
	}
	if output.Links[0].Code != first.Code || output.Links[1].Code != second.Code {
		t.Errorf("links = %#v, want creation order", output.Links)
	}
}

func TestListLinksReturnsInternalError(t *testing.T) {
	httpAPI := New(app.Application{
		Queries: app.Queries{ListLinks: listLinksHandlerStub{err: errors.New("repository failed")}},
	}, testLogger())
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links", nil))

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusInternalServerError, response.Body.String())
	}
}
