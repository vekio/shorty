package amigo

import "net/http"

// Problem is an RFC 9457 problem detail returned to an HTTP client.
type problem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`

	cause error
}

// newProblem creates a public HTTP problem.
func newProblem(status int, detail string) *problem {
	return &problem{
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	}
}

func (problem *problem) Error() string {
	return problem.Detail
}

func (problem *problem) Unwrap() error {
	return problem.cause
}
