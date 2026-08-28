package bootstrap

import (
	"log/slog"
	"os"

	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/infra/memory"
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func newLinkRepository() ports.LinkRepository {
	return memory.NewLinkRepository()
}
