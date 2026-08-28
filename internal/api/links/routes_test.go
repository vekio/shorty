package links

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/amigo"
	apivalidator "github.com/vekio/shorty/internal/api/validator"
)

func TestRegisterAppliesMiddlewareToLinksRouter(t *testing.T) {
	httpAPI := amigo.New(amigo.WithLogger(testLogger()))
	apivalidator.Register(httpAPI)
	Register(httpAPI, newTestApplication(), func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.Header().Set("X-Links-Middleware", "applied")
			next.ServeHTTP(w, request)
		})
	})

	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/links", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if got := response.Header().Get("X-Links-Middleware"); got != "applied" {
		t.Errorf("X-Links-Middleware = %q, want applied", got)
	}
}
