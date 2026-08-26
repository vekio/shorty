package api

import (
	"context"
	"time"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/visitlink"
)

type handlers struct {
	createLink app.CreateLinkHandler
	getLink    app.GetLinkHandler
	visitLink  app.VisitLinkHandler
}

func newHandlers(application app.Application) *handlers {
	return &handlers{
		createLink: application.Commands.CreateLink,
		getLink:    application.Queries.GetLink,
		visitLink:  application.Commands.VisitLink,
	}
}

type CreateLinkInput struct {
	Body struct {
		OriginURL string `json:"origin_url"`
	}
}

type CreateLinkOutput struct {
	Code string `json:"code"`
}

type GetLinkInput struct {
	ID string `path:"id"`
}

type GetLinkOutput struct {
	Code      string    `json:"code"`
	OriginURL string    `json:"origin_url"`
	CreatedAt time.Time `json:"created_at"`
	Visits    int       `json:"visits"`
}

type VisitLinkInput struct {
	ID string `path:"id"`
}

type VisitLinkOutput struct{}

// CreateLink creates a new short link from the request body.
func (h *handlers) CreateLink(ctx context.Context, input CreateLinkInput) (CreateLinkOutput, error) {
	result, err := h.createLink.Handle(ctx, createlink.CreateLinkCommand{
		URL: input.Body.OriginURL,
	})
	if err != nil {
		return CreateLinkOutput{}, err
	}

	return CreateLinkOutput{Code: result.Code}, nil
}

func (h *handlers) GetLink(ctx context.Context, input GetLinkInput) (GetLinkOutput, error) {
	result, err := h.getLink.Handle(ctx, getlink.GetLinkQuery{Code: input.ID})
	if err != nil {
		return GetLinkOutput{}, err
	}

	return GetLinkOutput{
		Code:      result.Code,
		OriginURL: result.OriginURL,
		CreatedAt: result.CreatedAt,
		Visits:    result.Visits,
	}, nil
}

func (h *handlers) VisitLink(ctx context.Context, input VisitLinkInput) (VisitLinkOutput, error) {
	_, err := h.visitLink.Handle(ctx, visitlink.VisitLinkCommand{Code: input.ID})
	if err != nil {
		return VisitLinkOutput{}, err
	}
	return VisitLinkOutput{}, nil
}
