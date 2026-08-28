package api

import (
	"log/slog"
	"net/http"
	"time"
)

// logRequest records one completion event after the endpoint and all inner
// middleware have finished. Keeping it at the links router boundary makes the
// log apply uniformly to every link endpoint.
func logRequest(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			startedAt := time.Now()
			next.ServeHTTP(w, request)
			logger.InfoContext(request.Context(), "request completed",
				"method", request.Method,
				"path", request.URL.Path,
				"duration", time.Since(startedAt),
			)
		})
	}
}
