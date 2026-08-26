package amigo

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerWritesTypedJSONResponse(t *testing.T) {
	type output struct {
		Name string `json:"name"`
	}
	app := New()
	app.POST("/things", func(context.Context, struct{}) (output, error) {
		return output{Name: "shorty"}, nil
	}, Status(http.StatusCreated))
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/things", nil))

	if response.Code != http.StatusCreated || response.Body.String() != `{"name":"shorty"}` {
		t.Errorf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q", response.Header().Get("Content-Type"))
	}
}

func TestHandlerSuppressesBodyForEmptySuccessStatus(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusResetContent} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			app := New()
			app.POST("/things", func(context.Context, struct{}) (string, error) {
				return "ignored", nil
			}, Status(status))
			response := httptest.NewRecorder()
			app.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/things", nil))
			if response.Code != status || response.Body.Len() != 0 {
				t.Errorf("status = %d, body = %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestHandlerReturnsHandlerError(t *testing.T) {
	errKnown := errors.New("known")
	app := New()
	app.MapErrors(func(err error) (*Problem, bool) {
		return BadRequest("mapped"), errors.Is(err, errKnown)
	})
	app.GET("/things", func(context.Context, struct{}) (struct{}, error) {
		return struct{}{}, errKnown
	})
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/things", nil))
	assertProblem(t, response, http.StatusBadRequest, "mapped")
}

func TestHandlerHidesResponseEncodingError(t *testing.T) {
	type output struct {
		Unsupported func() `json:"unsupported"`
	}
	app := New()
	app.GET("/things", func(context.Context, struct{}) (output, error) {
		return output{Unsupported: func() {}}, nil
	})
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/things", nil))
	assertProblem(t, response, http.StatusInternalServerError, "internal server error")
	if strings.Contains(response.Body.String(), "unsupported Go type") {
		t.Errorf("encoding error leaked: %s", response.Body.String())
	}
}
