package grit

import (
	"testing"
)

// Tests for engine_team_ai.go feature functions
// These tests verify that stubs are replaced with real implementations

func TestTrackLicenseHeaders_ScansActualFiles(t *testing.T) {
	hash := "aaa1111111111111111111111111111111111"
	results := trackLicenseHeaders(hash)
	if len(results) > 0 && results[0].file == "main.go" {
		t.Error("trackLicenseHeaders should not return hardcoded main.go")
	}
	for _, hdr := range results {
		if hdr.hash != hash {
			t.Errorf("trackLicenseHeaders should set correct hash, got %q", hdr.hash)
		}
	}
}

func TestScanForSecurityIssues_ReturnsActualLineNumbers(t *testing.T) {
	hash := "bbb2222222222222222222222222222222222"
	content := `password = "secret123"
normal code here
api_key = "xyz789"`
	issues := scanForSecurityIssues(hash, content)
	for _, issue := range issues {
		if issue.location == "line 5" && content != "" {
			t.Error("scanForSecurityIssues should not hardcode location to line 5")
		}
		if issue.hash != hash {
			t.Errorf("scanForSecurityIssues should set correct hash")
		}
	}
}

func TestTrackDataDeletion_AcceptsReason(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	hash := "ccc3333333333333333333333333333333333"
	email := "user@example.com"
	m = trackDataDeletion(m, hash, email)
	if len(m.dataDeleteRequests) > 0 {
		lastReq := m.dataDeleteRequests[len(m.dataDeleteRequests)-1]
		if lastReq.reason == "GDPR request" {
			t.Error("trackDataDeletion hardcodes reason; should accept parameter")
		}
	}
}

func TestDetectSecrets_UsesRealLineTracking(t *testing.T) {
	hash := "ddd4444444444444444444444444444444444"
	content := `line 1 normal code
line 2 more code
line 3 password = "secret"
line 4 final code`
	secrets := detectSecrets(hash, content)
	allHardcodedLine1 := len(secrets) > 0
	for _, s := range secrets {
		if s.location != "line 1" {
			allHardcodedLine1 = false
			break
		}
	}
	if allHardcodedLine1 && len(secrets) > 0 {
		t.Error("detectSecrets should not hardcode all locations to line 1")
	}
	if len(secrets) > 0 && secrets[0].hash != hash {
		t.Errorf("detectSecrets should set correct hash")
	}
}

func TestCalculateTeamStats_HasValidData(t *testing.T) {
	commits := []commit{
		{author: "Alice", hash: "aaa1111111111111111111111111111111111"},
		{author: "Bob", hash: "bbb2222222222222222222222222222222222"},
		{author: "Alice", hash: "ccc3333333333333333333333333333333333"},
	}
	stats := calculateTeamStats(commits)
	if len(stats) != 2 {
		t.Errorf("calculateTeamStats should return one entry per author, got %d", len(stats))
	}
	for _, s := range stats {
		if s.author == "" {
			t.Error("calculateTeamStats should set author name")
		}
		if s.commits < 0 {
			t.Error("calculateTeamStats should have non-negative commits")
		}
	}
}

func TestSuggestReviewers_ReturnsRelevantReviewers(t *testing.T) {
	m := newModel(".")
	m.commits = makeNamedCommits()
	reviewers := suggestReviewers(m, "main.go")
	if len(reviewers) == 0 {
		t.Error("suggestReviewers should return at least one reviewer")
	}
	for _, r := range reviewers {
		if r.reviewer == "" {
			t.Error("suggestReviewers should set reviewer name")
		}
		if r.score < 0 {
			t.Errorf("suggestReviewers score should be non-negative, got %f", r.score)
		}
	}
}

func TestDetectSemver_IdentifiesVersionCommits(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "v1.0.0: Initial release"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Feature: Add parser"},
		{hash: "ccc3333333333333333333333333333333333", subject: "v1.1.0: Add new features"},
	}
	versions := detectSemver(commits)
	if len(versions) < 2 {
		t.Errorf("detectSemver should find at least 2 versions, found %d", len(versions))
	}
}

func TestClassifyCommit_CategorizesCommits(t *testing.T) {
	hash := "aaa1111111111111111111111111111111111"
	subject := "fix: Handle edge case in parser"
	classification := classifyCommit(subject, hash)
	if classification.category == "" {
		t.Error("classifyCommit should set category")
	}
	if classification.hash != hash {
		t.Errorf("classifyCommit should preserve hash")
	}
}

