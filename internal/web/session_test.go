package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBrowserSessionCreatesAndReusesAnonymousCookie(t *testing.T) {
	handler := browserSession(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if browserSessionFromContext(request.Context()) == "" {
			t.Error("request context has no browser session")
		}
	}), false)

	firstResponse := httptest.NewRecorder()
	handler.ServeHTTP(firstResponse, httptest.NewRequest(http.MethodGet, "/", nil))
	cookies := firstResponse.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookies = %#v, want one session cookie", cookies)
	}
	if !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteLaxMode {
		t.Errorf("cookie security attributes = %#v", cookies[0])
	}

	secondRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	secondRequest.AddCookie(cookies[0])
	secondResponse := httptest.NewRecorder()
	handler.ServeHTTP(secondResponse, secondRequest)
	if got := secondResponse.Header().Get("Set-Cookie"); got != "" {
		t.Errorf("reused session emitted Set-Cookie %q", got)
	}
}

func TestBrowserProfilesReceiveDifferentSessions(t *testing.T) {
	handler := browserSession(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), false)
	first := httptest.NewRecorder()
	second := httptest.NewRecorder()

	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/", nil))
	handler.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/", nil))

	firstSession := first.Result().Cookies()[0].Value
	secondSession := second.Result().Cookies()[0].Value
	if firstSession == secondSession {
		t.Error("independent browser profiles received the same session")
	}
}
