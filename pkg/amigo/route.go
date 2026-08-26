package amigo

import (
	"net/http"
	"reflect"
)

type route struct {
	method string
	path   string
	status int
	input  inputMetadata
}

func (r route) pattern() string {
	return r.method + " " + r.path
}

// GET registers a typed GET handler.
func (app *App) GET[In, Out any](path string, handler Handler[In, Out], options ...RouteOption) {
	app.handle(http.MethodGet, path, handler, options...)
}

// POST registers a typed POST handler.
func (app *App) POST[In, Out any](path string, handler Handler[In, Out], options ...RouteOption) {
	app.handle(http.MethodPost, path, handler, options...)
}

// PUT registers a typed PUT handler.
func (app *App) PUT[In, Out any](path string, handler Handler[In, Out], options ...RouteOption) {
	app.handle(http.MethodPut, path, handler, options...)
}

// PATCH registers a typed PATCH handler.
func (app *App) PATCH[In, Out any](path string, handler Handler[In, Out], options ...RouteOption) {
	app.handle(http.MethodPatch, path, handler, options...)
}

// DELETE registers a typed DELETE handler.
func (app *App) DELETE[In, Out any](path string, handler Handler[In, Out], options ...RouteOption) {
	app.handle(http.MethodDelete, path, handler, options...)
}

// handle validates the route configuration once and registers its HTTP handler.
func (app *App) handle[In, Out any](
	method string,
	path string,
	handler Handler[In, Out],
	options ...RouteOption,
) {
	inputType := validateHandler(handler)
	route := newRoute(method, path, inputType, options...)
	app.mux.HandleFunc(
		route.pattern(),
		app.requestHandler(handler, route),
	)
}

func newRoute(
	method string,
	path string,
	inputType reflect.Type,
	options ...RouteOption,
) route {
	r := route{
		method: method,
		path:   path,
		status: http.StatusOK,
	}

	for _, option := range options {
		if option == nil {
			panic("amigo: route option cannot be nil")
		}
		option(&r)
	}

	r.input = buildInput(inputType, r.path)
	return r
}

func validateHandler[In, Out any](handler Handler[In, Out]) reflect.Type {
	if handler == nil {
		panic("amigo: handler cannot be nil")
	}

	inputType := reflect.TypeFor[In]()
	if inputType == nil || inputType.Kind() != reflect.Struct {
		panic("amigo: handler input must be a struct")
	}
	if body, exists := inputType.FieldByName("Body"); exists && !body.IsExported() {
		panic("amigo: Body field must be exported")
	}
	return inputType
}
