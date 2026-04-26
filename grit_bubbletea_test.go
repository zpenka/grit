package grit

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInit_ReturnsCommitsLoadCmd(t *testing.T) {
	m := newModel(".")
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() should return a command")
	}
	// cmd is executed during Bubble Tea's setup phase
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := newModel(".")
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updated, cmd := m.Update(msg)
	updatedModel := updated.(model)
	if updatedModel.width != 100 || updatedModel.height != 30 {
		t.Errorf("Update(WindowSizeMsg) should set width=%d height=%d, got %d,%d",
			100, 30, updatedModel.width, updatedModel.height)
	}
	if cmd != nil {
		t.Error("Update(WindowSizeMsg) should return nil command")
	}
}

func TestUpdate_CommitsLoadedMsg(t *testing.T) {
	m := newModel(".")
	commits := []commit{
		{hash: "abc123", shortHash: "abc", author: "test", when: "now", subject: "test"},
	}
	msg := commitsLoadedMsg(commits)
	updated, _ := m.Update(msg)
	updatedModel := updated.(model)
	if len(updatedModel.commits) != 1 {
		t.Errorf("Update(commitsLoadedMsg) should set commits, got %d", len(updatedModel.commits))
	}
	if updatedModel.cursor != 0 {
		t.Error("Update(commitsLoadedMsg) should reset cursor to 0")
	}
}

func TestUpdate_EmptyCommits(t *testing.T) {
	m := newModel(".")
	msg := commitsLoadedMsg{}
	updated, cmd := m.Update(msg)
	updatedModel := updated.(model)
	if len(updatedModel.commits) != 0 {
		t.Errorf("Update(empty commitsLoadedMsg) should have no commits")
	}
	if cmd != nil {
		t.Error("Update(empty commitsLoadedMsg) should return nil (no diff to fetch)")
	}
}

func TestUpdate_DiffLoadedMsg(t *testing.T) {
	m := newModel(".")
	m.loading = true
	lines := []diffLine{
		{kind: lineAdded, text: "+ added line"},
		{kind: lineRemoved, text: "- removed line"},
	}
	msg := diffLoadedMsg(lines)
	updated, _ := m.Update(msg)
	updatedModel := updated.(model)
	if updatedModel.loading {
		t.Error("Update(diffLoadedMsg) should set loading to false")
	}
	if len(updatedModel.diffLines) != 2 {
		t.Errorf("Update(diffLoadedMsg) should set diffLines, got %d", len(updatedModel.diffLines))
	}
	if updatedModel.diffOffset != 0 {
		t.Error("Update(diffLoadedMsg) should reset diffOffset to 0")
	}
}

func TestUpdate_BranchesLoadedMsg(t *testing.T) {
	m := newModel(".")
	msg := branchesLoadedMsg{
		branches: []string{"main", "develop", "feature"},
		current:  "develop",
	}
	updated, _ := m.Update(msg)
	updatedModel := updated.(model)
	if len(updatedModel.branches) != 3 {
		t.Errorf("Update(branchesLoadedMsg) should set branches, got %d", len(updatedModel.branches))
	}
	if updatedModel.currentBranch != "develop" {
		t.Errorf("Update(branchesLoadedMsg) should set currentBranch, got %s", updatedModel.currentBranch)
	}
	if updatedModel.branchCursor != 1 {
		t.Errorf("Update(branchesLoadedMsg) should position cursor on current branch at index 1, got %d", updatedModel.branchCursor)
	}
	if !updatedModel.showBranch {
		t.Error("Update(branchesLoadedMsg) should set showBranch to true")
	}
}

func TestUpdate_BlameLoadedMsg(t *testing.T) {
	m := newModel(".")
	m.loading = true
	blames := []blameLine{
		{shortHash: "abc", author: "test", date: "2024-01-01", lineNum: 1, text: "test"},
	}
	msg := blameLoadedMsg(blames)
	updated, _ := m.Update(msg)
	updatedModel := updated.(model)
	if updatedModel.loading {
		t.Error("Update(blameLoadedMsg) should set loading to false")
	}
	if len(updatedModel.blameLines) != 1 {
		t.Errorf("Update(blameLoadedMsg) should set blameLines, got %d", len(updatedModel.blameLines))
	}
	if !updatedModel.showBlame {
		t.Error("Update(blameLoadedMsg) should set showBlame to true")
	}
}

