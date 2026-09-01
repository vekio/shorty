package getlink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
)

type GetLinkHandler struct {
	repository ports.OwnedLinkFinder
}

func NewGetLinkHandler(repository ports.OwnedLinkFinder) *GetLinkHandler {
	return &GetLinkHandler{repository: repository}
}

func (h *GetLinkHandler) Handle(ctx context.Context, query GetLinkQuery) (GetLinkResult, error) {
	link, err := h.repository.FindOwnedByCode(ctx, query.OwnerID, query.Code)
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
