package amigo

import (
	"errors"
	"net/http"
)

const defaultMaxBodyBytes int64 = 1 << 20

type errorMapping struct {
	target error
	status int
}

type route struct {
	method       string
	path         string
	status       int
	maxBodyBytes int64
	errors       []errorMapping
}

func newRoute(method string, path string, options ...RouteOption) route {
	r := route{
		method:       method,
		path:         path,
		status:       http.StatusOK,
		maxBodyBytes: defaultMaxBodyBytes,
	}

	for _, option := range options {
		if option == nil {
			panic("amigo: route option cannot be nil")
		}
		option(&r)
	}

	return r
}

func (r route) pattern() string {
	return r.method + " " + r.path
}

func (r route) problem(err error) *problem {
	if direct, ok := errors.AsType[*problem](err); ok {
		clone := *direct
		return &clone
	}

	for _, mapping := range r.errors {
		if errors.Is(err, mapping.target) {
			problem := newProblem(mapping.status, mapping.target.Error())
			problem.cause = err
			return problem
		}
	}

	problem := newProblem(http.StatusInternalServerError, "internal server error")
	problem.cause = err
	return problem
}
