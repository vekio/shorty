package links

import (
	"context"

	"github.com/vekio/shorty/internal/app/createlink"
)

// CreateLink translates the HTTP input into the create-link command and maps
// its result back to the public representation.
func (h *handler) CreateLink(
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
