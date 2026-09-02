package resolvelink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
)

// ResolveLinkHandler resolves a destination and persists the resulting visit.
type ResolveLinkHandler struct {
	repository ports.LinkResolver
}

// NewResolveLinkHandler creates a resolve-link handler.
func NewResolveLinkHandler(repository ports.LinkResolver) *ResolveLinkHandler {
	return &ResolveLinkHandler{repository: repository}
}

func (h *ResolveLinkHandler) Handle(ctx context.Context, command ResolveLinkCommand) (ResolveLinkResult, error) {
	link, err := h.repository.ResolveByCode(ctx, command.Code)
	if err != nil {
		return ResolveLinkResult{}, err
	}

	originURL := link.OriginURL()
	return ResolveLinkResult{OriginURL: originURL.String()}, nil
}
