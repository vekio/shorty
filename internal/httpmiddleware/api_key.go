package httpmiddleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/vekio/shorty/internal/auth"
)

// APIKeyAuthenticator validates a clear-text API token.
type APIKeyAuthenticator interface {
	AuthenticateToken(context.Context, string) error
}

// RequireAPIKey protects a handler with an active opaque Bearer API key.
func RequireAPIKey(authenticator APIKeyAuthenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			token, valid := bearerToken(request.Header.Get("Authorization"))
			if !valid {
				writeUnauthorized(w)
				return
			}
			if err := authenticator.AuthenticateToken(request.Context(), token); err != nil {
				if errors.Is(err, auth.ErrInvalidAPIKey) {
					writeUnauthorized(w)
				} else {
					writeProblem(w, http.StatusInternalServerError, "API key authentication failed")
				}
				return
			}
			next.ServeHTTP(w, request)
		})
	}
}

func bearerToken(authorization string) (string, bool) {
	scheme, token, found := strings.Cut(authorization, " ")
	return token, found && strings.EqualFold(scheme, "Bearer") && token != ""
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	writeProblem(w, http.StatusUnauthorized, "a valid API key is required")
}

func writeProblem(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		Title  string `json:"title"`
		Status int    `json:"status"`
		Detail string `json:"detail"`
	}{
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	})
}
