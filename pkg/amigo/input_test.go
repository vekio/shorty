package amigo

import (
	"reflect"
	"testing"
)

func TestBuildInputFindsBodyAndPathParameters(t *testing.T) {
	type input struct {
		ID   string `path:"id"`
		Body struct{}
	}
	metadata := buildInput(reflect.TypeFor[input](), "/things/{id}")
	if metadata.bodyIndex != 1 || len(metadata.path) != 1 || metadata.path[0].name != "id" {
		t.Errorf("metadata = %#v", metadata)
	}
}

func TestBuildParametersRejectsInvalidDeclarations(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		typeOf  reflect.Type
	}{
		{name: "empty tag", pattern: "/things/{id}", typeOf: reflect.TypeFor[struct {
			ID string `path:""`
		}]()},
		{name: "missing wildcard", pattern: "/things", typeOf: reflect.TypeFor[struct {
			ID string `path:"id"`
		}]()},
		{name: "unbound wildcard", pattern: "/things/{id}", typeOf: reflect.TypeFor[struct{ ID string }]()},
		{name: "duplicate binding", pattern: "/things/{id}", typeOf: reflect.TypeFor[struct {
			ID   string `path:"id"`
			Code string `path:"id"`
		}]()},
		{name: "unsupported type", pattern: "/things/{id}", typeOf: reflect.TypeFor[struct {
			ID []string `path:"id"`
		}]()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertPanics(t, func() { buildInput(test.typeOf, test.pattern) })
		})
	}
}

func TestPathNames(t *testing.T) {
	names := pathNames("/users/{user}/files/{path...}/{$}")
	if len(names) != 2 {
		t.Fatalf("names = %#v", names)
	}
	if _, exists := names["user"]; !exists {
		t.Error("user wildcard is missing")
	}
	if _, exists := names["path"]; !exists {
		t.Error("catch-all wildcard is missing")
	}
}
