package api

import (
	"errors"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
	"github.com/vekio/shorty/pkg/amigo"
)

func mapApplicationError(err error) (*amigo.Problem, bool) {
	if errors.Is(err, domain.ErrOriginURLRequired) || errors.Is(err, domain.ErrOriginURLInvalid) {
		return amigo.BadRequest(err.Error()), true
	}
	if errors.Is(err, ports.ErrLinkNotFound) {
		return amigo.NotFound(err.Error()), true
	}

	return nil, false
}
