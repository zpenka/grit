package grit

import (
	"strings"
	"testing"
)

// TestRenderLostCommitsUI tests renderLostCommitsUI with different states
func TestRenderLostCommitsUI_EmptyList(t *testing.T) {
	m := model{
		lostCommits: []lostCommit{},
	}
	output := renderLostCommitsUI(m, 80)
	if output == "" {
		t.Error("renderLostCommitsUI should produce output even for empty list")
	}
	if !strings.Contains(output, "Lost Commits") {
		t.Error("output should contain 'Lost Commits' title")
	}
}

func TestRenderLostCommitsUI_WithCommits(t *testing.T) {
	m := model{
		lostCommits: []lostCommit{
			{hash: "aaa111", shortHash: "aaa111", author: "Alice", subject: "Lost feature"},
			{hash: "bbb222", shortHash: "bbb222", author: "Bob", subject: "Lost bugfix"},
		},
	}
	output := renderLostCommitsUI(m, 80)
	if output == "" {
		t.Error("renderLostCommitsUI should produce output for non-empty list")
	}
	if !strings.Contains(output, "aaa111") {
		t.Error("output should contain first lost commit hash")
	}
	if !strings.Contains(output, "bbb222") {
		t.Error("output should contain second lost commit hash")
	}
}

// TestRenderUndoMenu tests renderUndoMenu with various stack states
func TestRenderUndoMenu_EmptyStack(t *testing.T) {
	m := model{
		undoStack:    []string{},
		undoStackIdx: 0,
	}
	output := renderUndoMenu(m, 80)
	if output == "" {
		t.Error("renderUndoMenu should produce output")
	}
	if !strings.Contains(output, "Undo Stack") {
		t.Error("output should contain 'Undo Stack' title")
	}
}

func TestRenderUndoMenu_SingleItem(t *testing.T) {
	m := model{
		undoStack:    []string{"aaa111"},
		undoStackIdx: 1,
	}
	output := renderUndoMenu(m, 80)
	if !strings.Contains(output, "aaa111") {
		t.Error("output should contain undo stack item")
	}
}

func TestRenderUndoMenu_MultipleItems(t *testing.T) {
	m := model{
		undoStack:    []string{"aaa111", "bbb222", "ccc333"},
		undoStackIdx: 2,
	}
	output := renderUndoMenu(m, 80)
	if !strings.Contains(output, "aaa111") || !strings.Contains(output, "bbb222") || !strings.Contains(output, "ccc333") {
		t.Error("output should contain all undo stack items")
	}
}

// TestAnalyzeComplexity tests complexity analysis with different commit sets
func TestAnalyzeComplexity_EmptyCommits(t *testing.T) {
	m := model{
		commits: []commit{},
	}
	result := analyzeComplexity(m)
	if len(result.commitMetrics) != 0 {
		t.Error("analyzeComplexity should return empty metrics for empty commits")
	}
}

func TestAnalyzeComplexity_SingleSimpleCommit(t *testing.T) {
	m := model{
		commits: []commit{
			{hash: "aaa111", shortHash: "aaa111", subject: "Fix bug"},
		},
	}
	result := analyzeComplexity(m)
	if len(result.commitMetrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(result.commitMetrics))
	}
	metrics := result.commitMetrics[0]
	if metrics.hash != "aaa111" {
		t.Error("metric hash should match commit")
	}
	if metrics.complexity == 0 {
		t.Error("complexity should be calculated")
	}
}

func TestAnalyzeComplexity_LargeCommit(t *testing.T) {
	// Large commit message with many words
	m := model{
		commits: []commit{
			{hash: "aaa111", shortHash: "aaa111", subject: "Refactor module with new feature and update tests and fix bugs in processing"},
		},
	}
	result := analyzeComplexity(m)
	if len(result.commitMetrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(result.commitMetrics))
	}
	metrics := result.commitMetrics[0]
	if !metrics.isComplex {
		t.Error("large commit should be marked as complex")
	}
}

