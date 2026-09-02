package bootstrap

import (
	"log/slog"
	"net/http"
)

// Runtime is a fully composed HTTP process ready to be served by a command.
type Runtime struct {
	Handler http.Handler
	Logger  *slog.Logger
	close   func() error
}

// Close releases resources owned by the process.
func (runtime Runtime) Close() error {
	if runtime.close == nil {
		return nil
	}
	return runtime.close()
}
