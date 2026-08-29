package validator

import "github.com/vekio/shorty/internal/domain"

const urlValidator = "url"

func validateURL(rawURL string) error {
	return domain.ValidateOriginURL(rawURL)
}