func TestAnalyzeComplexity_MultipleCommits(t *testing.T) {
	m := model{
		commits: []commit{
			{hash: "aaa111", shortHash: "aaa111", subject: "Fix"},
			{hash: "bbb222", shortHash: "bbb222", subject: "Add feature with tests"},
			{hash: "ccc333", shortHash: "ccc333", subject: "Large refactor touching many files"},
		},
	}
	result := analyzeComplexity(m)
	if len(result.commitMetrics) != 3 {
		t.Errorf("expected 3 metrics, got %d", len(result.commitMetrics))
	}
}

// TestCalculateAuthorStats tests author statistics calculation
func TestCalculateAuthorStats_EmptyCommits(t *testing.T) {
	commits := []commit{}
	stats := calculateAuthorStats(commits)
	if len(stats) != 0 {
		t.Error("calculateAuthorStats should return empty map for empty commits")
	}
}

func TestCalculateAuthorStats_SingleAuthor(t *testing.T) {
	commits := []commit{
		{author: "Alice"},
		{author: "Alice"},
		{author: "Alice"},
	}
	stats := calculateAuthorStats(commits)
	if len(stats) != 1 {
		t.Errorf("expected 1 author, got %d", len(stats))
	}
	if stats["Alice"] != 3 {
		t.Errorf("expected 3 commits for Alice, got %d", stats["Alice"])
	}
}

func TestCalculateAuthorStats_MultipleAuthors(t *testing.T) {
	commits := []commit{
		{author: "Alice"},
		{author: "Bob"},
		{author: "Alice"},
		{author: "Charlie"},
		{author: "Bob"},
	}
	stats := calculateAuthorStats(commits)
	if len(stats) != 3 {
		t.Errorf("expected 3 authors, got %d", len(stats))
	}
	if stats["Alice"] != 2 || stats["Bob"] != 2 || stats["Charlie"] != 1 {
		t.Error("author counts don't match expected values")
	}
}

// TestCalculateTimeStats tests time-based statistics
func TestCalculateTimeStats_EmptyCommits(t *testing.T) {
	commits := []commit{}
	stats := calculateTimeStats(commits)
	if len(stats) != 0 {
		t.Error("calculateTimeStats should return empty map for empty commits")
	}
}

func TestCalculateTimeStats_RecentCommits(t *testing.T) {
	commits := []commit{
		{when: "1 day ago"},
		{when: "2 days ago"},
	}
	stats := calculateTimeStats(commits)
	if stats["recent"] != 2 {
		t.Errorf("expected 2 recent commits, got %d", stats["recent"])
	}
}

func TestCalculateTimeStats_WeekAgoCommits(t *testing.T) {
	commits := []commit{
		{when: "1 week ago"},
		{when: "2 weeks ago"},
	}
	stats := calculateTimeStats(commits)
	if stats["past_week"] != 2 {
		t.Errorf("expected 2 past_week commits, got %d", stats["past_week"])
	}
}

func TestCalculateTimeStats_OldCommits(t *testing.T) {
	commits := []commit{
		{when: "Jan 2024"},
		{when: "Dec 2023"},
	}
	stats := calculateTimeStats(commits)
	if stats["older"] != 2 {
		t.Errorf("expected 2 older commits, got %d", stats["older"])
	}
}

// TestAggregateByWeek tests weekly aggregation
func TestAggregateByWeek_EmptyCommits(t *testing.T) {
	commits := []commit{}
	weekly := aggregateByWeek(commits)
	if len(weekly) != 0 {
		t.Error("aggregateByWeek should return empty map for empty commits")
	}
}

func TestAggregateByWeek_WithAgoSuffix(t *testing.T) {
	commits := []commit{
		{when: "1 day ago"},
		{when: "3 days ago"},
	}
	weekly := aggregateByWeek(commits)
	if weekly["current"] != 2 {
		t.Errorf("expected 2 current week commits, got %d", weekly["current"])
	}
}

// TestExtractCoAuthors tests co-author extraction
func TestExtractCoAuthors_NoCoAuthors(t *testing.T) {
	message := "Fix critical bug"
	coAuthors := extractCoAuthors(message)
	if len(coAuthors) != 0 {
		t.Error("extractCoAuthors should return empty for no co-authors")
	}
}

