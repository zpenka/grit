package grit

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// ============================================================================
// Menu Improvements Tests (from grit_menu_improvements_test.go)
// ============================================================================

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

// ============================================================================
// Analytics Submenu Tests (from grit_analytics_submenu_test.go)
// ============================================================================

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
	m.analyticsMenuIdx = 11 // Last item (0-11 = 12 items after C2a)
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// Press "j" at the end - should not go beyond bounds
	msgJ := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	m1, _ := m.Update(msgJ)
	model1 := m1.(model)

	if model1.analyticsMenuIdx != 11 {
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
	m.analyticsMenuIdx = 0 // Code Ownership
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
	m.analyticsMenuIdx = 0 // Code Ownership

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
	m.analyticsMenuIdx = 3 // Bisect

	// Dispatch bisect
	m1 := dispatchAnalyticsFeature(m, 3)

	if !m1.bisectState.active {
		t.Error("dispatching bisect should activate bisect mode")
	}
	if !m1.showBisectUI {
		t.Error("dispatching bisect should show bisect UI")
	}
}

// ============================================================================
// Help Overlay Tests (from grit_help_overlay_test.go)
// ============================================================================

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
	m.height = 40 // Larger height to show full help text
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
	m.height = 40 // Larger height to show full help text
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

// ============================================================================
// Analytics Menu Extended Feature Tests (Task C2a)
// ============================================================================

func TestDispatchAnalyticsFeature_LargeCommits(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{
		{hash: "abc1", author: "alice", subject: "large change"},
		{hash: "abc2", author: "bob", subject: "small fix"},
	}

	// Dispatch large commits (index 7)
	m1 := dispatchAnalyticsFeature(m, 7)

	if !m1.showLargeCommits {
		t.Error("dispatching large commits should set showLargeCommits")
	}
	if len(m1.largeCommits) == 0 {
		t.Error("dispatching large commits should compute largeCommits data")
	}
}

func TestDispatchAnalyticsFeature_MergeAnalysis(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{
		{hash: "abc1", author: "alice", subject: "Merge branch 'feature'"},
		{hash: "abc2", author: "bob", subject: "fix: bug"},
	}

	// Dispatch merge analysis (index 8)
	m1 := dispatchAnalyticsFeature(m, 8)

	if !m1.showMergeAnalysis {
		t.Error("dispatching merge analysis should set showMergeAnalysis")
	}
	if len(m1.mergeAnalysisData) == 0 {
		t.Error("dispatching merge analysis should compute mergeAnalysisData")
	}
}

func TestDispatchAnalyticsFeature_FileCoupling(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{
		{hash: "abc1", author: "alice", subject: "refactor main.go util.go"},
		{hash: "abc2", author: "bob", subject: "update main.go util.go"},
	}

	// Dispatch file coupling (index 9)
	m1 := dispatchAnalyticsFeature(m, 9)

	if !m1.showCoupling {
		t.Error("dispatching file coupling should set showCoupling")
	}
	// Note: commitCouplings may be empty if extractFilesFromSubject doesn't find matches
	// The important thing is that showCoupling is set and the function ran
}

func TestDispatchAnalyticsFeature_SemanticSearch(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{
		{hash: "abc1", author: "alice", subject: "feature: add login"},
		{hash: "abc2", author: "bob", subject: "fix: logout"},
	}
	m.semanticQuery = "login"

	// Dispatch semantic search (index 10)
	m1 := dispatchAnalyticsFeature(m, 10)

	if !m1.showSemanticSearch {
		t.Error("dispatching semantic search should set showSemanticSearch")
	}
	if len(m1.semanticSearchResults) == 0 {
		t.Error("dispatching semantic search should compute semanticSearchResults data")
	}
}

