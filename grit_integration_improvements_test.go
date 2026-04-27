package grit

import (
	"strings"
	"testing"
)

// ============================================================================
// OPTION 1: Wire Integration Feature Displays
// ============================================================================

func TestShowPRLinksRendersUI(t *testing.T) {
	m := newModel(".")
	m.showPRLinks = true
	m.prReferences = []githubPRReference{
		{hash: "abc123", prNumber: 123, status: "merged"},
	}
	m.width = 80
	m.height = 20

	output := m.View()

	// renderPRLinksUI should be called and return something
	if output == "" {
		t.Errorf("PR Links UI should produce output when showPRLinks is true")
	}
}

func TestShowJiraLinksRendersUI(t *testing.T) {
	m := newModel(".")
	m.showJiraLinks = true
	m.jiraLinks = []jiraTicketLink{
		{hash: "abc123", ticket: "PROJ-123", status: "done"},
	}
	m.width = 80
	m.height = 20

	output := m.View()

	if output == "" {
		t.Errorf("Jira Links UI should produce output when showJiraLinks is true")
	}
}

func TestShowIssueRefsRendersUI(t *testing.T) {
	m := newModel(".")
	m.showIssueRefs = true
	m.issueReferences = []issueReference{
		{hash: "abc123", references: []string{"#123", "#456"}},
	}
	m.width = 80
	m.height = 20

	output := m.View()

	if output == "" {
		t.Errorf("Issue References UI should produce output when showIssueRefs is true")
	}
}

// ============================================================================
// OPTION 2: Git Operations Submenu
// ============================================================================

func TestGitOpsMenuItemsExist(t *testing.T) {
	if len(gitOpsMenuItems) == 0 {
		t.Error("gitOpsMenuItems should not be empty")
	}

	if len(gitOpsMenuItems) < 2 {
		t.Errorf("gitOpsMenuItems should have at least 2 items, got %d", len(gitOpsMenuItems))
	}

	if gitOpsMenuLen != len(gitOpsMenuItems) {
		t.Errorf("gitOpsMenuLen should equal len(gitOpsMenuItems): %d vs %d",
			gitOpsMenuLen, len(gitOpsMenuItems))
	}
}

func TestGitOpsMenuRenders(t *testing.T) {
	m := newModel(".")
	m.showGitOpsMenu = true
	m.gitOpsMenuIdx = 0

	output := renderGitOpsMenuOverlay(m, 80)

	if output == "" {
		t.Error("Git ops menu overlay should not be empty")
	}

	// Should show navigation hints
	if !strings.Contains(output, "j/k") && !strings.Contains(output, "Enter") {
		t.Errorf("Git ops menu should show navigation hints. Got: %s", output)
	}
}

func TestGKeyOpensGitOpsMenu(t *testing.T) {
	m := newModel(".")
	m.showGitOpsMenu = false

	// Simulate pressing 'g'
	km := createKeyMsg("g")
	updatedModel, _ := m.Update(km)
	m = updatedModel.(model)

	if !m.showGitOpsMenu {
		t.Error("'g' key should open git ops menu")
	}
}

func TestGKeyTogglesGitOpsMenu(t *testing.T) {
	m := newModel(".")
	m.showGitOpsMenu = true

	// Simulate pressing 'g'
	km := createKeyMsg("g")
	updatedModel, _ := m.Update(km)
	m = updatedModel.(model)

	if m.showGitOpsMenu {
		t.Error("'g' key should toggle off git ops menu")
	}
}

func TestGitOpsMenuNavigation(t *testing.T) {
	m := newModel(".")
	m.showGitOpsMenu = true
	m.gitOpsMenuIdx = 0

	// Test down
	km := createKeyMsg("j")
	updatedModel, _ := m.Update(km)
	m = updatedModel.(model)
	if m.gitOpsMenuIdx != 1 {
		t.Errorf("'j' should move cursor down, got index %d", m.gitOpsMenuIdx)
	}

	// Test up
	km = createKeyMsg("k")
	updatedModel, _ = m.Update(km)
	m = updatedModel.(model)
	if m.gitOpsMenuIdx != 0 {
		t.Errorf("'k' should move cursor up, got index %d", m.gitOpsMenuIdx)
	}
}

