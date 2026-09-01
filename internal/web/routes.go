package web

import (
	"log/slog"
	"net/http"
)

func registerRoutes(
	mux *http.ServeMux,
	api shortyAPI,
	logger *slog.Logger,
	shortURL string,
) {
	handler := newHandler(api, shortURL)

	// /{$} is ServeMux's exact-root pattern. Using / would also match every
	// otherwise unhandled path and would turn genuine 404s into the home page.
	mux.HandleFunc("GET /{$}", handleError(logger, handler.Home))
	mux.HandleFunc("POST /links", handleError(logger, handler.CreateLink))
	mux.HandleFunc("DELETE /links/{code}", handleError(logger, handler.DeleteLink))
	mux.HandleFunc("GET /r/{code}", handleError(logger, handler.ResolveLink))
	mux.HandleFunc("GET /{path...}", handleError(logger, handler.NotFound))
}
