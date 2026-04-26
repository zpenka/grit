package grit

import (
	"testing"
)

func TestRenderCommitRow_IndexOutOfBounds(t *testing.T) {
	m := newModel(".")
	vc := makeCommits(3)

	row := renderCommitRow(m, vc, -1, 80)
	if row == "" {
		t.Error("renderCommitRow with invalid index should return empty string")
	}
}

func TestRenderCommitRow_ValidIndex(t *testing.T) {
	m := newModel(".")
	m.cursor = 0
	m.focus = panelList
	vc := []commit{
		{hash: "abc123", shortHash: "abc", author: "Alice", when: "1d ago", subject: "Fix bug"},
	}

	row := renderCommitRow(m, vc, 0, 80)
	if row == "" {
		t.Error("renderCommitRow with valid index should return non-empty string")
	}
	if !contains(row, "abc") || !contains(row, "Alice") || !contains(row, "Fix bug") {
		t.Error("renderCommitRow should include hash, author, and subject")
	}
}

func TestRenderCommitRow_SelectedCursor(t *testing.T) {
	m := newModel(".")
	m.cursor = 0
	m.focus = panelList
	vc := []commit{
		{hash: "abc123", shortHash: "abc", author: "Alice", when: "1d ago", subject: "Fix"},
	}

	row := renderCommitRow(m, vc, 0, 80)
	if !contains(row, "▶") {
		t.Error("renderCommitRow should show cursor indicator for selected row")
	}
}

func TestRenderCommitRow_NarrowWidth(t *testing.T) {
	m := newModel(".")
	m.cursor = 0
	vc := makeCommits(1)

	row := renderCommitRow(m, vc, 0, 20)
	if len(row) == 0 {
		t.Error("renderCommitRow should handle narrow widths")
	}
}

func TestRenderFileRow_IndexOutOfBounds(t *testing.T) {
	m := newModel(".")

	row := renderFileRow(m, -1, 80)
	if row == "" {
		t.Error("renderFileRow with invalid index should return empty string")
	}
}

func TestRenderFileRow_ValidIndex(t *testing.T) {
	m := newModel(".")
	m.fileCursor = 0
	m.fileItems = []fileItem{
		{path: "src/main.go", diffIdx: 0},
	}

	row := renderFileRow(m, 0, 80)
	if row == "" {
		t.Error("renderFileRow with valid index should return non-empty string")
	}
	if !contains(row, "main.go") {
		t.Error("renderFileRow should include file path")
	}
}

func TestRenderFileRow_SelectedCursor(t *testing.T) {
	m := newModel(".")
	m.fileCursor = 0
	m.fileItems = []fileItem{
		{path: "file.txt", diffIdx: 0},
	}

	row := renderFileRow(m, 0, 80)
	if !contains(row, "▶") {
		t.Error("renderFileRow should show cursor for selected row")
	}
}

func TestRenderBranchRow_IndexOutOfBounds(t *testing.T) {
	m := newModel(".")

	row := renderBranchRow(m, -1, 80)
	if row == "" {
		t.Error("renderBranchRow with invalid index should return empty string")
	}
}

func TestRenderBranchRow_ValidIndex(t *testing.T) {
	m := newModel(".")
	m.branchCursor = 0
	m.branches = []string{"main", "develop"}
	m.currentBranch = "main"

	row := renderBranchRow(m, 0, 80)
	if row == "" {
		t.Error("renderBranchRow with valid index should return non-empty string")
	}
	if !contains(row, "main") {
		t.Error("renderBranchRow should include branch name")
	}
}

func TestRenderBranchRow_CurrentBranchMarker(t *testing.T) {
	m := newModel(".")
	m.branchCursor = 0
	m.branches = []string{"main", "develop"}
	m.currentBranch = "main"

	row := renderBranchRow(m, 0, 80)
	if !contains(row, "●") {
		t.Error("renderBranchRow should show marker for current branch")
	}
}

func TestRenderBlameRow_IndexOutOfBounds(t *testing.T) {
	m := newModel(".")

	row := renderBlameRow(m, -1, 80)
	if row != "" {
		t.Error("renderBlameRow with invalid index should return empty string")
	}
}

func TestRenderBlameRow_ValidIndex(t *testing.T) {
	m := newModel(".")
	m.blameLines = []blameLine{
		{shortHash: "abc", author: "Alice", date: "2024-01-01", lineNum: 1, text: "code"},
	}

	row := renderBlameRow(m, 0, 80)
	if row == "" {
		t.Error("renderBlameRow with valid index should return non-empty string")
	}
	if !contains(row, "abc") || !contains(row, "Alice") {
		t.Error("renderBlameRow should include hash and author")
	}
}

func TestRenderDiffRow_IndexOutOfBounds(t *testing.T) {
	m := newModel(".")

	row := renderDiffRow(m, -1, 80)
	if row != "" {
		t.Error("renderDiffRow with invalid index should return empty string")
	}
}

func TestRenderDiffRow_AddedLine(t *testing.T) {
	m := newModel(".")
	m.diffLines = []diffLine{
		{kind: lineAdded, text: "+ new code"},
	}

	row := renderDiffRow(m, 0, 80)
	if row == "" {
		t.Error("renderDiffRow with added line should return non-empty string")
	}
}

func TestRenderDiffRow_RemovedLine(t *testing.T) {
	m := newModel(".")
	m.diffLines = []diffLine{
		{kind: lineRemoved, text: "- old code"},
	}

	row := renderDiffRow(m, 0, 80)
	if row == "" {
		t.Error("renderDiffRow with removed line should return non-empty string")
	}
}

func TestRenderDiffRow_ContextLine(t *testing.T) {
	m := newModel(".")
	m.diffLines = []diffLine{
		{kind: lineContext, text: "  existing code"},
	}

	row := renderDiffRow(m, 0, 80)
	if row == "" {
		t.Error("renderDiffRow with context line should return non-empty string")
	}
}

func TestRenderDiffRow_HunkHeader(t *testing.T) {
	m := newModel(".")
	m.diffLines = []diffLine{
		{kind: lineHunk, text: "@@ -1,5 +1,6 @@"},
	}

	row := renderDiffRow(m, 0, 80)
	if row == "" {
		t.Error("renderDiffRow with hunk header should return non-empty string")
	}
}

func TestRenderDiffRow_MetaLine(t *testing.T) {
	m := newModel(".")
	m.diffLines = []diffLine{
		{kind: lineMeta, text: "index abc123..def456 100644"},
	}

	row := renderDiffRow(m, 0, 80)
	if row == "" {
		t.Error("renderDiffRow with meta line should return non-empty string")
	}
}

func TestRenderDiffRow_NarrowWidth(t *testing.T) {
	m := newModel(".")
	m.diffLines = []diffLine{
		{kind: lineAdded, text: "+ this is a very long line that needs to be truncated"},
	}

	row := renderDiffRow(m, 0, 20)
	if row == "" {
		t.Error("renderDiffRow should handle narrow widths")
	}
}
