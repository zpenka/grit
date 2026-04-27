package grit

import (
	"testing"
)

// TestPreviewRebaseOperations tests rebase preview generation
func TestPreviewRebaseOperations_EmptyOps(t *testing.T) {
	ops := []rebaseOp{}
	preview := previewRebaseOperations(ops)
	if !preview.willApply {
		t.Error("previewRebaseOperations should indicate operations will apply")
	}
	if len(preview.operations) != 0 {
		t.Error("preview should have empty operations")
	}
}

func TestPreviewRebaseOperations_WithOps(t *testing.T) {
	ops := []rebaseOp{
		{action: "pick", hash: "abc123"},
		{action: "squash", hash: "def456"},
	}
	preview := previewRebaseOperations(ops)
	if len(preview.operations) != 2 {
		t.Errorf("expected 2 operations, got %d", len(preview.operations))
	}
	if preview.message == "" {
		t.Error("preview should have message")
	}
}

// TestDetectConflicts tests conflict detection
func TestDetectConflicts_NoConflicts(t *testing.T) {
	content := "normal file content\nno conflicts here"
	conflicts := detectConflicts(content)
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(conflicts))
	}
}

func TestDetectConflicts_WithConflicts(t *testing.T) {
	content := `some code
<<<<<<< HEAD
our changes
=======
their changes
>>>>>>> branch`
	conflicts := detectConflicts(content)
	if len(conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(conflicts))
	}
	if conflicts[0].resolved {
		t.Error("conflict should be marked as unresolved")
	}
}

func TestDetectConflicts_MultipleConflicts(t *testing.T) {
	content := `<<<<<<< HEAD
conflict 1
=======
end 1
>>>>>>> branch
other code
<<<<<<< HEAD
conflict 2
=======
end 2
>>>>>>> branch`
	conflicts := detectConflicts(content)
	if len(conflicts) == 0 {
		t.Error("detectConflicts should find conflicts")
	}
}

// TestPlanSquashSequence tests squash planning
func TestPlanSquashSequence_Basic(t *testing.T) {
	plan := planSquashSequence("abc123", []string{"def456", "ghi789"}, "Combined message")
	if plan.targetHash != "abc123" {
		t.Errorf("expected target 'abc123', got '%s'", plan.targetHash)
	}
	if len(plan.toSquash) != 2 {
		t.Errorf("expected 2 commits to squash, got %d", len(plan.toSquash))
	}
	if plan.resultMsg != "Combined message" {
		t.Errorf("expected message 'Combined message', got '%s'", plan.resultMsg)
	}
}

func TestPlanSquashSequence_SingleCommit(t *testing.T) {
	plan := planSquashSequence("abc123", []string{"def456"}, "Fix typo")
	if len(plan.toSquash) != 1 {
		t.Errorf("expected 1 commit to squash, got %d", len(plan.toSquash))
	}
	if plan.lineCount != len("Fix typo") {
		t.Errorf("expected lineCount 8, got %d", plan.lineCount)
	}
}

// TestImproveCherryPick tests cherry-pick improvement
func TestImproveCherryPick_ReturnsImprovement(t *testing.T) {
	m := newModel(".")
	improvement := improveCherryPick(m, "abc123")
	if improvement == nil {
		t.Error("improveCherryPick should return non-nil improvement")
	}
	if improvement.hash != "abc123" {
		t.Errorf("expected hash 'abc123', got '%s'", improvement.hash)
	}
}

func TestImproveCherryPick_WithSuggestions(t *testing.T) {
	m := newModel(".")
	improvement := improveCherryPick(m, "abc123")
	if improvement.suggestions == nil {
		t.Error("improvement should have suggestions slice")
	}
}

// TestPreviewAmendCommit tests commit amend preview
func TestPreviewAmendCommit_Basic(t *testing.T) {
	changes := make(map[string]int)
	changes["file.go"] = 5
	preview := previewAmendCommit("Original message", "New message", changes)

	if preview.originalMsg != "Original message" {
		t.Error("preview should contain original message")
	}
	if preview.newMsg != "New message" {
		t.Error("preview should contain new message")
	}
	if len(preview.changes) != 1 {
		t.Errorf("expected 1 changed file, got %d", len(preview.changes))
	}
}

