package links

import (
	"net/http"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

var (
	linkNotFoundError = amigo.WithErrorMapping(
		ports.ErrLinkNotFound,
		http.StatusNotFound,
		ports.ErrLinkNotFound.Error(),
	)
	originURLErrors = []amigo.RouteOption{
		amigo.WithErrorMapping(
			domain.ErrOriginURLRequired,
			http.StatusUnprocessableEntity,
			domain.ErrOriginURLRequired.Error(),
		),
		amigo.WithErrorMapping(
			domain.ErrOriginURLInvalid,
			http.StatusUnprocessableEntity,
			domain.ErrOriginURLInvalid.Error(),
		),
		amigo.WithErrorMapping(
			domain.ErrOriginURLSelfReference,
			http.StatusUnprocessableEntity,
			domain.ErrOriginURLSelfReference.Error(),
		),
	}
	paginationErrors = []amigo.RouteOption{
		amigo.WithErrorMapping(
			listlinks.ErrInvalidLimit,
			http.StatusBadRequest,
			listlinks.ErrInvalidLimit.Error(),
		),
		amigo.WithErrorMapping(
			listlinks.ErrInvalidOffset,
			http.StatusBadRequest,
			listlinks.ErrInvalidOffset.Error(),
		),
	}
)
