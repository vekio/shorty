package links

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/amigo"
	apivalidator "github.com/vekio/shorty/internal/api/validator"
)

func TestRegisterUsesConfiguredRouter(t *testing.T) {
	httpAPI := amigo.New(amigo.WithLogger(testLogger()))
	apivalidator.Register(httpAPI)
	router := httpAPI.Group("/links", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.Header().Set("X-Links-Middleware", "applied")
			next.ServeHTTP(w, request)
		})
	})
	Register(router, newTestApplication())

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/links", nil)
	request.Header.Set("X-Shorty-Owner", "browser-a")
	httpAPI.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if got := response.Header().Get("X-Links-Middleware"); got != "applied" {
		t.Errorf("X-Links-Middleware = %q, want applied", got)
	}
}
