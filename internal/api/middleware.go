package api

import (
	"log/slog"
	"net/http"
	"time"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		next.ServeHTTP(w, request)
		slog.InfoContext(request.Context(), "request completed",
			"method", request.Method,
			"path", request.URL.Path,
			"duration", time.Since(startedAt),
		)
	})
}
