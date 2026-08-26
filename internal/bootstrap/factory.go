package bootstrap

import (
	"github.com/vekio/shorty/internal/app/ports"
	"github.com/vekio/shorty/internal/infra/memory"
)

func newLinkRepository() ports.LinkRepository {
	return memory.NewLinkRepository()
}
