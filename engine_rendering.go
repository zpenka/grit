// Package grit provides a terminal UI for exploring git history.
//
// The main rendering module (engine_rendering.go) is the heart of the UI,
// responsible for converting the model state into terminal output. It renders
// the four-panel interface: commit list, diff view, file tree, and feature panels.
//
// Key rendering groups:
//   - Commit list: Shows filtered commits with author, time, stats badges
//   - Diff panel: Displays unified diff with syntax highlighting and line numbers
//   - File panel: Shows files changed with hunks and change statistics
//   - Feature panels: Results for analytics, bisect, team stats, visualizations
//
// The rendering engine is performance-critical. It leverages Lipgloss for
// terminal styling and the consolidated templates from engine_render_consolidation.go.
// Large diffs are cached to prevent re-rendering on each frame.
package grit

import (
	"fmt"
	"strings"
)

func renderStatsBadgeInList(stats commitStatistics, maxWidth int) string {
	badge := diffStatBadge(stats)
	if len(badge) > maxWidth {
		badge = truncate(badge, maxWidth)
	}
	return badge
}

// formatFilterHeaderDisplay formats active filters for header display.
func formatFilterHeaderDisplay(m model) string {
	return formatActiveFilters(m)
}

// renderBookmarkMarker returns a visual marker for bookmarked commits.
func renderBookmarkMarker(m model, idx int) string {
	if isBookmarked(m, idx) {
		return "★"
	}
	return ""
}


// renderLineCommentMarker returns a visual marker for commented lines.
func renderLineCommentMarker(m model, lineIdx int) string {
	if m.comments != nil {
		if _, ok := m.comments[lineIdx]; ok {
			return "●"
		}
	}
	return ""
}

// ===== OPTION 2: COMMIT GRAPH =====

// parseCommitGraph builds a graph structure from commits.
func parseCommitGraph(commits []commit) []graphNode {
	var nodes []graphNode
	for i, c := range commits {
		node := graphNode{
			hash:    c.hash,
			depth:   0,
			isMerge: false,
		}
		// Simple heuristic: if subject contains "Merge", mark as merge
		if strings.Contains(strings.ToLower(c.subject), "merge") {
			node.isMerge = true
		}
		// Depth is based on position (linear for now)
		node.depth = i % 2
		nodes = append(nodes, node)
	}
	return nodes
}

// detectBranches identifies branches in commit graph.
func detectBranches(commits []commit) []string {
	// Simple implementation: treat as single branch unless merges detected
	var branches []string
	hasMerge := false
	for _, c := range commits {
		if strings.Contains(strings.ToLower(c.subject), "merge") {
			hasMerge = true
			break
		}
	}
	if hasMerge {
		branches = append(branches, "main", "feature")
	} else {
		branches = append(branches, "main")
	}
	return branches
}

// renderAsciiGraph renders ASCII art graph for commit history.
// renderAsciiGraph displays commits as an ASCII graph using the standard UI template.
func renderAsciiGraph(graph []graphNode) string {
	if len(graph) == 0 {
		return ""
	}
	var items []string
	for _, node := range graph {
		prefix := "* "
		if node.isMerge {
			prefix = "*   "
		}
		hash := node.hash
		if len(hash) > 7 {
			hash = hash[:7]
		}
		items = append(items, fmt.Sprintf("%s%s", prefix, hash))
	}
	return RenderStandardUI(RenderConfig{
		Title: "Commit Graph",
		Items: items,
	})
}

// navigateAlongGraph moves along the graph in a direction.
func navigateAlongGraph(graph []graphNode, currentIdx int, direction string) int {
	if len(graph) == 0 {
		return 0
	}
	switch direction {
	case "down":
		if currentIdx < len(graph)-1 {
			return currentIdx + 1
		}
	case "up":
		if currentIdx > 0 {
			return currentIdx - 1
		}
	}
	return currentIdx
}

// getCommitRelationships maps parent-child relationships.
func getCommitRelationships(commits []commit) map[string][]string {
	rels := make(map[string][]string)
	// Infrastructure: would populate from git data
	return rels
}

// ===== OPTION 3: FILE-CENTRIC VIEW =====

// buildFileHistory constructs commit history for a specific file.
func buildFileHistory(commits []commit, file string) []commit {
	if file == "" {
		return []commit{}
	}
	// Infrastructure: would query git for file history
	return []commit{}
}

// renderFileTimeline renders the evolution of a file over time using consolidated rendering.
func renderFileTimeline(commits []commit, file string, width int) string {
	if len(commits) == 0 {
		return RenderErrorList("File Timeline", []string{})
	}
	var items []string
	for _, c := range commits {
		items = append(items, fmt.Sprintf("%s - %s", c.shortHash, c.subject))
	}
	return RenderStandardUI(RenderConfig{
		Title: fmt.Sprintf("File Timeline: %s", file),
		Items: items,
	})
}

// getFileBlameContext gets blame information for a file.
func getFileBlameContext(lines []diffLine, file string) map[int]string {
	ctx := make(map[int]string)
	// Infrastructure: would populate from git blame
	return ctx
}

// filterCommitsByFileChange filters commits that modified a specific file.
func filterCommitsByFileChange(commits []commit, file string) []commit {
	if file == "" {
		return commits
	}
	var result []commit
	for _, c := range commits {
		if isFileModifiedInCommit(c.hash, file) {
			result = append(result, c)
		}
	}
	return result
}

