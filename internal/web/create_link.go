package web

import (
	"errors"
	"net/http"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

// CreateLink handles the HTMX form and returns the updated home fragment.
func (h *handler) CreateLink(w http.ResponseWriter, request *http.Request) error {
	limit, offset := paginationFromRequest(request)
	if err := request.ParseForm(); err != nil {
		return h.writeHomeFragment(w, request, pageData{
			FormError: "invalid form submission",
			Limit:     limit,
			Offset:    offset,
		})
	}

	originURL := request.PostForm.Get("origin_url")
	code, err := h.api.CreateLink(apiContext(request.Context()), originURL)
	if err != nil {
		if detail, known := createLinkErrorDetail(err); known {
			// This endpoint returns presentational HTML. A validation failure is
			// therefore a successful fragment update from HTMX's perspective.
			return h.writeHomeFragment(w, request, pageData{
				OriginURL: originURL,
				FormError: detail,
				Limit:     limit,
				Offset:    offset,
			})
		}
		return err
	}

	w.Header().Set("HX-Push-Url", "/")
	return h.writeHomeFragment(w, request, pageData{
		CreatedURL: h.shortURL(code),
		Limit:      limit,
	})
}

func createLinkErrorDetail(err error) (string, bool) {
	var problem *shortysdk.ProblemError
	if !errors.As(err, &problem) || problem.Status != http.StatusUnprocessableEntity {
		return "", false
	}
	for _, field := range problem.Errors {
		if field.Location == "body.origin_url" {
			return field.Message, true
		}
	}
	return problem.Detail, problem.Detail != ""
}

func isHTMX(request *http.Request) bool {
	return request.Header.Get("HX-Request") == "true"
}
