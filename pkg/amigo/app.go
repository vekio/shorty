// Package amigo provides a small typed layer over net/http.
package amigo

import (
	"net/http"
)

const defaultMaxBodyBytes int64 = 1 << 20

// App registers typed handlers in an http.ServeMux.
type App struct {
	mux          *http.ServeMux
	maxBodyBytes int64
	errorMappers []ErrorMapper
}

// New creates an empty application.
func New() *App {
	return &App{
		mux:          http.NewServeMux(),
		maxBodyBytes: defaultMaxBodyBytes,
	}
}

// SetMaxBodyBytes sets the maximum accepted JSON body size. Zero disables the limit.
func (app *App) SetMaxBodyBytes(limit int64) {
	if limit < 0 {
		panic("amigo: maximum body size cannot be negative")
	}
	app.maxBodyBytes = limit
}

// MapErrors adds application error mappings. The first matching mapper wins.
func (app *App) MapErrors(mappers ...ErrorMapper) {
	for _, mapper := range mappers {
		if mapper == nil {
			panic("amigo: error mapper cannot be nil")
		}
		app.errorMappers = append(app.errorMappers, mapper)
	}
}

// Handler exposes the application as a standard HTTP handler.
func (app *App) Handler() http.Handler {
	return app
}

// ServeHTTP implements http.Handler.
func (app *App) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	app.mux.ServeHTTP(w, request)
}
