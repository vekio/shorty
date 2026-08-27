package amigo

import (
	"fmt"
	"reflect"
	"strings"
)

type inputMetadata struct {
	path []pathParameter
}

type pathParameter struct {
	name  string
	index []int
}

type outputMetadata struct {
	headers []outputHeader
}

type outputHeader struct {
	name  string
	index []int
}

func buildInputMetadata[In any](path string) inputMetadata {
	inputType := reflect.TypeFor[In]()
	if inputType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("amigo: endpoint input must be a struct, got %s", inputType))
	}

	pathNames := routePathNames(path)
	boundNames := make(map[string]struct{}, len(pathNames))
	metadata := inputMetadata{}

	for field := range inputType.Fields() {
		name, tagged := field.Tag.Lookup("path")
		if !tagged {
			continue
		}
		if !field.IsExported() {
			panic(fmt.Sprintf("amigo: path field %s must be exported", field.Name))
		}
		if name == "" {
			panic(fmt.Sprintf("amigo: path field %s has an empty name", field.Name))
		}
		if _, exists := pathNames[name]; !exists {
			panic(fmt.Sprintf("amigo: path parameter %q is not present in route %q", name, path))
		}
		if _, exists := boundNames[name]; exists {
			panic(fmt.Sprintf("amigo: path parameter %q is bound more than once", name))
		}
		if !supportsPathType(field.Type) {
			panic(fmt.Sprintf("amigo: path parameter %q has unsupported type %s", name, field.Type))
		}
		if jsonName(field.Tag.Get("json")) != "-" {
			panic(fmt.Sprintf("amigo: path field %s must use json:\"-\"", field.Name))
		}

		boundNames[name] = struct{}{}
		metadata.path = append(metadata.path, pathParameter{
			name:  name,
			index: field.Index,
		})
	}

	for name := range pathNames {
		if _, exists := boundNames[name]; !exists {
			panic(fmt.Sprintf("amigo: path parameter %q is not bound by the endpoint input", name))
		}
	}

	return metadata
}

func buildOutputMetadata[Out any]() outputMetadata {
	outputType := reflect.TypeFor[Out]()
	if outputType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("amigo: endpoint output must be a struct, got %s", outputType))
	}

	metadata := outputMetadata{}
	for field := range outputType.Fields() {
		name, tagged := field.Tag.Lookup("header")
		if !tagged {
			continue
		}
		if !field.IsExported() {
			panic(fmt.Sprintf("amigo: header field %s must be exported", field.Name))
		}
		if name == "" {
			panic(fmt.Sprintf("amigo: header field %s has an empty name", field.Name))
		}
		if field.Type.Kind() != reflect.String {
			panic(fmt.Sprintf("amigo: response header %q must be a string", name))
		}
		if jsonName(field.Tag.Get("json")) != "-" {
			panic(fmt.Sprintf("amigo: header field %s must use json:\"-\"", field.Name))
		}
		metadata.headers = append(metadata.headers, outputHeader{name: name, index: field.Index})
	}
	return metadata
}

func jsonName(tag string) string {
	name, _, _ := strings.Cut(tag, ",")
	return name
}
