package links

import (
	"context"

	"github.com/vekio/shorty/internal/app/resolvelink"
)

// ResolveLink registers a visit and returns the destination as JSON for
// programmatic clients. The navigation endpoint exposes the same use case as a
// redirect under /r/{code}.
func (h *handler) ResolveLink(
	ctx context.Context,
	input ResolveLinkInput,
) (ResolveLinkOutput, error) {
	result, err := h.resolveLink.Handle(ctx, resolvelink.ResolveLinkCommand{Code: input.Code})
	if err != nil {
		return ResolveLinkOutput{}, err
	}
	return ResolveLinkOutput{OriginURL: result.OriginURL}, nil
}