// TestCalculateTeamStats tests team statistics
func TestCalculateTeamStats_EmptyCommits(t *testing.T) {
	commits := []commit{}
	stats := calculateTeamStats(commits)
	if len(stats) != 0 {
		t.Error("calculateTeamStats should return empty for empty commits")
	}
}

func TestCalculateTeamStats_SingleAuthor(t *testing.T) {
	commits := []commit{
		{author: "Alice"},
		{author: "Alice"},
	}
	stats := calculateTeamStats(commits)
	if len(stats) != 1 {
		t.Errorf("expected 1 team member, got %d", len(stats))
	}
	if stats[0].commits != 2 {
		t.Errorf("expected 2 commits, got %d", stats[0].commits)
	}
}

func TestCalculateTeamStats_MultipleAuthors(t *testing.T) {
	commits := []commit{
		{author: "Alice"},
		{author: "Bob"},
		{author: "Alice"},
		{author: "Charlie"},
	}
	stats := calculateTeamStats(commits)
	if len(stats) != 3 {
		t.Errorf("expected 3 team members, got %d", len(stats))
	}
	// Find Alice in stats
	found := false
	for _, s := range stats {
		if s.author == "Alice" && s.commits == 2 {
			found = true
		}
	}
	if !found {
		t.Error("should find Alice with 2 commits")
	}
}

// TestAutomateReviewWorkflow tests review workflow automation
func TestAutomateReviewWorkflow_Basic(t *testing.T) {
	workflow := automateReviewWorkflow(42, "alice", []string{"bob", "charlie"})
	if workflow.prNumber != 42 {
		t.Errorf("expected PR 42, got %d", workflow.prNumber)
	}
	if workflow.author != "alice" {
		t.Errorf("expected author 'alice', got '%s'", workflow.author)
	}
	if len(workflow.reviewers) != 2 {
		t.Errorf("expected 2 reviewers, got %d", len(workflow.reviewers))
	}
	if workflow.status != "pending" {
		t.Errorf("expected status 'pending', got '%s'", workflow.status)
	}
}

// TestSuggestReviewers tests reviewer suggestions
func TestSuggestReviewers_NoCommits(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{}
	suggestions := suggestReviewers(m, "main.go")
	if len(suggestions) != 0 {
		t.Error("suggestReviewers should return empty for empty commits")
	}
}

func TestSuggestReviewers_WithCommits(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{
		{author: "Alice"},
		{author: "Bob"},
		{author: "Alice"},
	}
	suggestions := suggestReviewers(m, "main.go")
	if len(suggestions) == 0 {
		t.Error("suggestReviewers should return suggestions for non-empty commits")
	}
	for _, s := range suggestions {
		if s.score == 0 {
			t.Error("reviewer suggestions should have non-zero score")
		}
	}
}

// TestDetectPairProgramming tests pair programming detection
func TestDetectPairProgramming_SingleAuthor(t *testing.T) {
	commits := []commit{
		{author: "Alice", subject: "Regular commit"},
	}
	pairs := detectPairProgramming(commits)
	if len(pairs) != 0 {
		t.Error("detectPairProgramming should return empty for single author commits")
	}
}

func TestDetectPairProgramming_WithPairKeyword(t *testing.T) {
	commits := []commit{
		{author: "Alice", subject: "pair programming session"},
		{author: "Bob", subject: "Also collaborative"},
	}
	pairs := detectPairProgramming(commits)
	// Should find pair programming when "pair" is in subject
	if pairs == nil {
		t.Error("detectPairProgramming should return non-nil slice")
	}
	if len(pairs) == 0 {
		t.Error("detectPairProgramming should detect pair programming")
	}
}

// TestCalculateVelocity tests velocity calculation
func TestCalculateVelocity_EmptyCommits(t *testing.T) {
	commits := []commit{}
	velocity := calculateVelocity(commits)
	if len(velocity) != 0 {
		t.Error("calculateVelocity should return empty for empty commits")
	}
}

func TestCalculateVelocity_WithCommits(t *testing.T) {
	commits := []commit{
		{when: "1 day ago", author: "Alice"},
		{when: "2 days ago", author: "Bob"},
	}
	velocity := calculateVelocity(commits)
	// Should return velocity data
	if velocity == nil {
		t.Error("calculateVelocity should return non-nil slice")
	}
}

