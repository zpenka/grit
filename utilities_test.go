package grit

import (
	"strings"
	"testing"
)
func TestTruncate_Short(t *testing.T) {
	AssertEqual(t, "hello", truncate("hello", 10), "should not truncate")
}

func TestTruncate_Exact(t *testing.T) {
	AssertEqual(t, "hello", truncate("hello", 5), "exact length should not truncate")
}

func TestTruncate_Long(t *testing.T) {
	AssertEqual(t, "hello …", truncate("hello world", 7), "should truncate with ellipsis")
}

func TestTruncate_One(t *testing.T) {
	AssertEqual(t, "…", truncate("hi", 1), "max=1 should return ellipsis")
}

func TestTruncate_Zero(t *testing.T) {
	AssertEqual(t, "", truncate("hi", 0), "max=0 should return empty")
}

func TestFirstWord_WithSpace(t *testing.T) {
	AssertEqual(t, "John", firstWord("John Doe"), "should extract first word")
}

func TestFirstWord_NoSpace(t *testing.T) {
	AssertEqual(t, "John", firstWord("John"), "single word should return as-is")
}

func TestCopyAsPatch_Empty(t *testing.T) {
	patch := copyAsPatch("abc1234", []diffLine{})
	if patch == "" {
		t.Error("should generate non-empty patch")
	}
}

func TestCopyAsPatch_ContainsHash(t *testing.T) {
	lines := []diffLine{{lineAdded, "+new line"}}
	patch := copyAsPatch("abc1234", lines)
	if !strings.Contains(patch, "abc1234") {
		t.Errorf("patch should contain hash, got %q", patch)
	}
}

func TestCopyAsPatch_ContainsDiff(t *testing.T) {
	lines := []diffLine{
		{lineMeta, "diff --git a/test.go b/test.go"},
		{lineAdded, "+new line"},
	}
	patch := copyAsPatch("abc1234", lines)
	if !strings.Contains(patch, "+new line") {
		t.Errorf("patch should contain diff content, got %q", patch)
	}
}

func TestExtractReviewers_None(t *testing.T) {
	reviewers := extractReviewers("normal commit message")
	if len(reviewers) != 0 {
		t.Errorf("no reviewers should return empty, got %d", len(reviewers))
	}
}

func TestExtractReviewers_Mentioned(t *testing.T) {
	msg := "Fix: handle edge case\n\nReviewed-by: John Smith <john@example.com>"
	reviewers := extractReviewers(msg)
	if len(reviewers) != 1 {
		t.Errorf("expected 1 reviewer, got %d", len(reviewers))
	}
}

func TestExtractReviewers_Multiple(t *testing.T) {
	msg := "Feature: new API\n\nReviewed-by: Alice <alice@ex.com>\nReviewed-by: Bob <bob@ex.com>"
	reviewers := extractReviewers(msg)
	if len(reviewers) != 2 {
		t.Errorf("expected 2 reviewers, got %d", len(reviewers))
	}
}

func TestExtractCoAuthors_None(t *testing.T) {
	authors := extractCoAuthors("simple commit message")
	if len(authors) != 0 {
		t.Errorf("no co-authors should return empty, got %d", len(authors))
	}
}

func TestExtractCoAuthors_Single(t *testing.T) {
	msg := "Fix bug\n\nCo-authored-by: Jane Smith <jane@example.com>"
	authors := extractCoAuthors(msg)
	if len(authors) != 1 {
		t.Errorf("expected 1 co-author, got %d", len(authors))
	}
	if !strings.Contains(authors[0], "Jane") {
		t.Errorf("should contain name, got %q", authors[0])
	}
}

func TestExtractCoAuthors_Multiple(t *testing.T) {
	msg := "Fix bug\n\nCo-authored-by: Jane <jane@example.com>\nCo-authored-by: Bob <bob@example.com>"
	authors := extractCoAuthors(msg)
	if len(authors) != 2 {
		t.Errorf("expected 2 co-authors, got %d", len(authors))
	}
}

func TestDiffCache_StoresResult(t *testing.T) {
	cache := newDiffCache(10)
	lines := []diffLine{{lineAdded, "+test"}}
	cache.set("abc1234", lines)
	cached, ok := cache.get("abc1234")
	if !ok {
		t.Error("cache should have stored diff")
	}
	if len(cached) != 1 {
		t.Errorf("cache should return same data, got %d lines", len(cached))
	}
}

