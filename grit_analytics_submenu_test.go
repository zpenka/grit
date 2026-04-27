package grit

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate_AnalyticsKeyTogglesMenu(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24

	// Press "a"
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if !model1.showAnalyticsMenu {
		t.Error("pressing a should toggle showAnalyticsMenu to true")
	}
	if model1.analyticsMenuIdx != 0 {
		t.Error("pressing a should reset analyticsMenuIdx to 0")
	}

	// Press "a" again to close
	m2, _ := model1.Update(msg)
	model2 := m2.(model)

	if model2.showAnalyticsMenu {
		t.Error("pressing a again should toggle showAnalyticsMenu to false")
	}
}

func TestUpdate_AnalyticsMenuNavigation(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showAnalyticsMenu = true
	m.analyticsMenuIdx = 0
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// Press "j" to move down
	msgJ := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	m1, _ := m.Update(msgJ)
	model1 := m1.(model)

	if model1.analyticsMenuIdx != 1 {
		t.Errorf("pressing j should move down, got idx %d", model1.analyticsMenuIdx)
	}

	// Press "k" to move up
	msgK := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	m2, _ := model1.Update(msgK)
	model2 := m2.(model)

	if model2.analyticsMenuIdx != 0 {
		t.Errorf("pressing k should move up, got idx %d", model2.analyticsMenuIdx)
	}
}

func TestUpdate_AnalyticsMenuBoundsCheck(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showAnalyticsMenu = true
	m.analyticsMenuIdx = 6  // Last item (0-6 = 7 items)
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// Press "j" at the end - should not go beyond bounds
	msgJ := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	m1, _ := m.Update(msgJ)
	model1 := m1.(model)

	if model1.analyticsMenuIdx != 6 {
		t.Error("pressing j at end should stay at last item")
	}

	// Now move to start and test k
	m.analyticsMenuIdx = 0
	msgK := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	m2, _ := m.Update(msgK)
	model2 := m2.(model)

	if model2.analyticsMenuIdx != 0 {
		t.Error("pressing k at start should stay at first item")
	}
}

func TestUpdate_AnalyticsMenuSelect(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showAnalyticsMenu = true
	m.analyticsMenuIdx = 0  // Code Ownership
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// Press "enter" to select
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if model1.showAnalyticsMenu {
		t.Error("pressing enter should close analytics menu")
	}
	if !model1.showCodeOwnership {
		t.Error("selecting item 0 should enable showCodeOwnership")
	}
}

func TestUpdate_AnalyticsMenuEscapeClosesMenu(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showAnalyticsMenu = true
	m.analyticsMenuIdx = 2

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if model1.showAnalyticsMenu {
		t.Error("pressing esc should close analytics menu")
	}
}

func TestView_AnalyticsMenuOverlay(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 40
	m.showAnalyticsMenu = true
	m.analyticsMenuIdx = 0
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	view := m.View()

	// Menu should contain analytics feature names
	expectedStrings := []string{"Code Ownership", "Hotspot", "Linting", "Bisect", "Heatmap", "Author Stats", "Complexity"}
	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("analytics menu should contain '%s'", expected)
		}
	}
}

func TestView_AnalyticsFeaturePanel(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 40
	m.showCodeOwnership = true
	m.codeOwnership = map[string]codeOwnershipData{
		"user1": {author: "user1", files: make(map[string]int), lines: 100, expertise: 0.8},
	}

	view := m.View()

	// Should show code ownership panel
	if !strings.Contains(view, "user1") {
		t.Error("code ownership view should show author data")
	}
}

func TestDispatchAnalyticsFeature_ComputesFeaturesOnDemand(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{
		{hash: "abc1", author: "alice", subject: "feature: add login"},
		{hash: "abc2", author: "bob", subject: "fix: typo in docs"},
	}
	m.analyticsMenuIdx = 0  // Code Ownership

	// Dispatch code ownership
	m1 := dispatchAnalyticsFeature(m, 0)

	if !m1.showCodeOwnership {
		t.Error("dispatching code ownership should set showCodeOwnership")
	}
	if len(m1.codeOwnership) == 0 {
		t.Error("dispatching code ownership should compute ownership data")
	}
}

func TestDispatchAnalyticsFeature_TogglesBisect(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{
		{hash: "abc1", author: "alice", subject: "feature: add login"},
		{hash: "abc2", author: "bob", subject: "fix: typo in docs"},
	}
	m.analyticsMenuIdx = 3  // Bisect

	// Dispatch bisect
	m1 := dispatchAnalyticsFeature(m, 3)

	if !m1.bisectState.active {
		t.Error("dispatching bisect should activate bisect mode")
	}
	if !m1.showBisectUI {
		t.Error("dispatching bisect should show bisect UI")
	}
}
