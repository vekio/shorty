package amigo

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
)

// Handler turns a typed request into a typed response.
type Handler[In, Out any] func(context.Context, In) (Out, error)

func (app *App) requestHandler[In, Out any](
	handler Handler[In, Out],
	configuredRoute route,
) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		app.processRequest(w, request, handler, configuredRoute)
	}
}

func (app *App) processRequest[In, Out any](
	w http.ResponseWriter,
	request *http.Request,
	handler Handler[In, Out],
	configuredRoute route,
) {
	input, err := bindInput[In](w, request, configuredRoute.input, app.maxBodyBytes)
	if err != nil {
		app.writeError(w, request, err)
		return
	}

	output, err := handler(request.Context(), input)
	if err != nil {
		app.writeError(w, request, err)
		return
	}

	if err := writeOutput(w, configuredRoute.status, output); err != nil {
		app.writeError(w, request, err)
	}
}

func writeOutput[Out any](w http.ResponseWriter, status int, output Out) error {
	if status == http.StatusNoContent || status == http.StatusResetContent {
		w.WriteHeader(status)
		return nil
	}

	data, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
	return nil
}
