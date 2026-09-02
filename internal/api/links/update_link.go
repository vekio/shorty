package links

import (
	"context"

	"github.com/vekio/shorty/internal/app/updatelink"
)

// UpdateLink changes a link destination and returns the updated resource.
func (h *handler) UpdateLink(ctx context.Context, input UpdateLinkInput) (LinkOutput, error) {
	result, err := h.updateLink.Handle(ctx, updatelink.UpdateLinkCommand{
		Code:      input.Code,
		OriginURL: input.OriginURL,
	})
	if err != nil {
		return LinkOutput{}, err
	}
	return LinkOutput{
		Code:      result.Code,
		OriginURL: result.OriginURL,
		CreatedAt: result.CreatedAt,
		Visits:    result.Visits,
	}, nil
}
