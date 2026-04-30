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
	m.height = 60  // Increased height to ensure Integration section is visible

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

// TEST: Help overlay lists all menu names
func TestRenderHelpOverlay_ListsAllMenus(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 80
	m.height = 100  // Large height to ensure no clipping

	output := renderHelpOverlay(m)

	// Should contain all submenu headers
	menus := []string{"Analytics", "Visualizations", "Team", "Integration", "Git Ops", "Workflows"}
	for _, menu := range menus {
		if !strings.Contains(output, menu) {
			t.Errorf("Help overlay should contain menu name: %s", menu)
		}
	}
}

// TEST: Help overlay lists entries from each menu
func TestRenderHelpOverlay_ListsMenuEntries(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 80
	m.height = 100  // Large height to ensure no clipping

	output := renderHelpOverlay(m)

	// Sample entries from each menu that should appear
	entries := []string{
		"Code Ownership",        // Analytics
		"Hotspot Detection",     // Analytics
		"Contributor Flamegraph", // Visualizations
		"Timeline Slider",       // Visualizations
		"Team Statistics",       // Team
		"Reviewer Suggestions",  // Team
		"GitHub PR Links",       // Integration
		"Jira Tickets",         // Integration
		"Interactive Rebase",    // Git Ops
		"Cherry-pick Mode",     // Git Ops
		"Worktrees",            // Workflows
		"Named stashes",        // Workflows
	}

	for _, entry := range entries {
		if !strings.Contains(output, entry) {
			t.Errorf("Help overlay should contain menu entry: %s", entry)
		}
	}
}

// TEST: Help overlay stays in sync with analyticsMenuItems
func TestRenderHelpOverlay_StaysInSyncAnalytics(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 100
	m.height = 100  // Large height to ensure no clipping

	output := renderHelpOverlay(m)

	for _, item := range analyticsMenuItems {
		if !strings.Contains(output, item) {
			t.Errorf("Help overlay missing Analytics entry: %s", item)
		}
	}
}

// TEST: Help overlay stays in sync with vizMenuItems
func TestRenderHelpOverlay_StaysInSyncVisualization(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 100
	m.height = 100  // Large height to ensure no clipping

	output := renderHelpOverlay(m)

	for _, item := range vizMenuItems {
		if !strings.Contains(output, item) {
			t.Errorf("Help overlay missing Visualization entry: %s", item)
		}
	}
}

// TEST: Help overlay stays in sync with teamMenuItems
func TestRenderHelpOverlay_StaysInSyncTeam(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 100
	m.height = 100  // Large height to ensure no clipping

	output := renderHelpOverlay(m)

	for _, item := range teamMenuItems {
		if !strings.Contains(output, item) {
			t.Errorf("Help overlay missing Team entry: %s", item)
		}
	}
}

// TEST: Help overlay stays in sync with integrationMenuItems
func TestRenderHelpOverlay_StaysInSyncIntegration(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 100
	m.height = 100  // Large height to ensure no clipping

	output := renderHelpOverlay(m)

	for _, item := range integrationMenuItems {
		if !strings.Contains(output, item) {
			t.Errorf("Help overlay missing Integration entry: %s", item)
		}
	}
}

// TEST: Help overlay stays in sync with gitOpsMenuItems
func TestRenderHelpOverlay_StaysInSyncGitOps(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 100
	m.height = 100  // Large height to ensure no clipping

	output := renderHelpOverlay(m)

	for _, item := range gitOpsMenuItems {
		if !strings.Contains(output, item) {
			t.Errorf("Help overlay missing Git Ops entry: %s", item)
		}
	}
}

// TEST: Help overlay stays in sync with workflowsMenuItems
func TestRenderHelpOverlay_StaysInSyncWorkflows(t *testing.T) {
	m := newModel(".")
	m.showHelp = true
	m.width = 100
	m.height = 100  // Large height to ensure no clipping

	output := renderHelpOverlay(m)

	for _, item := range workflowsMenuItems {
		if !strings.Contains(output, item) {
			t.Errorf("Help overlay missing Workflows entry: %s", item)
		}
	}
}
