package amigo

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func bindInput[In any](request *http.Request, metadata inputMetadata) (In, error) {
	var input In
	if err := bindJSONBody(request, &input); err != nil {
		return input, err
	}
	if err := bindPathParameters(reflect.ValueOf(&input).Elem(), request, metadata.path); err != nil {
		return input, err
	}
	return input, nil
}

func bindJSONBody(request *http.Request, destination any) error {
	if request.Body == nil || request.Body == http.NoBody || request.ContentLength == 0 {
		return nil
	}

	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return newProblem(http.StatusUnsupportedMediaType, "content type must be application/json")
	}

	if err := json.UnmarshalRead(request.Body, destination, json.RejectUnknownMembers(true)); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			return newProblem(http.StatusRequestEntityTooLarge, "request body exceeds the maximum allowed size")
		}
		return newProblem(http.StatusBadRequest, "invalid JSON request body")
	}
	return nil
}

func bindPathParameters(value reflect.Value, request *http.Request, parameters []pathParameter) error {
	for _, parameter := range parameters {
		if err := setPathParameter(value.FieldByIndex(parameter.index), request.PathValue(parameter.name)); err != nil {
			return newProblem(http.StatusBadRequest, fmt.Sprintf("invalid path parameter %q", parameter.name))
		}
	}
	return nil
}

func setPathParameter(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(value, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(value, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetUint(parsed)
	default:
		panic("amigo: unsupported path parameter type reached binding")
	}
	return nil
}

func supportsPathType(type_ reflect.Type) bool {
	switch type_.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func routePathNames(path string) map[string]struct{} {
	names := make(map[string]struct{})
	for segment := range strings.SplitSeq(path, "/") {
		if len(segment) < 3 || segment[0] != '{' || segment[len(segment)-1] != '}' {
			continue
		}
		name := strings.TrimSuffix(segment[1:len(segment)-1], "...")
		if name != "$" {
			names[name] = struct{}{}
		}
	}
	return names
}
