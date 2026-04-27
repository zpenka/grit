package grit

import (
	"strings"
	"testing"
)

// TestParseStashList tests stash parsing from git output
func TestParseStashList_EmptyOutput(t *testing.T) {
	output := ""
	stashes := parseStashList(output)
	if len(stashes) != 0 {
		t.Error("parseStashList should return empty slice for empty output")
	}
}

func TestParseStashList_SingleStash(t *testing.T) {
	output := "stash@{0}: WIP on main: abc1234 work in progress"
	stashes := parseStashList(output)
	if len(stashes) != 1 {
		t.Errorf("expected 1 stash, got %d", len(stashes))
	}
	if stashes[0].name != "stash@{0}" {
		t.Errorf("expected stash name 'stash@{0}', got '%s'", stashes[0].name)
	}
	if stashes[0].branch != "main" {
		t.Errorf("expected branch 'main', got '%s'", stashes[0].branch)
	}
}

func TestParseStashList_MultipleStashes(t *testing.T) {
	output := `stash@{0}: WIP on main: abc1234 work1
stash@{1}: WIP on develop: def5678 work2
stash@{2}: WIP on feature: ghi9012 work3`
	stashes := parseStashList(output)
	if len(stashes) != 3 {
		t.Errorf("expected 3 stashes, got %d", len(stashes))
	}
	if stashes[0].name != "stash@{0}" || stashes[1].name != "stash@{1}" || stashes[2].name != "stash@{2}" {
		t.Error("stash names don't match expected pattern")
	}
}

func TestParseStashList_WithWhitespace(t *testing.T) {
	output := `  stash@{0}: WIP on main: abc1234 work

  stash@{1}: WIP on develop: def5678 work2  `
	stashes := parseStashList(output)
	if len(stashes) != 2 {
		t.Errorf("expected 2 stashes, got %d", len(stashes))
	}
}

// TestParseReflog tests reflog parsing from git output
func TestParseReflog_EmptyOutput(t *testing.T) {
	output := ""
	entries := parseReflog(output)
	if len(entries) != 0 {
		t.Error("parseReflog should return empty slice for empty output")
	}
}

func TestParseReflog_SingleEntry(t *testing.T) {
	output := "abc1234 HEAD@{0}: commit: Add new feature"
	entries := parseReflog(output)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].hash != "abc1234" {
		t.Errorf("expected hash 'abc1234', got '%s'", entries[0].hash)
	}
	if entries[0].action != "commit" {
		t.Errorf("expected action 'commit', got '%s'", entries[0].action)
	}
	if entries[0].message != "Add new feature" {
		t.Errorf("expected message 'Add new feature', got '%s'", entries[0].message)
	}
}

func TestParseReflog_MultipleEntries(t *testing.T) {
	output := `abc1234 HEAD@{0}: commit: Add feature
def5678 HEAD@{1}: commit: Fix bug
ghi9012 HEAD@{2}: checkout: moving to main`
	entries := parseReflog(output)
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestParseReflog_WithWhitespace(t *testing.T) {
	output := `  abc1234 HEAD@{0}: commit: message

  def5678 HEAD@{1}: rebase: continue  `
	entries := parseReflog(output)
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

// TestRenderStashView tests stash view rendering
func TestRenderStashView_HeadingAndContent(t *testing.T) {
	stashes := []stashEntry{
		{name: "stash@{0}", branch: "main", subject: "WIP on main"},
	}
	output := renderStashView(stashes, 80)
	if !strings.Contains(output, "stash@{0}") {
		t.Error("renderStashView should contain stash name")
	}
	if !strings.Contains(output, "main") {
		t.Error("renderStashView should contain branch name")
	}
}

// TestRenderReflogView tests reflog view rendering
func TestRenderReflogView_HeadingAndContent(t *testing.T) {
	entries := []reflogEntry{
		{hash: "abc1234567890", action: "commit", message: "Add feature"},
	}
	output := renderReflogView(entries, 80)
	if !strings.Contains(output, "abc1234") {
		t.Error("renderReflogView should contain commit hash")
	}
	if !strings.Contains(output, "Add feature") {
		t.Error("renderReflogView should contain message")
	}
}

// TestParseCommitGraph tests graph parsing
func TestParseCommitGraph_EmptyList(t *testing.T) {
	commits := []commit{}
	graph := parseCommitGraph(commits)
	if len(graph) != 0 {
		t.Error("parseCommitGraph should return empty graph for empty commits")
	}
}

func TestParseCommitGraph_SingleCommit(t *testing.T) {
	commits := []commit{
		{hash: "abc1234", subject: "Initial commit"},
	}
	graph := parseCommitGraph(commits)
	if len(graph) != 1 {
		t.Errorf("expected 1 node, got %d", len(graph))
	}
	if graph[0].hash != "abc1234" {
		t.Errorf("expected hash 'abc1234', got '%s'", graph[0].hash)
	}
}

func TestParseCommitGraph_MergeDetection(t *testing.T) {
	commits := []commit{
		{hash: "abc1234", subject: "Merge branch develop"},
	}
	graph := parseCommitGraph(commits)
	if len(graph) != 1 || !graph[0].isMerge {
		t.Error("parseCommitGraph should detect merge commits")
	}
}

func TestParseCommitGraph_DepthAssignment(t *testing.T) {
	commits := []commit{
		{hash: "aaa", subject: "First"},
		{hash: "bbb", subject: "Second"},
		{hash: "ccc", subject: "Third"},
	}
	graph := parseCommitGraph(commits)
	if len(graph) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(graph))
	}
	// Verify depth is assigned (pattern: 0, 1, 0)
	for i, node := range graph {
		expectedDepth := i % 2
		if node.depth != expectedDepth {
			t.Errorf("node %d: expected depth %d, got %d", i, expectedDepth, node.depth)
		}
	}
}

