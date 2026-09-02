package bootstrap

import "net/http"

// NewServer mounts public redirects, the management API, and the admin panel.
func NewServer(
	redirectHandler http.Handler,
	apiHandler http.Handler,
	adminHandler http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/r/", redirectHandler)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiHandler))
	mux.Handle("/_/", http.StripPrefix("/_", adminHandler))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, request *http.Request) {
		http.Redirect(w, request, "/_/", http.StatusTemporaryRedirect)
	})

	return mux
}
