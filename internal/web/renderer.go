package web

import (
	"bytes"
	"html/template"
)

type renderer struct {
	templates *template.Template
}

func newRenderer() renderer {
	return renderer{
		templates: template.Must(template.ParseFS(
			content,
			"layouts/*.html",
			"pages/*.html",
			"components/*.html",
		)),
	}
}

func (renderer renderer) render(
	templateName string,
	data any,
) (string, error) {
	// Render into memory so callers can still turn failures into a safe 500.
	var body bytes.Buffer
	if err := renderer.templates.ExecuteTemplate(&body, templateName, data); err != nil {
		return "", err
	}
	return body.String(), nil
}
