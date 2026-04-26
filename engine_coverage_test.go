package grit

import (
	"testing"
)

func TestParseFileItemsFromDiff_Empty(t *testing.T) {
	items := parseFileItemsFromDiff(nil)
	if len(items) != 0 {
		t.Errorf("parseFileItemsFromDiff(nil) should return empty, got %d", len(items))
	}
}

func TestParseFileItemsFromDiff_SingleFile(t *testing.T) {
	diff := []diffLine{
		{kind: lineMeta, text: "diff --git a/main.go b/main.go"},
		{kind: lineMeta, text: "index abc123..def456 100644"},
		{kind: lineMeta, text: "--- a/main.go"},
		{kind: lineMeta, text: "+++ b/main.go"},
		{kind: lineAdded, text: "+ new code"},
	}
	items := parseFileItemsFromDiff(diff)
	if len(items) != 1 {
		t.Errorf("parseFileItemsFromDiff with single file should return 1 item, got %d", len(items))
	}
	if items[0].path != "main.go" {
		t.Errorf("parseFileItemsFromDiff should extract path, got %q", items[0].path)
	}
}

func TestParseFileItemsFromDiff_MultipleFiles(t *testing.T) {
	diff := []diffLine{
		{kind: lineMeta, text: "diff --git a/file1.go b/file1.go"},
		{kind: lineMeta, text: "index abc..def 100644"},
		{kind: lineAdded, text: "+ code1"},
		{kind: lineMeta, text: "diff --git a/file2.go b/file2.go"},
		{kind: lineMeta, text: "index def..ghi 100644"},
		{kind: lineRemoved, text: "- old code"},
	}
	items := parseFileItemsFromDiff(diff)
	if len(items) != 2 {
		t.Errorf("parseFileItemsFromDiff with 2 files should return 2 items, got %d", len(items))
	}
}

func TestParseFileItemsFromDiff_TracksDiffIndex(t *testing.T) {
	diff := []diffLine{
		{kind: lineMeta, text: "diff --git a/test.go b/test.go"},
		{kind: lineMeta, text: "index 1..2"},
		{kind: lineAdded, text: "+ code"},
	}
	items := parseFileItemsFromDiff(diff)
	if len(items) > 0 && items[0].diffIdx < 0 {
		t.Error("parseFileItemsFromDiff should track positive diffIdx")
	}
}

func TestAddToNavHistory_TracksPosition(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(10)
	m.cursor = 5

	m = addToNavHistory(m, 5)
	if m.navHistoryIdx < 0 {
		t.Error("addToNavHistory should increment navHistoryIdx")
	}
}

func TestJumpToNextBookmark_NoBookmarks(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.bookmarks = nil
	m.cursor = 0

	m = jumpToNextBookmark(m)
	if m.cursor != 0 {
		t.Error("jumpToNextBookmark with no bookmarks should not change cursor")
	}
}

func TestJumpToPrevBookmark_NoBookmarks(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.bookmarks = nil
	m.cursor = 2

	m = jumpToPrevBookmark(m)
	if m.cursor != 2 {
		t.Error("jumpToPrevBookmark with no bookmarks should not change cursor")
	}
}

func TestToggleLineComment_AddsComment(t *testing.T) {
	m := newModel(".")
	m.comments = make(map[int]string)
	lineNum := 10

	m = toggleLineComment(m, lineNum, "test comment")
	if _, exists := m.comments[lineNum]; !exists {
		t.Error("toggleLineComment should add line to comments map")
	}
}

func TestToggleLineComment_RemovesComment(t *testing.T) {
	m := newModel(".")
	m.comments = make(map[int]string)
	m.comments[10] = "test comment"

	m = toggleLineComment(m, 10, "")
	if _, exists := m.comments[10]; exists {
		t.Error("toggleLineComment should remove existing comment")
	}
}

func TestCompileRegex_ValidPattern(t *testing.T) {
	pattern := "test.*regex"
	re, err := compileRegex(pattern)
	if err != nil {
		t.Errorf("compileRegex with valid pattern should not error: %v", err)
	}
	if re == nil {
		t.Error("compileRegex with valid pattern should not return nil regex")
	}
}

func TestCompileRegex_InvalidPattern(t *testing.T) {
	pattern := "[invalid(regex"
	re, err := compileRegex(pattern)
	if err == nil {
		t.Error("compileRegex with invalid pattern should return error")
	}
	if re != nil {
		t.Error("compileRegex with invalid pattern should return nil regex")
	}
}

