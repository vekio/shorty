package amigo

import (
	"encoding/json/v2"
	"errors"
	"log/slog"
	"net/http"
)

// ErrorMapper translates an application error into a safe public problem.
// Returning false delegates the error to the next mapper.
type ErrorMapper func(error) (*Problem, bool)

// Problem is an RFC 9457 problem detail returned to an HTTP client.
type Problem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`

	err error
}

// NewProblem creates a public HTTP problem.
func NewProblem(status int, detail string) *Problem {
	return &Problem{
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	}
}

func BadRequest(detail string) *Problem {
	return NewProblem(http.StatusBadRequest, detail)
}

func NotFound(detail string) *Problem {
	return NewProblem(http.StatusNotFound, detail)
}

func UnsupportedMediaType(detail string) *Problem {
	return NewProblem(http.StatusUnsupportedMediaType, detail)
}

func ContentTooLarge(detail string) *Problem {
	return NewProblem(http.StatusRequestEntityTooLarge, detail)
}

func InternalServerError(detail string) *Problem {
	return NewProblem(http.StatusInternalServerError, detail)
}

// WrapProblem associates a private cause with a public problem.
func WrapProblem(err error, problem *Problem) *Problem {
	if problem == nil {
		problem = InternalServerError("internal server error")
	}

	wrapped := cloneProblem(problem)
	wrapped.err = err
	return wrapped
}

func (problem *Problem) Error() string {
	if problem == nil {
		return http.StatusText(http.StatusInternalServerError)
	}
	if problem.Detail != "" {
		return problem.Detail
	}
	return problem.Title
}

// Unwrap returns the private cause without exposing it in the JSON response.
func (problem *Problem) Unwrap() error {
	if problem == nil {
		return nil
	}
	return problem.err
}

func (app *App) writeError(w http.ResponseWriter, request *http.Request, err error) {
	problem := app.problemFor(err)
	if problem.Status >= http.StatusInternalServerError {
		slog.ErrorContext(request.Context(), "request failed",
			"method", request.Method,
			"path", request.URL.Path,
			"error", err,
		)
	}

	if problem.Instance == "" {
		problem.Instance = request.URL.Path
	}
	data, marshalErr := json.Marshal(problem)
	if marshalErr != nil {
		problem = InternalServerError("internal server error")
		problem.Instance = request.URL.Path
		data, _ = json.Marshal(problem)
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(problem.Status)
	_, _ = w.Write(data)
}

func (app *App) problemFor(err error) *Problem {
	var problem *Problem
	if errors.As(err, &problem) && validProblem(problem) {
		return cloneProblem(problem)
	}

	for _, mapper := range app.errorMappers {
		mapped, ok := mapper(err)
		if ok && validProblem(mapped) {
			return WrapProblem(err, mapped)
		}
	}

	return WrapProblem(err, InternalServerError("internal server error"))
}

func validProblem(problem *Problem) bool {
	return problem != nil && problem.Status >= 400 && problem.Status <= 599
}

func cloneProblem(problem *Problem) *Problem {
	clone := *problem
	if clone.Title == "" {
		clone.Title = http.StatusText(clone.Status)
	}
	return &clone
}