// isFileModifiedInCommit checks if a file was modified in a commit.
func isFileModifiedInCommit(hash, file string) bool {
	// Infrastructure: would query git
	return false
}

// ===== OPTION 4: STASH & REFLOG =====

// parseStashList parses git stash output.
func parseStashList(output string) []stashEntry {
	var stashes []stashEntry
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Format: "stash@{0}: WIP on main: abc1234 message"
		parts := strings.SplitN(line, ":", 3)
		if len(parts) >= 2 {
			name := strings.TrimSpace(parts[0])
			rest := strings.TrimSpace(parts[1])
			// Extract branch name
			branchParts := strings.Fields(rest)
			var branch string
			if len(branchParts) >= 3 {
				branch = branchParts[2]
			}
			stashes = append(stashes, stashEntry{
				name:    name,
				branch:  branch,
				subject: line,
			})
		}
	}
	return stashes
}

// parseReflog parses git reflog output.
func parseReflog(output string) []reflogEntry {
	var entries []reflogEntry
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Format: "abc1234 HEAD@{0}: commit: message"
		parts := strings.SplitN(line, " ", 3)
		if len(parts) >= 3 {
			hash := parts[0]
			rest := parts[2]
			actionParts := strings.SplitN(rest, ":", 2)
			action := ""
			message := ""
			if len(actionParts) >= 1 {
				action = actionParts[0]
			}
			if len(actionParts) >= 2 {
				message = strings.TrimSpace(actionParts[1])
			}
			entries = append(entries, reflogEntry{
				hash:    hash,
				action:  action,
				message: message,
			})
		}
	}
	return entries
}

// renderStashView renders the stash browser view.
// renderStashView renders the stash list using consolidated rendering.
func renderStashView(stashes []stashEntry, width int) string {
	var items []string
	for _, s := range stashes {
		items = append(items, fmt.Sprintf("%s - %s", s.name, s.branch))
	}

	config := RenderConfig{
		Title:       "Stashes",
		Items:       items,
		ShowIndices: true,
	}
	return RenderStandardUI(config)
}

// renderReflogView renders the reflog browser view using consolidated rendering.
func renderReflogView(entries []reflogEntry, width int) string {
	var items []string
	for _, e := range entries {
		items = append(items, fmt.Sprintf("%s - %s: %s", e.hash[:7], e.action, e.message))
	}

	config := RenderConfig{
		Title: "Reflog",
		Items: items,
	}
	return RenderStandardUI(config)
}

// stashToCommitLike converts a stash entry to a commit-like structure.
func stashToCommitLike(stash stashEntry) commit {
	return commit{
		shortHash: stash.hash,
		hash:      stash.hash,
		subject:   stash.subject,
		author:    "stash",
		when:      "stash",
	}
}

// reflogToCommitLike converts a reflog entry to a commit-like structure.
func reflogToCommitLike(entry reflogEntry) commit {
	return commit{
		shortHash: entry.hash[:7],
		hash:      entry.hash,
		subject:   entry.message,
		author:    entry.action,
		when:      entry.date,
	}
}

// switchViewMode switches between log, stash, and reflog views.
func switchViewMode(m model, newMode string) model {
	m.viewMode = newMode
	m.cursor = 0
	return m
}

// findStashByIndex finds a stash by its index.
func findStashByIndex(stashes []stashEntry, idx int) *stashEntry {
	if idx < 0 || idx >= len(stashes) {
		return nil
	}
	return &stashes[idx]
}

// ===== UI INTEGRATION: RENDERING =====

// renderBookmarkList renders list of bookmarked commits.
// renderBookmarkList renders bookmarked commits using consolidated rendering.
func renderBookmarkList(m model, width int) string {
	var items []string
	for _, hash := range m.bookmarks {
		for _, c := range m.commits {
			if c.shortHash == hash {
				items = append(items, fmt.Sprintf("%s - %s", hash, c.subject))
				break
			}
		}
	}

	config := RenderConfig{
		Title:       "Bookmarks",
		Items:       items,
		ShowIndices: true,
	}
	return RenderStandardUI(config)
}

// renderGraphView renders the commit graph.
func renderGraphView(m model, width int) string {
	if len(m.commitGraph) == 0 {
		return ""
	}
	return renderAsciiGraph(m.commitGraph)
}

// renderViewMode renders current view (log, stash, or reflog).
func renderViewMode(m model, width int) string {
	switch m.viewMode {
	case "stash":
		return renderStashView(m.stashes, width)
	case "reflog":
		return renderReflogView(m.reflogEntries, width)
	default:
		return ""
	}
}

// renderDiffWithComments renders diff with comment markers using the standard UI template.
func renderDiffWithComments(m model, panelHeight, width int) string {
	var items []string
	for i := 0; i < panelHeight && m.diffOffset+i < len(m.diffLines); i++ {
		marker := renderLineCommentMarker(m, m.diffOffset+i)
		text := m.diffLines[m.diffOffset+i].text
		if marker != "" {
			text = marker + " " + text
		}
		items = append(items, text)
	}
	return RenderStandardUI(RenderConfig{
		Title: "Diff View",
		Items: items,
	})
}

// enterCommentMode enters line comment mode.
func enterCommentMode(m model) model {
	m.inCommentMode = true
	m.commentInput = ""
	return m
}

// exitCommentMode exits line comment mode.
func exitCommentMode(m model) model {
	m.inCommentMode = false
	m.commentInput = ""
	return m
}
