package visitlink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
)

type VisitLinkHandler struct {
	repository ports.LinkRepository
}

func NewVisitLinkHandler(repository ports.LinkRepository) *VisitLinkHandler {
	return &VisitLinkHandler{repository: repository}
}

func (h *VisitLinkHandler) Handle(ctx context.Context, command VisitLinkCommand) (VisitLinkResult, error) {
	link, err := h.repository.FindByCode(ctx, command.Code)
	if err != nil {
		return VisitLinkResult{}, err
	}

	link.RegisterVisit()
	if err := h.repository.UpdateLinkVisits(ctx, link); err != nil {
		return VisitLinkResult{}, err
	}
	return VisitLinkResult{}, nil
}