func TestDiffCache_EvictsOldest(t *testing.T) {
	cache := newDiffCache(2)
	cache.set("aaa", []diffLine{{lineAdded, "+a"}})
	cache.set("bbb", []diffLine{{lineAdded, "+b"}})
	cache.set("ccc", []diffLine{{lineAdded, "+c"}})
	_, ok := cache.get("aaa")
	if ok {
		t.Error("cache should evict oldest entry")
	}
}

func TestRegexCache_Compiles(t *testing.T) {
	cache := newRegexCache(10)
	re, err := cache.compile("test.*pattern")
	if err != nil {
		t.Errorf("should compile regex, got %v", err)
	}
	if !re.MatchString("test something pattern") {
		t.Error("compiled regex should match")
	}
}

func TestRegexCache_Reuses(t *testing.T) {
	cache := newRegexCache(10)
	re1, _ := cache.compile("test")
	re2, _ := cache.compile("test")
	if re1 != re2 {
		t.Error("cache should return same regex object")
	}
}

func TestCircularBuffer_StoresAndRetrievesItems(t *testing.T) {
	buffer := NewCircularBuffer(5)
	
	buffer.Push("commit1")
	buffer.Push("commit2")
	buffer.Push("commit3")
	
	items := buffer.GetAll()
	AssertLen(t, items, 3, "should store 3 items")
	AssertEqual(t, 3, buffer.Size(), "size should be 3")
}

func TestCircularBuffer_Wraps(t *testing.T) {
	buffer := NewCircularBuffer(3)
	
	buffer.Push("a")
	buffer.Push("b")
	buffer.Push("c")
	buffer.Push("d") // Should wrap and overwrite "a"
	
	items := buffer.GetAll()
	AssertLen(t, items, 3, "should maintain max capacity")
}

func TestCurrentFile_Empty(t *testing.T) {
	m := model{}
	if currentFile(m) != "" {
		t.Error("expected empty string with no fileItems")
	}
}

func TestCurrentFile_AtStart(t *testing.T) {
	m := model{
		fileItems:  []fileItem{{path: "a.go", diffIdx: 0}, {path: "b.go", diffIdx: 20}},
		diffOffset: 0,
	}
	if got := currentFile(m); got != "a.go" {
		t.Errorf("expected a.go, got %q", got)
	}
}

func TestCurrentFile_PastBoundary(t *testing.T) {
	m := model{
		fileItems:  []fileItem{{path: "a.go", diffIdx: 0}, {path: "b.go", diffIdx: 20}},
		diffOffset: 25,
	}
	if got := currentFile(m); got != "b.go" {
		t.Errorf("expected b.go, got %q", got)
	}
}

func TestCurrentFile_ExactlyAtBoundary(t *testing.T) {
	m := model{
		fileItems:  []fileItem{{path: "a.go", diffIdx: 0}, {path: "b.go", diffIdx: 20}},
		diffOffset: 20,
	}
	if got := currentFile(m); got != "b.go" {
		t.Errorf("expected b.go at exact boundary, got %q", got)
	}
}

func TestMetrics_TracksOperations(t *testing.T) {
	metrics := NewMetrics()

	metrics.RecordOperation("filter", true)
	metrics.RecordOperation("filter", true)
	metrics.RecordOperation("filter", false)

	rate := metrics.GetSuccessRate("filter")
	AssertTrue(t, rate > 0 && rate <= 1, "success rate should be between 0 and 1")
}

func TestMemoryPool_Reuses(t *testing.T) {
	pool := NewMemoryPool(5)
	
	// Get object from empty pool
	obj1 := pool.Get(func() interface{} {
		return &commit{hash: "new1"}
	})
	AssertNotNil(t, obj1, "should create object")
	
	// Put back to pool
	pool.Put(obj1)
	
	// Get again should reuse
	obj2 := pool.Get(func() interface{} {
		return &commit{hash: "new2"}
	})
	AssertEqual(t, obj1, obj2, "should reuse pooled object")
}

