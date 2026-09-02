// Package validator registers validation rules shared by Shorty's HTTP
// resources.
package validator

import (
	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/domain"
)

const (
	urlValidator        = "url"
	pageLimitValidator  = "page_limit"
	pageOffsetValidator = "page_offset"
)

// Register adds Shorty's shared request validators to api.
func Register(api *amigo.API) {
	api.Validator(urlValidator, domain.ValidateOriginURL)
	api.Validator(pageLimitValidator, listlinks.ValidateLimit)
	api.Validator(pageOffsetValidator, listlinks.ValidateOffset)
}
