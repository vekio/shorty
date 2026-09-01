package web

import (
	"log/slog"
	"net/http"
)

type endpoint func(http.ResponseWriter, *http.Request) error

// handleError keeps endpoint failures out of public responses while retaining
// their cause in server logs. Endpoints return errors before writing headers.
func handleError(logger *slog.Logger, endpoint endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if err := endpoint(w, request); err != nil {
			logger.ErrorContext(request.Context(), "request failed",
				"method", request.Method,
				"path", request.URL.Path,
				"error", err,
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
