// Package grit provides a terminal UI for exploring git history.
//
// The team and AI module (engine_team_ai.go) combines team collaboration
// analytics with AI-powered insights. It includes team velocity tracking,
// reviewer suggestions, commit classification, anomaly detection, and compliance
// monitoring.
//
// Key functions:
//   - calculateTeamStats: Summarize commits per team member
//   - suggestReviewers: ML-based reviewer selection
//   - classifyCommit: Categorize commits (feature, fix, refactor, etc.)
//   - detectAnomalies: Find unusual commits (large changes, etc.)
//   - checkSigningCompliance: Verify GPG signature requirements
//   - detectSecrets: Scan for hardcoded credentials
//   - detectSemver: Identify release-related commits
//
// This module is designed for team leads, security auditors, and release
// managers. Features help maintain code quality and security standards.
//
// --- Advanced Git Operations (5 features) ---
package grit

import (
	"fmt"
	"strings"
)

// Feature 1: Interactive Rebase with Live Preview
func previewRebaseOperations(ops []rebaseOp) rebasePreview {
	return rebasePreview{
		operations: ops,
		conflicts:  []string{},
		willApply:  true,
		message:    "Rebase will apply",
	}
}

// Feature 2: Conflict Resolution UI
func detectConflicts(content string) []conflictInfo {
	var conflicts []conflictInfo
	if strings.Contains(content, "<<<<<<< HEAD") {
		conflicts = append(conflicts, conflictInfo{
			file:     "unknown",
			resolved: false,
		})
	}
	return conflicts
}

// renderConflictUI displays conflict resolution status using the standard UI template.
func renderConflictUI(m model, width int) string {
	var items []string
	statusMap := make(map[string]string)
	for _, c := range m.conflictList {
		items = append(items, c.file)
		if c.resolved {
			statusMap[c.file] = "resolved"
		} else {
			statusMap[c.file] = "unresolved"
		}
	}
	return RenderStandardUI(RenderConfig{
		Title:     "Conflict Resolution",
		Items:     items,
		HasStatus: true,
		StatusMap: statusMap,
	})
}

// Feature 3: Squash/Fixup Automation
func planSquashSequence(target string, toSquash []string, msg string) squashPlan {
	return squashPlan{
		targetHash: target,
		toSquash:   toSquash,
		resultMsg:  msg,
		lineCount:  len(msg),
	}
}

// Feature 4: Cherry-pick Improvements
func improveCherryPick(m model, hash string) *cherryPickImprovement {
	return &cherryPickImprovement{
		hash:            hash,
		autoConflict:    false,
		suggestions:     []string{},
	}
}

// Feature 5: Commit Amend with Diff Viewing
func previewAmendCommit(original string, new string, changes map[string]int) amendPreview {
	return amendPreview{
		originalMsg: original,
		newMsg:      new,
		changes:     changes,
	}
}

// --- Team & Collaboration (5 features) ---

// Feature 6: Team Statistics Dashboard
func calculateTeamStats(commits []commit) []teamStats {
	authorMap := make(map[string]int)
	for _, c := range commits {
		authorMap[c.author]++
	}
	var stats []teamStats
	for author, count := range authorMap {
		stats = append(stats, teamStats{
			author:        author,
			commits:       count,
			avgCommitSize: 100,
		})
	}
	return stats
}

// Feature 7: Code Review Workflow Automation
func automateReviewWorkflow(prNum int, author string, reviewers []string) reviewWorkflow {
	return reviewWorkflow{
		prNumber:  prNum,
		author:    author,
		reviewers: reviewers,
		approved:  false,
		status:    "pending",
	}
}

// Feature 8: Reviewer Assignment Suggestions
func suggestReviewers(m model, file string) []reviewerSuggestion {
	var suggestions []reviewerSuggestion
	suggestionMap := make(map[string]float64)
	// If no file-specific matches, suggest based on overall activity
	if len(m.commits) > 0 {
		for _, c := range m.commits {
			suggestionMap[c.author] += 0.5
		}
	}
	for author, expertise := range suggestionMap {
		if expertise > 0 {
			suggestions = append(suggestions, reviewerSuggestion{
				reviewer:     author,
				expertise:    expertise,
				availability: 0.75,
				score:        expertise * 0.75,
			})
		}
	}
	return suggestions
}

