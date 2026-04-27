package grit

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate_HelpKeyToggles(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24

	// Press "?"
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if !model1.showHelp {
		t.Error("pressing ? should toggle showHelp to true")
	}

	// Press "?" again
	m2, _ := model1.Update(msg)
	model2 := m2.(model)

	if model2.showHelp {
		t.Error("pressing ? again should toggle showHelp to false")
	}
}

func TestUpdate_HelpKeyInGlobalBindings(t *testing.T) {
	// "?" should work from any context, not just commit list
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test commit"}}
	m.showBranch = true // even with branch picker open
	m.searching = false

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if !model1.showHelp {
		t.Error("? should work even with branch picker open")
	}
}

func TestView_HelpOverlayShown(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 40  // Larger height to show full help text
	m.showHelp = true

	view := m.View()

	// Help overlay should be the main content
	hasQuestion := strings.Contains(view, "?")
	hasHelp := strings.Contains(view, "help") || strings.Contains(view, "HELP")
	if !hasQuestion {
		t.Error("help overlay should contain '?'")
	}
	if !hasHelp {
		t.Error("help overlay should contain 'help' or 'HELP'")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestView_HelpOverlayContainsKeybindings(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 40  // Larger height to show full help text
	m.showHelp = true

	view := m.View()

	expectedKeys := []string{"j", "k", "g", "G", "/", "f", "b", "B", "y", "e", "q", "?"}
	for _, key := range expectedKeys {
		if !strings.Contains(view, key) {
			t.Errorf("help overlay should contain key '%s'", key)
		}
	}
}

func TestUpdate_EscClosesHelpOverlay(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showHelp = true

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if model1.showHelp {
		t.Error("pressing Esc should close help overlay")
	}
}

func TestUpdate_EscClearsAllFeaturePanels(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showCodeOwnership = true
	m.showHotspots = true
	m.showAnalytics = true
	m.showFlamegraph = true
	m.showTimeline = true
	m.showTeamStats = true

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if model1.showCodeOwnership || model1.showHotspots || model1.showAnalytics ||
		model1.showFlamegraph || model1.showTimeline || model1.showTeamStats {
		t.Error("pressing Esc should clear all feature panels")
	}

	// Structural panels should NOT be cleared
	if model1.showBranch != m.showBranch || model1.showFiles != m.showFiles || model1.showBlame != m.showBlame {
		t.Error("pressing Esc should not clear structural panels (branch, files, blame)")
	}
}

func TestView_NormalViewWithoutHelpOverlay(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test commit"}}
	m.diffLines = []diffLine{{kind: lineContext, text: "some code"}}
	m.showHelp = false

	view := m.View()

	// Should render normal two-panel UI, not help
	if strings.Contains(view, "HELP") && strings.Contains(view, "Navigation") {
		t.Error("normal view should not show help overlay")
	}
}