func TestExtractCoAuthors_SingleCoAuthor(t *testing.T) {
	message := "Fix bug\n\nCo-authored-by: Jane Doe <jane@example.com>"
	coAuthors := extractCoAuthors(message)
	if len(coAuthors) != 1 {
		t.Errorf("expected 1 co-author, got %d", len(coAuthors))
	}
	if coAuthors[0] != "Jane Doe" {
		t.Errorf("expected 'Jane Doe', got '%s'", coAuthors[0])
	}
}

func TestExtractCoAuthors_MultipleCoAuthors(t *testing.T) {
	message := `Fix bug

Co-authored-by: Jane Doe <jane@example.com>
Co-authored-by: John Smith <john@example.com>`
	coAuthors := extractCoAuthors(message)
	if len(coAuthors) != 2 {
		t.Errorf("expected 2 co-authors, got %d", len(coAuthors))
	}
}

// TestExtractReviewers tests reviewer extraction
func TestExtractReviewers_NoReviewers(t *testing.T) {
	message := "Fix critical bug"
	reviewers := extractReviewers(message)
	if len(reviewers) != 0 {
		t.Error("extractReviewers should return empty for no reviewers")
	}
}

func TestExtractReviewers_SingleReviewer(t *testing.T) {
	message := "Fix bug\n\nReviewed-by: Bob Johnson <bob@example.com>"
	reviewers := extractReviewers(message)
	if len(reviewers) != 1 {
		t.Errorf("expected 1 reviewer, got %d", len(reviewers))
	}
	if reviewers[0] != "Bob Johnson" {
		t.Errorf("expected 'Bob Johnson', got '%s'", reviewers[0])
	}
}

func TestExtractReviewers_MultipleReviewers(t *testing.T) {
	message := `Fix bug

Reviewed-by: Alice Smith <alice@example.com>
Reviewed-by: Bob Jones <bob@example.com>`
	reviewers := extractReviewers(message)
	if len(reviewers) != 2 {
		t.Errorf("expected 2 reviewers, got %d", len(reviewers))
	}
}

// TestCalculateProductivity tests productivity metrics
func TestCalculateProductivity_EmptyCommits(t *testing.T) {
	commits := []commit{}
	metrics := calculateProductivity(commits)
	if len(metrics) != 0 {
		t.Error("calculateProductivity should return empty for empty commits")
	}
}

func TestCalculateProductivity_SingleCommit(t *testing.T) {
	commits := []commit{
		{author: "Alice"},
	}
	metrics := calculateProductivity(commits)
	if metrics["commits"] != 1 {
		t.Errorf("expected 1 commit, got %v", metrics["commits"])
	}
	if metrics["unique_authors"] != 1 {
		t.Errorf("expected 1 unique author, got %v", metrics["unique_authors"])
	}
}

func TestCalculateProductivity_MultipleAuthors(t *testing.T) {
	commits := []commit{
		{author: "Alice"},
		{author: "Bob"},
		{author: "Alice"},
	}
	metrics := calculateProductivity(commits)
	if metrics["commits"] != 3 {
		t.Errorf("expected 3 commits, got %v", metrics["commits"])
	}
	if metrics["unique_authors"] != 2 {
		t.Errorf("expected 2 unique authors, got %v", metrics["unique_authors"])
	}
}

// TestBisectMarkGood tests marking commits as good in bisect
func TestBisectMarkGood_AddsHashToBisectGood(t *testing.T) {
	m := newModel(".")
	m.bisectState.active = true
	m.bisectState.current = "abc123"
	m.bisectState.good = []string{}

	result := bisectMarkGood(m)
	if len(result.bisectState.good) != 1 {
		t.Errorf("expected 1 good commit, got %d", len(result.bisectState.good))
	}
	if result.bisectState.good[0] != "abc123" {
		t.Errorf("expected 'abc123', got '%s'", result.bisectState.good[0])
	}
}

