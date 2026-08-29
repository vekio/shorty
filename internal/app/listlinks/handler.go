package listlinks

import (
	"context"

	"github.com/vekio/shorty/internal/app/ports"
)

type ListLinksHandler struct {
	repository ports.LinkLister
}

func NewListLinksHandler(repository ports.LinkLister) *ListLinksHandler {
	return &ListLinksHandler{repository: repository}
}

func (h *ListLinksHandler) Handle(ctx context.Context, query ListLinksQuery) (ListLinksResult, error) {
	if query.Limit == 0 {
		query.Limit = DefaultLimit
	}
	if err := ValidateLimit(query.Limit); err != nil {
		return ListLinksResult{}, err
	}
	if err := ValidateOffset(query.Offset); err != nil {
		return ListLinksResult{}, err
	}

	page, err := h.repository.FindPage(ctx, query.Limit, query.Offset)
	if err != nil {
		return ListLinksResult{}, err
	}

	result := ListLinksResult{
		Links:  make([]LinkResult, 0, len(page.Links)),
		Total:  page.Total,
		Limit:  query.Limit,
		Offset: query.Offset,
	}
	for _, link := range page.Links {
		originURL := link.OriginURL()
		result.Links = append(result.Links, LinkResult{
			Code:      link.Code(),
			OriginURL: originURL.String(),
			CreatedAt: link.CreatedAt(),
			Visits:    link.Visits(),
		})
	}
	return result, nil
}
