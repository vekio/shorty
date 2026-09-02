// Package web exposes Shorty's server-rendered management interface.
package web

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/deletelink"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/auth"
)

//go:embed static templates
var files embed.FS

var dashboardTemplate = template.Must(template.ParseFS(files, "templates/dashboard.html"))

type handler struct {
	application   app.Application
	apiKeys       *auth.Service
	workspaceID   string
	workspaceName string
}

type pageData struct {
	WorkspaceID   string
	WorkspaceName string
	Links         []listlinks.LinkResult
	APIKeys       []auth.APIKey
	NewToken      string
	Error         string
}

// New builds the management Web handler using application use cases directly.
func New(
	application app.Application,
	apiKeys *auth.Service,
	workspaceID string,
	workspaceName string,
) http.Handler {
	handler := &handler{
		application:   application,
		apiKeys:       apiKeys,
		workspaceID:   workspaceID,
		workspaceName: workspaceName,
	}
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(files))
	mux.HandleFunc("GET /{$}", handler.dashboard)
	mux.HandleFunc("POST /links", handler.createLink)
	mux.HandleFunc("POST /links/{code}/delete", handler.deleteLink)
	mux.HandleFunc("POST /api-keys", handler.createAPIKey)
	mux.HandleFunc("POST /api-keys/{id}/revoke", handler.revokeAPIKey)
	return mux
}

func (handler *handler) dashboard(w http.ResponseWriter, request *http.Request) {
	handler.render(w, request, http.StatusOK, "", "")
}

func (handler *handler) createLink(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		handler.render(w, request, http.StatusBadRequest, "", "invalid form")
		return
	}
	_, err := handler.application.Commands.CreateLink.Handle(
		request.Context(),
		createlink.CreateLinkCommand{OriginURL: request.FormValue("origin_url")},
	)
	if err != nil {
		handler.render(w, request, http.StatusUnprocessableEntity, "", err.Error())
		return
	}
	http.Redirect(w, request, "/_/", http.StatusSeeOther)
}

func (handler *handler) deleteLink(w http.ResponseWriter, request *http.Request) {
	_, err := handler.application.Commands.DeleteLink.Handle(
		request.Context(),
		deletelink.DeleteLinkCommand{Code: request.PathValue("code")},
	)
	if err != nil {
		handler.render(w, request, http.StatusUnprocessableEntity, "", err.Error())
		return
	}
	http.Redirect(w, request, "/_/", http.StatusSeeOther)
}

func (handler *handler) createAPIKey(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		handler.render(w, request, http.StatusBadRequest, "", "invalid form")
		return
	}
	_, token, err := handler.apiKeys.Create(
		request.Context(), handler.workspaceID, request.FormValue("name"),
	)
	if err != nil {
		handler.render(w, request, http.StatusUnprocessableEntity, "", err.Error())
		return
	}
	handler.render(w, request, http.StatusCreated, token, "")
}

func (handler *handler) revokeAPIKey(w http.ResponseWriter, request *http.Request) {
	if err := handler.apiKeys.Revoke(
		request.Context(), handler.workspaceID, request.PathValue("id"),
	); err != nil {
		handler.render(w, request, http.StatusUnprocessableEntity, "", err.Error())
		return
	}
	http.Redirect(w, request, "/_/", http.StatusSeeOther)
}

func (handler *handler) render(
	w http.ResponseWriter,
	request *http.Request,
	status int,
	newToken string,
	publicError string,
) {
	links, err := handler.application.Queries.ListLinks.Handle(
		request.Context(), listlinks.ListLinksQuery{Limit: listlinks.MaximumLimit},
	)
	if err != nil {
		http.Error(w, "load links", http.StatusInternalServerError)
		return
	}
	keys, err := handler.apiKeys.List(request.Context(), handler.workspaceID)
	if err != nil {
		http.Error(w, "load API keys", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if newToken != "" {
		w.Header().Set("Cache-Control", "no-store")
	}
	w.WriteHeader(status)
	_ = dashboardTemplate.Execute(w, pageData{
		WorkspaceID:   handler.workspaceID,
		WorkspaceName: handler.workspaceName,
		Links:         links.Links,
		APIKeys:       keys,
		NewToken:      newToken,
		Error:         strings.TrimSpace(publicError),
	})
}
