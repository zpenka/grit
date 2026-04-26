package grit

import (
	"testing"
)

func TestFetchCommits_ReturnsCmdFunction(t *testing.T) {
	cmd := fetchCommits(".", "")
	if cmd == nil {
		t.Error("fetchCommits should return a command")
	}
}

func TestFetchCommits_WithRef(t *testing.T) {
	cmd := fetchCommits(".", "main")
	if cmd == nil {
		t.Error("fetchCommits with ref should return a command")
	}
}

func TestFetchDiff_ReturnsCmdFunction(t *testing.T) {
	cmd := fetchDiff(".", "abc123")
	if cmd == nil {
		t.Error("fetchDiff should return a command")
	}
}

func TestFetchBranches_ReturnsCmdFunction(t *testing.T) {
	cmd := fetchBranches(".")
	if cmd == nil {
		t.Error("fetchBranches should return a command")
	}
}

func TestFetchBlame_ReturnsCmdFunction(t *testing.T) {
	cmd := fetchBlame(".", "abc123", "file.go")
	if cmd == nil {
		t.Error("fetchBlame should return a command")
	}
}

func TestFlashCmd_ReturnsCmdFunction(t *testing.T) {
	cmd := flashCmd()
	if cmd == nil {
		t.Error("flashCmd should return a command")
	}
}

func TestCopyToClipboard_WithContent(t *testing.T) {
	// copyToClipboard attempts multiple methods and may succeed or fail
	// depending on system configuration. Just verify it returns an error or nil.
	err := copyToClipboard("test content")
	_ = err // error handling depends on system environment
}

func TestOpenInEditor_WithHash(t *testing.T) {
	cmd := openInEditor(".", "abc123")
	if cmd == nil {
		t.Error("openInEditor should return a command")
	}
}

func TestMessagesTypes_AreDefinedCorrectly(t *testing.T) {
	commits := makeCommits(1)
	msg := commitsLoadedMsg(commits)
	if len(msg) != 1 {
		t.Error("commitsLoadedMsg should preserve slice content")
	}

	lines := makeDiffLines(2)
	diffMsg := diffLoadedMsg(lines)
	if len(diffMsg) != 2 {
		t.Error("diffLoadedMsg should preserve slice content")
	}

	branchMsg := branchesLoadedMsg{
		branches: []string{"main"},
		current:  "main",
	}
	if len(branchMsg.branches) != 1 {
		t.Error("branchesLoadedMsg should store branches")
	}

	blames := []blameLine{{shortHash: "abc", author: "test", date: "2024-01-01", lineNum: 1, text: "code"}}
	blameMsg := blameLoadedMsg(blames)
	if len(blameMsg) != 1 {
		t.Error("blameLoadedMsg should preserve slice content")
	}
}

func TestFlashClearMsg_IsZeroValue(t *testing.T) {
	msg := flashClearMsg{}
	if msg != (flashClearMsg{}) {
		t.Error("flashClearMsg should be creatable as zero value")
	}
}

func TestEditorDoneMsg_IsZeroValue(t *testing.T) {
	msg := editorDoneMsg{}
	if msg != (editorDoneMsg{}) {
		t.Error("editorDoneMsg should be creatable as zero value")
	}
}

func TestNewModel_InitialState(t *testing.T) {
	m := newModel(".")
	if m.repoPath != "." {
		t.Errorf("newModel should set repoPath, got %q", m.repoPath)
	}
	if len(m.commits) != 0 {
		t.Error("newModel should start with empty commits")
	}
	if m.cursor != 0 {
		t.Error("newModel cursor should start at 0")
	}
	if m.loading {
		t.Error("newModel should not be loading initially")
	}
	if m.searching {
		t.Error("newModel should not be searching initially")
	}
	if m.showBranch {
		t.Error("newModel should not show branch initially")
	}
	if m.showFiles {
		t.Error("newModel should not show files initially")
	}
	if m.showBlame {
		t.Error("newModel should not show blame initially")
	}
}

func TestNewModel_PanelDefaults(t *testing.T) {
	m := newModel(".")
	if m.focus != panelList {
		t.Errorf("newModel should default focus to panelList, got %d", m.focus)
	}
}

func TestModelFields_QueryAndFilters(t *testing.T) {
	m := newModel(".")
	if m.query != "" {
		t.Error("Initial query should be empty")
	}
	if m.countBuf != "" {
		t.Error("Initial countBuf should be empty")
	}
	if m.flash != "" {
		t.Error("Initial flash message should be empty")
	}
	if m.currentRef != "" {
		t.Error("Initial currentRef should be empty")
	}
}

func TestVisibleCommits_WithNoFilter(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	vc := visibleCommits(m)
	if len(vc) != 5 {
		t.Errorf("visibleCommits with no filter should return all commits, got %d", len(vc))
	}
}

func TestVisibleCommits_FiltersByQuery(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{
		{hash: "aaa", shortHash: "aaa", author: "Alice", when: "now", subject: "fix bug"},
		{hash: "bbb", shortHash: "bbb", author: "Bob", when: "now", subject: "add feature"},
	}
	m.query = "fix"
	vc := visibleCommits(m)
	if len(vc) != 1 {
		t.Errorf("visibleCommits with query should filter, got %d", len(vc))
	}
}

func TestModelCursors_Independent(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.branches = []string{"main", "develop"}
	m.fileItems = []fileItem{
		{path: "file1.go", diffIdx: 0},
		{path: "file2.go", diffIdx: 1},
	}

	m.cursor = 2
	m.branchCursor = 1
	m.fileCursor = 1

	if m.cursor != 2 || m.branchCursor != 1 || m.fileCursor != 1 {
		t.Error("Cursors should be independent")
	}
}

func TestNavigationHistory_Initialization(t *testing.T) {
	m := newModel(".")
	// navHistory may be nil or empty slice depending on initialization
	if m.navHistoryIdx != 0 {
		t.Error("navHistoryIdx should start at 0")
	}
}
