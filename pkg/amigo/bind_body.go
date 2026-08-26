package amigo

import (
	"encoding/json/v2"
	"errors"
	"mime"
	"net/http"
	"reflect"
)

func bindBody(
	w http.ResponseWriter,
	request *http.Request,
	input reflect.Value,
	bodyIndex int,
	maxBodyBytes int64,
) error {
	if bodyIndex < 0 {
		return nil
	}

	if request.Body == nil || request.Body == http.NoBody || request.ContentLength == 0 {
		return BadRequest("request body is required")
	}

	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return UnsupportedMediaType("content type must be application/json")
	}

	reader := request.Body
	if maxBodyBytes > 0 {
		reader = http.MaxBytesReader(w, request.Body, maxBodyBytes)
	}
	body := input.Field(bodyIndex)
	if err := json.UnmarshalRead(reader, body.Addr().Interface(), json.RejectUnknownMembers(true)); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			return ContentTooLarge("request body exceeds the maximum allowed size")
		}
		return BadRequest("invalid JSON request body")
	}

	return nil
}
