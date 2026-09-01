package links

import (
	"context"

	"github.com/vekio/shorty/internal/app/getlink"
)

// GetLink returns a link without changing its visit count.
func (h *handler) GetLink(
	ctx context.Context,
	input LinkByCodeInput,
) (LinkOutput, error) {
	result, err := h.getLink.Handle(ctx, getlink.GetLinkQuery{
		OwnerID: input.OwnerID,
		Code:    input.Code,
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
