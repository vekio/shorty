package listlinks

import (
	"context"
	"sort"

	"github.com/vekio/shorty/internal/app/ports"
)

type ListLinksHandler struct {
	repository ports.LinkLister
}

func NewListLinksHandler(repository ports.LinkLister) *ListLinksHandler {
	return &ListLinksHandler{repository: repository}
}

func (h *ListLinksHandler) Handle(ctx context.Context, _ ListLinksQuery) (ListLinksResult, error) {
	links, err := h.repository.FindAll(ctx)
	if err != nil {
		return ListLinksResult{}, err
	}

	result := ListLinksResult{Links: make([]LinkResult, 0, len(links))}
	for _, link := range links {
		originURL := link.OriginURL()
		result.Links = append(result.Links, LinkResult{
			Code:      link.Code(),
			OriginURL: originURL.String(),
			CreatedAt: link.CreatedAt(),
			Visits:    link.Visits(),
		})
	}

	sort.Slice(result.Links, func(i, j int) bool {
		if result.Links[i].CreatedAt.Equal(result.Links[j].CreatedAt) {
			return result.Links[i].Code < result.Links[j].Code
		}
		return result.Links[i].CreatedAt.Before(result.Links[j].CreatedAt)
	})
	return result, nil
}
