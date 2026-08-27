package amigo

import (
	"context"
	"net/http"
)

// EndpointFunc receives a typed input and returns the typed response body.
type EndpointFunc[In, Out any] func(context.Context, In) (Out, error)

// RawEndpointFunc is the escape hatch for endpoints that need direct access to
// net/http, such as streaming responses or file downloads.
type RawEndpointFunc func(http.ResponseWriter, *http.Request) error

func (app *Api) handler[In, Out any](
	route route,
	inputMetadata inputMetadata,
	outputMetadata outputMetadata,
	endpoint EndpointFunc[In, Out],
) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		limitRequestBody(w, request, route.maxBodyBytes)

		input, err := bindInput[In](request, inputMetadata)
		if err != nil {
			writeError(w, request, route, err)
			return
		}

		output, err := endpoint(request.Context(), input)
		if err != nil {
			writeError(w, request, route, err)
			return
		}

		if err := writeOutput(w, route.status, output, outputMetadata); err != nil {
			writeError(w, request, route, err)
		}
	}
}

func (app *Api) rawHandler(route route, endpoint RawEndpointFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		limitRequestBody(w, request, route.maxBodyBytes)
		if err := endpoint(w, request); err != nil {
			writeError(w, request, route, err)
		}
	}
}

func limitRequestBody(w http.ResponseWriter, request *http.Request, limit int64) {
	if limit > 0 && request.Body != nil {
		request.Body = http.MaxBytesReader(w, request.Body, limit)
	}
}

// GET registers a typed GET endpoint.
func (app *Api) GET[In, Out any](path string, endpoint EndpointFunc[In, Out], options ...RouteOption) {
	app.handle(http.MethodGet, path, endpoint, options...)
}

// POST registers a typed POST endpoint.
func (app *Api) POST[In, Out any](path string, endpoint EndpointFunc[In, Out], options ...RouteOption) {
	app.handle(http.MethodPost, path, endpoint, options...)
}

// RAW registers an endpoint that owns its complete net/http response.
func (app *Api) RAW(method string, path string, endpoint RawEndpointFunc, options ...RouteOption) {
	if endpoint == nil {
		panic("amigo: raw endpoint cannot be nil")
	}
	route := newRoute(method, path, options...)
	app.mux.HandleFunc(route.pattern(), app.rawHandler(route, endpoint))
}

func (app *Api) handle[In, Out any](
	method string,
	path string,
	endpoint EndpointFunc[In, Out],
	options ...RouteOption,
) {
	if endpoint == nil {
		panic("amigo: endpoint cannot be nil")
	}
	route := newRoute(method, path, options...)
	inputMetadata := buildInputMetadata[In](route.path)
	outputMetadata := buildOutputMetadata[Out]()
	app.mux.HandleFunc(route.pattern(), app.handler(route, inputMetadata, outputMetadata, endpoint))
}
