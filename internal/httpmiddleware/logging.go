// Package httpmiddleware contains transport middleware shared by Shorty's
// independent HTTP servers.
package httpmiddleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LogRequests records one completion event after the wrapped handler finishes.
func LogRequests(logger *slog.Logger) func(http.Handler) http.Handler {
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