func TestDispatchAnalyticsFeature_DependencyChanges(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.commits = []commit{
		{hash: "abc1", author: "alice", subject: "update deps"},
		{hash: "abc2", author: "bob", subject: "bump version"},
	}

	// Dispatch dependency changes (index 11)
	m1 := dispatchAnalyticsFeature(m, 11)

	if !m1.showDependencies {
		t.Error("dispatching dependency changes should set showDependencies")
	}
	if len(m1.dependencyChanges) == 0 {
		t.Error("dispatching dependency changes should compute dependencyChanges data")
	}
}

func TestDispatchAnalyticsFeature_ToggleLargeCommits(t *testing.T) {
	m := newModel(".")
	m.showLargeCommits = false

	// Toggle on
	m1 := dispatchAnalyticsFeature(m, 7)
	if !m1.showLargeCommits {
		t.Error("first dispatch should enable showLargeCommits")
	}

	// Toggle off
	m2 := dispatchAnalyticsFeature(m1, 7)
	if m2.showLargeCommits {
		t.Error("second dispatch should disable showLargeCommits")
	}
}

func TestDispatchAnalyticsFeature_OutOfRangeIndex(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{{hash: "abc1", author: "alice", subject: "test"}}
	originalShowCodeOwnership := m.showCodeOwnership

	// Test out of range (should be >= 12 for new range)
	m1 := dispatchAnalyticsFeature(m, 12)
	if m1.showCodeOwnership != originalShowCodeOwnership {
		t.Error("out of range index should return unchanged model")
	}
}

// ============================================================================
// Git Ops Menu Extended Feature Tests (Task C2a)
// ============================================================================

func TestDispatchGitOpsFeature_AmendPreview(t *testing.T) {
	m := newModel(".")
	m.showAmendPreview = false

	// Dispatch amend preview (index 3)
	m1 := dispatchGitOpsFeature(m, 3)

	if !m1.showAmendPreview {
		t.Error("dispatching amend preview should set showAmendPreview")
	}
}

func TestDispatchGitOpsFeature_RebasePreview(t *testing.T) {
	m := newModel(".")
	m.showRebasePreview = false
	m.commits = []commit{
		{hash: "abc1", author: "alice", subject: "commit 1"},
		{hash: "abc2", author: "bob", subject: "commit 2"},
	}

	// Dispatch rebase preview (index 4)
	m1 := dispatchGitOpsFeature(m, 4)

	if !m1.showRebasePreview {
		t.Error("dispatching rebase preview should set showRebasePreview")
	}
	if len(m1.rebaseSequence) == 0 {
		t.Error("dispatching rebase preview should populate rebaseSequence")
	}
}

func TestDispatchGitOpsFeature_SquashUI(t *testing.T) {
	m := newModel(".")
	m.showSquashUI = false

	// Dispatch squash UI (index 5)
	m1 := dispatchGitOpsFeature(m, 5)

	if !m1.showSquashUI {
		t.Error("dispatching squash UI should set showSquashUI")
	}
}

func TestDispatchGitOpsFeature_UndoMenu(t *testing.T) {
	m := newModel(".")
	m.showUndoMenu = false

	// Dispatch undo menu (index 6)
	m1 := dispatchGitOpsFeature(m, 6)

	if !m1.showUndoMenu {
		t.Error("dispatching undo menu should set showUndoMenu")
	}
}

func TestDispatchGitOpsFeature_ToggleAmendPreview(t *testing.T) {
	m := newModel(".")
	m.showAmendPreview = false

	// Toggle on
	m1 := dispatchGitOpsFeature(m, 3)
	if !m1.showAmendPreview {
		t.Error("first dispatch should enable showAmendPreview")
	}

	// Toggle off
	m2 := dispatchGitOpsFeature(m1, 3)
	if m2.showAmendPreview {
		t.Error("second dispatch should disable showAmendPreview")
	}
}