func TestGitOpsMenuDispatchesFeatures(t *testing.T) {
	m := newModel(".")
	m.showGitOpsMenu = true
	m.showRebaseUI = false
	m.gitOpsMenuIdx = 0 // Rebase

	// Simulate pressing Enter
	km := createKeyMsg("enter")
	updatedModel, _ := m.Update(km)
	m = updatedModel.(model)

	// Menu should close after selection
	if m.showGitOpsMenu {
		t.Error("Git ops menu should close after selecting a feature")
	}

	// First feature (rebase) should be toggled
	if !m.showRebaseUI {
		t.Error("Rebase feature should be activated when dispatched")
	}
}

// ============================================================================
// OPTION 3: Menu Rendering Refactor Tests
// ============================================================================

func TestMenuRenderingFunctionExists(t *testing.T) {
	// Verify that renderMenuOverlay helper function exists and works
	m := newModel(".")
	m.analyticsMenuIdx = 0

	config := RenderConfig{
		Title: "TEST MENU",
		Items: []string{"Item 1", "Item 2", "Item 3"},
	}

	output := renderMenuOverlay(config)

	if output == "" {
		t.Error("renderMenuOverlay should produce non-empty output")
	}

	if !strings.Contains(output, "TEST MENU") {
		t.Error("Menu should contain title")
	}

	if !strings.Contains(output, "j/k") {
		t.Error("Menu should contain navigation hints")
	}
}

func TestAllMenusUseConsistentHints(t *testing.T) {
	menus := []struct {
		name   string
		render func(m model, w int) string
	}{
		{"Analytics", renderAnalyticsMenuOverlay},
		{"Visualization", renderVisualizationMenuOverlay},
		{"Team", renderTeamMenuOverlay},
		{"Integration", renderIntegrationMenuOverlay},
		{"GitOps", renderGitOpsMenuOverlay},
	}

	m := newModel(".")
	m.width = 80

	for _, menu := range menus {
		output := menu.render(m, 80)

		if !strings.Contains(output, "j/k") {
			t.Errorf("%s menu missing j/k hint", menu.name)
		}
		if !strings.Contains(output, "Enter") {
			t.Errorf("%s menu missing Enter hint", menu.name)
		}
		if !strings.Contains(output, "Esc") {
			t.Errorf("%s menu missing Esc hint", menu.name)
		}
	}
}

func TestViewMethodIncludesIntegrationFeatures(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 20

	// Test each integration feature is rendered when shown
	tests := []struct {
		flag   bool
		ptr    *bool
		name   string
		check  func(string) bool
	}{
		{true, &m.showPRLinks, "PR Links", func(s string) bool {
			return strings.Contains(s, "PR") || strings.Contains(s, "Link")
		}},
		{true, &m.showJiraLinks, "Jira Links", func(s string) bool {
			return strings.Contains(s, "JIRA") || strings.Contains(s, "Jira")
		}},
		{true, &m.showIssueRefs, "Issue Refs", func(s string) bool {
			return strings.Contains(s, "Issue") || strings.Contains(s, "Reference")
		}},
	}

	for _, test := range tests {
		*test.ptr = true
		output := m.View()
		*test.ptr = false

		if !test.check(output) {
			t.Errorf("%s feature not rendered when flag is true", test.name)
		}
	}
}

// ============================================================================
// Git Ops Menu Structure Tests
// ============================================================================

func TestGitOpsMenuHasRebaseOption(t *testing.T) {
	found := false
	for _, item := range gitOpsMenuItems {
		if strings.Contains(item, "Rebase") || strings.Contains(item, "rebase") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Git ops menu should include rebase option")
	}
}

func TestGitOpsMenuHasCherryPickOption(t *testing.T) {
	found := false
	for _, item := range gitOpsMenuItems {
		if strings.Contains(item, "Cherry") || strings.Contains(item, "cherry") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Git ops menu should include cherry-pick option")
	}
}

func TestGitOpsMenuHasResetOption(t *testing.T) {
	found := false
	for _, item := range gitOpsMenuItems {
		if strings.Contains(item, "Reset") || strings.Contains(item, "reset") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Git ops menu should include reset option")
	}
}
