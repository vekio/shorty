package web

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
)

// New builds Shorty's server-rendered interface and public redirect routes.
func New(api shortyAPI, logger *slog.Logger, shortURL string) http.Handler {
	mux := http.NewServeMux()
	staticFiles, err := fs.Sub(content, "static")
	if err != nil {
		panic("web: open embedded static files: " + err.Error())
	}
	mux.Handle("GET /assets/", http.StripPrefix("/assets", http.FileServerFS(staticFiles)))
	registerRoutes(mux, api, logger, shortURL)
	publicURL, err := url.Parse(shortURL)
	if err != nil {
		panic("web: parse public short URL: " + err.Error())
	}
	return browserSession(mux, publicURL.Scheme == "https")
}
