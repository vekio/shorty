package amigo

import (
	"encoding/json/v2"
	"net/http/httptest"
	"testing"
)

func assertProblem(t *testing.T, response *httptest.ResponseRecorder, status int, detail string) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, status, response.Body.String())
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Errorf("Content-Type = %q, want application/problem+json", contentType)
	}

	var problem Problem
	if err := json.Unmarshal(response.Body.Bytes(), &problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Status != status || problem.Detail != detail {
		t.Errorf("problem = %#v, want status %d and detail %q", problem, status, detail)
	}
}

func assertPanics(t *testing.T, action func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	action()
}
