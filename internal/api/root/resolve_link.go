package root

import (
	"context"

	"github.com/vekio/shorty/internal/app/resolvelink"
)

// ResolveLink follows a shortened code and returns its redirect location.
func (h *handler) ResolveLink(
	ctx context.Context,
	input ResolveLinkInput,
) (ResolveLinkOutput, error) {
	result, err := h.resolveLink.Handle(ctx, resolvelink.ResolveLinkCommand{
		Code: input.Code,
	})
	if err != nil {
		return ResolveLinkOutput{}, err
	}

	return ResolveLinkOutput{Location: result.OriginURL}, nil
}