// TestDetectBranches tests branch detection
func TestDetectBranches_NoMerges(t *testing.T) {
	commits := []commit{
		{hash: "aaa", subject: "First commit"},
		{hash: "bbb", subject: "Second commit"},
	}
	branches := detectBranches(commits)
	if len(branches) != 1 || branches[0] != "main" {
		t.Error("detectBranches should return only 'main' for non-merge history")
	}
}

func TestDetectBranches_WithMerge(t *testing.T) {
	commits := []commit{
		{hash: "aaa", subject: "First commit"},
		{hash: "bbb", subject: "Merge branch feature"},
	}
	branches := detectBranches(commits)
	if len(branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(branches))
	}
	found := false
	for _, b := range branches {
		if b == "main" && found {
			found = true
		}
	}
}

// TestNavigateAlongGraph tests graph navigation
func TestNavigateAlongGraph_EmptyGraph(t *testing.T) {
	graph := []graphNode{}
	result := navigateAlongGraph(graph, 0, "down")
	if result != 0 {
		t.Error("navigateAlongGraph should return 0 for empty graph")
	}
}

func TestNavigateAlongGraph_MoveDown(t *testing.T) {
	graph := []graphNode{
		{hash: "aaa"},
		{hash: "bbb"},
		{hash: "ccc"},
	}
	result := navigateAlongGraph(graph, 0, "down")
	if result != 1 {
		t.Errorf("navigateAlongGraph down should return 1, got %d", result)
	}
}

func TestNavigateAlongGraph_MoveUp(t *testing.T) {
	graph := []graphNode{
		{hash: "aaa"},
		{hash: "bbb"},
		{hash: "ccc"},
	}
	result := navigateAlongGraph(graph, 1, "up")
	if result != 0 {
		t.Errorf("navigateAlongGraph up should return 0, got %d", result)
	}
}

func TestNavigateAlongGraph_BoundaryDown(t *testing.T) {
	graph := []graphNode{
		{hash: "aaa"},
		{hash: "bbb"},
	}
	result := navigateAlongGraph(graph, 1, "down")
	if result != 1 {
		t.Errorf("navigateAlongGraph should not go past end, got %d", result)
	}
}

func TestNavigateAlongGraph_BoundaryUp(t *testing.T) {
	graph := []graphNode{
		{hash: "aaa"},
		{hash: "bbb"},
	}
	result := navigateAlongGraph(graph, 0, "up")
	if result != 0 {
		t.Errorf("navigateAlongGraph should not go before start, got %d", result)
	}
}

// TestRenderAsciiGraph tests ASCII graph rendering
func TestRenderAsciiGraph_EmptyGraph(t *testing.T) {
	graph := []graphNode{}
	output := renderAsciiGraph(graph)
	if output != "" {
		t.Error("renderAsciiGraph should return empty string for empty graph")
	}
}

func TestRenderAsciiGraph_SingleNode(t *testing.T) {
	graph := []graphNode{
		{hash: "abc1234567890", isMerge: false},
	}
	output := renderAsciiGraph(graph)
	if !strings.Contains(output, "abc1234") {
		t.Error("renderAsciiGraph should contain commit hash")
	}
}