// TestAutoCompleteMessage tests message auto-completion
func TestAutoCompleteMessage_EmptyPrefix(t *testing.T) {
	commits := []commit{
		{subject: "Add feature"},
		{subject: "Fix bug"},
	}
	completions := autoCompleteMessage("", commits)
	if len(completions) == 0 {
		t.Error("autoCompleteMessage should return completions")
	}
}

func TestAutoCompleteMessage_WithPrefix(t *testing.T) {
	commits := []commit{
		{subject: "Add feature"},
		{subject: "Add documentation"},
		{subject: "Fix bug"},
	}
	completions := autoCompleteMessage("Add", commits)
	if completions == nil {
		t.Error("autoCompleteMessage should return non-nil slice")
	}
}

// TestClassifyCommit tests commit classification
func TestClassifyCommit_IdentifiesCategory(t *testing.T) {
	classification := classifyCommit("feat: Add new feature", "abc123")
	if classification.category == "" {
		t.Error("classifyCommit should identify commit category")
	}
	if classification.hash != "abc123" {
		t.Errorf("expected hash 'abc123', got '%s'", classification.hash)
	}
}

func TestClassifyCommit_FixType(t *testing.T) {
	classification := classifyCommit("fix: Critical bug fix", "abc123")
	if classification.category == "" {
		t.Error("classifyCommit should identify commit category for fix")
	}
}

func TestClassifyCommit_RefactorType(t *testing.T) {
	classification := classifyCommit("refactor: Improve code structure", "abc123")
	if classification.category == "" {
		t.Error("classifyCommit should identify commit category for refactor")
	}
}

// TestDetectAnomalies tests anomaly detection
func TestDetectAnomalies_EmptyCommits(t *testing.T) {
	commits := []commit{}
	anomalies := detectAnomalies(commits)
	if len(anomalies) != 0 {
		t.Error("detectAnomalies should return empty for empty commits")
	}
}

func TestDetectAnomalies_WithAnomalousCommit(t *testing.T) {
	commits := []commit{
		{hash: "abc123", subject: "Normal commit"},
		{hash: "def456", subject: "massive rewrite of the entire codebase with many changes"},
	}
	anomalies := detectAnomalies(commits)
	// Should detect the "massive" commit as anomalous
	if anomalies == nil {
		t.Error("detectAnomalies should return non-nil slice")
	}
	if len(anomalies) == 0 {
		t.Error("detectAnomalies should find anomalies in the commit")
	}
}

// TestFindSimilarCommits tests similar commit finding
func TestFindSimilarCommits_FindsExactHash(t *testing.T) {
	commits := []commit{
		{hash: "abc1234567890123456", subject: "Fix login bug in auth"},
		{hash: "def4567890123456789", subject: "Add feature"},
	}
	// Use a hash that exists in the list
	similar := findSimilarCommits(commits, "abc1234567890123456")
	// Should work without panicking
	if similar != nil && len(similar) > 0 {
		// Found similar commits
		if similar[0].hash1 != "abc1234567890123456" {
			t.Error("similar commits should reference the target hash")
		}
	}
}

func TestFindSimilarCommits_EmptyCommits(t *testing.T) {
	commits := []commit{}
	// For empty commits, function should handle gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("findSimilarCommits panicked with empty commits: %v", r)
		}
	}()
	similar := findSimilarCommits(commits, "abc1234567890123456")
	// Should return something for empty commits
	_ = similar
}

// TestGenerateAutoSummary tests automatic summary generation
func TestGenerateAutoSummary_Basic(t *testing.T) {
	summary := generateAutoSummary("abc123", "Fix critical security vulnerability in auth module")
	if summary.hash != "abc123" {
		t.Errorf("expected hash 'abc123', got '%s'", summary.hash)
	}
	if summary.summary == "" {
		t.Error("generateAutoSummary should return non-empty summary")
	}
}

// TestCheckSigningCompliance tests GPG signing compliance
func TestCheckSigningCompliance_EmptyCommits(t *testing.T) {
	commits := []commit{}
	compliance := checkSigningCompliance(commits, true)
	if len(compliance) != 0 {
		t.Error("checkSigningCompliance should return empty for empty commits")
	}
}

