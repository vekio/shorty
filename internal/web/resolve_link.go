package web

import (
	"errors"
	"net/http"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

// ResolveLink asks the API to register the visit, then redirects the browser
// directly to the returned destination.
func (h *handler) ResolveLink(
	w http.ResponseWriter,
	request *http.Request,
) error {
	code := request.PathValue("code")
	destination, err := h.api.ResolveLink(request.Context(), code)
	if err != nil {
		var problem *shortysdk.ProblemError
		if errors.As(err, &problem) && problem.Status == http.StatusNotFound {
			return h.writeNotFound(w, notFoundView{
				Heading: "Short code not found",
				Detail:  "The code “" + code + "” does not exist.",
			})
		}
		return err
	}
	w.Header().Set("Location", destination)
	w.WriteHeader(http.StatusFound)
	return nil
}
