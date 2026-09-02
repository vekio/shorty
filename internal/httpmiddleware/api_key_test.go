package httpmiddleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/shorty/internal/auth"
)

type apiKeyAuthenticatorFunc func(context.Context, string) error

func (authenticate apiKeyAuthenticatorFunc) AuthenticateToken(ctx context.Context, token string) error {
	return authenticate(ctx, token)
}

func TestRequireAPIKeyAuthenticatesBearerToken(t *testing.T) {
	var authenticatedToken string
	handler := RequireAPIKey(apiKeyAuthenticatorFunc(func(_ context.Context, token string) error {
		authenticatedToken = token
		return nil
	}))(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
	request.Header.Set("Authorization", "Bearer shorty_secret")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent || authenticatedToken != "shorty_secret" {
		t.Errorf("status = %d, token = %q", response.Code, authenticatedToken)
	}
}

func TestRequireAPIKeyMapsAuthenticationErrors(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
		authError     error
		wantStatus    int
	}{
		{name: "missing header", wantStatus: http.StatusUnauthorized},
		{name: "wrong scheme", authorization: "Basic abc", wantStatus: http.StatusUnauthorized},
		{name: "invalid key", authorization: "Bearer invalid", authError: auth.ErrInvalidAPIKey, wantStatus: http.StatusUnauthorized},
		{name: "repository failure", authorization: "Bearer shorty_secret", authError: errors.New("database unavailable"), wantStatus: http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := RequireAPIKey(apiKeyAuthenticatorFunc(func(context.Context, string) error {
				return test.authError
			}))(http.NotFoundHandler())
			request := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
			request.Header.Set("Authorization", test.authorization)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Errorf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if response.Header().Get("Content-Type") != "application/problem+json" {
				t.Errorf("Content-Type = %q", response.Header().Get("Content-Type"))
			}
			if test.wantStatus == http.StatusUnauthorized && response.Header().Get("WWW-Authenticate") != "Bearer" {
				t.Errorf("WWW-Authenticate = %q", response.Header().Get("WWW-Authenticate"))
			}
		})
	}
}