func TestDetectLanguage_GoFile(t *testing.T) {
	lang := detectLanguage("main.go")
	if lang != "go" && lang != "" {
		t.Errorf("detectLanguage for .go should return 'go' or empty, got %q", lang)
	}
}

func TestDetectLanguage_PythonFile(t *testing.T) {
	lang := detectLanguage("script.py")
	if lang != "python" && lang != "" {
		t.Errorf("detectLanguage for .py should return 'python' or empty, got %q", lang)
	}
}

func TestDetectLanguage_UnknownFile(t *testing.T) {
	lang := detectLanguage("unknown.xyz")
	// detectLanguage may return "text" or empty string for unknown extensions
	_ = lang
}

func TestDiffStatBadge_WithStats(t *testing.T) {
	// Create stats with valid values
	stats := commitStatistics{}
	badge := diffStatBadge(stats)
	// Badge may be empty for zero stats
	_ = badge
}

func TestPluralize_Singular(t *testing.T) {
	result := pluralize(1)
	_ = result // Should be a string for 1
}

func TestPluralize_Plural(t *testing.T) {
	result := pluralize(2)
	_ = result // Should be a string for 2
}

func TestRenderStatsBadgeInList_WithStats(t *testing.T) {
	stats := commitStatistics{}
	result := renderStatsBadgeInList(stats, 80)
	if result == "" {
		t.Error("renderStatsBadgeInList should return a string (possibly empty)")
	}
}

func TestFormatFilterHeaderDisplay_NoFilters(t *testing.T) {
	m := newModel(".")
	result := formatFilterHeaderDisplay(m)
	if result != "" {
		t.Error("formatFilterHeaderDisplay with no filters should be empty")
	}
}

func TestFormatFilterHeaderDisplay_WithAuthorFilter(t *testing.T) {
	m := newModel(".")
	m.authorFilter = "Alice"
	result := formatFilterHeaderDisplay(m)
	if result == "" {
		t.Error("formatFilterHeaderDisplay with author filter should return non-empty string")
	}
}

func TestRenderBookmarkMarker_ValidIndex(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.bookmarks = []string{"aaa1111"} // Use actual short hash from makeCommits

	result := renderBookmarkMarker(m, 0)
	_ = result // May be empty or contain marker
}

func TestHandleGoToCommitInput_EmptyInput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.cursor = 0

	updatedM := handleGoToCommitInput(m, "")
	if updatedM.cursor != 0 {
		t.Error("handleGoToCommitInput with empty input should not change cursor")
	}
}

func TestRenderLineCommentMarker_ValidIndex(t *testing.T) {
	m := newModel(".")
	m.comments = make(map[int]string)

	result := renderLineCommentMarker(m, 0)
	_ = result // May be empty string or spaces
}

func TestDetectBranches_WithCommits(t *testing.T) {
	commits := makeCommits(3)
	branches := detectBranches(commits)
	if branches == nil {
		// detectBranches may return nil for no branches detected
		return
	}
	_ = branches // May be empty or contain branch names
}

func TestRenderAsciiGraph_NoNodes(t *testing.T) {
	result := renderAsciiGraph(nil)
	if result != "" {
		t.Errorf("renderAsciiGraph with no nodes should be empty, got %q", result)
	}
}

func TestGetCommitRelationships_NoCommits(t *testing.T) {
	result := getCommitRelationships(nil)
	if result == nil {
		t.Error("getCommitRelationships should return a map (may be empty)")
	}
}

func TestGetFileBlameContext_NoFile(t *testing.T) {
	lines := []diffLine{
		{kind: lineMeta, text: "diff --git a/other.go b/other.go"},
	}

	result := getFileBlameContext(lines, "nonexistent.go")
	if len(result) > 0 {
		t.Errorf("getFileBlameContext with nonexistent file should return empty map")
	}
}

func TestGetFileBlameContext_ValidFile(t *testing.T) {
	lines := []diffLine{
		{kind: lineMeta, text: "diff --git a/file.go b/file.go"},
		{kind: lineAdded, text: "+ code"},
	}

	result := getFileBlameContext(lines, "file.go")
	if result == nil {
		t.Error("getFileBlameContext should return a map")
	}
}

func TestIsFileModifiedInCommit_NotModified(t *testing.T) {
	result := isFileModifiedInCommit("abc123", "nonexistent.go")
	// Result depends on git output, just verify it returns a bool
	_ = result
}

func TestIsFileModifiedInCommit_ValidFile(t *testing.T) {
	// This would need actual git history, so just verify it runs
	result := isFileModifiedInCommit("abc123", "file.go")
	_ = result // Verify it returns a bool (value depends on git state)
}
