package domain

import (
	"strings"
	"testing"
)

func TestNewWorkspaceUsesRandomIDAndProvidedName(t *testing.T) {
	first, err := NewWorkspace("default")
	if err != nil {
		t.Fatalf("NewWorkspace() error = %v", err)
	}
	second, err := NewWorkspace("default")
	if err != nil {
		t.Fatalf("NewWorkspace() second error = %v", err)
	}
	if !strings.HasPrefix(first.ID(), "ws_") || first.ID() == second.ID() || first.Name() != "default" {
		t.Errorf("workspaces = %#v, %#v", first, second)
	}
}
