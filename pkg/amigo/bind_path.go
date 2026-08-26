package amigo

import (
	"encoding"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

var textUnmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()

func supportsPathType(typeOf reflect.Type) bool {
	if reflect.PointerTo(typeOf).Implements(textUnmarshalerType) {
		return true
	}

	switch typeOf.Kind() {
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func bindPathParameters(input reflect.Value, request *http.Request, parameters []pathParameter) error {
	for _, parameter := range parameters {
		field := input.Field(parameter.fieldIndex)
		if err := setPathValue(field, request.PathValue(parameter.name)); err != nil {
			return BadRequest(fmt.Sprintf(
				"path parameter %q must be a valid %s",
				parameter.name,
				parameter.typeOf,
			))
		}
	}
	return nil
}

func setPathValue(field reflect.Value, raw string) error {
	if field.Addr().Type().Implements(textUnmarshalerType) {
		return field.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(raw))
	}

	var err error
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Bool:
		var value bool
		value, err = strconv.ParseBool(raw)
		field.SetBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var value int64
		value, err = strconv.ParseInt(raw, 10, field.Type().Bits())
		field.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var value uint64
		value, err = strconv.ParseUint(raw, 10, field.Type().Bits())
		field.SetUint(value)
	case reflect.Float32, reflect.Float64:
		var value float64
		value, err = strconv.ParseFloat(raw, field.Type().Bits())
		field.SetFloat(value)
	}
	return err
}