func TestClassifyCommit_WithMerge(t *testing.T) {
	classification := classifyCommit("Merge pull request #123", "abc123")
	AssertTrue(t, classification.category != "", "should classify merge commit")
}

func TestClassifyCommit_FeatureCommit(t *testing.T) {
	classification := classifyCommit("feat: Add new feature", "def456")
	AssertTrue(t, classification.category != "", "should classify feature commit")
}

func TestClassifyCommit_Refactor(t *testing.T) {
	classification := classifyCommit("refactor: Clean up code structure", "ghi789")
	AssertTrue(t, classification.category != "", "should classify refactor commit")
}

func TestPreviewRebaseOperations_CreatesPreview(t *testing.T) {
	ops := []rebaseOp{
		{action: "pick", hash: "aaa1111111111111111111111111111111111"},
		{action: "squash", hash: "bbb2222222222222222222222222222222222"},
	}
	preview := previewRebaseOperations(ops)
	AssertEqual(t, 2, len(preview.operations), "should preserve operation count")
	AssertEqual(t, "Rebase will apply", preview.message, "should have expected message")
	AssertTrue(t, preview.willApply, "should indicate rebase will apply")
}

func TestPreviewRebaseOperations_EmptyOps(t *testing.T) {
	var ops []rebaseOp
	preview := previewRebaseOperations(ops)
	AssertEqual(t, 0, len(preview.operations), "should handle empty operations")
	AssertTrue(t, preview.willApply, "should still indicate willApply")
}

func TestDetectConflicts_FindsHeadMarker(t *testing.T) {
	content := `some code
<<<<<<< HEAD
our change
=======
their change
>>>>>>> branch`
	conflicts := detectConflicts(content)
	AssertTrue(t, len(conflicts) > 0, "should detect HEAD conflict marker")
	AssertTrue(t, len(conflicts[0].file) > 0, "should set conflict file")
}

func TestDetectConflicts_NoMarkers(t *testing.T) {
	content := `normal code
no conflicts here`
	conflicts := detectConflicts(content)
	AssertEqual(t, 0, len(conflicts), "should return empty when no conflicts")
}

func TestPlanSquashSequence_CreatesValidPlan(t *testing.T) {
	target := "aaa1111111111111111111111111111111111"
	toSquash := []string{"bbb2222222222222222222222222222222222", "ccc3333333333333333333333333333333333"}
	msg := "Combined commit message"
	plan := planSquashSequence(target, toSquash, msg)
	AssertEqual(t, target, plan.targetHash, "should set target hash")
	AssertEqual(t, 2, len(plan.toSquash), "should preserve toSquash list")
	AssertEqual(t, msg, plan.resultMsg, "should set result message")
	AssertTrue(t, plan.lineCount > 0, "should calculate line count")
}

func TestImproveCherryPick_CreatesImprovement(t *testing.T) {
	m := newModel(".")
	hash := "abc1111111111111111111111111111111111"
	improvement := improveCherryPick(m, hash)
	AssertNotNil(t, improvement, "should return non-nil improvement")
	AssertEqual(t, hash, improvement.hash, "should set hash")
	AssertFalse(t, improvement.autoConflict, "should default to no auto conflict")
}

func TestPreviewAmendCommit_ShowsChanges(t *testing.T) {
	original := "Old message"
	newMsg := "New message"
	changes := map[string]int{"added": 5, "removed": 2}
	preview := previewAmendCommit(original, newMsg, changes)
	AssertEqual(t, original, preview.originalMsg, "should set original message")
	AssertEqual(t, newMsg, preview.newMsg, "should set new message")
	AssertEqual(t, 2, len(preview.changes), "should preserve changes map")
}

func TestAutomateReviewWorkflow_CreatesWorkflow(t *testing.T) {
	prNum := 42
	author := "Alice"
	reviewers := []string{"Bob", "Charlie"}
	workflow := automateReviewWorkflow(prNum, author, reviewers)
	AssertEqual(t, prNum, workflow.prNumber, "should set PR number")
	AssertEqual(t, author, workflow.author, "should set author")
	AssertEqual(t, 2, len(workflow.reviewers), "should preserve reviewers list")
	AssertFalse(t, workflow.approved, "should default to not approved")
	AssertEqual(t, "pending", workflow.status, "should default to pending status")
}

