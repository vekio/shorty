package deletelink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
)

// DeleteLinkHandler removes links through the configured repository.
type DeleteLinkHandler struct {
	repository ports.LinkDeleter
}

// NewDeleteLinkHandler creates a delete-link handler.
func NewDeleteLinkHandler(repository ports.LinkDeleter) *DeleteLinkHandler {
	return &DeleteLinkHandler{repository: repository}
}

func (h *DeleteLinkHandler) Handle(ctx context.Context, command DeleteLinkCommand) (DeleteLinkResult, error) {
	if err := h.repository.Delete(ctx, command.Code); err != nil {
		return DeleteLinkResult{}, err
	}
	return DeleteLinkResult{}, nil
}
