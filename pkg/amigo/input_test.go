package amigo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBindInputDecodesJSON(t *testing.T) {
	type inputBody struct {
		Name string `json:"name"`
	}
	request := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(`{"name":"shorty"}`))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	input, err := bindInput[inputBody](request, buildInputMetadata[inputBody]("/things"))

	if err != nil {
		t.Fatalf("bindInput() error = %v", err)
	}
	if input.Name != "shorty" {
		t.Errorf("name = %q, want %q", input.Name, "shorty")
	}
}

func TestBindInputAllowsMissingBody(t *testing.T) {
	type inputBody struct{ Name string }
	request := httptest.NewRequest(http.MethodGet, "/things", nil)

	input, err := bindInput[inputBody](request, buildInputMetadata[inputBody]("/things"))

	if err != nil {
		t.Fatalf("bindInput() error = %v", err)
	}
	if input.Name != "" {
		t.Errorf("name = %q, want empty", input.Name)
	}
}

func TestBindInputRejectsInvalidJSONRequests(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		status      int
	}{
		{name: "missing content type", body: `{}`, status: http.StatusUnsupportedMediaType},
		{name: "wrong content type", body: `{}`, contentType: "text/plain", status: http.StatusUnsupportedMediaType},
		{name: "malformed JSON", body: `{"name":`, contentType: "application/json", status: http.StatusBadRequest},
		{name: "unknown member", body: `{"unknown":true}`, contentType: "application/json", status: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			type inputBody struct {
				Name string `json:"name"`
			}
			request := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(test.body))
			if test.contentType != "" {
				request.Header.Set("Content-Type", test.contentType)
			}

			_, err := bindInput[inputBody](request, buildInputMetadata[inputBody]("/things"))
			problem, ok := errors.AsType[*problem](err)
			if !ok {
				t.Fatalf("error = %T, want *problem", err)
			}
			if problem.Status != test.status {
				t.Errorf("status = %d, want %d", problem.Status, test.status)
			}
		})
	}
}

func TestBindInputBindsPathParameter(t *testing.T) {
	type inputBody struct {
		ID int `path:"id" json:"-"`
	}
	request := httptest.NewRequest(http.MethodGet, "/things/42", nil)
	request.SetPathValue("id", "42")

	input, err := bindInput[inputBody](request, buildInputMetadata[inputBody]("/things/{id}"))

	if err != nil {
		t.Fatalf("bindInput() error = %v", err)
	}
	if input.ID != 42 {
		t.Errorf("ID = %d, want %d", input.ID, 42)
	}
}

func TestBuildInputMetadataRejectsMissingPathBinding(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("buildInputMetadata() did not panic")
		}
	}()

	_ = buildInputMetadata[struct{}]("/things/{id}")
}
