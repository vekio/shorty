// Package validator registers validation rules shared by Shorty's HTTP
// resources.
package validator

import (
	"strings"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/domain"
)

const urlValidator = "url"

// Register adds Shorty's shared request validators to api.
func Register(api *amigo.API) {
	api.Validator(urlValidator, validateURL)
}

func validateURL(rawURL string) error {
	// Required owns the empty-value error, so an empty value does not produce
	// both a presence issue and a URL-format issue.
	if strings.TrimSpace(rawURL) == "" {
		return nil
	}
	return domain.ValidateOriginURL(rawURL)
}