// TestBisectMarkBad tests marking commits as bad in bisect
func TestBisectMarkBad_AddsHashToBisectBad(t *testing.T) {
	m := newModel(".")
	m.bisectState.active = true
	m.bisectState.current = "def456"
	m.bisectState.bad = []string{}

	result := bisectMarkBad(m)
	if len(result.bisectState.bad) != 1 {
		t.Errorf("expected 1 bad commit, got %d", len(result.bisectState.bad))
	}
	if result.bisectState.bad[0] != "def456" {
		t.Errorf("expected 'def456', got '%s'", result.bisectState.bad[0])
	}
}

// TestBisectFindCulprit tests finding culprit in bisect
func TestBisectFindCulprit_RequiresGoodAndBad(t *testing.T) {
	commits := makeCommits(5)
	for i := range commits {
		commits[i].hash = string(byte(65 + i)) + "bc"
	}
	good := []string{}
	bad := []string{"abc"}

	culprit := bisectFindCulprit(commits, good, bad)
	if culprit != "" {
		t.Error("bisectFindCulprit should return empty when good list is empty")
	}
}

func TestBisectFindCulprit_WithGoodAndBad(t *testing.T) {
	commits := makeCommits(5)
	for i := range commits {
		commits[i].hash = string(byte(65+i)) + string(byte(65+i)) + string(byte(65+i))
	}
	good := []string{"aaa"}
	bad := []string{"ccc"}

	culprit := bisectFindCulprit(commits, good, bad)
	if culprit == "" {
		t.Error("bisectFindCulprit should return result")
	}
}

// TestInitiateBisect tests starting a bisect operation
func TestInitiateBisect_InitializesBisectState(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)

	result := initiateBisect(m)
	if !result.showBisectUI {
		t.Error("initiateBisect should set showBisectUI to true")
	}
}

// TestCalculateCommitMetrics tests commit metrics calculation
func TestCalculateCommitMetrics_BasicMetrics(t *testing.T) {
	metrics := calculateCommitMetrics("abc123", 10, 2)
	if metrics.linesChanged != 10 {
		t.Errorf("expected 10 lines changed, got %d", metrics.linesChanged)
	}
	if metrics.filesChanged != 2 {
		t.Errorf("expected 2 files changed, got %d", metrics.filesChanged)
	}
}

// TestLintCommitMessage tests message linting
func TestLintCommitMessage_ValidMessage(t *testing.T) {
	result := lintCommitMessage("Add feature", "abc123")
	if result.score < 0 {
		t.Errorf("lintCommitMessage should return non-negative score, got %d", result.score)
	}
}

func TestLintCommitMessage_EmptyMessage(t *testing.T) {
	result := lintCommitMessage("", "abc123")
	if len(result.issues) == 0 {
		t.Error("lintCommitMessage should find issues with empty message")
	}
}

// TestValidateCommitFormat tests commit format validation
func TestValidateCommitFormat_ValidFormat(t *testing.T) {
	errors := validateCommitFormat("Add feature description")
	if len(errors) != 0 {
		t.Errorf("expected no errors, got %d", len(errors))
	}
}

func TestValidateCommitFormat_EmptySubject(t *testing.T) {
	errors := validateCommitFormat("")
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors))
	}
}

// TestAnalyzeCodeOwnership tests code ownership analysis
func TestAnalyzeCodeOwnership_EmptyCommits(t *testing.T) {
	commits := []commit{}
	ownership := analyzeCodeOwnership(commits)
	if len(ownership) != 0 {
		t.Error("analyzeCodeOwnership should return empty for empty commits")
	}
}

func TestAnalyzeCodeOwnership_WithCommits(t *testing.T) {
	commits := []commit{
		{author: "Alice", hash: "abc"},
		{author: "Bob", hash: "def"},
	}
	ownership := analyzeCodeOwnership(commits)
	// Implementation may be a stub, so just verify it returns a map
	if ownership == nil {
		t.Error("analyzeCodeOwnership should return non-nil map")
	}
}

