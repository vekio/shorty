package links

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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
	output := decodeResponse[ListLinksOutput](t, response)
	if output.Links == nil || len(output.Links) != 0 {
		t.Errorf("links = %#v, want non-nil empty collection", output.Links)
	}
	if output.Total != 0 || output.Limit != listlinks.DefaultLimit || output.Offset != 0 {
		t.Errorf("pagination = total %d, limit %d, offset %d", output.Total, output.Limit, output.Offset)
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
	output := decodeResponse[ListLinksOutput](t, response)
	if len(output.Links) != 2 {
		t.Fatalf("links length = %d, want 2", len(output.Links))
	}
	if output.Links[0].Code != first.Code || output.Links[1].Code != second.Code {
		t.Errorf("links = %#v, want creation order", output.Links)
	}
	if output.Total != 2 || output.Limit != listlinks.DefaultLimit || output.Offset != 0 {
		t.Errorf("pagination = total %d, limit %d, offset %d", output.Total, output.Limit, output.Offset)
	}
}

func TestListLinksReturnsRequestedPage(t *testing.T) {
	httpAPI := newTestAPI()
	createLink(t, httpAPI, "https://example.com/first")
	second := createLink(t, httpAPI, "https://example.com/second")
	createLink(t, httpAPI, "https://example.com/third")

	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links?limit=1&offset=1", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	output := decodeResponse[ListLinksOutput](t, response)
	if len(output.Links) != 1 || output.Links[0].Code != second.Code {
		t.Errorf("links = %#v, want second created link", output.Links)
	}
	if output.Total != 3 || output.Limit != 1 || output.Offset != 1 {
		t.Errorf("pagination = total %d, limit %d, offset %d", output.Total, output.Limit, output.Offset)
	}
}

func TestListLinksValidatesPagination(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		location string
		message  string
	}{
		{name: "zero limit", query: "limit=0", location: "query.limit", message: listlinks.ErrInvalidLimit.Error()},
		{name: "limit above maximum", query: "limit=101", location: "query.limit", message: listlinks.ErrInvalidLimit.Error()},
		{name: "negative offset", query: "offset=-1", location: "query.offset", message: listlinks.ErrInvalidOffset.Error()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			newTestAPI().ServeHTTP(
				response,
				httptest.NewRequest(http.MethodGet, "/links?"+test.query, nil),
			)

			problem := assertProblem(t, response, http.StatusUnprocessableEntity, "request validation failed", "/links")
			want := fieldErrorResponse{Location: test.location, Message: test.message}
			if len(problem.Errors) != 1 || problem.Errors[0] != want {
				t.Errorf("errors = %#v, want %#v", problem.Errors, []fieldErrorResponse{want})
			}
		})
	}
}

func TestListLinksRejectsMalformedPagination(t *testing.T) {
	response := httptest.NewRecorder()
	newTestAPI().ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/links?limit=many", nil),
	)

	assertProblem(t, response, http.StatusBadRequest, `invalid query parameter "limit"`, "/links")
}

func TestListLinksReturnsInternalError(t *testing.T) {
	httpAPI := newTestAPIWithApplication(app.Application{
		Queries: app.Queries{ListLinks: listLinksHandlerStub{err: errors.New("repository failed")}},
	})
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links", nil))

	problem := assertProblem(t, response, http.StatusInternalServerError, "internal server error", "/links")
	if strings.Contains(response.Body.String(), "repository failed") || len(problem.Errors) != 0 {
		t.Errorf("problem exposes private error information: %s", response.Body.String())
	}
}

func TestListLinksMapsApplicationPaginationError(t *testing.T) {
	httpAPI := newTestAPIWithApplication(app.Application{
		Queries: app.Queries{ListLinks: listLinksHandlerStub{err: listlinks.ErrInvalidLimit}},
	})
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links", nil))

	assertProblem(
		t,
		response,
		http.StatusBadRequest,
		listlinks.ErrInvalidLimit.Error(),
		"/links",
	)
}
