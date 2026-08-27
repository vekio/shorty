package amigo

import (
	"fmt"
	"net/http"
)

// RouteOption configures a route.
type RouteOption func(*route)

// Status sets the successful response status.
func Status(status int) RouteOption {
	if status < http.StatusOK || status > 299 {
		panic(fmt.Sprintf("amigo: success status must be between 200 and 299, got %d", status))
	}

	return func(route *route) {
		route.status = status
	}
}

// MaxBodyBytes sets the maximum accepted JSON body size. Zero disables the limit.
func MaxBodyBytes(limit int64) RouteOption {
	if limit < 0 {
		panic("amigo: maximum body size cannot be negative")
	}

	return func(route *route) {
		route.maxBodyBytes = limit
	}
}

// MapError associates an application error with an HTTP error status.
func MapError(target error, status int) RouteOption {
	if target == nil {
		panic("amigo: mapped error cannot be nil")
	}
	if status < http.StatusBadRequest || status > 599 {
		panic(fmt.Sprintf("amigo: error status must be between 400 and 599, got %d", status))
	}

	return func(route *route) {
		route.errors = append(route.errors, errorMapping{
			target: target,
			status: status,
		})
	}
}
