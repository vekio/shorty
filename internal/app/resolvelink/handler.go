package resolvelink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
)

// ResolveLinkHandler resolves a destination and persists the resulting visit.
type ResolveLinkHandler struct {
	repository ports.LinkVisitor
}

// NewResolveLinkHandler creates a resolve-link handler.
func NewResolveLinkHandler(repository ports.LinkVisitor) *ResolveLinkHandler {
	return &ResolveLinkHandler{repository: repository}
}

func (h *ResolveLinkHandler) Handle(ctx context.Context, command ResolveLinkCommand) (ResolveLinkResult, error) {
	link, err := h.repository.FindByCode(ctx, command.Code)
	if err != nil {
		return ResolveLinkResult{}, err
	}

	link.RegisterVisit()
	if err := h.repository.UpdateLinkVisits(ctx, link); err != nil {
		return ResolveLinkResult{}, err
	}
	originURL := link.OriginURL()
	return ResolveLinkResult{OriginURL: originURL.String()}, nil
}
