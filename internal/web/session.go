package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

const (
	browserSessionCookie = "shorty_session"
	browserSessionBytes  = 24
	browserSessionMaxAge = 365 * 24 * time.Hour
)

type browserSessionContextKey struct{}

// browserSession gives each browser profile an anonymous ownership boundary.
// Incognito profiles receive another cookie and therefore another link list.
func browserSession(next http.Handler, secure bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sessionID, ok := readBrowserSession(request)
		if !ok {
			var err error
			sessionID, err = newBrowserSessionID()
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     browserSessionCookie,
				Value:    sessionID,
				Path:     "/",
				Expires:  time.Now().Add(browserSessionMaxAge),
				MaxAge:   int(browserSessionMaxAge.Seconds()),
				HttpOnly: true,
				Secure:   secure,
				SameSite: http.SameSiteLaxMode,
			})
		}

		ctx := context.WithValue(request.Context(), browserSessionContextKey{}, sessionID)
		next.ServeHTTP(w, request.WithContext(ctx))
	})
}

func readBrowserSession(request *http.Request) (string, bool) {
	cookie, err := request.Cookie(browserSessionCookie)
	if err != nil {
		return "", false
	}
	decoded, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	return cookie.Value, err == nil && len(decoded) == browserSessionBytes
}

func newBrowserSessionID() (string, error) {
	random := make([]byte, browserSessionBytes)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(random), nil
}

func browserSessionFromContext(ctx context.Context) string {
	sessionID, _ := ctx.Value(browserSessionContextKey{}).(string)
	return sessionID
}
