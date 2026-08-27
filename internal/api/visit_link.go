package api

import (
	"context"

	"github.com/vekio/shorty/internal/app/visitlink"
)

func (h *handlers) VisitLink(
	ctx context.Context,
	input LinkByIDInput,
) (Empty, error) {
	_, err := h.visitLink.Handle(ctx, visitlink.VisitLinkCommand{
		Code: input.Code,
	})
	if err != nil {
		return Empty{}, err
	}
	return Empty{}, nil
}