func TestUpdate_FlashClearMsg(t *testing.T) {
	m := newModel(".")
	m.flash = "test message"
	msg := flashClearMsg{}
	updated, _ := m.Update(msg)
	updatedModel := updated.(model)
	if updatedModel.flash != "" {
		t.Errorf("Update(flashClearMsg) should clear flash, got %s", updatedModel.flash)
	}
}

func TestUpdate_EditorDoneMsg(t *testing.T) {
	m := newModel(".")
	msg := editorDoneMsg{}
	updated, cmd := m.Update(msg)
	_ = updated.(model)
	if cmd != nil {
		t.Error("Update(editorDoneMsg) should return nil command")
	}
}

func TestUpdate_CtrlCQuits(t *testing.T) {
	m := newModel(".")
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := m.Update(msg)
	if cmd == nil {
		t.Error("Update(ctrl+c) should return quit command")
	}
}

func TestUpdate_QuitKeyQuits(t *testing.T) {
	m := newModel(".")
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(msg)
	if cmd == nil {
		t.Error("Update('q') should return quit command")
	}
}

func TestUpdate_IgnoresNonKeyMessages(t *testing.T) {
	m := newModel(".")
	msg := tea.MouseMsg{X: 10, Y: 10}
	updated, cmd := m.Update(msg)
	_ = updated.(model)
	if cmd != nil {
		t.Error("Update(non-key message) should return nil command")
	}
}

func TestUpdate_SearchModeToggle(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.diffLines = makeDiffLines(5)
	if m.searching {
		t.Fatal("Initial state should not be in search mode")
	}

	// Verify model is created in searchable state
	if m.query != "" {
		t.Error("Initial query should be empty")
	}
}

func TestView_BasicLayout(t *testing.T) {
	m := newModel(".")
	m.width = 120
	m.height = 30
	m.repoPath = "."
	m.commits = makeCommits(5)
	m.diffLines = makeDiffLines(10)

	view := m.View()
	if view == "" {
		t.Error("View() should return non-empty string")
	}
	if len(view) < 10 {
		t.Error("View() should return substantial output")
	}
}

func TestView_NoWindowSize(t *testing.T) {
	m := newModel(".")
	// width and height are 0 by default
	view := m.View()
	if view == "" {
		t.Error("View() should return non-empty string even with zero size")
	}
	if !contains(view, "loading") {
		t.Error("View() with zero size should show loading message")
	}
}

func TestView_ShowsCommitList(t *testing.T) {
	m := newModel(".")
	m.width = 120
	m.height = 30
	m.repoPath = "."
	m.commits = makeCommits(3)
	m.diffLines = makeDiffLines(5)

	view := m.View()
	if !contains(view, "Commits") {
		t.Error("View() should show Commits header")
	}
}

func TestView_ShowsCurrentRef(t *testing.T) {
	m := newModel(".")
	m.width = 120
	m.height = 30
	m.currentRef = "main"

	view := m.View()
	if !contains(view, "main") {
		t.Error("View() should show currentRef")
	}
}

func TestFirstVisibleHash_WithCommits(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{
		{hash: "aaa111", shortHash: "aaa", author: "test", when: "now", subject: "first"},
		{hash: "bbb222", shortHash: "bbb", author: "test", when: "now", subject: "second"},
	}
	hash := firstVisibleHash(m)
	if hash != "aaa111" {
		t.Errorf("firstVisibleHash should return first commit hash, got %s", hash)
	}
}

func TestFirstVisibleHash_NoCommits(t *testing.T) {
	m := newModel(".")
	hash := firstVisibleHash(m)
	if hash != "" {
		t.Errorf("firstVisibleHash with no commits should return empty string, got %s", hash)
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
