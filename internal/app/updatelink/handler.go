package updatelink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

type UpdateLinkHandler struct {
	repository        ports.LinkEditor
	originURLPolicies []domain.OriginURLPolicy
}

func NewUpdateLinkHandler(repository ports.LinkEditor, policies ...domain.OriginURLPolicy) *UpdateLinkHandler {
	return &UpdateLinkHandler{repository: repository, originURLPolicies: policies}
}

func (handler *UpdateLinkHandler) Handle(ctx context.Context, command UpdateLinkCommand) (UpdateLinkResult, error) {
	link, err := handler.repository.FindByCode(ctx, command.Code)
	if err != nil {
		return UpdateLinkResult{}, err
	}
	if err := link.ChangeOriginURL(command.OriginURL, handler.originURLPolicies...); err != nil {
		return UpdateLinkResult{}, err
	}
	if err := handler.repository.UpdateLinkOrigin(ctx, link); err != nil {
		return UpdateLinkResult{}, err
	}
	originURL := link.OriginURL()
	return UpdateLinkResult{
		Code:      link.Code(),
		OriginURL: originURL.String(),
		CreatedAt: link.CreatedAt(),
		Visits:    link.Visits(),
	}, nil
}
