package amigo

import (
	"net/http"
	"reflect"
)

func bindInput[In any](
	w http.ResponseWriter,
	request *http.Request,
	metadata inputMetadata,
	maxBodyBytes int64,
) (In, error) {
	var input In
	value := reflect.ValueOf(&input).Elem()

	if err := bindPathParameters(value, request, metadata.path); err != nil {
		return input, err
	}
	if err := bindBody(w, request, value, metadata.bodyIndex, maxBodyBytes); err != nil {
		return input, err
	}

	return input, nil
}
