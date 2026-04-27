package grit

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TEST 1: Menu overlays contain navigation hints
func TestAnalyticsMenuShowsNavigationHints(t *testing.T) {
	m := newModel(".")
	m.showAnalyticsMenu = true
	m.analyticsMenuIdx = 0

	output := renderAnalyticsMenuOverlay(m, 80)

	// Should contain hints for navigation
	if !strings.Contains(output, "j/k") && !strings.Contains(output, "Enter") {
		t.Errorf("Analytics menu should show navigation hints. Got: %s", output)
	}
}

func TestVisualizationMenuShowsNavigationHints(t *testing.T) {
	m := newModel(".")
	m.showVisualizationMenu = true
	m.vizMenuIdx = 0

	output := renderVisualizationMenuOverlay(m, 80)

	// Should contain hints for navigation
	if !strings.Contains(output, "j/k") && !strings.Contains(output, "Enter") {
		t.Errorf("Visualization menu should show navigation hints. Got: %s", output)
	}
}

func TestTeamMenuShowsNavigationHints(t *testing.T) {
	m := newModel(".")
	m.showTeamMenu = true
	m.teamMenuIdx = 0

	output := renderTeamMenuOverlay(m, 80)

	// Should contain hints for navigation
	if !strings.Contains(output, "j/k") && !strings.Contains(output, "Enter") {
		t.Errorf("Team menu should show navigation hints. Got: %s", output)
	}
}

// TEST 2: Integration menu exists and is properly defined
func TestIntegrationMenuItemsExist(t *testing.T) {
	if len(integrationMenuItems) == 0 {
		t.Error("integrationMenuItems should not be empty")
	}

	// Should have reasonable count of items
	if len(integrationMenuItems) < 2 {
		t.Errorf("integrationMenuItems should have at least 2 items, got %d", len(integrationMenuItems))
	}

	if integrationMenuLen != len(integrationMenuItems) {
		t.Errorf("integrationMenuLen should equal len(integrationMenuItems): %d vs %d",
			integrationMenuLen, len(integrationMenuItems))
	}
}

// TEST 3: Integration menu overlay renders correctly
func TestIntegrationMenuRenders(t *testing.T) {
	m := newModel(".")
	m.showIntegrationMenu = true
	m.integrationMenuIdx = 0

	output := renderIntegrationMenuOverlay(m, 80)

	if output == "" {
		t.Error("Integration menu overlay should not be empty")
	}

	// Should show navigation hints
	if !strings.Contains(output, "j/k") && !strings.Contains(output, "Enter") {
		t.Errorf("Integration menu should show navigation hints. Got: %s", output)
	}
}

// TEST 4: 'i' key opens integration menu
func TestIKeyOpensIntegrationMenu(t *testing.T) {
	m := newModel(".")
	m.showIntegrationMenu = false

	// Simulate pressing 'i'
	km := createKeyMsg("i")
	updatedModel, _ := m.Update(km)
	m = updatedModel.(model)

	if !m.showIntegrationMenu {
		t.Error("'i' key should open integration menu")
	}
}

// TEST 5: 'i' key closes integration menu when open
func TestIKeyTogglesIntegrationMenu(t *testing.T) {
	m := newModel(".")
	m.showIntegrationMenu = true

	// Simulate pressing 'i'
	km := createKeyMsg("i")
	updatedModel, _ := m.Update(km)
	m = updatedModel.(model)

	if m.showIntegrationMenu {
		t.Error("'i' key should toggle off integration menu")
	}
}

// TEST 6: Integration menu navigation works
func TestIntegrationMenuNavigation(t *testing.T) {
	m := newModel(".")
	m.showIntegrationMenu = true
	m.integrationMenuIdx = 0

	// Test down
	km := createKeyMsg("j")
	updatedModel, _ := m.Update(km)
	m = updatedModel.(model)
	if m.integrationMenuIdx != 1 {
		t.Errorf("'j' should move cursor down, got index %d", m.integrationMenuIdx)
	}

	// Test up
	km = createKeyMsg("k")
	updatedModel, _ = m.Update(km)
	m = updatedModel.(model)
	if m.integrationMenuIdx != 0 {
		t.Errorf("'k' should move cursor up, got index %d", m.integrationMenuIdx)
	}
}

// TEST 7: Integration menu dispatches features
func TestIntegrationMenuDispatchesFeatures(t *testing.T) {
	m := newModel(".")
	m.showIntegrationMenu = true
	m.integrationMenuIdx = 0

	// Simulate pressing Enter
	km := createKeyMsg("enter")
	updatedModel, _ := m.Update(km)
	m = updatedModel.(model)

	// Menu should close after selection
	if m.showIntegrationMenu {
		t.Error("Integration menu should close after selecting a feature")
	}
}

// TEST 8: Help hint appears in main UI footer
func TestMainUIShowsHelpHint(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{
		{hash: "abc123", author: "test", subject: "test commit", when: "1 hour ago"},
	}
	m.width = 80
	m.height = 20

	output := m.View()

	// Should contain hint about ? key
	if !strings.Contains(output, "?") {
		t.Errorf("Main UI should show help hint with '?' key")
	}
}

// Helper function to create key messages
func createKeyMsg(s string) interface{} {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}