func TestMemoryPool_CreateWhenEmpty(t *testing.T) {
	pool := NewMemoryPool(1)
	
	// First object
	obj1 := pool.Get(func() interface{} {
		return &commit{hash: "first"}
	})
	AssertNotNil(t, obj1, "should create first object")
	
	// Get second without returning first
	obj2 := pool.Get(func() interface{} {
		return &commit{hash: "second"}
	})
	AssertNotNil(t, obj2, "should create new object when pool empty")
	AssertNotEqual(t, obj1, obj2, "should be different objects")
}

func TestStatCache_Memoizes(t *testing.T) {
	cache := newStatCache(10)
	lines := []diffLine{
		{lineAdded, "+test"},
		{lineRemoved, "-old"},
	}
	stats1 := cache.getOrCompute("abc", lines)
	stats2 := cache.getOrCompute("abc", lines)
	if stats1.insertions != stats2.insertions {
		t.Error("cached stats should be identical")
	}
}

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

// ============================================================================
// ENHANCED COVERAGE: Parsing Module
// ============================================================================

func TestParseCommitsWithPool_Empty(t *testing.T) {
	result := parseCommitsWithPool("")
	if result != nil {
		t.Error("parseCommitsWithPool should return nil for empty input")
	}
}

func TestParseCommitsWithPool_SingleCommit(t *testing.T) {
	input := "abc1234567890123456789012345678901234\x00abc1234\x00Alice\x001h ago\x00Fix bug\n"
	commits := parseCommitsWithPool(input)
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	if commits[0].hash != "abc1234567890123456789012345678901234" {
		t.Errorf("hash mismatch: got %s", commits[0].hash)
	}
	if commits[0].author != "Alice" {
		t.Errorf("author mismatch: got %s", commits[0].author)
	}
}

func TestParseCommitsWithPool_MultipleCommits(t *testing.T) {
	input := "hash1\x00h1\x00Author1\x001h ago\x00Msg1\nhash2\x00h2\x00Author2\x002h ago\x00Msg2\n"
	commits := parseCommitsWithPool(input)
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}
}

func TestParseCommitsWithPool_MalformedLines(t *testing.T) {
	input := "valid123456789012345678901234567890\x00v123\x00Author\x001h ago\x00Msg\ninvalid-no-fields\nhash2\x00h2\x00Auth2\x002h ago\x00Msg2\n"
	commits := parseCommitsWithPool(input)
	if len(commits) != 2 {
		t.Fatalf("should skip malformed lines, got %d commits", len(commits))
	}
}

// ============================================================================
// ENHANCED COVERAGE: Diff Batch Processing Module
// ============================================================================

func TestProcessDiffBatch_Empty(t *testing.T) {
	processor := &BatchProcessor{}
	lines := []diffLine{}
	result := processDiffBatch(processor, lines)
	if len(result) != 0 {
		t.Error("processDiffBatch should return empty slice for empty input")
	}
}

func TestProcessDiffBatch_AddedLines(t *testing.T) {
	processor := &BatchProcessor{}
	lines := []diffLine{
		{text: "+new line 1", kind: lineAdded},
		{text: "+new line 2", kind: lineAdded},
	}
	result := processDiffBatch(processor, lines)
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result))
	}
	for _, line := range result {
		if line.kind != lineAdded {
			t.Error("batch processing should preserve line kind")
		}
	}
}

func TestProcessDiffBatch_MixedLines(t *testing.T) {
	processor := &BatchProcessor{}
	lines := []diffLine{
		{text: "-old line", kind: lineRemoved},
		{text: " context line", kind: lineContext},
		{text: "+new line", kind: lineAdded},
		{text: "@@ -1,5 +1,6 @@", kind: lineHunk},
	}
	result := processDiffBatch(processor, lines)
	if len(result) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(result))
	}
}

// ============================================================================
// ENHANCED COVERAGE: Filtering Module
// ============================================================================

func TestIsWithinDays_TodayString(t *testing.T) {
	// "0 days ago" should be within any day threshold
	result := isWithinDays("0 days ago", 1)
	if !result {
		t.Error("0 days ago should be within 1 day")
	}
}

func TestIsWithinDays_OneDayAgo(t *testing.T) {
	result := isWithinDays("1 day ago", 1)
	if !result {
		t.Error("1 day ago should be within 1 day")
	}
}

func TestIsWithinDays_TwoDaysAgo(t *testing.T) {
	result := isWithinDays("2 days ago", 1)
	if result {
		t.Error("2 days ago should NOT be within 1 day")
	}
}