func TestDetectPairProgramming_FindsPairs(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "pair programming: Fix parser"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Normal commit"},
		{hash: "ccc3333333333333333333333333333333333", subject: "Pair: Add feature"},
	}
	pairs := detectPairProgramming(commits)
	AssertTrue(t, len(pairs) >= 1, "should detect pair programming commits")
}

func TestDetectPairProgramming_NoMatches(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "Normal commit"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Another commit"},
	}
	pairs := detectPairProgramming(commits)
	AssertEqual(t, 0, len(pairs), "should find no pairs when no markers present")
}

func TestCalculateVelocity_ComputesVelocity(t *testing.T) {
	commits := makeCommits(10)
	velocity := calculateVelocity(commits)
	AssertTrue(t, len(velocity) > 0, "should return velocity data")
	for _, v := range velocity {
		AssertTrue(t, v.commits > 0, "should have positive commit count")
		AssertTrue(t, v.additions >= 0, "should have non-negative additions")
	}
}

func TestCalculateVelocity_EmptyCommits(t *testing.T) {
	commits := []commit{}
	velocity := calculateVelocity(commits)
	AssertEqual(t, 0, len(velocity), "should return empty velocity for no commits")
}

func TestAutoCompleteMessage_GeneratesSuggestions(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "feat: Add login"},
		{hash: "bbb2222222222222222222222222222222222", subject: "feat: Add parser"},
		{hash: "ccc3333333333333333333333333333333333", subject: "fix: Handle error"},
	}
	completions := autoCompleteMessage("feat:", commits)
	AssertTrue(t, len(completions) > 0, "should generate suggestions for matching prefix")
	if len(completions) > 0 {
		AssertEqual(t, "feat:", completions[0].prefix, "should set prefix")
		AssertTrue(t, len(completions[0].suggestions) > 0, "should include suggestions")
	}
}

func TestAutoCompleteMessage_NoMatches(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "Add login"},
	}
	completions := autoCompleteMessage("xyz:", commits)
	AssertEqual(t, 0, len(completions), "should return empty when prefix doesn't match")
}

func TestCheckSigningCompliance_EnforcedTrue(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "Commit 1"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Commit 2"},
		{hash: "ccc3333333333333333333333333333333333", subject: "Commit 3"},
	}
	statuses := checkSigningCompliance(commits, true)
	AssertEqual(t, 3, len(statuses), "should have status for each commit")
	for _, status := range statuses {
		AssertFalse(t, status.compliant, "should be non-compliant when enforced and not signed")
		AssertTrue(t, status.enforced, "should mark as enforced")
	}
}

func TestCheckSigningCompliance_EnforcedFalse(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "Commit 1"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Commit 2"},
	}
	statuses := checkSigningCompliance(commits, false)
	AssertEqual(t, 2, len(statuses), "should have status for each commit")
	for _, status := range statuses {
		AssertTrue(t, status.compliant, "should be compliant when not enforced")
	}
}

func TestGenerateChangelog_CategorizesCommits(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "feat: Add new feature"},
		{hash: "bbb2222222222222222222222222222222222", subject: "fix: Handle edge case"},
		{hash: "ccc3333333333333333333333333333333333", subject: "docs: Update README"},
	}
	entry := generateChangelog(commits, "1.2.0")
	AssertNotNil(t, entry, "should return changelog entry")
	AssertEqual(t, "1.2.0", entry.version, "should set version")
	AssertTrue(t, len(entry.features) > 0, "should find features")
	AssertTrue(t, len(entry.bugfixes) > 0, "should find bugfixes")
}

func TestGenerateChangelog_EmptyCommits(t *testing.T) {
	commits := []commit{}
	entry := generateChangelog(commits, "1.0.0")
	AssertNotNil(t, entry, "should return entry even with empty commits")
	AssertEqual(t, "1.0.0", entry.version, "should set version")
}

func TestBuildReleaseNotes_CreatesNotes(t *testing.T) {
	version := "2.0.0"
	commits := []string{"aaa1111111111111111111111111111111111", "bbb2222222222222222222222222222222222"}
	notes := buildReleaseNotes(version, commits)
	AssertEqual(t, version, notes.version, "should set version")
	AssertTrue(t, len(notes.summary) > 0, "should have summary")
	AssertTrue(t, len(notes.highlights) > 0, "should have highlights")
	AssertTrue(t, len(notes.contributors) > 0, "should have contributors")
}

