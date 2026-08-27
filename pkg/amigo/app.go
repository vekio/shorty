// Package amigo contains small HTTP abstractions extracted from Shorty.
package amigo

import "net/http"

// Api wraps an HTTP multiplexer with shared endpoint error handling.
type Api struct {
	mux *http.ServeMux
}

// New creates an HTTP application with an optional application error mapper.
func New() *Api {
	return &Api{
		mux: http.NewServeMux(),
	}
}

// Handler exposes the application as a standard HTTP handler.
func (app *Api) Handler() http.Handler {
	return app
}

// ServeHTTP implements http.Handler.
func (app *Api) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	app.mux.ServeHTTP(w, request)
}