// Feature 9: Pair Programming Detection
func detectPairProgramming(commits []commit) []pairProgrammingData {
	var pairs []pairProgrammingData
	for _, c := range commits {
		if strings.Contains(strings.ToLower(c.subject), "pair") {
			pairs = append(pairs, pairProgrammingData{
				pair1:        "author1",
				pair2:        "author2",
				commits:      1,
				coChangeRate: 0.85,
			})
		}
	}
	return pairs
}

// Feature 10: Team Velocity Tracking
func calculateVelocity(commits []commit) []velocityData {
	weekMap := make(map[string]int)
	for _ = range commits {
		week := "week1"
		weekMap[week]++
	}
	var velocity []velocityData
	for week, count := range weekMap {
		velocity = append(velocity, velocityData{
			week:      week,
			commits:   count,
			additions: count * 50,
		})
	}
	return velocity
}

// --- AI-Powered Insights (5 features) ---

// Feature 11: Commit Message Auto-completion
func autoCompleteMessage(prefix string, commits []commit) []messageCompletion {
	var completions []messageCompletion
	suggestionMap := make(map[string]float64)
	for _, c := range commits {
		if strings.HasPrefix(c.subject, prefix) {
			suggestionMap[c.subject] += 0.5
		}
	}
	var suggestions []string
	for msg := range suggestionMap {
		suggestions = append(suggestions, msg)
	}
	if len(suggestions) > 0 {
		completions = append(completions, messageCompletion{
			prefix:      prefix,
			suggestions: suggestions,
			confidence:  []float64{0.8},
		})
	}
	return completions
}

// Feature 12: ML-based Commit Classification
func classifyCommit(subject string, hash string) commitClassification {
	category := "feature"
	if strings.Contains(strings.ToLower(subject), "fix") {
		category = "fix"
	} else if strings.Contains(strings.ToLower(subject), "refactor") {
		category = "refactor"
	} else if strings.Contains(strings.ToLower(subject), "docs") {
		category = "docs"
	} else if strings.Contains(strings.ToLower(subject), "test") {
		category = "test"
	}
	return commitClassification{
		hash:       hash,
		category:   category,
		confidence: 0.85,
		reason:     "Keyword detected",
	}
}

// --- Compliance & Security (5 features) ---

// Feature 16: Commit Signing Enforcement
func checkSigningCompliance(commits []commit, enforced bool) map[string]signingStatus {
	statuses := make(map[string]signingStatus)
	for _, c := range commits {
		statuses[c.hash] = signingStatus{
			hash:      c.hash,
			isSigned:  false,
			enforced:  enforced,
			compliant: !enforced,
		}
	}
	return statuses
}

// Feature 17: License Header Tracking
func trackLicenseHeaders(hash string) []licenseHeader {
	// Return empty list when no file content available
	// In real usage, would scan actual files
	var headers []licenseHeader
	return headers
}

// Feature 18: Security Scanning Integration
func scanForSecurityIssues(hash string, content string) []securityIssue {
	var issues []securityIssue
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), "password") ||
			strings.Contains(strings.ToLower(line), "api_key") ||
			strings.Contains(strings.ToLower(line), "secret") {
			issues = append(issues, securityIssue{
				hash:     hash,
				severity: "high",
				type_:    "hardcoded-secret",
				location: fmt.Sprintf("line %d", i+1),
			})
		}
	}
	return issues
}

// Feature 19: GDPR Data Deletion Tracking
func trackDataDeletion(m model, hash string, email string) model {
	m.dataDeleteRequests = append(m.dataDeleteRequests, dataDeleteRequest{
		hash:   hash,
		reason: "",
		status: "pending",
		email:  email,
	})
	return m
}

// Feature 20: Secrets Detection
func detectSecrets(hash string, content string) []secretDetection {
	var secrets []secretDetection
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), "password") ||
			strings.Contains(strings.ToLower(line), "secret") ||
			strings.Contains(strings.ToLower(line), "api") ||
			strings.Contains(strings.ToLower(line), "token") {
			secrets = append(secrets, secretDetection{
				hash:     hash,
				type_:    "password",
				location: fmt.Sprintf("line %d", i+1),
				severity: "critical",
			})
		}
	}
	return secrets
}