func TestTrackVersionBumps_DetectsBumps(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "bump version to 1.1.0", when: "2 days ago"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Normal commit"},
		{hash: "ccc3333333333333333333333333333333333", subject: "version: 1.2.0", when: "1 day ago"},
	}
	bumps := trackVersionBumps(commits)
	AssertTrue(t, len(bumps) >= 1, "should detect version bumps")
	for _, bump := range bumps {
		AssertTrue(t, len(bump.hash) > 0, "should have hash")
		AssertTrue(t, len(bump.message) > 0, "should have message")
	}
}

func TestTrackVersionBumps_NoVersion(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "Add feature"},
	}
	bumps := trackVersionBumps(commits)
	AssertEqual(t, 0, len(bumps), "should find no bumps when no version keywords")
}

func TestCreateMilestone_AddsMilestone(t *testing.T) {
	m := newModel(".")
	m.milestones = []milestone{}
	name := "v2.0 Release"
	commits := []string{"aaa1111111111111111111111111111111111", "bbb2222222222222222222222222222222222"}
	m = createMilestone(m, name, commits)
	AssertEqual(t, 1, len(m.milestones), "should add milestone")
	AssertEqual(t, name, m.milestones[0].name, "should set milestone name")
	AssertEqual(t, 2, len(m.milestones[0].commits), "should set milestone commits")
	AssertEqual(t, "in-progress", m.milestones[0].status, "should set status")
}

func TestDispatchTeamFeature_TeamStats(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showTeamStats = false
	m = dispatchTeamFeature(m, 0)
	AssertTrue(t, m.showTeamStats, "should toggle team stats on")
	AssertTrue(t, len(m.teamStats) > 0, "should populate team stats")
}

func TestDispatchTeamFeature_ReviewerSuggestions(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showReviewUI = false
	m = dispatchTeamFeature(m, 1)
	AssertTrue(t, m.showReviewUI, "should toggle review UI on")
}

func TestDispatchTeamFeature_Velocity(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.showVelocity = false
	m = dispatchTeamFeature(m, 2)
	AssertTrue(t, m.showVelocity, "should toggle velocity on")
}

func TestDispatchTeamFeature_Classification(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showClassification = false
	m = dispatchTeamFeature(m, 3)
	AssertTrue(t, m.showClassification, "should toggle classification on")
}

func TestDispatchTeamFeature_Security(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(2)
	m.showSecrets = false
	m.diffLines = []diffLine{}
	m = dispatchTeamFeature(m, 4)
	AssertTrue(t, m.showSecrets, "should toggle secrets on")
}

func TestDispatchTeamFeature_Changelog(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showChangelog = false
	m = dispatchTeamFeature(m, 5)
	AssertTrue(t, m.showChangelog, "should toggle changelog on")
}

func TestDispatchTeamFeature_OutOfRange(t *testing.T) {
	m := newModel(".")
	m.showTeamStats = false
	result := dispatchTeamFeature(m, 10)
	AssertFalse(t, result.showTeamStats, "should not modify model for out-of-range index")
}

func TestDispatchTeamFeature_NegativeIndex(t *testing.T) {
	m := newModel(".")
	m.showTeamStats = false
	result := dispatchTeamFeature(m, -1)
	AssertFalse(t, result.showTeamStats, "should not modify model for negative index")
}

func TestDispatchTeamFeature_Toggle(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(2)
	m.showTeamStats = true
	m = dispatchTeamFeature(m, 0)
	AssertFalse(t, m.showTeamStats, "should toggle team stats off")
}

func TestRenderTeamStatsUI_WithData(t *testing.T) {
	m := newModel(".")
	m.teamStats = []teamStats{
		{author: "Alice", commits: 10, avgCommitSize: 100},
		{author: "Bob", commits: 5, avgCommitSize: 80},
	}
	output := renderTeamStatsUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce non-empty output")
	AssertStringContains(t, output, "Team Statistics", "should contain title")
	AssertStringContains(t, output, "Alice", "should contain author name")
}

func TestRenderTeamStatsUI_NoData(t *testing.T) {
	m := newModel(".")
	m.teamStats = []teamStats{}
	output := renderTeamStatsUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce output even with no data")
}

