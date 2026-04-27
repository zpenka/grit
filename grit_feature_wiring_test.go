package grit

import (
	"strings"
	"testing"
)

// ============================================================================
// Feature Wiring Tests - Ensure all render functions are called from View()
// ============================================================================

// TEST 1: Export UI Rendering
func TestShowExportUIRendersExportUI(t *testing.T) {
	m := newModel(".")
	m.showExportUI = true
	m.pendingExports = []exportData{
		{format: "markdown", filename: "commits.md", content: "# Commits"},
	}
	m.width = 80
	m.height = 20

	output := m.View()

	if output == "" {
		t.Error("Export UI should render when showExportUI is true")
	}
}

// TEST 2: Conflict UI Rendering
func TestShowConflictUIRendersConflictUI(t *testing.T) {
	m := newModel(".")
	m.showConflictUI = true
	m.width = 80
	m.height = 20

	output := m.View()

	if output == "" {
		t.Error("Conflict UI should render when showConflictUI is true")
	}
}


// ============================================================================
// Help Documentation Tests - Verify keybindings are documented
// ============================================================================

func TestHelpOverlayShowsAllMainKeybindings(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 80
	m.height = 30

	output := renderHelpOverlay(m)

	// Should document main navigation keys
	requiredKeys := []string{"j/k", "f", "b", "B", "/", "?", "q"}
	for _, key := range requiredKeys {
		if !strings.Contains(output, key) {
			t.Errorf("Help should document '%s' key", key)
		}
	}
}

func TestHelpOverlayShowsMenuKeybindings(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 80
	m.height = 30

	output := renderHelpOverlay(m)

	// Should document submenu keys
	menuKeys := []string{"a", "v", "t", "i", "g"}
	for _, key := range menuKeys {
		if !strings.Contains(output, key) {
			t.Errorf("Help should document '%s' submenu key", key)
		}
	}
}

func TestHelpOverlayDocumentsGitOpsKey(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 80
	m.height = 30

	output := renderHelpOverlay(m)

	// Should specifically mention git operations
	if !strings.Contains(output, "Git") && !strings.Contains(output, "git") &&
		!strings.Contains(output, "Rebase") && !strings.Contains(output, "rebase") {
		t.Error("Help should document git operations (g key)")
	}
}

func TestHelpOverlayDocumentsIntegrationKey(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 80
	m.height = 30

	output := renderHelpOverlay(m)

	// Should mention integration/export
	if !strings.Contains(output, "Integration") && !strings.Contains(output, "integration") &&
		!strings.Contains(output, "Export") && !strings.Contains(output, "export") {
		t.Error("Help should document integration features (i key)")
	}
}

// TEST: Help mentions menu navigation hints
func TestHelpOverlayDocumentsMenuNavigation(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 80
	m.height = 30

	output := renderHelpOverlay(m)

	// Should mention that menus use j/k/Enter
	if !strings.Contains(output, "j/k") {
		t.Error("Help should document j/k menu navigation")
	}
}
