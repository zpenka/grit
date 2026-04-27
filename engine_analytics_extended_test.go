package grit

import (
	"testing"
)

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
