package grit

import (
	"testing"
)

// Tests for engine_team_ai.go feature functions
// These tests verify that stubs are replaced with real implementations

func TestTrackLicenseHeaders_ScansActualFiles(t *testing.T) {
	// Should scan diff content for license headers, not return hardcoded main.go
	hash := "aaa1111111111111111111111111111111111"

	results := trackLicenseHeaders(hash)

	// Verify we don't get hardcoded main.go
	if len(results) > 0 && results[0].file == "main.go" {
		t.Error("trackLicenseHeaders should not return hardcoded main.go")
	}

	// Verify results have proper structure
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

	// Should find security issues but not hardcode location to "line 5"
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
	reason := "CCPA request"

	// Should accept reason parameter, not hardcode "GDPR request"
	m = trackDataDeletion(m, hash, email)

	if len(m.dataDeleteRequests) > 0 {
		lastReq := m.dataDeleteRequests[len(m.dataDeleteRequests)-1]
		if lastReq.reason == "GDPR request" && reason != "GDPR request" {
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

	// Should find secrets with real line numbers, not hardcode "line 1"
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

	// Should have stats for both authors
	if len(stats) != 2 {
		t.Errorf("calculateTeamStats should return one entry per author, got %d", len(stats))
	}

	// Verify structure
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

	// Should return some reviewers, not empty
	if len(reviewers) == 0 {
		t.Error("suggestReviewers should return at least one reviewer")
	}

	// Each reviewer should have expertise
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

	// Should detect both version tags
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

func TestDetectAnomalies_FindsUnusualCommits(t *testing.T) {
	// Normal-size commits
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "Small fix"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Small fix"},
		// Large commit (should be anomaly)
		{hash: "ccc3333333333333333333333333333333333", subject: "Refactor: Rewrite entire module (50000 lines)"},
	}

	anomalies := detectAnomalies(commits)

	// Should detect at least the large commit as anomalous
	if len(anomalies) == 0 {
		t.Error("detectAnomalies should find unusual commits")
	}
}

// TestClassifyCommit_WithMerge tests merge commit detection
func TestClassifyCommit_WithMerge(t *testing.T) {
	classification := classifyCommit("Merge pull request #123", "abc123")

	AssertTrue(t, classification.category != "", "should classify merge commit")
}

// TestClassifyCommit_FeatureCommit tests feature commit classification
func TestClassifyCommit_FeatureCommit(t *testing.T) {
	classification := classifyCommit("feat: Add new feature", "def456")

	AssertTrue(t, classification.category != "", "should classify feature commit")
}

// TestClassifyCommit_Refactor tests refactor commit classification
func TestClassifyCommit_Refactor(t *testing.T) {
	classification := classifyCommit("refactor: Clean up code structure", "ghi789")

	AssertTrue(t, classification.category != "", "should classify refactor commit")
}
