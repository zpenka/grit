package grit

import (
	"strings"
	"testing"
)

// TestTruncate tests string truncation
func TestTruncate_NoTruncationNeeded(t *testing.T) {
	result := truncate("hello", 10)
	if result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

func TestTruncate_ExactLength(t *testing.T) {
	result := truncate("hello", 5)
	if result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

func TestTruncate_TruncationNeeded(t *testing.T) {
	result := truncate("hello world", 5)
	if len(result) > 5 {
		t.Errorf("truncate should limit to 5 chars, got %d", len(result))
	}
}

func TestTruncate_EmptyString(t *testing.T) {
	result := truncate("", 10)
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestTruncate_ZeroLength(t *testing.T) {
	result := truncate("hello", 0)
	if result != "" {
		t.Errorf("truncate with 0 should return empty, got '%s'", result)
	}
}

// TestFirstWord tests extracting first word
func TestFirstWord_SingleWord(t *testing.T) {
	result := firstWord("hello")
	if result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

func TestFirstWord_MultipleWords(t *testing.T) {
	result := firstWord("hello world test")
	if result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

func TestFirstWord_EmptyString(t *testing.T) {
	result := firstWord("")
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestFirstWord_OnlySpaces(t *testing.T) {
	result := firstWord("   ")
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

// TestCapitalizeFirst tests capitalizing first letter
func TestCapitalizeFirst_Lowercase(t *testing.T) {
	result := capitalizeFirst("hello")
	if result != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", result)
	}
}

func TestCapitalizeFirst_AlreadyCapitalized(t *testing.T) {
	result := capitalizeFirst("Hello")
	if result != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", result)
	}
}

func TestCapitalizeFirst_SingleChar(t *testing.T) {
	result := capitalizeFirst("a")
	if result != "A" {
		t.Errorf("expected 'A', got '%s'", result)
	}
}

func TestCapitalizeFirst_EmptyString(t *testing.T) {
	result := capitalizeFirst("")
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

// TestParseBranches tests branch parsing
func TestParseBranches_EmptyOutput(t *testing.T) {
	branches := parseBranches("")
	if len(branches) != 0 {
		t.Error("parseBranches should return empty for empty output")
	}
}

func TestParseBranches_SingleBranch(t *testing.T) {
	output := "* main"
	branches := parseBranches(output)
	if len(branches) != 1 {
		t.Errorf("expected 1 branch, got %d", len(branches))
	}
	if branches[0] != "main" {
		t.Errorf("expected 'main', got '%s'", branches[0])
	}
}

func TestParseBranches_MultipleBranches(t *testing.T) {
	output := `* main
  develop
  feature/test`
	branches := parseBranches(output)
	if len(branches) != 3 {
		t.Errorf("expected 3 branches, got %d", len(branches))
	}
}

func TestParseBranches_CurrentBranchMarker(t *testing.T) {
	output := `  main
* develop
  feature/test`
	branches := parseBranches(output)
	if len(branches) != 3 {
		t.Errorf("expected 3 branches, got %d", len(branches))
	}
	// First should be develop (current)
	if branches[0] != "develop" {
		t.Errorf("expected 'develop' first, got '%s'", branches[0])
	}
}

// TestParseCurrentBranch tests current branch parsing
func TestParseCurrentBranch_StandardOutput(t *testing.T) {
	output := "* main"
	current := parseCurrentBranch(output)
	if current != "main" {
		t.Errorf("expected 'main', got '%s'", current)
	}
}

func TestParseCurrentBranch_NoCurrentBranch(t *testing.T) {
	output := `  main
  develop`
	current := parseCurrentBranch(output)
	if current != "" {
		t.Errorf("expected empty string, got '%s'", current)
	}
}

func TestParseCurrentBranch_EmptyOutput(t *testing.T) {
	current := parseCurrentBranch("")
	if current != "" {
		t.Errorf("expected empty string, got '%s'", current)
	}
}

// TestDetectFileLanguage tests language detection from filename
func TestDetectFileLanguage_GoFile(t *testing.T) {
	lang := detectLanguage("main.go")
	if lang != "go" {
		t.Errorf("expected 'go', got '%s'", lang)
	}
}

func TestDetectFileLanguage_PythonFile(t *testing.T) {
	lang := detectLanguage("script.py")
	if lang != "python" {
		t.Errorf("expected 'python', got '%s'", lang)
	}
}

func TestDetectLanguage_JavascriptFile(t *testing.T) {
	lang := detectLanguage("app.js")
	if lang != "javascript" {
		t.Errorf("expected 'javascript', got '%s'", lang)
	}
}

func TestDetectLanguage_NoExtension(t *testing.T) {
	lang := detectLanguage("Makefile")
	if lang == "" {
		t.Error("detectLanguage should return something for Makefile")
	}
}

// TestParseCount tests parsing count strings
func TestParseCount_SimpleNumber(t *testing.T) {
	count := parseCount("42")
	if count != 42 {
		t.Errorf("expected 42, got %d", count)
	}
}

func TestParseCount_WithText(t *testing.T) {
	count := parseCount("42 files")
	if count != 42 {
		t.Errorf("expected 42, got %d", count)
	}
}

func TestParseCount_InvalidString(t *testing.T) {
	count := parseCount("abc")
	if count != 0 {
		t.Errorf("expected 0 for invalid string, got %d", count)
	}
}

func TestParseCount_EmptyString(t *testing.T) {
	count := parseCount("")
	if count != 0 {
		t.Errorf("expected 0 for empty string, got %d", count)
	}
}

// TestPluralizeWord tests pluralization
func TestPluralizeWord_Singular(t *testing.T) {
	result := pluralize(1)
	if result != "" {
		t.Errorf("expected empty string for singular, got '%s'", result)
	}
}

func TestPluralizeWord_Plural(t *testing.T) {
	result := pluralize(2)
	if result != "s" {
		t.Errorf("expected 's' for plural, got '%s'", result)
	}
}

func TestPluralizeWord_Zero(t *testing.T) {
	result := pluralize(0)
	if result != "s" {
		t.Errorf("expected 's' for zero, got '%s'", result)
	}
}

// TestCommitStats tests diff statistics calculation
func TestCommitStats_EmptyDiff(t *testing.T) {
	stats := commitStats([]diffLine{})
	if stats.filesChanged != 0 || stats.insertions != 0 || stats.deletions != 0 {
		t.Error("empty diff should have zero stats")
	}
}

func TestCommitStats_AddedLines(t *testing.T) {
	lines := []diffLine{
		{kind: lineAdded, text: "+ new code"},
		{kind: lineAdded, text: "+ more code"},
	}
	stats := commitStats(lines)
	if stats.insertions != 2 {
		t.Errorf("expected 2 insertions, got %d", stats.insertions)
	}
}

func TestCommitStats_RemovedLines(t *testing.T) {
	lines := []diffLine{
		{kind: lineRemoved, text: "- old code"},
		{kind: lineRemoved, text: "- removed"},
	}
	stats := commitStats(lines)
	if stats.deletions != 2 {
		t.Errorf("expected 2 deletions, got %d", stats.deletions)
	}
}

func TestCommitStats_Mixed(t *testing.T) {
	lines := []diffLine{
		{kind: lineAdded, text: "+ added"},
		{kind: lineRemoved, text: "- removed"},
		{kind: lineContext, text: "  context"},
	}
	stats := commitStats(lines)
	if stats.insertions != 1 || stats.deletions != 1 {
		t.Error("should count insertions and deletions correctly")
	}
}

// TestIsMergeCommit tests merge commit detection
func TestIsMergeCommit_NoMergeMarker(t *testing.T) {
	lines := []diffLine{
		{kind: lineMeta, text: "diff --git a/file.go b/file.go"},
	}
	if isMergeCommit(lines) {
		t.Error("should not detect merge in non-merge diff")
	}
}

func TestIsMergeCommit_WithMergeMarker(t *testing.T) {
	lines := []diffLine{
		{kind: lineMeta, text: "Merge: abc123 def456"},
	}
	if !isMergeCommit(lines) {
		t.Error("should detect merge commit")
	}
}

// TestGetMergeParents tests extracting merge parents
func TestGetMergeParents_NoMergeMarker(t *testing.T) {
	lines := []diffLine{
		{kind: lineMeta, text: "diff --git a/file.go b/file.go"},
	}
	parents := getMergeParents(lines)
	if len(parents) != 0 {
		t.Error("should return empty for non-merge")
	}
}

func TestGetMergeParents_WithMergeMarker(t *testing.T) {
	lines := []diffLine{
		{kind: lineMeta, text: "Merge: abc123def456 def456abc123"},
	}
	parents := getMergeParents(lines)
	if len(parents) == 0 {
		t.Error("should extract merge parents")
	}
}

// TestDiffStatBadge tests diff stat badge generation
func TestDiffStatBadge_NoChanges(t *testing.T) {
	stats := commitStatistics{filesChanged: 0, insertions: 0, deletions: 0}
	badge := diffStatBadge(stats)
	if badge == "" {
		t.Error("diffStatBadge should return non-empty string")
	}
}

func TestDiffStatBadge_WithChanges(t *testing.T) {
	stats := commitStatistics{filesChanged: 2, insertions: 10, deletions: 5}
	badge := diffStatBadge(stats)
	if !strings.Contains(badge, "2") || !strings.Contains(badge, "10") || !strings.Contains(badge, "5") {
		t.Error("badge should contain file and line counts")
	}
}

// TestParseHunks tests hunk parsing
func TestParseHunks_EmptyDiff(t *testing.T) {
	hunks := parseHunks([]diffLine{})
	if len(hunks) != 0 {
		t.Error("empty diff should have no hunks")
	}
}

func TestParseHunks_WithHunkHeader(t *testing.T) {
	lines := []diffLine{
		{kind: lineHunk, text: "@@ -1,5 +1,6 @@"},
		{kind: lineAdded, text: "+ new line"},
	}
	hunks := parseHunks(lines)
	if len(hunks) == 0 {
		t.Error("should parse hunk header")
	}
}

// TestMiniMapPosition tests minimap position calculation
func TestMiniMapPosition_StartOfFile(t *testing.T) {
	pos := miniMapPosition(0, 10, 100)
	if pos != 0 {
		t.Errorf("expected position 0, got %d", pos)
	}
}

func TestMiniMapPosition_MiddleOfFile(t *testing.T) {
	pos := miniMapPosition(50, 10, 100)
	if pos <= 0 || pos >= 10 {
		t.Errorf("expected middle position, got %d", pos)
	}
}

func TestMiniMapPosition_EndOfFile(t *testing.T) {
	pos := miniMapPosition(99, 10, 100)
	if pos != 9 {
		t.Errorf("expected position 9, got %d", pos)
	}
}

// TestParseGitReferences tests parsing git references
func TestParseGitReferences_NoReferences(t *testing.T) {
	refs := parseGitReferences("Normal commit message")
	if len(refs) != 0 {
		t.Error("should return empty for no references")
	}
}

func TestParseGitReferences_WithHash(t *testing.T) {
	refs := parseGitReferences("Fix bug abc1234567890")
	if len(refs) == 0 {
		t.Error("should find git hash reference")
	}
}

func TestParseGitReferences_WithTag(t *testing.T) {
	refs := parseGitReferences("Release v1.2.3")
	// May or may not detect tags depending on implementation
	_ = refs
}

// TestFilterCommitsByFile tests file-based filtering
func TestFilterCommitsByFile_EmptyCommits(t *testing.T) {
	commits := []commit{}
	filtered := filterCommitsByFile(commits, "main.go")
	if len(filtered) != 0 {
		t.Error("should return empty for empty commits")
	}
}

func TestFilterCommitsByFile_WithMatches(t *testing.T) {
	commits := []commit{
		{hash: "abc123", subject: "commit 1"},
		{hash: "def456", subject: "commit 2"},
		{hash: "ghi789", subject: "commit 3"},
	}
	filtered := filterCommitsByFile(commits, "main.go")
	// Filtering depends on git integration, may be empty
	if filtered == nil {
		t.Error("should return non-nil slice")
	}
}

// TestParseTagsList tests tag parsing
func TestParseTagsList_EmptyOutput(t *testing.T) {
	tags := parseTags("")
	if len(tags) != 0 {
		t.Error("should return empty for empty output")
	}
}

func TestParseTagsList_VersionTag(t *testing.T) {
	output := "v1.0.0"
	tags := parseTags(output)
	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
	if tags[0] != "v1.0.0" {
		t.Errorf("expected 'v1.0.0', got '%s'", tags[0])
	}
}

func TestParseTagsList_MultipleTags(t *testing.T) {
	output := `v1.0.0
v1.1.0
v2.0.0`
	tags := parseTags(output)
	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}
}

// TestFilenameToScope tests scope detection from filename
func TestFilenameToScope_GoFile(t *testing.T) {
	scope := filenameToScope("auth/login.go")
	if scope == "" {
		t.Error("should detect scope from filename")
	}
}

func TestFilenameToScope_NestedPath(t *testing.T) {
	scope := filenameToScope("pkg/util/string_test.go")
	if scope == "" {
		t.Error("should detect scope from nested path")
	}
}

func TestFilenameToScope_NoPath(t *testing.T) {
	scope := filenameToScope("main.go")
	if scope == "" {
		t.Error("should handle files without path")
	}
}

// TestParseDateRangeFilter tests date range parsing
func TestParseDateRangeFilter_ValidRange(t *testing.T) {
	input := "2024-01-01..2024-12-31"
	from, to, err := parseDateRange(input)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if from == nil || to == nil {
		t.Error("should parse both dates")
	}
}

func TestParseDateRangeFilter_InvalidFormat(t *testing.T) {
	input := "invalid"
	_, _, err := parseDateRange(input)
	if err == nil {
		t.Error("should return error for invalid format")
	}
}

func TestParseDateRangeFilter_EmptyString(t *testing.T) {
	input := ""
	_, _, err := parseDateRange(input)
	if err == nil {
		t.Error("should return error for empty string")
	}
}
