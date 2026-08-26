package getlink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
)

type GetLinkHandler struct {
	repository ports.LinkRepository
}

func NewGetLinkHandler(repository ports.LinkRepository) *GetLinkHandler {
	return &GetLinkHandler{repository: repository}
}

func (h *GetLinkHandler) Handle(ctx context.Context, query GetLinkQuery) (GetLinkResult, error) {
	link, err := h.repository.FindByCode(ctx, query.Code)
	if err != nil {
		return GetLinkResult{}, err
	}

	originURL := link.OriginURL()
	return GetLinkResult{
		Code:      link.Code(),
		OriginURL: originURL.String(),
		CreatedAt: link.CreatedAt(),
		Visits:    link.Visits(),
	}, nil
}