func TestRenderAsciiGraph_MergeNode(t *testing.T) {
	graph := []graphNode{
		{hash: "abc1234567890", isMerge: true},
	}
	output := renderAsciiGraph(graph)
	if !strings.Contains(output, "*  ") && !strings.Contains(output, "*   ") {
		t.Error("renderAsciiGraph should show merge prefix for merge commits")
	}
}

// TestFilterCommitsByFileChange tests file-based filtering
func TestFilterCommitsByFileChange_EmptyFile(t *testing.T) {
	commits := makeCommits(3)
	result := filterCommitsByFileChange(commits, "")
	if len(result) != len(commits) {
		t.Error("filterCommitsByFileChange should return all commits for empty file")
	}
}

func TestFilterCommitsByFileChange_NoMatches(t *testing.T) {
	commits := makeCommits(3)
	result := filterCommitsByFileChange(commits, "nonexistent.go")
	if len(result) != 0 {
		t.Errorf("expected 0 commits, got %d", len(result))
	}
}

// TestBuildFileHistory tests file history construction
func TestBuildFileHistory_EmptyFile(t *testing.T) {
	commits := makeCommits(3)
	history := buildFileHistory(commits, "")
	if len(history) != 0 {
		t.Error("buildFileHistory should return empty for empty file")
	}
}

func TestBuildFileHistory_WithFile(t *testing.T) {
	commits := makeCommits(3)
	history := buildFileHistory(commits, "main.go")
	// Expected to return empty since implementation doesn't query git
	if len(history) != 0 {
		t.Error("buildFileHistory should return empty (infrastructure not implemented)")
	}
}

// TestRenderFileTimeline tests file timeline rendering
func TestRenderFileTimeline_EmptyAndContainsFile(t *testing.T) {
	commits := []commit{
		{shortHash: "abc", subject: "Add feature"},
		{shortHash: "def", subject: "Fix bug"},
	}
	output := renderFileTimeline(commits, "main.go", 80)
	if !strings.Contains(output, "File Timeline") {
		t.Error("renderFileTimeline should contain title")
	}
	if !strings.Contains(output, "main.go") {
		t.Error("renderFileTimeline should contain filename")
	}
}

// TestStashToCommitLike tests stash conversion
func TestStashToCommitLike_Conversion(t *testing.T) {
	stash := stashEntry{
		hash:    "abc123",
		name:    "stash@{0}",
		branch:  "main",
		subject: "WIP on main",
	}
	c := stashToCommitLike(stash)
	if c.author != "stash" {
		t.Errorf("stash commit should have author 'stash', got '%s'", c.author)
	}
	if c.subject == "" {
		t.Error("stash commit should have non-empty subject")
	}
}

// TestReflogToCommitLike tests reflog conversion
func TestReflogToCommitLike_Conversion(t *testing.T) {
	entry := reflogEntry{
		hash:    "abc1234567890abc",
		action:  "commit",
		message: "Add feature",
	}
	c := reflogToCommitLike(entry)
	if c.author != "commit" {
		t.Errorf("reflog commit author should equal action 'commit', got '%s'", c.author)
	}
	if c.subject != "Add feature" {
		t.Errorf("reflog commit subject should be 'Add feature', got '%s'", c.subject)
	}
}

// TestFindStashByIndex tests stash lookup
func TestFindStashByIndex_NotFound(t *testing.T) {
	stashes := []stashEntry{
		{name: "stash@{0}"},
		{name: "stash@{1}"},
	}
	result := findStashByIndex(stashes, 5)
	if result != nil {
		t.Error("findStashByIndex should return nil for out-of-bounds index")
	}
}

func TestFindStashByIndex_Found(t *testing.T) {
	stashes := []stashEntry{
		{name: "stash@{0}", hash: "abc"},
		{name: "stash@{1}", hash: "def"},
	}
	result := findStashByIndex(stashes, 1)
	if result == nil {
		t.Error("findStashByIndex should find valid stash")
	}
	if result.hash != "def" {
		t.Errorf("expected hash 'def', got '%s'", result.hash)
	}
}

// TestSwitchViewMode tests view mode switching
func TestSwitchViewMode_Basic(t *testing.T) {
	m := newModel(".")
	result := switchViewMode(m, "reflog")
	if result.viewMode == "" {
		t.Error("switchViewMode should set view mode")
	}
}