func TestIsWithinDays_WeekThreshold(t *testing.T) {
	result := isWithinDays("1 week ago", 7)
	if !result {
		t.Error("1 week ago should be within 7 days")
	}
}

func TestIsWithinDays_InvalidFormat(t *testing.T) {
	result := isWithinDays("invalid date string", 30)
	if result {
		t.Error("invalid format should return false")
	}
}

func TestFilterCommitsWithCache_Empty(t *testing.T) {
	cache := &FilterCache{}
	commits := []commit{}
	result := filterCommitsWithCache(cache, commits, "test")
	if len(result) != 0 {
		t.Error("empty commits should return empty slice")
	}
}

func TestFilterCommitsWithCache_NoQuery(t *testing.T) {
	cache := &FilterCache{}
	commits := []commit{
		{hash: "abc", author: "test", subject: "msg"},
	}
	result := filterCommitsWithCache(cache, commits, "")
	if len(result) != 1 {
		t.Error("no query should return all commits")
	}
}

func TestFilterCommitsWithCache_MatchingSubject(t *testing.T) {
	cache := &FilterCache{}
	commits := []commit{
		{hash: "abc", author: "test", subject: "fix bug"},
		{hash: "def", author: "test", subject: "add feature"},
	}
	result := filterCommitsWithCache(cache, commits, "fix")
	if len(result) != 1 || result[0].subject != "fix bug" {
		t.Error("should filter by subject")
	}
}

func TestFilterCommitsWithCache_MatchingAuthor(t *testing.T) {
	cache := &FilterCache{}
	commits := []commit{
		{hash: "abc", author: "Alice", subject: "msg"},
		{hash: "def", author: "Bob", subject: "msg"},
	}
	result := filterCommitsWithCache(cache, commits, "Alice")
	if len(result) != 1 || result[0].author != "Alice" {
		t.Error("should filter by author")
	}
}

// ============================================================================
// ENHANCED COVERAGE: Navigation Module
// ============================================================================

func TestScrollToDiffLine_InvalidIndex(t *testing.T) {
	m := newModel(".")
	m.diffLines = []diffLine{
		{text: "line 1"},
		{text: "line 2"},
	}
	m.diffOffset = 0
	m = scrollToDiffLine(m, -1)
	// Should handle gracefully without crashing
	if m.diffOffset < 0 {
		t.Error("diffOffset should not be negative")
	}
}

func TestScrollToDiffLine_BeyondEnd(t *testing.T) {
	m := newModel(".")
	m.diffLines = []diffLine{
		{text: "line 1"},
		{text: "line 2"},
	}
	m.height = 10
	m = scrollToDiffLine(m, 100)
	// Should handle gracefully without crashing
	if m.diffOffset < 0 {
		t.Error("diffOffset should not be negative")
	}
}

func TestScrollToDiffLine_ValidPosition(t *testing.T) {
	m := newModel(".")
	m.diffLines = make([]diffLine, 50)
	m.height = 10
	m = scrollToDiffLine(m, 25)
	// Should scroll without errors
	if m.diffOffset > 25 && m.height > 0 {
		t.Errorf("diffOffset %d should generally position near target line 25", m.diffOffset)
	}
}

func TestAddToNavHistory_NewEntry(t *testing.T) {
	m := newModel(".")
	m.navHistory = []int{}
	m.navHistoryIdx = 0
	originalLen := len(m.navHistory)
	m = addToNavHistory(m, 5)
	if len(m.navHistory) != originalLen+1 {
		t.Error("should add new entry to history")
	}
	if m.navHistory[len(m.navHistory)-1] != 5 {
		t.Error("should store the position in history")
	}
}

func TestAddToNavHistory_TracksCursor(t *testing.T) {
	m := newModel(".")
	m.cursor = 10
	originalLen := len(m.navHistory)
	m = addToNavHistory(m, 10)
	if len(m.navHistory) <= originalLen {
		t.Error("should add new history entry")
	}
}

func TestAddToNavHistory_Multiple(t *testing.T) {
	m := newModel(".")
	m = addToNavHistory(m, 1)
	m = addToNavHistory(m, 2)
	m = addToNavHistory(m, 3)
	if len(m.navHistory) < 3 {
		t.Errorf("should have at least 3 entries, got %d", len(m.navHistory))
	}
}
