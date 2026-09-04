package shorty

import (
	"fmt"
	"net/http"
)

// FieldError identifies one invalid request value.
type FieldError struct {
	Location string `json:"location"`
	Message  string `json:"message"`
}

// Problem is the error representation returned by the Shorty API.
type Problem struct {
	Title    string       `json:"title"`
	Status   int          `json:"status"`
	Detail   string       `json:"detail"`
	Instance string       `json:"instance"`
	Errors   []FieldError `json:"errors,omitempty"`
}

// APIError reports a non-successful response from the Shorty API.
type APIError struct {
	StatusCode int
	Problem    Problem
}

func (err *APIError) Error() string {
	detail := err.Problem.Detail
	if detail == "" {
		detail = http.StatusText(err.StatusCode)
	}
	return fmt.Sprintf("Shorty API returned HTTP %d: %s", err.StatusCode, detail)
}
