package links

import (
	"context"

	"github.com/vekio/shorty/internal/app/listlinks"
)

// ListLinks returns the current link collection in application order.
func (h *handler) ListLinks(
	ctx context.Context,
	_ Empty,
) (ListLinksOutput, error) {
	result, err := h.listLinks.Handle(ctx, listlinks.ListLinksQuery{})
	if err != nil {
		return ListLinksOutput{}, err
	}

	output := ListLinksOutput{Links: make([]LinkOutput, 0, len(result.Links))}
	for _, link := range result.Links {
		output.Links = append(output.Links, LinkOutput{
			Code:      link.Code,
			OriginURL: link.OriginURL,
			CreatedAt: link.CreatedAt,
			Visits:    link.Visits,
		})
	}
	return output, nil
}
