package web

import (
	"net/http"
)

// DeleteLink removes a link owned by this browser session and returns the
// refreshed workspace for HTMX to swap into the page.
func (h *handler) DeleteLink(w http.ResponseWriter, request *http.Request) error {
	limit, offset := paginationFromRequest(request)
	ctx := apiContext(request.Context())
	if err := h.api.DeleteLink(ctx, request.PathValue("code")); err != nil {
		return err
	}

	return h.writeHomeFragment(w, request, pageData{Limit: limit, Offset: offset})
}