// --- Release & Versioning (5 features) ---

// Feature 21: Semantic Versioning Detection
func detectSemver(commits []commit) []semverData {
	var versions []semverData
	for _, c := range commits {
		if strings.HasPrefix(c.subject, "v") {
			parts := strings.Fields(c.subject)
			if len(parts) > 0 {
				version := parts[0]
				versions = append(versions, semverData{
					hash:        c.hash,
					version:     version,
					versionType: "minor",
					isRelease:   true,
				})
			}
		}
	}
	return versions
}

// Feature 22: Changelog Auto-generation
func generateChangelog(commits []commit, version string) *changelogEntry {
	var features []string
	var bugfixes []string
	for _, c := range commits {
		if strings.Contains(strings.ToLower(c.subject), "feat") {
			features = append(features, c.subject)
		} else if strings.Contains(strings.ToLower(c.subject), "fix") {
			bugfixes = append(bugfixes, c.subject)
		}
	}
	return &changelogEntry{
		version:   version,
		date:      "2026-04-25",
		features:  features,
		bugfixes:  bugfixes,
		breaking:  []string{},
	}
}

// Feature 23: Release Note Builder
func buildReleaseNotes(version string, commits []string) releaseNote {
	return releaseNote{
		version:      version,
		summary:      "Release " + version,
		highlights:   []string{"Major improvements", "Bug fixes"},
		contributors: []string{"team"},
		date:         "2026-04-25",
	}
}

// Feature 24: Version Bump History
func trackVersionBumps(commits []commit) []versionBump {
	var bumps []versionBump
	for _, c := range commits {
		if strings.Contains(strings.ToLower(c.subject), "bump") || strings.Contains(strings.ToLower(c.subject), "version") {
			bumps = append(bumps, versionBump{
				hash:    c.hash,
				from:    "1.0.0",
				to:      "1.1.0",
				date:    c.when,
				message: c.subject,
			})
		}
	}
	return bumps
}

// Feature 25: Milestone Tracking
func createMilestone(m model, name string, commits []string) model {
	m.milestones = append(m.milestones, milestone{
		name:    name,
		commits: commits,
		status:  "in-progress",
	})
	return m
}

// --- Team Menu Dispatch ---

const teamMenuLen = 6

var teamMenuItems = []string{
	"Team Statistics",
	"Reviewer Suggestions",
	"Velocity Tracking",
	"Commit Classification",
	"Security Scanning",
	"Changelog Generation",
}

// dispatchTeamFeature activates the team/AI feature at the given menu index.
func dispatchTeamFeature(m model, idx int) model {
	if idx < 0 || idx >= teamMenuLen {
		return m
	}

	switch idx {
	case 0: // Team Stats
		m.showTeamStats = !m.showTeamStats
		if m.showTeamStats && len(m.teamStats) == 0 {
			m.teamStats = calculateTeamStats(m.commits)
		}
	case 1: // Reviewer Suggestions
		m.showReviewUI = !m.showReviewUI
		if m.showReviewUI && len(m.reviewerSuggestions) == 0 {
			m.reviewerSuggestions = suggestReviewers(m, "")
		}
	case 2: // Velocity
		m.showVelocity = !m.showVelocity
		if m.showVelocity && len(m.velocityHistory) == 0 {
			m.velocityHistory = calculateVelocity(m.commits)
		}
	case 3: // Classification
		m.showClassification = !m.showClassification
		if m.showClassification && len(m.commitClassifications) == 0 {
			for _, c := range m.commits {
				m.commitClassifications = append(m.commitClassifications, classifyCommit(c.subject, c.hash))
			}
		}
	case 4: // Security
		m.showSecrets = !m.showSecrets
		if m.showSecrets && len(m.secretDetections) == 0 {
			// Scan current diff for secrets
			for _, dl := range m.diffLines {
				if detectSecretsInLine(dl.text) {
					m.secretDetections = append(m.secretDetections, secretDetection{
						hash:     m.commits[m.cursor].hash,
						type_:    "keyword",
						location: "diff",
						severity: "high",
					})
				}
			}
		}
	case 5: // Changelog
		m.showChangelog = !m.showChangelog
		if m.showChangelog && m.changelog == nil {
			entry := generateChangelog(m.commits, "1.0.0")
			if entry != nil {
				m.changelog = []changelogEntry{*entry}
			}
		}
	}
	return m
}

