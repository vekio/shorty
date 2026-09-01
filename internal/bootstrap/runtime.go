package bootstrap

import (
	"log/slog"
	"net/http"
)

// Runtime is a fully composed HTTP process ready to be served by a command.
type Runtime struct {
	Handler http.Handler
	Logger  *slog.Logger
}