// TestDetectHotspots tests hotspot detection
func TestDetectHotspots_EmptyCommits(t *testing.T) {
	commits := []commit{}
	hotspots := detectHotspots(commits)
	if len(hotspots) != 0 {
		t.Error("detectHotspots should return empty for empty commits")
	}
}

func TestDetectHotspots_WithCommits(t *testing.T) {
	commits := makeCommits(5)
	hotspots := detectHotspots(commits)
	// Verify it returns a slice (may be empty if implementation is stub)
	if hotspots == nil {
		t.Error("detectHotspots should return non-nil slice")
	}
}

// TestAssessRiskLevel tests risk assessment
func TestAssessRiskLevel_LowRisk(t *testing.T) {
	hotspot := hotspotData{
		path:          "README.md",
		changeFrequency: 2,
		recentChanges: 1,
		collaborators: 1,
	}
	risk := assessRiskLevel(hotspot)
	if risk == "" {
		t.Error("assessRiskLevel should return non-empty string")
	}
}

// TestSemanticSearch tests semantic search functionality
func TestSemanticSearch_EmptyCommits(t *testing.T) {
	commits := []commit{}
	results := semanticSearch(commits, "feature")
	if len(results) != 0 {
		t.Error("semanticSearch should return empty for empty commits")
	}
}

func TestSemanticSearch_EmptyQuery(t *testing.T) {
	commits := makeCommits(3)
	results := semanticSearch(commits, "")
	// Should handle empty query gracefully
	if results == nil {
		t.Error("semanticSearch should return non-nil slice")
	}
}

// TestCalculateComplexityScore tests complexity scoring
func TestCalculateComplexityScore_BasicMetrics(t *testing.T) {
	metrics := commitMetrics{
		linesChanged: 10,
		filesChanged: 2,
		complexity:   3,
	}
	score := calculateComplexityScore(metrics)
	if score < 0 {
		t.Errorf("complexity score should be non-negative, got %d", score)
	}
}

func TestCalculateComplexityScore_HighComplexity(t *testing.T) {
	metrics := commitMetrics{
		linesChanged: 500,
		filesChanged: 20,
		complexity:   100,
	}
	score := calculateComplexityScore(metrics)
	if score <= 0 {
		t.Error("high complexity metrics should produce positive score")
	}
}
func TestRenderAuthorStats_Valid(t *testing.T) {
	stats := make(map[string]int)
	stats["Alice"] = 5
	_ = renderAuthorStats(stats, 80)
}

func TestRenderTimeStats_Valid(t *testing.T) {
	stats := make(map[string]int)
	stats["recent"] = 10
	_ = renderTimeStats(stats, 80)
}

func TestRenderProductivityMetrics_Valid(t *testing.T) {
	metrics := make(map[string]interface{})
	metrics["commits"] = 42
	_ = renderProductivityMetrics(metrics, 80)
}

func TestRenderRebaseUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderRebaseUI(m, 80)
}

func TestRenderAnalyticsPanel_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderAnalyticsPanel(m, 80)
}

func TestRenderBisectUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderBisectUI(m, 80)
}

func TestExtractReflogEntries_Valid(t *testing.T) {
	output := "abc1234 HEAD@{0}: commit: test\ndef5678 HEAD@{1}: commit: another"
	_ = extractReflogEntries(output)
}

func TestEnableReflogRecovery_Valid(t *testing.T) {
	m := newModel(".")
	_ = enableReflogRecovery(m)
}

func TestFindLostCommits_Valid(t *testing.T) {
	fsckOutput := "dangling commit abc1234\ndangling commit def5678"
	_ = findLostCommits(fsckOutput)
}

func TestRenderLostCommitsUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderLostCommitsUI(m, 80)
}

func TestPushUndo_Valid(t *testing.T) {
	m := newModel(".")
	_ = pushUndo(m, "abc1234")
}

func TestPerformUndo_Valid(t *testing.T) {
	m := newModel(".")
	_ = performUndo(m)
}

func TestRenderUndoMenu_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderUndoMenu(m, 80)
}

