package links

import (
	"context"

	"github.com/vekio/shorty/internal/app/visitlink"
)

// VisitLink registers one visit and intentionally returns no representation.
func (h *handler) VisitLink(
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
