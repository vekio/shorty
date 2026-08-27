package api

import (
	"context"

	"github.com/vekio/shorty/internal/app/createlink"
)

func (h *handlers) CreateLink(
	ctx context.Context,
	input CreateLinkInput,
) (CreateLinkOutput, error) {
	result, err := h.createLink.Handle(ctx, createlink.CreateLinkCommand{URL: input.OriginURL})
	if err != nil {
		return CreateLinkOutput{}, err
	}

	return CreateLinkOutput{
		Location: "/links/" + result.Code,
		Code:     result.Code,
	}, nil
}
