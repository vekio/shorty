package links

import (
	"context"

	"github.com/vekio/shorty/internal/app/listlinks"
)

// ListLinks returns the current link collection in application order.
func (h *handler) ListLinks(
	ctx context.Context,
	input ListLinksInput,
) (ListLinksOutput, error) {
	result, err := h.listLinks.Handle(ctx, listlinks.ListLinksQuery{
		OwnerID: input.OwnerID,
		Limit:   input.Limit,
		Offset:  input.Offset,
	})
	if err != nil {
		return ListLinksOutput{}, err
	}

	output := ListLinksOutput{
		Links:  make([]LinkOutput, 0, len(result.Links)),
		Total:  result.Total,
		Limit:  result.Limit,
		Offset: result.Offset,
	}
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
