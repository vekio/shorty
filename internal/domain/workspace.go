package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

// Workspace is an isolated container for links and API keys.
type Workspace struct {
	id   string
	name string
}

// NewWorkspace creates a workspace with a random identifier.
func NewWorkspace(name string) (Workspace, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Workspace{}, errors.New("workspace name is required")
	}
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return Workspace{}, err
	}
	return Workspace{id: "ws_" + hex.EncodeToString(random), name: name}, nil
}

func (workspace Workspace) ID() string {
	return workspace.id
}

func (workspace Workspace) Name() string {
	return workspace.name
}
