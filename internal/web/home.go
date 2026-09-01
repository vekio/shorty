package web

import (
	"context"
	"net/http"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

type pageData struct {
	Title       string
	NotFound    *notFoundView
	OriginURL   string
	FormError   string
	CreatedURL  string
	Links       []linkView
	Total       int
	Limit       int
	Offset      int
	PageStart   int
	PageEnd     int
	Previous    int
	Next        int
	HasPrevious bool
	HasNext     bool
}

type linkView struct {
	shortysdk.Link
	ShortURL string
}

// Home renders the full document or its HTMX fragment.
func (h *handler) Home(w http.ResponseWriter, request *http.Request) error {
	limit, offset := paginationFromRequest(request)
	content, err := h.renderHome(request.Context(), isHTMX(request), pageData{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}
	writeHTML(w, http.StatusOK, content)
	return nil
}

func (h *handler) writeHomeFragment(
	w http.ResponseWriter,
	request *http.Request,
	data pageData,
) error {
	output, err := h.renderHome(request.Context(), true, data)
	if err != nil {
		return err
	}
	writeHTML(w, http.StatusOK, output)
	return nil
}

func (h *handler) renderHome(
	ctx context.Context,
	fragment bool,
	data pageData,
) (string, error) {
	if data.Limit == 0 {
		data.Limit = defaultLinksPageLimit
	}
	result, err := h.api.ListLinks(apiContext(ctx), shortysdk.ListOptions{
		Limit:  data.Limit,
		Offset: data.Offset,
	})
	if err != nil {
		return "", err
	}
	if data.Offset > 0 && len(result.Links) == 0 {
		data.Offset = lastPageOffset(result.Total, data.Limit)
		result, err = h.api.ListLinks(apiContext(ctx), shortysdk.ListOptions{
			Limit:  data.Limit,
			Offset: data.Offset,
		})
		if err != nil {
			return "", err
		}
	}

	listed := h.pageData(result)
	data.Links = listed.Links
	data.Total = listed.Total
	data.PageStart = 0
	data.PageEnd = data.Offset + len(data.Links)
	if len(data.Links) > 0 {
		data.PageStart = data.Offset + 1
	}
	data.HasPrevious = data.Offset > 0
	data.Previous = max(0, data.Offset-data.Limit)
	data.HasNext = data.PageEnd < data.Total
	data.Next = data.Offset + data.Limit
	data.Title = "Shorty"
	templateName := "layout"
	if fragment {
		templateName = "content"
	}
	content, err := h.renderer.render(templateName, data)
	if err != nil {
		return "", err
	}
	return content, nil
}

func lastPageOffset(total int, limit int) int {
	if total == 0 {
		return 0
	}
	return ((total - 1) / limit) * limit
}

func (h *handler) pageData(result shortysdk.LinkPage) pageData {
	data := pageData{
		Links: make([]linkView, 0, len(result.Links)),
		Total: result.Total,
	}
	for _, link := range result.Links {
		data.Links = append(data.Links, linkView{
			Link:     link,
			ShortURL: h.shortURL(link.Code),
		})
	}
	return data
}