func TestDetectCodeOwners_Valid(t *testing.T) {
	ownership := make(map[string]codeOwnershipData)
	_ = detectCodeOwners(ownership)
}

func TestRenderCodeOwnershipUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderCodeOwnershipUI(m, 80)
}

func TestRenderHotspotsUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderHotspotsUI(m, 80)
}

func TestRenderLintingUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderLintingUI(m, 80)
}

func TestAnalyzeCommitSize_Valid(t *testing.T) {
	m := newModel(".")
	_ = analyzeCommitSize(m)
}

func TestRenderLargeCommitsUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderLargeCommitsUI(m, 80)
}

func TestAnalyzeComplexity_Valid(t *testing.T) {
	m := newModel(".")
	_ = analyzeComplexity(m)
}

func TestRenderComplexityUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderComplexityUI(m, 80)
}

func TestRenderSemanticSearchUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderSemanticSearchUI(m, 80)
}

func TestBuildActivityHeatmap_Valid(t *testing.T) {
	_ = buildActivityHeatmap([]commit{{when: "2024-01-01"}})
}

func TestFindPeakHour_Valid(t *testing.T) {
	data := authorActivityData{}
	_ = findPeakHour(data)
}

func TestRenderActivityHeatmapUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderActivityHeatmapUI(m, 80)
}

func TestAnalyzeMerges_Valid(t *testing.T) {
	_ = analyzeMerges([]commit{{hash: "a"}})
}

func TestRenderMergeAnalysisUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderMergeAnalysisUI(m, 80)
}

func TestAnalyzeCommitCoupling_Valid(t *testing.T) {
	_ = analyzeCommitCoupling([]commit{{hash: "a"}})
}

func TestRenderCouplingAnalysisUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderCouplingAnalysisUI(m, 80)
}

func TestToggleExtensionFilter_Valid(t *testing.T) {
	m := newModel(".")
	_ = toggleExtensionFilter(m, "go")
}

func TestRenderExtensionFilterUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderExtensionFilterUI(m, 80)
}

func TestFilterByExtension_Valid(t *testing.T) {
	_ = filterByExtension([]commit{{hash: "a"}}, "go")
}

func TestTrackDependencyChanges_Valid(t *testing.T) {
	_ = trackDependencyChanges([]commit{{hash: "a"}})
}

func TestRenderDependenciesUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderDependenciesUI(m, 80)
}

func TestExtractFilesFromSubject_Valid(t *testing.T) {
	_ = extractFilesFromSubject("update file.go and main.go")
}

func TestDetectFastForwardMerges_Valid(t *testing.T) {
	_ = detectFastForwardMerges([]commit{{hash: "a"}})
}

func TestRenderFastForwardsUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderFastForwardsUI(m, 80)
}

func TestRenderCommitGroupsUI_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderCommitGroupsUI(m, 80)
}

// TestCalculateBisectProgress_SingleCommit tests progress with single commit
func TestCalculateBisectProgress_SingleCommit(t *testing.T) {
	state := bisectState{
		candidates: []string{"hash1"},
	}

	progress := calculateBisectProgress(state)

	AssertEqual(t, 1, progress, "single commit should return 1")
}

// TestCalculateBisectProgress_MultipleCommits tests progress calculation
func TestCalculateBisectProgress_MultipleCommits(t *testing.T) {
	state := bisectState{
		candidates: []string{"h1", "h2", "h3", "h4", "h5"},
	}

	progress := calculateBisectProgress(state)

	AssertTrue(t, progress > 0, "progress should be positive")
}

// TestGroupCommits tests the groupCommits function with various grouping modes
func TestGroupCommits_EmptyInput(t *testing.T) {
	commits := []commit{}
	groups := groupCommits(commits, "date")

	AssertEqual(t, 0, len(groups), "empty input should return empty slice")
}

