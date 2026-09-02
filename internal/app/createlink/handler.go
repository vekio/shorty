package createlink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

type CreateLinkHandler struct {
	repository        ports.LinkSaver
	originURLPolicies []domain.OriginURLPolicy
}

func NewCreateLinkHandler(
	repository ports.LinkSaver,
	originURLPolicies ...domain.OriginURLPolicy,
) *CreateLinkHandler {
	return &CreateLinkHandler{
		repository:        repository,
		originURLPolicies: originURLPolicies,
	}
}

func (h *CreateLinkHandler) Handle(ctx context.Context, command CreateLinkCommand) (CreateLinkResult, error) {
	link, err := domain.NewLink(command.OriginURL, h.originURLPolicies...)
	if err != nil {
		return CreateLinkResult{}, err
	}
	if err := h.repository.Save(ctx, link); err != nil {
		return CreateLinkResult{}, err
	}

	return CreateLinkResult{
		Code: link.Code(),
	}, nil
}