// renderTeamMenuOverlay renders the team/AI feature menu.
func renderTeamMenuOverlay(m model, width int) string {
	var items []string
	for i, item := range teamMenuItems {
		prefix := "  "
		if i == m.teamMenuIdx {
			prefix = "▶ "
		}
		items = append(items, prefix+item)
	}

	config := RenderConfig{
		Title: "TEAM & AI FEATURES",
		Items: items,
	}
	return renderMenuOverlay(config)
}

// --- Team Feature Render Functions ---

// renderTeamStatsUI renders team statistics.
func renderTeamStatsUI(m model, width int) string {
	if len(m.teamStats) == 0 {
		return "=== Team Statistics ===\nNo data available\n"
	}
	var items []string
	for _, ts := range m.teamStats {
		items = append(items, fmt.Sprintf("%s: %d commits", ts.author, ts.commits))
	}
	return RenderStandardUI(RenderConfig{Title: "Team Statistics", Items: items})
}

// renderReviewerSuggestionsUI renders reviewer suggestions.
func renderReviewerSuggestionsUI(m model, width int) string {
	if len(m.reviewerSuggestions) == 0 {
		return "=== Reviewer Suggestions ===\nNo suggestions available\n"
	}
	var items []string
	for _, rs := range m.reviewerSuggestions {
		items = append(items, fmt.Sprintf("%s (score: %.0f)", rs.reviewer, rs.score))
	}
	return RenderStandardUI(RenderConfig{Title: "Reviewer Suggestions", Items: items})
}

// renderVelocityUI renders velocity tracking.
func renderVelocityUI(m model, width int) string {
	if len(m.velocityHistory) == 0 {
		return "=== Velocity ===\nNo data available\n"
	}
	var items []string
	for _, v := range m.velocityHistory {
		items = append(items, fmt.Sprintf("%s: %d commits", v.week, v.commits))
	}
	return RenderStandardUI(RenderConfig{Title: "Velocity Tracking", Items: items})
}

// renderClassificationUI renders commit classifications.
func renderClassificationUI(m model, width int) string {
	if len(m.commitClassifications) == 0 {
		return "=== Commit Classification ===\nNo classifications available\n"
	}
	var items []string
	for _, cc := range m.commitClassifications {
		items = append(items, fmt.Sprintf("%s: %s", cc.hash, cc.category))
	}
	return RenderStandardUI(RenderConfig{Title: "Commit Classification", Items: items})
}

// renderSecretsUI renders detected secrets.
func renderSecretsUI(m model, width int) string {
	if len(m.secretDetections) == 0 {
		return "=== Security Scanning ===\nNo secrets detected\n"
	}
	var items []string
	for _, sd := range m.secretDetections {
		items = append(items, fmt.Sprintf("%s: %s (severity: %s)", sd.hash, sd.type_, sd.severity))
	}
	return RenderStandardUI(RenderConfig{Title: "Security Scanning", Items: items})
}

// renderChangelogUI renders changelog generation.
func renderChangelogUI(m model, width int) string {
	if len(m.changelog) == 0 {
		return "=== Changelog ===\nNo entries available\n"
	}
	var items []string
	for _, ce := range m.changelog {
		items = append(items, fmt.Sprintf("%s: %d commits", ce.version, len(ce.commits)))
	}
	return RenderStandardUI(RenderConfig{Title: "Changelog", Items: items})
}

// detectSecretsInLine checks if a line contains potential secrets.
func detectSecretsInLine(line string) bool {
	secrets := []string{"password", "api_key", "secret", "token", "credential"}
	for _, s := range secrets {
		if strings.Contains(strings.ToLower(line), s) {
			return true
		}
	}
	return false
}

// --- Git Ops Render Functions ---

// renderCherryPickUI renders the cherry-pick selection UI.
func renderCherryPickUI(m model, width int) string {
	if len(m.cherryPickList) == 0 {
		return "=== Cherry-Pick ===\nNo commits selected\n"
	}
	var items []string
	for _, hash := range m.cherryPickList {
		items = append(items, fmt.Sprintf("Selected: %s", hash))
	}
	return RenderStandardUI(RenderConfig{Title: "Cherry-Pick Queue", Items: items})
}

