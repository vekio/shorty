// Package web exposes Shorty's server-rendered management interface.
package web

import (
	"bytes"
	"embed"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/vekio/shorty/internal/auth"
	"github.com/vekio/shorty/internal/web/components"
	"github.com/vekio/shorty/internal/web/pages"
)

//go:embed static
var files embed.FS

type handler struct {
	apiKeys       *auth.Service
	workspaceID   string
	workspaceName string
}

// New builds the minimal administration handler for API key management.
func New(
	apiKeys *auth.Service,
	workspaceID string,
	workspaceName string,
) http.Handler {
	handler := &handler{
		apiKeys:       apiKeys,
		workspaceID:   workspaceID,
		workspaceName: workspaceName,
	}
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(files))
	mux.HandleFunc("GET /{$}", handler.dashboard)
	mux.HandleFunc("POST /api-keys", handler.createAPIKey)
	mux.HandleFunc("POST /api-keys/{id}/revoke", handler.revokeAPIKey)
	return mux
}

func (handler *handler) dashboard(w http.ResponseWriter, request *http.Request) {
	data, ok := handler.apiKeyPanelData(w, request, "", "")
	if !ok {
		return
	}
	handler.render(w, request, http.StatusOK, false, pages.Dashboard(pages.DashboardData{
		WorkspaceName: handler.workspaceName,
		APIKeys:       data,
	}))
}

func (handler *handler) createAPIKey(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		handler.renderAPIKeyPanel(w, request, http.StatusBadRequest, "", "invalid form")
		return
	}
	_, token, err := handler.apiKeys.Create(
		request.Context(), handler.workspaceID, request.FormValue("name"),
	)
	if err != nil {
		handler.renderAPIKeyPanel(w, request, http.StatusUnprocessableEntity, "", err.Error())
		return
	}
	handler.renderAPIKeyPanel(w, request, http.StatusCreated, token, "")
}

func (handler *handler) revokeAPIKey(w http.ResponseWriter, request *http.Request) {
	if err := handler.apiKeys.Revoke(
		request.Context(), handler.workspaceID, request.PathValue("id"),
	); err != nil {
		handler.renderAPIKeyPanel(w, request, http.StatusUnprocessableEntity, "", err.Error())
		return
	}
	handler.renderAPIKeyPanel(w, request, http.StatusOK, "", "")
}

func (handler *handler) renderAPIKeyPanel(
	w http.ResponseWriter,
	request *http.Request,
	status int,
	newToken string,
	publicError string,
) {
	data, ok := handler.apiKeyPanelData(w, request, newToken, publicError)
	if !ok {
		return
	}
	if status >= http.StatusBadRequest {
		// HTMX does not swap error responses by default. The validation message
		// is part of the successful fragment response instead.
		status = http.StatusOK
	}
	handler.render(w, request, status, newToken != "", components.APIKeyPanel(data))
}

func (handler *handler) apiKeyPanelData(
	w http.ResponseWriter,
	request *http.Request,
	newToken string,
	publicError string,
) (components.APIKeyPanelData, bool) {
	keys, err := handler.apiKeys.List(request.Context(), handler.workspaceID)
	if err != nil {
		http.Error(w, "load API keys", http.StatusInternalServerError)
		return components.APIKeyPanelData{}, false
	}
	return components.APIKeyPanelData{
		APIKeys:  keys,
		NewToken: newToken,
		Error:    strings.TrimSpace(publicError),
	}, true
}

func (handler *handler) render(
	w http.ResponseWriter,
	request *http.Request,
	status int,
	private bool,
	component templ.Component,
) {
	if err := renderComponent(w, request, status, component, private); err != nil {
		http.Error(w, "render page", http.StatusInternalServerError)
	}
}

func renderComponent(
	w http.ResponseWriter,
	request *http.Request,
	status int,
	component templ.Component,
	private bool,
) error {
	var content bytes.Buffer
	if err := component.Render(request.Context(), &content); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if private {
		w.Header().Set("Cache-Control", "no-store")
	}
	w.WriteHeader(status)
	_, err := content.WriteTo(w)
	return err
}