func TestGroupCommits_ByDate(t *testing.T) {
	commits := []commit{
		{hash: "abc123", when: "2024-01-15", subject: "Feature A"},
		{hash: "def456", when: "2024-01-15", subject: "Feature B"},
		{hash: "ghi789", when: "2024-01-15", subject: "Feature C"},
		{hash: "jkl012", when: "2024-01-16", subject: "Fix bug"},
		{hash: "mno345", when: "2024-01-16", subject: "Docs update"},
	}

	groups := groupCommits(commits, "date")

	AssertEqual(t, 2, len(groups), "should have 2 date groups")

	// Verify group contents
	dateGroups := make(map[string]int)
	for _, g := range groups {
		dateGroups[g.name] = len(g.commits)
		AssertEqual(t, "date", g.groupBy, "groupBy should be set to 'date'")
	}

	AssertEqual(t, 3, dateGroups["2024-01-15"], "should have 3 commits on 2024-01-15")
	AssertEqual(t, 2, dateGroups["2024-01-16"], "should have 2 commits on 2024-01-16")
}

func TestGroupCommits_ByBranch(t *testing.T) {
	commits := []commit{
		{hash: "abc123", subject: "feature: add login"},
		{hash: "def456", subject: "feature: add signup"},
		{hash: "ghi789", subject: "feature: add dashboard"},
		{hash: "jkl012", subject: "bugfix: fix timeout"},
		{hash: "mno345", subject: "bugfix: fix crash"},
		{hash: "pqr678", subject: "docs: update readme"},
	}

	groups := groupCommits(commits, "branch")

	AssertEqual(t, 3, len(groups), "should have 3 branch groups (feature, bugfix, docs)")

	// Verify group contents
	branchGroups := make(map[string]int)
	for _, g := range groups {
		branchGroups[g.name] = len(g.commits)
		AssertEqual(t, "branch", g.groupBy, "groupBy should be set to 'branch'")
	}

	AssertEqual(t, 3, branchGroups["feature:"], "should have 3 'feature:' commits")
	AssertEqual(t, 2, branchGroups["bugfix:"], "should have 2 'bugfix:' commits")
	AssertEqual(t, 1, branchGroups["docs:"], "should have 1 'docs:' commit")
}

func TestGroupCommits_UnknownGroupMode(t *testing.T) {
	commits := []commit{
		{hash: "abc123", when: "2024-01-15", subject: "Feature A"},
		{hash: "def456", when: "2024-01-15", subject: "Feature B"},
		{hash: "ghi789", when: "2024-01-16", subject: "Fix bug"},
	}

	groups := groupCommits(commits, "unknown")

	AssertEqual(t, 1, len(groups), "unknown mode should fall back to default grouping")
	AssertEqual(t, "default", groups[0].name, "should have 'default' group name")
	AssertEqual(t, 3, len(groups[0].commits), "should have all 3 commits in default group")
	AssertEqual(t, "unknown", groups[0].groupBy, "groupBy should preserve the mode string")
}

func TestGroupCommits_DefaultGrouping(t *testing.T) {
	commits := []commit{
		{hash: "abc123", when: "2024-01-15", subject: "Feature A"},
		{hash: "def456", when: "2024-01-16", subject: "Feature B"},
		{hash: "ghi789", when: "2024-01-17", subject: "Fix bug"},
	}

	groups := groupCommits(commits, "")

	AssertEqual(t, 1, len(groups), "empty mode should use default grouping")
	AssertEqual(t, "default", groups[0].name, "should have 'default' group name")
	AssertEqual(t, 3, len(groups[0].commits), "should have all 3 commits in default group")
}

func TestGroupCommits_PreservesCommitHashes(t *testing.T) {
	commits := []commit{
		{hash: "hash1", when: "2024-01-15", subject: "Commit 1"},
		{hash: "hash2", when: "2024-01-15", subject: "Commit 2"},
		{hash: "hash3", when: "2024-01-16", subject: "Commit 3"},
	}

	groups := groupCommits(commits, "date")

	// Collect all hashes from groups
	hashMap := make(map[string]bool)
	for _, g := range groups {
		for _, hash := range g.commits {
			hashMap[hash] = true
		}
	}

	AssertTrue(t, hashMap["hash1"], "should preserve hash1")
	AssertTrue(t, hashMap["hash2"], "should preserve hash2")
	AssertTrue(t, hashMap["hash3"], "should preserve hash3")
	AssertEqual(t, 3, len(hashMap), "should have all 3 hashes")
}

