package links

import (
	"context"

	"github.com/vekio/shorty/internal/app/deletelink"
)

// DeleteLink removes the link identified by the route path.
func (h *handler) DeleteLink(ctx context.Context, input LinkByCodeInput) (Empty, error) {
	_, err := h.deleteLink.Handle(ctx, deletelink.DeleteLinkCommand{
		Code: input.Code,
	})
	if err != nil {
		return Empty{}, err
	}
	return Empty{}, nil
}
