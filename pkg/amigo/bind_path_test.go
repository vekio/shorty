package amigo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type upperText string

func (value *upperText) UnmarshalText(text []byte) error {
	*value = upperText(strings.ToUpper(string(text)))
	return nil
}

func TestBindPathParametersConvertsSupportedTypes(t *testing.T) {
	type input struct {
		Text  string    `path:"text"`
		Int   int       `path:"int"`
		Uint  uint      `path:"uint"`
		Bool  bool      `path:"bool"`
		Float float64   `path:"float"`
		Upper upperText `path:"upper"`
	}

	var received input
	app := New()
	app.GET("/values/{text}/{int}/{uint}/{bool}/{float}/{upper}", func(_ context.Context, input input) (struct{}, error) {
		received = input
		return struct{}{}, nil
	})
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/values/go/-2/3/true/1.5/amigo", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", response.Code, response.Body.String())
	}
	if received.Text != "go" || received.Int != -2 || received.Uint != 3 || !received.Bool || received.Float != 1.5 || received.Upper != "AMIGO" {
		t.Errorf("input = %#v", received)
	}
}

func TestBindPathParametersRejectsWrongType(t *testing.T) {
	type input struct {
		ID int `path:"id"`
	}
	app := New()
	app.GET("/things/{id}", func(context.Context, input) (struct{}, error) { return struct{}{}, nil })
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/things/not-an-int", nil))
	assertProblem(t, response, http.StatusBadRequest, `path parameter "id" must be a valid int`)
}
