package root

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

func TestRegisterAppliesMiddlewareToRootRoute(t *testing.T) {
	httpAPI := amigo.New()
	Register(
		httpAPI,
		app.Application{Commands: app.Commands{
			ResolveLink: &resolveLinkHandlerStub{
				result: resolvelink.ResolveLinkResult{OriginURL: "https://example.com"},
			},
		}},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				w.Header().Set("X-Root-Middleware", "applied")
				next.ServeHTTP(w, request)
			})
		},
	)
	response := httptest.NewRecorder()
	httpAPI.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/abc123", nil))

	if got := response.Header().Get("X-Root-Middleware"); got != "applied" {
		t.Errorf("X-Root-Middleware = %q, want applied", got)
	}
}
