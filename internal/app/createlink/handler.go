package createlink

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/domain"
)

type CreateLinkHandler struct {
	repository ports.LinkRepository
}

func NewCreateLinkHandler(
	repository ports.LinkRepository,
) *CreateLinkHandler {
	return &CreateLinkHandler{
		repository: repository,
	}
}

func (h *CreateLinkHandler) Handle(ctx context.Context, command CreateLinkCommand) (CreateLinkResult, error) {
	link, err := domain.New(command.URL)
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
