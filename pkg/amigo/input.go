package amigo

import (
	"fmt"
	"reflect"
	"strings"
)

type inputMetadata struct {
	bodyIndex int
	path      []pathParameter
}

type pathParameter struct {
	name       string
	fieldIndex int
	typeOf     reflect.Type
}

func buildInput(inputType reflect.Type, pattern string) inputMetadata {
	metadata := inputMetadata{
		bodyIndex: -1,
		path:      buildParameters(inputType, pattern),
	}
	for index := range inputType.NumField() {
		if inputType.Field(index).Name == "Body" {
			metadata.bodyIndex = index
		}
	}
	return metadata
}

func buildParameters(inputType reflect.Type, pattern string) []pathParameter {
	routeNames := pathNames(pattern)
	boundNames := make(map[string]struct{}, len(routeNames))
	parameters := make([]pathParameter, 0, len(routeNames))

	for index := range inputType.NumField() {
		field := inputType.Field(index)
		// Lookup distinguishes a missing tag from an explicitly empty path tag.
		name, tagged := field.Tag.Lookup("path")
		if !tagged {
			continue
		}

		parameter := buildPathParam(field, index, name, pattern, routeNames, boundNames)
		boundNames[name] = struct{}{}
		parameters = append(parameters, parameter)
	}

	// Every route wildcard must be represented in the typed input.
	for name := range routeNames {
		if _, exists := boundNames[name]; !exists {
			panic(fmt.Sprintf("amigo: route path parameter %q has no matching input field", name))
		}
	}
	return parameters
}

func buildPathParam(
	field reflect.StructField,
	fieldIndex int,
	name string,
	pattern string,
	routeNames map[string]struct{},
	boundNames map[string]struct{},
) pathParameter {
	if !field.IsExported() {
		panic(fmt.Sprintf("amigo: path field %s must be exported", field.Name))
	}
	if name == "" {
		panic(fmt.Sprintf("amigo: path field %s has an empty name", field.Name))
	}
	if _, exists := routeNames[name]; !exists {
		panic(fmt.Sprintf("amigo: path parameter %q is not present in route %q", name, pattern))
	}
	if _, exists := boundNames[name]; exists {
		panic(fmt.Sprintf("amigo: path parameter %q is bound more than once", name))
	}
	if !supportsPathType(field.Type) {
		panic(fmt.Sprintf("amigo: path parameter %q has unsupported type %s", name, field.Type))
	}

	return pathParameter{
		name:       name,
		fieldIndex: fieldIndex,
		typeOf:     field.Type,
	}
}

// pathNames extracts named wildcards from an http.ServeMux route pattern.
func pathNames(pattern string) map[string]struct{} {
	names := make(map[string]struct{})
	for segment := range strings.SplitSeq(pattern, "/") {
		if len(segment) < 3 || segment[0] != '{' || segment[len(segment)-1] != '}' {
			continue
		}

		name := segment[1 : len(segment)-1]
		// A catch-all wildcard such as {path...} is bound using the name "path".
		name = strings.TrimSuffix(name, "...")
		// {$} is ServeMux's end-of-path marker, not an input parameter.
		if name != "$" && name != "" {
			names[name] = struct{}{}
		}
	}
	return names
}