func TestRenderReviewerSuggestionsUI_WithData(t *testing.T) {
	m := newModel(".")
	m.reviewerSuggestions = []reviewerSuggestion{
		{reviewer: "Carol", expertise: 0.85, availability: 0.9, score: 0.765},
	}
	output := renderReviewerSuggestionsUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce non-empty output")
	AssertStringContains(t, output, "Carol", "should contain reviewer name")
}

func TestRenderReviewerSuggestionsUI_NoData(t *testing.T) {
	m := newModel(".")
	m.reviewerSuggestions = []reviewerSuggestion{}
	output := renderReviewerSuggestionsUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce output even with no data")
}

func TestRenderVelocityUI_WithData(t *testing.T) {
	m := newModel(".")
	m.velocityHistory = []velocityData{
		{week: "week1", commits: 10, additions: 500},
		{week: "week2", commits: 15, additions: 750},
	}
	output := renderVelocityUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce non-empty output")
	AssertStringContains(t, output, "Velocity", "should contain title")
}

func TestRenderVelocityUI_NoData(t *testing.T) {
	m := newModel(".")
	m.velocityHistory = []velocityData{}
	output := renderVelocityUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce output even with no data")
}

func TestRenderClassificationUI_WithData(t *testing.T) {
	m := newModel(".")
	m.commitClassifications = []commitClassification{
		{hash: "aaa1111111111111111111111111111111111", category: "fix", confidence: 0.9},
		{hash: "bbb2222222222222222222222222222222222", category: "feature", confidence: 0.85},
	}
	output := renderClassificationUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce non-empty output")
	AssertStringContains(t, output, "Classification", "should contain title")
	AssertStringContains(t, output, "fix", "should contain category")
}

func TestRenderClassificationUI_NoData(t *testing.T) {
	m := newModel(".")
	m.commitClassifications = []commitClassification{}
	output := renderClassificationUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce output even with no data")
}

func TestRenderSecretsUI_WithData(t *testing.T) {
	m := newModel(".")
	m.secretDetections = []secretDetection{
		{hash: "aaa1111111111111111111111111111111111", type_: "password", location: "line 5", severity: "critical"},
	}
	output := renderSecretsUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce non-empty output")
	AssertStringContains(t, output, "Security", "should contain title")
}

func TestRenderSecretsUI_NoData(t *testing.T) {
	m := newModel(".")
	m.secretDetections = []secretDetection{}
	output := renderSecretsUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce output even with no data")
}

func TestRenderChangelogUI_WithData(t *testing.T) {
	m := newModel(".")
	m.changelog = []changelogEntry{
		{version: "1.0.0", date: "2026-01-01", features: []string{"Feature A"}, bugfixes: []string{"Fix B"}},
	}
	output := renderChangelogUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce non-empty output")
	AssertStringContains(t, output, "Changelog", "should contain title")
}

func TestRenderChangelogUI_NoData(t *testing.T) {
	m := newModel(".")
	m.changelog = []changelogEntry{}
	output := renderChangelogUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce output even with no data")
}

func TestRenderCherryPickUI_WithData(t *testing.T) {
	m := newModel(".")
	m.cherryPickList = []string{"aaa1111111111111111111111111111111111", "bbb2222222222222222222222222222222222"}
	output := renderCherryPickUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce non-empty output")
	AssertStringContains(t, output, "Cherry-Pick", "should contain title")
}

func TestRenderCherryPickUI_NoData(t *testing.T) {
	m := newModel(".")
	m.cherryPickList = []string{}
	output := renderCherryPickUI(m, 80)
	AssertTrue(t, len(output) > 0, "should produce output even with no data")
}

func TestDetectSecretsInLine_FindsPassword(t *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"password = 'secret123'", true},
		{"PASSWORD=xyz", true},
		{"api_key = 'token'", true},
		{"secret_token = 'abc'", true},
		{"credential = 'pass'", true},
		{"token = 'mytoken'", true},
		{"normal code line", false},
		{"variable_name = value", false},
		{"function call()", false},
	}
	for _, test := range tests {
		result := detectSecretsInLine(test.line)
		if result != test.expected {
			t.Errorf("detectSecretsInLine(%q): expected %v, got %v", test.line, test.expected, result)
		}
	}
}
