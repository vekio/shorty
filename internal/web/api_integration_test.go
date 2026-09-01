package web

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	apibootstrap "github.com/vekio/shorty/internal/bootstrap/api"
	shortyconfig "github.com/vekio/shorty/internal/config"
	apiconfig "github.com/vekio/shorty/internal/config/api"
	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

func TestWebCreatesAndResolvesLinkThroughAPI(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	apiRuntime, err := apibootstrap.New(apiconfig.Config{
		Address:  ":8080",
		ShortURL: "https://sho.rt",
		Logger:   shortyconfig.DefaultLoggerConfig(),
	})
	if err != nil {
		t.Fatalf("bootstrap API: %v", err)
	}
	apiServer := httptest.NewServer(apiRuntime.Handler)
	defer apiServer.Close()

	client, err := shortysdk.NewClient(apiServer.URL, apiServer.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	httpWeb := New(client, logger, "https://sho.rt")

	form := url.Values{"origin_url": {"https://example.com/docs"}}
	createRequest := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(form.Encode()))
	createRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createRequest.Header.Set("HX-Request", "true")
	createResponse := httptest.NewRecorder()
	httpWeb.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusOK || !strings.Contains(createResponse.Body.String(), "https://example.com/docs") {
		t.Fatalf("create response = %d %s", createResponse.Code, createResponse.Body.String())
	}
	sessionCookie := createResponse.Result().Cookies()[0]
	ctx := shortysdk.WithOwner(t.Context(), sessionCookie.Value)

	normalRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	normalRequest.AddCookie(sessionCookie)
	normalResponse := httptest.NewRecorder()
	httpWeb.ServeHTTP(normalResponse, normalRequest)
	if !strings.Contains(normalResponse.Body.String(), "https://example.com/docs") {
		t.Errorf("normal browser does not see its created link")
	}
	incognitoResponse := httptest.NewRecorder()
	httpWeb.ServeHTTP(incognitoResponse, httptest.NewRequest(http.MethodGet, "/", nil))
	if strings.Contains(incognitoResponse.Body.String(), "https://example.com/docs") ||
		!strings.Contains(incognitoResponse.Body.String(), "No links yet.") {
		t.Errorf("incognito browser sees another session's links: %s", incognitoResponse.Body.String())
	}

	page, err := client.ListLinks(ctx, shortysdk.ListOptions{})
	if err != nil || len(page.Links) != 1 {
		t.Fatalf("API page = %#v, error = %v", page, err)
	}
	code := page.Links[0].Code

	if !strings.Contains(createResponse.Body.String(), "https://sho.rt/r/"+code) {
		t.Errorf("Web response does not contain the public short URL")
	}
	redirectResponse := httptest.NewRecorder()
	httpWeb.ServeHTTP(
		redirectResponse,
		httptest.NewRequest(http.MethodGet, "/r/"+code, nil),
	)
	if redirectResponse.Code != http.StatusFound ||
		redirectResponse.Header().Get("Location") != "https://example.com/docs" {
		t.Fatalf("Web redirect = %d, Location = %q", redirectResponse.Code, redirectResponse.Header().Get("Location"))
	}

	link, err := client.GetLink(ctx, code)
	if err != nil {
		t.Fatalf("GetLink() error = %v", err)
	}
	if link.Visits != 1 {
		t.Errorf("visits = %d, want 1", link.Visits)
	}

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/links/"+code, nil)
	deleteRequest.Header.Set("HX-Request", "true")
	deleteRequest.AddCookie(sessionCookie)
	deleteResponse := httptest.NewRecorder()
	httpWeb.ServeHTTP(deleteResponse, deleteRequest)
	if deleteResponse.Code != http.StatusOK || !strings.Contains(deleteResponse.Body.String(), "No links yet.") {
		t.Fatalf("delete response = %d %s", deleteResponse.Code, deleteResponse.Body.String())
	}
}
