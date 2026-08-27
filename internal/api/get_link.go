package api

import (
	"context"

	"github.com/vekio/shorty/internal/app/getlink"
)

func (h *handlers) GetLink(
	ctx context.Context,
	input LinkByIDInput,
) (LinkOutput, error) {
	result, err := h.getLink.Handle(ctx, getlink.GetLinkQuery{
		Code: input.Code,
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
