package web

import "net/http"

func writeHTML(w http.ResponseWriter, status int, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(content))
}