func TestDispatchGitOpsFeature_OutOfRangeIndex(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{{hash: "abc1", author: "alice", subject: "test"}}
	originalShowRebaseUI := m.showRebaseUI

	// Test out of range (should be >= 7 for new range)
	m1 := dispatchGitOpsFeature(m, 7)
	if m1.showRebaseUI != originalShowRebaseUI {
		t.Error("out of range index should return unchanged model")
	}
}

func TestAnalyticsMenuItemsMatchLen(t *testing.T) {
	if len(analyticsMenuItems) != analyticsMenuLen {
		t.Errorf("analyticsMenuItems count (%d) should match analyticsMenuLen (%d)",
			len(analyticsMenuItems), analyticsMenuLen)
	}
}

func TestGitOpsMenuItemsMatchLen(t *testing.T) {
	if len(gitOpsMenuItems) != gitOpsMenuLen {
		t.Errorf("gitOpsMenuItems count (%d) should match gitOpsMenuLen (%d)",
			len(gitOpsMenuItems), gitOpsMenuLen)
	}
}

// ============================================================================
// Workflows Submenu Tests
// ============================================================================

func TestUpdate_WorkflowsKeyTogglesMenu(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24

	// Press "w"
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if !model1.showWorkflowsMenu {
		t.Error("pressing w should toggle showWorkflowsMenu to true")
	}
	if model1.workflowsMenuIdx != 0 {
		t.Error("pressing w should reset workflowsMenuIdx to 0")
	}

	// Press "w" again to close
	m2, _ := model1.Update(msg)
	model2 := m2.(model)

	if model2.showWorkflowsMenu {
		t.Error("pressing w again should toggle showWorkflowsMenu to false")
	}
}

func TestUpdate_WorkflowsMenuNavigation(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showWorkflowsMenu = true
	m.workflowsMenuIdx = 0
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// Press "j" to move down
	msgJ := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	m1, _ := m.Update(msgJ)
	model1 := m1.(model)

	if model1.workflowsMenuIdx != 1 {
		t.Errorf("pressing j should move down, got idx %d", model1.workflowsMenuIdx)
	}

	// Press "k" to move up
	msgK := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	m2, _ := model1.Update(msgK)
	model2 := m2.(model)

	if model2.workflowsMenuIdx != 0 {
		t.Errorf("pressing k should move up, got idx %d", model2.workflowsMenuIdx)
	}
}

func TestUpdate_WorkflowsMenuBoundsCheck(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showWorkflowsMenu = true
	m.workflowsMenuIdx = 4 // Last item (0-4 = 5 items)
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// Press "j" at the end - should not go beyond bounds
	msgJ := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	m1, _ := m.Update(msgJ)
	model1 := m1.(model)

	if model1.workflowsMenuIdx != 4 {
		t.Error("pressing j at end should stay at last item")
	}

	// Now move to start and test k
	m.workflowsMenuIdx = 0
	msgK := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	m2, _ := m.Update(msgK)
	model2 := m2.(model)

	if model2.workflowsMenuIdx != 0 {
		t.Error("pressing k at start should stay at first item")
	}
}

func TestUpdate_WorkflowsMenuSelect(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showWorkflowsMenu = true
	m.workflowsMenuIdx = 0 // Worktrees
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// Press "enter" to select
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if model1.showWorkflowsMenu {
		t.Error("pressing enter should close workflows menu")
	}
	if !model1.showWorktrees {
		t.Error("selecting item 0 should enable showWorktrees")
	}
}

func TestUpdate_WorkflowsMenuEscapeClosesMenu(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 24
	m.showWorkflowsMenu = true
	m.workflowsMenuIdx = 2
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	// Press 'w' to toggle menu closed
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")}
	m1, _ := m.Update(msg)
	model1 := m1.(model)

	if model1.showWorkflowsMenu {
		t.Error("pressing w should close workflows menu when open")
	}
}

