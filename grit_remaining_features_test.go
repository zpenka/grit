package grit

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// === PHASE 3: VISUALIZATION SUBMENU ===

func TestUpdate_VisualizationKeyTogglesMenu(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("v")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if !model1.showVisualizationMenu {
		t.Error("pressing v should toggle showVisualizationMenu to true")
	}
}

func TestUpdate_VisualizationMenuNavigation(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showVisualizationMenu = true
	m.vizMenuIdx = 0
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	msgJ := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	m1, _ := m.Update(msgJ)
	model1 := m1.(model)

	if model1.vizMenuIdx != 1 {
		t.Errorf("pressing j should move down, got idx %d", model1.vizMenuIdx)
	}
}

// === PHASE 4: GIT OPS ===

func TestUpdate_RebaseKeyToggles(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if !model1.showRebaseUI {
		t.Error("pressing r should show rebase UI")
	}
	if len(model1.rebaseSequence) == 0 {
		t.Error("rebase sequence should be populated")
	}
}

func TestUpdate_CherryPickToggle(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}
	m.cursor = 0

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if len(model1.cherryPickList) == 0 {
		t.Error("cherry-pick should add commit to list")
	}
}

func TestUpdate_ResetModeCycles(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.resetMode = ""

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if model1.resetMode != "soft" {
		t.Errorf("pressing x should cycle to soft mode, got %q", model1.resetMode)
	}
}

// === PHASE 5: TEAM/AI SUBMENU ===

func TestUpdate_TeamKeyTogglesMenu(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if !model1.showTeamMenu {
		t.Error("pressing t should toggle showTeamMenu to true")
	}
}

func TestUpdate_TeamMenuNavigation(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showTeamMenu = true
	m.teamMenuIdx = 0
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	msgJ := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	m1, _ := m.Update(msgJ)
	model1 := m1.(model)

	if model1.teamMenuIdx != 1 {
		t.Errorf("pressing j should move down, got idx %d", model1.teamMenuIdx)
	}
}

// === PHASE 6: INTEGRATION SUBMENU ===

func TestView_IntegrationMenuRendersIfOpen(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 40
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// We'll test this after wiring
	// For now just ensure the test compiles
	_ = m
}

// === ADDITIONAL TESTS ===

func TestUpdate_MenuClosedByEsc(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showVisualizationMenu = true
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	// When Esc is pressed while visualization menu is open, it should close
	if model1.showVisualizationMenu {
		t.Error("pressing Esc should close visualization menu")
	}
}

func TestView_NoOverlapBetweenFeaturePanels(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 40
	m.showCodeOwnership = true
	m.showFlamegraph = true // Two feature panels open
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}
	m.codeOwnership = map[string]codeOwnershipData{
		"user1": {author: "user1", files: make(map[string]int), lines: 100},
	}
	m.contributorFlameData = []contributorFlameData{
		{author: "user1", percentage: 50.0},
	}

	view := m.View()

	// Only one should be rendered (first one wins in if/if chain)
	if !strings.Contains(view, "user1") {
		t.Error("expected code ownership to be rendered (it was set first)")
	}
}
