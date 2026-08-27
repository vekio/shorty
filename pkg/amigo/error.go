package amigo

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"
)

func writeError(w http.ResponseWriter, request *http.Request, route route, err error) {
	problem := route.problem(err)
	problem.Instance = request.URL.Path

	if problem.Status >= http.StatusInternalServerError {
		slog.ErrorContext(request.Context(), "request failed",
			"method", request.Method,
			"path", request.URL.Path,
			"error", err,
		)
	}

	writeProblem(w, problem)
}

func writeProblem(w http.ResponseWriter, problem *problem) {
	data, err := json.Marshal(problem)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(problem.Status)
	_, _ = w.Write(data)
}
