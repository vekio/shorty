package web

import (
	"net/http"
)

type notFoundView struct {
	Heading string
	Detail  string
}

// NotFound renders an HTML response for GET paths not owned by the Web
// application. More specific routes, such as /r/{code}, take precedence.
func (h *handler) NotFound(w http.ResponseWriter, request *http.Request) error {
	return h.writeNotFound(w, notFoundView{
		Heading: "Page not found",
		Detail:  "The page “" + request.URL.Path + "” does not exist.",
	})
}

func (h *handler) writeNotFound(w http.ResponseWriter, view notFoundView) error {
	content, err := h.renderer.render("layout", pageData{
		Title:    "Not found · Shorty",
		NotFound: &view,
	})
	if err != nil {
		return err
	}
	writeHTML(w, http.StatusNotFound, content)
	return nil
}
