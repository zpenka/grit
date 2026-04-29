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

// ============================================================================
// Export to Markdown Tests
// ============================================================================

func TestExportToMarkdownBasic(t *testing.T) {
	commits := []commit{
		{hash: "abc123def456", shortHash: "abc123", author: "Alice", subject: "Add feature X"},
		{hash: "def456ghi789", shortHash: "def456", author: "Bob", subject: "Fix bug Y"},
	}

	result := exportToMarkdown(commits)

	// Check format
	if result.format != "markdown" {
		t.Errorf("Expected format 'markdown', got %q", result.format)
	}

	// Check filename
	if result.filename != "commits.md" {
		t.Errorf("Expected filename 'commits.md', got %q", result.filename)
	}

	// Check commits are preserved
	if len(result.commits) != len(commits) {
		t.Errorf("Expected %d commits, got %d", len(commits), len(result.commits))
	}

	// Check content structure
	content := result.content
	if !strings.Contains(content, "# Commit History") {
		t.Error("Content should contain markdown heading '# Commit History'")
	}

	// Check content contains commit data
	if !strings.Contains(content, "abc123") {
		t.Error("Content should contain first commit hash")
	}
	if !strings.Contains(content, "def456") {
		t.Error("Content should contain second commit hash")
	}
	if !strings.Contains(content, "Add feature X") {
		t.Error("Content should contain first commit subject")
	}
	if !strings.Contains(content, "Fix bug Y") {
		t.Error("Content should contain second commit subject")
	}
	if !strings.Contains(content, "Alice") {
		t.Error("Content should contain first commit author")
	}
	if !strings.Contains(content, "Bob") {
		t.Error("Content should contain second commit author")
	}
}

func TestExportToMarkdownEmpty(t *testing.T) {
	commits := []commit{}

	result := exportToMarkdown(commits)

	// Check format and filename still set
	if result.format != "markdown" {
		t.Errorf("Expected format 'markdown', got %q", result.format)
	}
	if result.filename != "commits.md" {
		t.Errorf("Expected filename 'commits.md', got %q", result.filename)
	}

	// Check content has heading even for empty list
	content := result.content
	if !strings.Contains(content, "# Commit History") {
		t.Error("Content should contain markdown heading even for empty commits")
	}

	// Should not panic
	if len(result.commits) != 0 {
		t.Error("Empty input should result in empty commits slice")
	}
}

func TestExportToMarkdownSingleCommit(t *testing.T) {
	commits := []commit{
		{hash: "abc123def456", shortHash: "abc123", author: "Charlie", subject: "Initial commit"},
	}

	result := exportToMarkdown(commits)

	content := result.content
	if !strings.Contains(content, "# Commit History") {
		t.Error("Content should contain markdown heading")
	}
	if !strings.Contains(content, "abc123") {
		t.Error("Content should contain commit hash")
	}
	if !strings.Contains(content, "Charlie") {
		t.Error("Content should contain author")
	}
	if !strings.Contains(content, "Initial commit") {
		t.Error("Content should contain subject")
	}

	// Verify format
	if !strings.Contains(content, "-") {
		t.Error("Markdown should use bullet list format with dashes")
	}
}

func TestExportToMarkdownPreservesCommitOrder(t *testing.T) {
	commits := []commit{
		{hash: "111", shortHash: "111", author: "Author1", subject: "First"},
		{hash: "222", shortHash: "222", author: "Author2", subject: "Second"},
		{hash: "333", shortHash: "333", author: "Author3", subject: "Third"},
	}

	result := exportToMarkdown(commits)
	content := result.content

	// Find positions of each hash in content
	pos1 := strings.Index(content, "111")
	pos2 := strings.Index(content, "222")
	pos3 := strings.Index(content, "333")

	if pos1 == -1 || pos2 == -1 || pos3 == -1 {
		t.Error("Content should contain all commit hashes")
		return
	}

	if pos1 >= pos2 || pos2 >= pos3 {
		t.Error("Commits should be in order in the content")
	}
}

