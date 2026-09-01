package resolve

import (
	"context"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

type resolveLinkInput struct {
	Code string `path:"code" json:"-"`
}

type resolveLinkOutput struct {
	Location string `header:"Location" json:"-"`
}

type handler struct {
	resolveLink app.ResolveLinkHandler
}

func (h *handler) ResolveLink(
	ctx context.Context,
	input resolveLinkInput,
) (resolveLinkOutput, error) {
	result, err := h.resolveLink.Handle(ctx, resolvelink.ResolveLinkCommand{Code: input.Code})
	if err != nil {
		return resolveLinkOutput{}, err
	}
	return resolveLinkOutput{Location: result.OriginURL}, nil
}
