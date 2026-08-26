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