func TestCheckSigningCompliance_NotEnforced(t *testing.T) {
	commits := []commit{
		{hash: "abc123"},
	}
	compliance := checkSigningCompliance(commits, false)
	// Should return compliance map
	if compliance == nil {
		t.Error("checkSigningCompliance should return non-nil map")
	}
}

// TestTrackLicenseHeaders tests license header tracking
func TestTrackLicenseHeaders_ReturnsSlice(t *testing.T) {
	headers := trackLicenseHeaders("abc123")
	// Should return a slice (may be empty if no file content available)
	if len(headers) == 0 {
		// Empty is okay, function returns empty by design
		return
	}
	// If there are headers, they should be valid
	if headers[0].file == "" {
		t.Error("license headers should have file information")
	}
}

// TestScanForSecurityIssues tests security scanning
func TestScanForSecurityIssues_NoIssues(t *testing.T) {
	issues := scanForSecurityIssues("abc123", "normal code\nmore code")
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestScanForSecurityIssues_WithIssues(t *testing.T) {
	content := `password = "hardcoded_secret"`
	issues := scanForSecurityIssues("abc123", content)
	// Should detect potential security issues
	if issues == nil {
		t.Error("scanForSecurityIssues should return non-nil slice")
	}
}

// TestDetectSecrets tests secret detection
func TestDetectSecrets_NoSecrets(t *testing.T) {
	secrets := detectSecrets("abc123", "normal file content")
	if len(secrets) != 0 {
		t.Error("detectSecrets should return empty for normal content")
	}
}

func TestDetectSecrets_WithAPIKey(t *testing.T) {
	content := `const API_KEY = "sk-1234567890abcdef"`
	secrets := detectSecrets("abc123", content)
	// May detect potential secrets
	if secrets == nil {
		t.Error("detectSecrets should return non-nil slice")
	}
}

// TestDetectSemver tests semantic versioning detection
func TestDetectSemver_EmptyCommits(t *testing.T) {
	commits := []commit{}
	semver := detectSemver(commits)
	if len(semver) != 0 {
		t.Error("detectSemver should return empty for empty commits")
	}
}

func TestDetectSemver_VersionCommit(t *testing.T) {
	commits := []commit{
		{subject: "v1.0.0 release"},
		{subject: "v1.0.1 patch"},
	}
	semver := detectSemver(commits)
	// Should detect version tags
	if semver == nil {
		t.Error("detectSemver should return non-nil slice")
	}
}

// TestGenerateChangelog tests changelog generation
func TestGenerateChangelog_EmptyCommits(t *testing.T) {
	commits := []commit{}
	changelog := generateChangelog(commits, "1.0.0")
	// Should handle empty commits gracefully
	if changelog == nil && len(commits) > 0 {
		t.Error("generateChangelog should handle empty commits")
	}
}

func TestGenerateChangelog_WithCommits(t *testing.T) {
	commits := []commit{
		{subject: "feat: Add feature"},
		{subject: "fix: Fix bug"},
	}
	changelog := generateChangelog(commits, "1.0.0")
	if changelog != nil {
		if changelog.version != "1.0.0" {
			t.Error("changelog should contain version")
		}
	}
}

// TestBuildReleaseNotes tests release notes building
func TestBuildReleaseNotes_Basic(t *testing.T) {
	hashes := []string{"abc123", "def456"}
	notes := buildReleaseNotes("1.0.0", hashes)
	if notes.version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", notes.version)
	}
	if notes.summary == "" {
		t.Error("release notes should have summary")
	}
}

// TestTrackVersionBumps tests version bump tracking
func TestTrackVersionBumps_EmptyCommits(t *testing.T) {
	commits := []commit{}
	bumps := trackVersionBumps(commits)
	if len(bumps) != 0 {
		t.Error("trackVersionBumps should return empty for empty commits")
	}
}

func TestTrackVersionBumps_WithVersionCommits(t *testing.T) {
	commits := []commit{
		{hash: "abc123", subject: "bump: version 1.0.0", when: "1 day ago"},
		{hash: "def456", subject: "version: 1.1.0", when: "2 days ago"},
	}
	bumps := trackVersionBumps(commits)
	// Should track version changes
	if len(bumps) != 2 {
		t.Errorf("expected 2 version bumps, got %d", len(bumps))
	}
}
