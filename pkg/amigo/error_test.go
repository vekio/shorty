package amigo

import (
	"context"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProblemConstructors(t *testing.T) {
	tests := []struct {
		name    string
		problem *Problem
		status  int
	}{
		{name: "bad request", problem: BadRequest("detail"), status: http.StatusBadRequest},
		{name: "not found", problem: NotFound("detail"), status: http.StatusNotFound},
		{name: "unsupported media", problem: UnsupportedMediaType("detail"), status: http.StatusUnsupportedMediaType},
		{name: "content too large", problem: ContentTooLarge("detail"), status: http.StatusRequestEntityTooLarge},
		{name: "internal", problem: InternalServerError("detail"), status: http.StatusInternalServerError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.problem.Status != test.status || test.problem.Title != http.StatusText(test.status) || test.problem.Error() != "detail" {
				t.Errorf("problem = %#v", test.problem)
			}
		})
	}
}

func TestProblemErrorFallsBackToTitle(t *testing.T) {
	problem := NewProblem(http.StatusConflict, "")
	if problem.Error() != http.StatusText(http.StatusConflict) {
		t.Errorf("Error() = %q", problem.Error())
	}
	var nilProblem *Problem
	if nilProblem.Error() != http.StatusText(http.StatusInternalServerError) || nilProblem.Unwrap() != nil {
		t.Error("nil problem fallback is invalid")
	}
}

func TestWrapProblemPreservesPrivateCause(t *testing.T) {
	cause := errors.New("database failed")
	problem := WrapProblem(cause, BadRequest("public detail"))
	if !errors.Is(problem, cause) {
		t.Error("wrapped problem does not preserve its cause")
	}
	data, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("marshal problem: %v", err)
	}
	if strings.Contains(string(data), cause.Error()) {
		t.Errorf("private cause leaked in JSON: %s", data)
	}
}

func TestWrapProblemUsesInternalFallbackForNilProblem(t *testing.T) {
	cause := errors.New("failure")
	problem := WrapProblem(cause, nil)
	if problem.Status != http.StatusInternalServerError || !errors.Is(problem, cause) {
		t.Errorf("problem = %#v", problem)
	}
}

func TestErrorMapperPublishesKnownError(t *testing.T) {
	errKnown := errors.New("known failure")
	app := New()
	app.MapErrors(
		func(error) (*Problem, bool) { return nil, false },
		func(err error) (*Problem, bool) {
			return BadRequest(err.Error()), errors.Is(err, errKnown)
		},
	)
	app.GET("/failure", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, errKnown })
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/failure", nil))
	assertProblem(t, response, http.StatusBadRequest, errKnown.Error())
}

func TestDirectProblemIsPublic(t *testing.T) {
	app := New()
	app.GET("/failure", func(context.Context, struct{}) (struct{}, error) {
		return struct{}{}, NotFound("missing")
	})
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/failure", nil))
	assertProblem(t, response, http.StatusNotFound, "missing")
}

func TestUnexpectedErrorIsNotExposed(t *testing.T) {
	app := New()
	app.GET("/failure", func(context.Context, struct{}) (struct{}, error) {
		return struct{}{}, errors.New("database password is secret")
	})
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/failure", nil))
	assertProblem(t, response, http.StatusInternalServerError, "internal server error")
	if strings.Contains(response.Body.String(), "database password") {
		t.Errorf("private error leaked in body: %s", response.Body.String())
	}
}

func TestInvalidMappedProblemFallsBackToInternalError(t *testing.T) {
	app := New()
	app.MapErrors(func(error) (*Problem, bool) { return NewProblem(http.StatusOK, "invalid"), true })
	app.GET("/failure", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, errors.New("failure") })
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/failure", nil))
	assertProblem(t, response, http.StatusInternalServerError, "internal server error")
}