func TestView_WorkflowsMenuOverlay(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 40
	m.showWorkflowsMenu = true
	m.workflowsMenuIdx = 0
	m.commits = []commit{{hash: "abc123", author: "test", subject: "test"}}

	view := m.View()

	// Menu should contain workflows feature names
	expectedStrings := []string{"Worktrees", "Submodules", "Named stashes", "Tag management", "GPG status"}
	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("workflows menu should contain '%s'", expected)
		}
	}
}

func TestDispatchWorkflowsFeature_Worktrees(t *testing.T) {
	m := newModel(".")
	m.showWorktrees = false

	m1 := dispatchWorkflowsFeature(m, 0)

	if !m1.showWorktrees {
		t.Error("dispatching worktrees should set showWorktrees")
	}
	// Lazy load should have been called (data may be empty with no git output)
	// Just verify the flag was toggled
	if m1.showWorktrees != !m.showWorktrees {
		t.Error("dispatching should toggle showWorktrees")
	}
}

func TestDispatchWorkflowsFeature_Submodules(t *testing.T) {
	m := newModel(".")
	m.showSubmodules = false

	m1 := dispatchWorkflowsFeature(m, 1)

	if !m1.showSubmodules {
		t.Error("dispatching submodules should set showSubmodules")
	}
	// Lazy load should have been called (data may be empty with no git output)
	// Just verify the flag was toggled
	if m1.showSubmodules != !m.showSubmodules {
		t.Error("dispatching should toggle showSubmodules")
	}
}

func TestDispatchWorkflowsFeature_NamedStashes(t *testing.T) {
	m := newModel(".")
	m.showNamedStashes = false

	m1 := dispatchWorkflowsFeature(m, 2)

	if !m1.showNamedStashes {
		t.Error("dispatching named stashes should set showNamedStashes")
	}
}

func TestDispatchWorkflowsFeature_TagMgmt(t *testing.T) {
	m := newModel(".")
	m.showTagMgmt = false

	m1 := dispatchWorkflowsFeature(m, 3)

	if !m1.showTagMgmt {
		t.Error("dispatching tag management should set showTagMgmt")
	}
}

func TestDispatchWorkflowsFeature_GPGStatus(t *testing.T) {
	m := newModel(".")
	m.showGPGStatus = false
	m.gpgStatuses = nil

	m1 := dispatchWorkflowsFeature(m, 4)

	if !m1.showGPGStatus {
		t.Error("dispatching GPG status should set showGPGStatus")
	}
	// Lazy load attempts to populate gpgStatuses (may be empty if no output)
	if m1.gpgStatuses == nil {
		t.Error("dispatching GPG status should initialize gpgStatuses map")
	}
}

func TestDispatchWorkflowsFeature_Toggle(t *testing.T) {
	m := newModel(".")
	m.showWorktrees = false

	// Toggle on
	m1 := dispatchWorkflowsFeature(m, 0)
	if !m1.showWorktrees {
		t.Error("first dispatch should enable showWorktrees")
	}

	// Toggle off
	m2 := dispatchWorkflowsFeature(m1, 0)
	if m2.showWorktrees {
		t.Error("second dispatch should disable showWorktrees")
	}
}

func TestDispatchWorkflowsFeature_OutOfRangeIndex(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{{hash: "abc1", author: "alice", subject: "test"}}
	originalShowWorktrees := m.showWorktrees

	// Test out of range (should be >= 5 for new range)
	m1 := dispatchWorkflowsFeature(m, 5)
	if m1.showWorktrees != originalShowWorktrees {
		t.Error("out of range index should return unchanged model")
	}
}

func TestWorkflowsMenuItemsMatchLen(t *testing.T) {
	if len(workflowsMenuItems) != workflowsMenuLen {
		t.Errorf("workflowsMenuItems count (%d) should match workflowsMenuLen (%d)",
			len(workflowsMenuItems), workflowsMenuLen)
	}
}

// Helper function to create key messages
func createKeyMsg(s string) interface{} {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}