// ============================================================================
// Export Patch Series Tests
// ============================================================================

func TestExportPatchSeriesBasic(t *testing.T) {
	commits := []commit{
		{hash: "abc123def456", shortHash: "abc123", author: "Alice", subject: "Add feature X"},
		{hash: "def456ghi789", shortHash: "def456", author: "Bob", subject: "Fix bug Y"},
	}

	result := exportPatchSeries(commits)

	// Check format
	if result.format != "patch" {
		t.Errorf("Expected format 'patch', got %q", result.format)
	}

	// Check filename
	if result.filename != "series.patch" {
		t.Errorf("Expected filename 'series.patch', got %q", result.filename)
	}

	// Check commits are preserved
	if len(result.commits) != len(commits) {
		t.Errorf("Expected %d commits, got %d", len(commits), len(result.commits))
	}

	// Check content structure
	content := result.content
	if !strings.Contains(content, "From:") {
		t.Error("Patch content should contain 'From:' header")
	}

	// Check for Subject headers for each commit
	if !strings.Contains(content, "Subject: Add feature X") {
		t.Error("Content should contain first commit subject in patch format")
	}
	if !strings.Contains(content, "Subject: Fix bug Y") {
		t.Error("Content should contain second commit subject in patch format")
	}
}

func TestExportPatchSeriesEmpty(t *testing.T) {
	commits := []commit{}

	result := exportPatchSeries(commits)

	// Check format and filename still set
	if result.format != "patch" {
		t.Errorf("Expected format 'patch', got %q", result.format)
	}
	if result.filename != "series.patch" {
		t.Errorf("Expected filename 'series.patch', got %q", result.filename)
	}

	// Check content has From header even for empty list
	content := result.content
	if !strings.Contains(content, "From:") {
		t.Error("Patch content should contain 'From:' header even for empty commits")
	}

	// Should not panic
	if len(result.commits) != 0 {
		t.Error("Empty input should result in empty commits slice")
	}
}

func TestExportPatchSeriesSingleCommit(t *testing.T) {
	commits := []commit{
		{hash: "abc123def456", shortHash: "abc123", author: "Charlie", subject: "Initial commit"},
	}

	result := exportPatchSeries(commits)

	content := result.content
	if !strings.Contains(content, "From:") {
		t.Error("Content should contain From header")
	}
	if !strings.Contains(content, "Subject: Initial commit") {
		t.Error("Content should contain commit subject in patch format")
	}

	// Should have only one Subject line for single commit
	subjectCount := strings.Count(content, "Subject:")
	if subjectCount != 1 {
		t.Errorf("Expected 1 Subject header for single commit, got %d", subjectCount)
	}
}

func TestExportPatchSeriesPreservesCommitOrder(t *testing.T) {
	commits := []commit{
		{hash: "111", shortHash: "111", author: "Author1", subject: "First"},
		{hash: "222", shortHash: "222", author: "Author2", subject: "Second"},
		{hash: "333", shortHash: "333", author: "Author3", subject: "Third"},
	}

	result := exportPatchSeries(commits)
	content := result.content

	// Find positions of each subject in content
	posFirst := strings.Index(content, "Subject: First")
	posSecond := strings.Index(content, "Subject: Second")
	posThird := strings.Index(content, "Subject: Third")

	if posFirst == -1 || posSecond == -1 || posThird == -1 {
		t.Error("Content should contain all commit subjects")
		return
	}

	if posFirst >= posSecond || posSecond >= posThird {
		t.Error("Commits should be in order in the patch content")
	}
}

func TestExportPatchSeriesFromHeader(t *testing.T) {
	commits := []commit{
		{hash: "abc123", shortHash: "abc123", author: "Test Author", subject: "Test"},
	}

	result := exportPatchSeries(commits)
	content := result.content

	// From header should be present at start
	if !strings.HasPrefix(strings.TrimSpace(content), "From:") {
		t.Error("Patch content should start with From header")
	}

	// Should contain the default email
	if !strings.Contains(content, "user@example.com") {
		t.Error("From header should contain user@example.com")
	}
}