func TestGroupCommits_SingleCommit(t *testing.T) {
	commits := []commit{
		{hash: "abc123", when: "2024-01-15", subject: "Single feature"},
	}

	groups := groupCommits(commits, "date")

	AssertEqual(t, 1, len(groups), "should have 1 group")
	AssertEqual(t, "2024-01-15", groups[0].name, "group name should be the date")
	AssertEqual(t, 1, len(groups[0].commits), "group should have 1 commit")
	AssertEqual(t, "abc123", groups[0].commits[0], "should preserve the commit hash")
}

// TestDispatchAnalyticsFeature tests the dispatch function for analytics features
func TestDispatchAnalyticsFeature_InvalidIndex(t *testing.T) {
	m := newModel(".")

	// Test negative index
	result := dispatchAnalyticsFeature(m, -1)
	AssertEqual(t, false, result.showCodeOwnership, "should not change state for invalid index")

	// Test index >= analyticsMenuLen
	result = dispatchAnalyticsFeature(m, 100)
	AssertEqual(t, false, result.showCodeOwnership, "should not change state for out-of-range index")
}

func TestDispatchAnalyticsFeature_ToggleCodeOwnership(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)

	// Enable code ownership
	result := dispatchAnalyticsFeature(m, 0)
	AssertEqual(t, true, result.showCodeOwnership, "should enable code ownership")

	// Disable code ownership
	result = dispatchAnalyticsFeature(result, 0)
	AssertEqual(t, false, result.showCodeOwnership, "should disable code ownership")
}

func TestDispatchAnalyticsFeature_ToggleHotspots(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)

	// Enable hotspots
	result := dispatchAnalyticsFeature(m, 1)
	AssertEqual(t, true, result.showHotspots, "should enable hotspots")

	// Disable hotspots
	result = dispatchAnalyticsFeature(result, 1)
	AssertEqual(t, false, result.showHotspots, "should disable hotspots")
}

func TestDispatchAnalyticsFeature_ToggleLinting(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(2)
	m.commits[0].subject = "fix: bad commit message with very long text that should be linted"
	m.commits[1].subject = "feat: good commit message"

	// Enable linting
	result := dispatchAnalyticsFeature(m, 2)
	AssertEqual(t, true, result.showLinting, "should enable linting")
	AssertTrue(t, len(result.lintingResults) > 0, "should populate linting results")
}

func TestDispatchAnalyticsFeature_ToggleBisect(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)

	// Enable bisect
	result := dispatchAnalyticsFeature(m, 3)
	AssertEqual(t, true, result.bisectState.active, "should activate bisect")

	// Disable bisect
	result = dispatchAnalyticsFeature(result, 3)
	AssertEqual(t, false, result.bisectState.active, "should deactivate bisect")
	AssertEqual(t, false, result.showBisectUI, "should hide bisect UI")
}

func TestDispatchAnalyticsFeature_ToggleActivityHeatmap(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)

	// Enable activity heatmap
	result := dispatchAnalyticsFeature(m, 4)
	AssertEqual(t, true, result.showActivityHeatmap, "should enable activity heatmap")

	// Disable activity heatmap
	result = dispatchAnalyticsFeature(result, 4)
	AssertEqual(t, false, result.showActivityHeatmap, "should disable activity heatmap")
}

func TestDispatchAnalyticsFeature_ToggleAuthorStats(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)

	// Enable author stats
	result := dispatchAnalyticsFeature(m, 5)
	AssertEqual(t, true, result.showAnalytics, "should enable analytics")
	AssertTrue(t, len(result.authorStats) > 0, "should populate author stats")

	// Disable author stats
	result = dispatchAnalyticsFeature(result, 5)
	AssertEqual(t, false, result.showAnalytics, "should disable analytics")
}
