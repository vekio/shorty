package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	shortyconfig "github.com/vekio/shorty/internal/config"
	webconfig "github.com/vekio/shorty/internal/config/web"
)

func TestNewComposesWebProcessWithAPIClient(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/links" {
			t.Errorf("API path = %q, want /links", request.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"links":[],"total":0,"limit":20,"offset":0}`))
	}))
	defer apiServer.Close()

	runtime, err := New(webconfig.Config{
		Address:  ":3000",
		ShortURL: "https://sho.rt",
		APIURL:   apiServer.URL,
		Logger:   shortyconfig.DefaultLoggerConfig(),
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if runtime.Handler == nil || runtime.Logger == nil {
		t.Fatalf("runtime = %#v, want complete Web process", runtime)
	}

	response := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "Short links") {
		t.Errorf("response = %d %s", response.Code, response.Body.String())
	}
}
