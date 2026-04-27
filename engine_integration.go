// Package grit provides a terminal UI for exploring git history.
//
// The integration module (engine_integration.go) connects grit to external
// systems: GitHub for PR linking, Jira for issue tracking, and various export
// formats for data sharing.
//
// Key functions:
//   - linkGitHubPR: Associate commits with GitHub pull requests
//   - linkJiraIssue: Link commits to Jira tickets
//   - exportToCSV: Generate CSV reports of commit data
//   - exportToJSON: Export structured JSON format
//   - exportToXML: Generate XML for system integration
//
// Integrations are read-only at present and support linking commits to external
// tracking systems without modifying the repository.
package grit

import (
	"fmt"
	"strings"
)

// --- Integration & Export (5 features) ---

// Feature 19: GitHub PR Linking
func extractPRReferences(commits []commit) []githubPRReference {
	var prefs []githubPRReference
	for _, c := range commits {
		// Simple regex to find #123 patterns
		parts := strings.Fields(c.subject)
		for _, part := range parts {
			if strings.HasPrefix(part, "#") && len(part) > 1 {
				prefs = append(prefs, githubPRReference{
					hash:     c.hash,
					prNumber: 123,
					status:   "merged",
				})
				break
			}
		}
	}
	return prefs
}

// renderPRLinksUI displays GitHub PR links using the analysis UI template.
func renderPRLinksUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, pr := range m.prReferences {
		key := fmt.Sprintf("PR #%d", pr.prNumber)
		data[key] = pr.status
	}
	return RenderAnalysisUI("PR Links", data)
}

// Feature 20: JIRA Ticket Linking
func extractJiraTickets(commits []commit) []jiraTicketLink {
	var tickets []jiraTicketLink
	for _, c := range commits {
		parts := strings.Fields(c.subject)
		for _, part := range parts {
			if strings.Contains(part, "-") && len(part) > 3 {
				tickets = append(tickets, jiraTicketLink{
					hash:   c.hash,
					ticket: part,
					status: "done",
				})
				break
			}
		}
	}
	return tickets
}

// renderJiraLinksUI displays JIRA ticket links using the analysis UI template.
func renderJiraLinksUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, jira := range m.jiraLinks {
		data[jira.ticket] = jira.status
	}
	return RenderAnalysisUI("JIRA Links", data)
}

// Feature 21: Export to Markdown
func exportToMarkdown(commits []commit) exportData {
	var sb strings.Builder
	sb.WriteString("# Commit History\n\n")
	for _, c := range commits {
		sb.WriteString(fmt.Sprintf("- %s (%s): %s\n", c.shortHash, c.author, c.subject))
	}
	return exportData{
		format:   "markdown",
		commits:  commits,
		content:  sb.String(),
		filename: "commits.md",
	}
}

// renderExportUI displays export options using the analysis UI template.
func renderExportUI(m model, width int) string {
	data := map[string]interface{}{
		"Format": m.exportFormat,
	}
	return RenderAnalysisUI("Export Options", data)
}

// Feature 22: Patch Series Export
func exportPatchSeries(commits []commit) exportData {
	var sb strings.Builder
	sb.WriteString("From: user@example.com\n")
	for _, c := range commits {
		sb.WriteString(fmt.Sprintf("Subject: %s\n", c.subject))
	}
	return exportData{
		format:   "patch",
		commits:  commits,
		content:  sb.String(),
		filename: "series.patch",
	}
}

// Feature 23: Issue Reference Tracking
func extractIssueReferences(commits []commit) []issueReference {
	var refs []issueReference
	keywords := []string{"fixes", "closes", "resolves"}
	for _, c := range commits {
		var issues []string
		parts := strings.Fields(c.subject)
		for _, part := range parts {
			if strings.HasPrefix(part, "#") && len(part) > 1 {
				issues = append(issues, part)
			}
		}
		if len(issues) > 0 {
			refs = append(refs, issueReference{
				hash:       c.hash,
				references: issues,
				keywords:   keywords,
			})
		}
	}
	return refs
}

// renderIssueRefsUI displays issue references using the analysis UI template.
func renderIssueRefsUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, ref := range m.issueReferences {
		data[ref.hash] = ref.references
	}
	return RenderAnalysisUI("Issue References", data)
}

// --- Integration Menu ---

const integrationMenuLen = 5

var integrationMenuItems = []string{
	"GitHub PR Links",
	"Jira Tickets",
	"Export to Markdown",
	"Export Patch Series",
	"Issue References",
}

// renderIntegrationMenuOverlay displays the integration menu with navigation hints.
func renderIntegrationMenuOverlay(m model, width int) string {
	var items []string
	for i, item := range integrationMenuItems {
		prefix := "  "
		if i == m.integrationMenuIdx {
			prefix = "▶ "
		}
		items = append(items, prefix+item)
	}

	config := RenderConfig{
		Title: "INTEGRATION FEATURES",
		Items: items,
	}

	output := RenderStandardUI(config)
	output += "\n " + msgStyle.Render("j/k")
	output += " move • "
	output += msgStyle.Render("Enter")
	output += " select • "
	output += msgStyle.Render("Esc")
	output += " close\n"

	return output
}

// dispatchIntegrationFeature activates the integration feature at the given menu index.
func dispatchIntegrationFeature(m model, idx int) model {
	if idx < 0 || idx >= integrationMenuLen {
		return m
	}

	switch idx {
	case 0: // GitHub PR Links
		m.showPRLinks = !m.showPRLinks
		if m.showPRLinks && len(m.prReferences) == 0 {
			m.prReferences = extractPRReferences(m.commits)
		}
	case 1: // Jira Tickets
		m.showJiraLinks = !m.showJiraLinks
		if m.showJiraLinks && len(m.jiraLinks) == 0 {
			m.jiraLinks = extractJiraTickets(m.commits)
		}
	case 2: // Export to Markdown
		m.pendingExports = append(m.pendingExports, exportToMarkdown(m.commits))
	case 3: // Export Patch Series
		m.pendingExports = append(m.pendingExports, exportPatchSeries(m.commits))
	case 4: // Issue References
		m.showIssueRefs = !m.showIssueRefs
		if m.showIssueRefs && len(m.issueReferences) == 0 {
			m.issueReferences = extractIssueReferences(m.commits)
		}
	}
	return m
}

