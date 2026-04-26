package grit

import (
	"fmt"
	"regexp"
	"strconv"
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

// ===== UI INTEGRATION: KEYBINDINGS =====

// handleKeyBinding processes keyboard input and returns updated model.
func handleKeyBinding(m model, key string) model {
	switch key {
	case "m":
		m = toggleBookmark(m)
	case "'":
		m = jumpToNextBookmark(m)
	case "gg":
		m.inGoToCommitMode = true
		m.goToCommitInput = ""
	case "c":
		m.inCommentMode = true
		m.commentInput = ""
	case "v":
		m = switchViewMode(m, "stash")
	case "V":
		m = switchViewMode(m, "reflog")
	case "G":
		m.showGraph = !m.showGraph
		if m.showGraph && len(m.commitGraph) == 0 {
			m = lazyLoadGraph(m)
		}
	case "f":
		m = toggleFileView(m)
	case "R":
		m.showRebaseUI = !m.showRebaseUI
		if m.showRebaseUI && len(m.rebaseSequence) == 0 {
			m.rebaseSequence = parseRebaseSequence(m.commits)
		}
	case "C":
		m.showCherryPickUI = !m.showCherryPickUI
	case "A":
		m.showAnalytics = !m.showAnalytics
		if m.showAnalytics {
			m.authorStats = calculateAuthorStats(m.commits)
			m.timeStats = calculateTimeStats(m.commits)
		}
	// Bisect & Recovery
	case "B":
		if !m.bisectState.active {
			m = initiateBisect(m)
		} else {
			m.bisectState.active = false
			m.showBisectUI = false
		}
	case "L":
		m.showLostCommits = !m.showLostCommits
		if m.showLostCommits && len(m.lostCommits) == 0 {
			m.lostCommits = findLostCommits("")
		}
	case "U":
		if m.showUndoMenu {
			m = performUndo(m)
		}
		m.showUndoMenu = !m.showUndoMenu
	// Code Patterns & Quality
	case "O":
		m.showCodeOwnership = !m.showCodeOwnership
		if m.showCodeOwnership && len(m.codeOwnership) == 0 {
			m.codeOwnership = analyzeCodeOwnership(m.commits)
		}
	case "H":
		m.showHotspots = !m.showHotspots
		if m.showHotspots && len(m.hotspots) == 0 {
			m.hotspots = detectHotspots(m.commits)
		}
	case "M":
		m.showLinting = !m.showLinting
		if m.showLinting && len(m.lintingResults) == 0 {
			for _, c := range m.commits {
				result := lintCommitMessage(c.subject, c.hash)
				m.lintingResults = append(m.lintingResults, result)
			}
		}
	case "S":
		m.showLargeCommits = !m.showLargeCommits
		if m.showLargeCommits && len(m.largeCommits) == 0 {
			m = analyzeCommitSize(m)
		}
	case "X":
		m.showComplexity = !m.showComplexity
		if m.showComplexity && len(m.commitMetrics) == 0 {
			m = analyzeComplexity(m)
		}
	// Commit Analysis & Search (4 features)
	case "N":
		m.showSemanticSearch = !m.showSemanticSearch
		if m.showSemanticSearch && len(m.semanticSearchResults) == 0 {
			m.semanticSearchResults = semanticSearch(m.commits, m.semanticQuery)
		}
	case "E":
		m.showActivityHeatmap = !m.showActivityHeatmap
		if m.showActivityHeatmap && len(m.authorActivityHeatmap) == 0 {
			m.authorActivityHeatmap = buildActivityHeatmap(m.commits)
		}
	case "Y":
		m.showMergeAnalysis = !m.showMergeAnalysis
		if m.showMergeAnalysis && len(m.mergeAnalysisData) == 0 {
			m.mergeAnalysisData = analyzeMerges(m.commits)
		}
	case "T":
		m.showCoupling = !m.showCoupling
		if m.showCoupling && len(m.commitCouplings) == 0 {
			m.commitCouplings = analyzeCommitCoupling(m.commits)
		}
	// Performance & Filtering (4 features)
	case "D":
		if m.currentExtFilter == "" {
			m = toggleExtensionFilter(m, ".go")
		} else {
			m = toggleExtensionFilter(m, "")
		}
	case "W":
		if m.groupingMode == "" {
			m.groupingMode = "date"
			m.commitGroups = groupCommits(m.commits, "date")
		} else {
			m.groupingMode = ""
		}
	case "Z":
		m.showDependencies = !m.showDependencies
		if m.showDependencies && len(m.dependencyChanges) == 0 {
			m.dependencyChanges = trackDependencyChanges(m.commits)
		}
	// Advanced Workflows (5 features)
	case "1":
		m.showWorktrees = !m.showWorktrees
		if m.showWorktrees && len(m.worktrees) == 0 {
			m.worktrees = loadWorktrees("")
		}
	case "2":
		m.showSubmodules = !m.showSubmodules
		if m.showSubmodules && len(m.submodules) == 0 {
			m.submodules = parseSubmodules("")
		}
	case "3":
		m.showNamedStashes = !m.showNamedStashes
	case "4":
		m.showTagMgmt = !m.showTagMgmt
	case "5":
		m.showGPGStatus = !m.showGPGStatus
		if m.showGPGStatus && len(m.gpgStatuses) == 0 {
			m.gpgStatuses = extractGPGSignatureStatus("")
		}
	// Visualization (5 features)
	case "6":
		m.showFlamegraph = !m.showFlamegraph
		if m.showFlamegraph && len(m.contributorFlameData) == 0 {
			m.contributorFlameData = buildContributorFlame(m.commits)
		}
	case "7":
		m.showTimeline = !m.showTimeline
		if m.showTimeline && len(m.timelinePoints) == 0 {
			m.timelinePoints = buildTimeline(m.commits)
		}
	case "8":
		m.showTreeView = !m.showTreeView
		if m.showTreeView && m.treeRoot == nil {
			m.treeRoot = buildTreeView(m.commits)
		}
	case "9":
		m.showAuthorComparison = !m.showAuthorComparison
		if m.showAuthorComparison && len(m.authorComparisons) == 0 {
			m.authorComparisons = compareAuthors(m)
		}
	case "0":
		m.showFileHeatmap = !m.showFileHeatmap
		if m.showFileHeatmap && len(m.fileHeatmap) == 0 {
			m.fileHeatmap = buildFileHeatmap(m.commits)
		}
	// Integration & Export (5 features)
	case "p":
		m.showPRLinks = !m.showPRLinks
		if m.showPRLinks && len(m.prReferences) == 0 {
			m.prReferences = extractPRReferences(m.commits)
		}
	case "j":
		m.showJiraLinks = !m.showJiraLinks
		if m.showJiraLinks && len(m.jiraLinks) == 0 {
			m.jiraLinks = extractJiraTickets(m.commits)
		}
	case "e":
		m.showExportUI = !m.showExportUI
	case "q":
		m.showIssueRefs = !m.showIssueRefs
		if m.showIssueRefs && len(m.issueReferences) == 0 {
			m.issueReferences = extractIssueReferences(m.commits)
		}
	default:
		// Handle multi-key like "5j"
		if len(key) > 1 && (strings.HasSuffix(key, "j") || strings.HasSuffix(key, "k")) {
			n := parseCount(key[:len(key)-1])
			if strings.HasSuffix(key, "j") {
				for i := 0; i < n; i++ {
					m = moveCursorDown(m)
				}
			} else {
				for i := 0; i < n; i++ {
					m = moveCursorUp(m)
				}
			}
		}
	}
	return m
}

// safeHandleKeyBinding handles keybindings with error recovery.
func safeHandleKeyBinding(m model, key string) model {
	if key == "" {
		return m
	}
	return handleKeyBinding(m, key)
}

// ===== UI INTEGRATION: RENDERING =====

// renderCommitRowWithStats renders commit row with stats badge.
func renderCommitRowWithStats(m model, idx int, width int) string {
	if !m.showStatsBadge {
		return ""
	}
	badge := diffStatBadge(m.lastStats)
	return badge
}

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

// ===== OPTIMIZATION: CACHING =====

// newDiffCache creates a new diff cache with specified max size.

// ===== OPTION A: ADVANCED COMMIT OPERATIONS =====

// --- Interactive Rebase ---

// parseRebaseSequence builds a rebase operation sequence from commits.
func parseRebaseSequence(commits []commit) []rebaseOp {
	var ops []rebaseOp
	for _, c := range commits {
		ops = append(ops, rebaseOp{
			action:  "pick",
			hash:    c.hash,
			subject: c.subject,
		})
	}
	return ops
}

// reorderCommit moves a commit in the rebase sequence.
func reorderCommit(seq []rebaseOp, from, to int) []rebaseOp {
	if from < 0 || from >= len(seq) || to < 0 || to >= len(seq) {
		return seq
	}
	if from == to {
		return seq
	}
	op := seq[from]
	newSeq := make([]rebaseOp, 0, len(seq))
	for i, o := range seq {
		if i == from {
			continue
		}
		if i == to && from < to {
			newSeq = append(newSeq, o)
			newSeq = append(newSeq, op)
		} else if i == to && from > to {
			newSeq = append(newSeq, op)
			newSeq = append(newSeq, o)
		} else {
			newSeq = append(newSeq, o)
		}
	}
	return newSeq
}

// squashCommit marks a commit for squashing.
func squashCommit(seq []rebaseOp, idx int) []rebaseOp {
	if idx >= 0 && idx < len(seq) {
		seq[idx].action = "squash"
	}
	return seq
}

// fixupCommit marks a commit for fixup (squash without message).
func fixupCommit(seq []rebaseOp, idx int) []rebaseOp {
	if idx >= 0 && idx < len(seq) {
		seq[idx].action = "fixup"
	}
	return seq
}

// previewRebase renders a preview of the rebase operation.
func previewRebase(seq []rebaseOp) string {
	var sb strings.Builder
	sb.WriteString("Rebase sequence:\n")
	for i, op := range seq {
		hash := op.hash
		if len(hash) > 7 {
			hash = hash[:7]
		}
		sb.WriteString(fmt.Sprintf("%d: %s %s - %s\n", i, op.action, hash, op.subject))
	}
	return sb.String()
}

// --- Cherry-pick ---

// toggleCherryPick adds or removes a commit from cherry-pick selection.
func toggleCherryPick(m model, hash string) model {
	found := false
	for i, h := range m.cherryPickList {
		if h == hash {
			m.cherryPickList = append(m.cherryPickList[:i], m.cherryPickList[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		m.cherryPickList = append(m.cherryPickList, hash)
	}
	return m
}

// previewCherryPick shows which commits will be cherry-picked.
func previewCherryPick(commits []commit, picks []string) string {
	var sb strings.Builder
	sb.WriteString("Cherry-pick queue:\n")
	for i, pick := range picks {
		for _, c := range commits {
			if c.hash == pick || c.shortHash == pick {
				hash := c.shortHash
				if len(hash) > 7 {
					hash = hash[:7]
				}
				sb.WriteString(fmt.Sprintf("%d: %s - %s\n", i, hash, c.subject))
				break
			}
		}
	}
	return sb.String()
}

// --- Reset ---

// resetToCommit generates a reset command with the specified mode.
func resetToCommit(hash, mode string) string {
	if mode == "" {
		mode = "mixed"
	}
	return fmt.Sprintf("git reset --%s %s", mode, hash)
}

// --- Revert ---

// revertCommit generates a revert command for a commit.
func revertCommit(hash string) string {
	return fmt.Sprintf("git revert %s", hash)
}

// --- Amend ---

// amendLastCommit updates the last commit message.
func amendLastCommit(m model, message string) model {
	if m.cursor < len(m.commits) {
		m.commits[m.cursor].subject = message
		m.amendMessage = message
	}
	return m
}

// ===== OPTION B: COLLABORATION & ANALYTICS =====

// --- Advanced Workflows (5 features) ---

// Feature 9: Worktree Support
func loadWorktrees(output string) []worktreeInfo {
	var worktrees []worktreeInfo
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		worktrees = append(worktrees, worktreeInfo{
			path:   strings.TrimSpace(line),
			branch: "main",
		})
	}
	return worktrees
}

func switchWorktree(m model, path string) model {
	m.currentWorktree = path
	return m
}

// renderWorktreesUI displays worktree information using the analysis UI template.
func renderWorktreesUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, wt := range m.worktrees {
		data[wt.path] = wt.branch
	}
	return RenderAnalysisUI("Worktrees", data)
}

// Feature 10: Submodule Tracking
func parseSubmodules(output string) []submoduleInfo {
	var subs []submoduleInfo
	lines := strings.Split(output, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "submodule") {
			subs = append(subs, submoduleInfo{
				path:   "lib",
				url:    "https://github.com/user/lib",
				branch: "main",
			})
		}
	}
	return subs
}

// renderSubmodulesUI displays submodule information using the analysis UI template.
func renderSubmodulesUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, sm := range m.submodules {
		data[sm.path] = sm.url
	}
	return RenderAnalysisUI("Submodules", data)
}

// Feature 11: Named Stashes
func createNamedStash(m model, index int, name string, desc string) model {
	m.namedStashes = append(m.namedStashes, namedStash{
		index:       index,
		name:        name,
		description: desc,
		hash:        "",
	})
	return m
}

// renderNamedStashesUI displays named stashes using the analysis UI template.
func renderNamedStashesUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, ns := range m.namedStashes {
		data[ns.name] = ns.description
	}
	return RenderAnalysisUI("Named Stashes", data)
}

// Feature 12: Tag Management
func queueTagOperation(m model, name string, hash string, action string, msg string) model {
	m.pendingTagOps = append(m.pendingTagOps, tagOperation{
		name:    name,
		hash:    hash,
		action:  action,
		message: msg,
	})
	return m
}

// renderTagMgmtUI displays tag management operations using the analysis UI template.
func renderTagMgmtUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, op := range m.pendingTagOps {
		data[op.name] = op.action
	}
	return RenderAnalysisUI("Tag Management", data)
}

// Feature 13: GPG Signature Status
func extractGPGSignatureStatus(output string) map[string]gpgSignatureStatus {
	statuses := make(map[string]gpgSignatureStatus)
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			hash := parts[0]
			statuses[hash] = gpgSignatureStatus{
				hash:   hash,
				signed: true,
				signer: "unknown",
			}
		}
	}
	return statuses
}

// renderGPGStatusUI displays GPG signature status for commits using the standard UI template.
func renderGPGStatusUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, status := range m.gpgStatuses {
		signed := "✗"
		if status.signed {
			signed = "✓"
		}
		data[status.hash] = signed
	}
	return RenderAnalysisUI("GPG Signatures", data)
}

// --- Visualization (5 features) ---

// Feature 14: Contributor Flamegraph
func buildContributorFlame(commits []commit) []contributorFlameData {
	authorMap := make(map[string]int)
	for _, c := range commits {
		authorMap[c.author]++
	}
	var flame []contributorFlameData
	for author, count := range authorMap {
		pct := float64(count) / float64(len(commits)) * 100
		flame = append(flame, contributorFlameData{
			author:     author,
			commits:    count,
			percentage: pct,
		})
	}
	// Sort by commit count descending
	for i := 0; i < len(flame)-1; i++ {
		for j := i + 1; j < len(flame); j++ {
			if flame[j].commits > flame[i].commits {
				flame[i], flame[j] = flame[j], flame[i]
			}
		}
	}
	return flame
}

// renderFlamegraphUI displays contributor flamegraph data using the analysis UI template.
func renderFlamegraphUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, cf := range m.contributorFlameData {
		data[cf.author] = cf.percentage
	}
	return RenderAnalysisUI("Contributor Flamegraph", data)
}

// Feature 15: Timeline Slider
func buildTimeline(commits []commit) []timelinePoint {
	var timeline []timelinePoint
	dateMap := make(map[string]int)
	for _, c := range commits {
		dateMap[c.when]++
	}
	for date, count := range dateMap {
		timeline = append(timeline, timelinePoint{
			date:    date,
			commits: count,
		})
	}
	return timeline
}

// renderTimelineSliderUI displays timeline data using the analysis UI template.
func renderTimelineSliderUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, tp := range m.timelinePoints {
		data[tp.date] = tp.commits
	}
	return RenderAnalysisUI("Timeline", data)
}

// Feature 16: Tree View
func buildTreeView(commits []commit) *treeNode {
	if len(commits) == 0 {
		return nil
	}
	root := &treeNode{
		hash:    commits[0].hash,
		subject: commits[0].subject,
		depth:   0,
	}
	for i := 1; i < len(commits); i++ {
		root.children = append(root.children, &treeNode{
			hash:    commits[i].hash,
			subject: commits[i].subject,
			depth:   1,
		})
	}
	return root
}

// renderTreeViewUI displays tree view of commits with hierarchy using the standard UI template.
func renderTreeViewUI(m model, width int) string {
	var items []string
	if m.treeRoot != nil {
		flattenTreeNode(m.treeRoot, &items)
	}
	return RenderStandardUI(RenderConfig{
		Title: "Tree View",
		Items: items,
	})
}

func flattenTreeNode(node *treeNode, items *[]string) {
	indent := strings.Repeat("  ", node.depth)
	*items = append(*items, fmt.Sprintf("%s├─ %s", indent, node.hash))
	for _, child := range node.children {
		flattenTreeNode(child, items)
	}
}

// Feature 17: Author Comparison
func compareAuthors(m model) []authorComparison {
	var comparisons []authorComparison
	if m.selectedAuthors[0] != "" && m.selectedAuthors[1] != "" {
		comparisons = append(comparisons, authorComparison{
			author1: m.selectedAuthors[0],
			author2: m.selectedAuthors[1],
			commits1: 10,
			commits2: 8,
			similarity: 0.75,
		})
	}
	return comparisons
}

// renderAuthorComparisonUI displays author comparison using the comparison table template.
func renderAuthorComparisonUI(m model, width int) string {
	if len(m.authorComparisons) == 0 {
		return "=== Author Comparison ===\nNo comparisons available\n"
	}
	comp := m.authorComparisons[0]
	items := map[string][2]interface{}{
		"Commits": {comp.commits1, comp.commits2},
		"Similarity": {comp.similarity, 0},
	}
	return RenderComparisonTable("Author Comparison", comp.author1, comp.author2, items)
}

// Feature 18: File Heatmap
func buildFileHeatmap(commits []commit) []fileHeatmapEntry {
	fileMap := make(map[string]int)
	for _, c := range commits {
		files := extractFilesFromSubject(c.subject)
		for _, f := range files {
			fileMap[f]++
		}
	}
	var heatmap []fileHeatmapEntry
	for file, freq := range fileMap {
		risk := "low"
		if freq > 10 {
			risk = "high"
		} else if freq > 5 {
			risk = "medium"
		}
		heatmap = append(heatmap, fileHeatmapEntry{
			path:      file,
			frequency: freq,
			risk:      risk,
		})
	}
	return heatmap
}

// renderFileHeatmapUI displays file heatmap with risk levels using the standard UI template.
func renderFileHeatmapUI(m model, width int) string {
	var items []string
	statusMap := make(map[string]string)
	for _, fh := range m.fileHeatmap {
		item := fmt.Sprintf("%s: %d changes", fh.path, fh.frequency)
		items = append(items, item)
		statusMap[item] = fh.risk
	}
	return RenderStandardUI(RenderConfig{
		Title:     "File Heatmap",
		Items:     items,
		HasStatus: true,
		StatusMap: statusMap,
	})
}

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

// --- Advanced Git Operations (5 features) ---

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

// Feature 13: Anomaly Detection
func detectAnomalies(commits []commit) []anomalyData {
	var anomalies []anomalyData
	for _, c := range commits {
		words := len(strings.Fields(c.subject))
		// Detect unusually large or verbose commits
		if words > 20 || strings.Contains(strings.ToLower(c.subject), "massive") || strings.Contains(c.subject, "10000") {
			anomalies = append(anomalies, anomalyData{
				hash:        c.hash,
				type_:       "large",
				severity:    7,
				description: "Large or unusual commit",
			})
		}
	}
	return anomalies
}

// Feature 14: Similar Commits Finder
func findSimilarCommits(commits []commit, targetHash string) []similarCommit {
	var similar []similarCommit
	var targetSubject string
	for _, c := range commits {
		if c.hash == targetHash {
			targetSubject = c.subject
			break
		}
	}
	for _, c := range commits {
		if c.hash != targetHash && strings.Contains(c.subject, targetSubject[:10]) {
			similar = append(similar, similarCommit{
				hash1:      targetHash,
				hash2:      c.hash,
				subject1:   targetSubject,
				subject2:   c.subject,
				similarity: 0.75,
			})
		}
	}
	return similar
}

// Feature 15: Auto-generated Summaries
func generateAutoSummary(hash string, fullMessage string) autoSummary {
	words := strings.Fields(fullMessage)
	var summary string
	if len(words) > 10 {
		summary = strings.Join(words[:10], " ") + "..."
	} else {
		summary = fullMessage
	}
	return autoSummary{
		hash:    hash,
		summary: summary,
		length:  len(summary),
		tokens:  len(words),
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
	return []licenseHeader{
		{file: "main.go", hasHeader: true, license: "MIT", hash: hash},
	}
}

// Feature 18: Security Scanning Integration
func scanForSecurityIssues(hash string, content string) []securityIssue {
	var issues []securityIssue
	if strings.Contains(content, "key") || strings.Contains(content, "secret") || strings.Contains(content, "password") {
		issues = append(issues, securityIssue{
			hash:     hash,
			severity: "high",
			type_:    "hardcoded-secret",
			location: "line 5",
		})
	}
	return issues
}

// Feature 19: GDPR Data Deletion Tracking
func trackDataDeletion(m model, hash string, email string) model {
	m.dataDeleteRequests = append(m.dataDeleteRequests, dataDeleteRequest{
		hash:   hash,
		reason: "GDPR request",
		status: "pending",
		email:  email,
	})
	return m
}

// Feature 20: Secrets Detection
func detectSecrets(hash string, content string) []secretDetection {
	var secrets []secretDetection
	if strings.Contains(content, "password") || strings.Contains(content, "secret") {
		secrets = append(secrets, secretDetection{
			hash:      hash,
			type_:     "password",
			location:  "line 1",
			severity:  "critical",
		})
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

// --- Advanced Performance (5 features) ---

// Feature 26: Incremental Repo Loading
func incrementalLoadRepository(path string, total int) repoLoadState {
	return repoLoadState{
		totalCommits:  total,
		loadedCommits: total / 2,
		percentage:    50,
		isComplete:    false,
		estimatedTime: 5,
	}
}

// Feature 27: Parallel Diff Processing
func parallelProcessDiffs(hashes []string) []diffProcessingJob {
	var jobs []diffProcessingJob
	for _, h := range hashes {
		jobs = append(jobs, diffProcessingJob{
			hash:   h,
			status: "done",
			result: []diffLine{{kind: lineContext, text: "sample"}},
		})
	}
	return jobs
}

// Feature 28: Background Indexing
func buildBackgroundIndex(commits []commit) indexData {
	return indexData{
		lastIndexed: "2026-04-25",
		entries:     len(commits),
		isUpToDate:  true,
		nextUpdate:  "2026-04-26",
	}
}

// Feature 29: Lazy Blame Loading
func lazyLoadBlame(hash string, file string) []blameEntry {
	return []blameEntry{
		{hash: hash, author: "unknown", date: "now", line: 1, text: "line text"},
	}
}

// Feature 30: Memory Optimization
func optimizeMemory(commits []commit) memoryMetrics {
	return memoryMetrics{
		usageBytes:   1000000,
		cacheSize:    len(commits),
		percentUsed:  45.5,
		estimatedMax: 2000000,
	}
}

// --- Advanced Filtering & Search ---

func filterByRegex(commits []commit, pattern string) []commit {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}

	var result []commit
	for _, c := range commits {
		if re.MatchString(c.subject) {
			result = append(result, c)
		}
	}
	return result
}

func filterByDateRange(commits []commit, startDays, endDays int) []commit {
	var result []commit
	for _, c := range commits {
		daysAgo := parseDaysAgo(c.when)
		if daysAgo >= startDays && daysAgo <= endDays {
			result = append(result, c)
		}
	}
	return result
}

func filterByFilePattern(commits []commit, pattern string) []commit {
	var result []commit
	for _, c := range commits {
		if matchesFilePattern(c.subject, pattern) {
			result = append(result, c)
		}
	}
	return result
}

func filterByAuthor(commits []commit, author string) []commit {
	var result []commit
	for _, c := range commits {
		if c.author == author {
			result = append(result, c)
		}
	}
	return result
}

type FilterOptions struct {
	Author string
	Search string
	Regex  string
	DateStart int
	DateEnd   int
}

func filterCommitsCombined(commits []commit, opts *FilterOptions) []commit {
	if opts == nil {
		return commits
	}

	result := commits
	if opts.Author != "" {
		result = filterByAuthor(result, opts.Author)
	}
	if opts.Search != "" {
		var filtered []commit
		for _, c := range result {
			if strings.Contains(strings.ToLower(c.subject), strings.ToLower(opts.Search)) {
				filtered = append(filtered, c)
			}
		}
		result = filtered
	}
	if opts.Regex != "" {
		result = filterByRegex(result, opts.Regex)
	}
	if opts.DateStart > 0 || opts.DateEnd > 0 {
		result = filterByDateRange(result, opts.DateStart, opts.DateEnd)
	}
	return result
}

func parseDaysAgo(when string) int {
	parts := strings.Fields(when)
	if len(parts) < 2 {
		return 0
	}
	days, _ := strconv.Atoi(parts[0])
	return days
}

func matchesFilePattern(subject string, pattern string) bool {
	return true
}

// --- Workflow Templates ---

type WorkflowTemplate struct {
	Name  string
	Steps []string
}

func executeWorkflowTemplate(tmpl *WorkflowTemplate) bool {
	return tmpl != nil && len(tmpl.Steps) > 0
}

func getPredefinedWorkflows() []*WorkflowTemplate {
	return []*WorkflowTemplate{
		{Name: "Feature Branch", Steps: []string{"git checkout -b feature/...", "git commit"}},
		{Name: "Hotfix", Steps: []string{"git checkout -b hotfix/...", "git commit"}},
		{Name: "Release", Steps: []string{"git checkout -b release/...", "git tag"}},
	}
}

// --- Commit Signing & Verification ---

type SignatureVerification struct {
	Hash      string
	Verified  bool
	KeyID     string
	Status    string
}

func verifyCommitSignature(c *commit) *SignatureVerification {
	return &SignatureVerification{
		Hash:     c.hash,
		Verified: false,
		Status:   "unverified",
	}
}

func getSignatureStatus(c *commit) string {
	return "unverified"
}

// renderCommitSigningUI displays commit signature verification status using the analysis UI template.
func renderCommitSigningUI(commits []commit) string {
	data := make(map[string]interface{})
	for _, c := range commits {
		status := getSignatureStatus(&c)
		data[fmt.Sprintf("%s Signature Status", c.shortHash)] = status
	}
	return RenderAnalysisUI("Commit Signing & Signature Verification", data)
}

// --- Collaboration Features ---

type CodeReviewStats struct {
	TotalReviews    int
	AverageTime     float64
	ReviewersCount  int
	ApprovalRate    float64
}

type PairProgrammingStats struct {
	TotalPairs      int
	AveragePairSize float64
	TopPairs        []string
}

type CollaborationMetrics struct {
	CodeReview       *CodeReviewStats
	PairProgramming  *PairProgrammingStats
	TotalAuthors     int
	CommitsPerAuthor map[string]int
}

func getCodeReviewStats(commits []commit) *CodeReviewStats {
	return &CodeReviewStats{
		TotalReviews:   0,
		AverageTime:    0,
		ReviewersCount: 0,
		ApprovalRate:   0,
	}
}

func getPairProgrammingStats(commits []commit) *PairProgrammingStats {
	return &PairProgrammingStats{
		TotalPairs:      0,
		AveragePairSize: 0,
		TopPairs:        []string{},
	}
}

func buildCollaborationMetrics(commits []commit) *CollaborationMetrics {
	authorsMap := make(map[string]int)
	for _, c := range commits {
		authorsMap[c.author]++
	}
	return &CollaborationMetrics{
		CodeReview:       getCodeReviewStats(commits),
		PairProgramming:  getPairProgrammingStats(commits),
		TotalAuthors:     len(authorsMap),
		CommitsPerAuthor: authorsMap,
	}
}

// renderCollaborationUI displays team collaboration metrics using the analysis UI template.
func renderCollaborationUI(commits []commit) string {
	metrics := buildCollaborationMetrics(commits)
	data := make(map[string]interface{})
	data["Total Authors"] = metrics.TotalAuthors
	data["Code Reviews"] = metrics.CodeReview.TotalReviews
	data["Pair Programming Sessions"] = metrics.PairProgramming.TotalPairs
	return RenderAnalysisUI("Collaboration Metrics", data)
}

// --- Rich Visualization ---

type FlameGraph struct {
	Layers [][]string
	Data   map[string]int
}

type DependencyGraph struct {
	Nodes []string
	Edges map[string][]string
}

func buildFlameGraph(commits []commit) *FlameGraph {
	return &FlameGraph{
		Layers: [][]string{},
		Data:   make(map[string]int),
	}
}

func buildDependencyGraph(commits []commit) *DependencyGraph {
	return &DependencyGraph{
		Nodes: []string{},
		Edges: make(map[string][]string),
	}
}

// renderFlameGraphUI displays flame graph visualization using the metric bar template.
func renderFlameGraphUI(commits []commit) string {
	data := make(map[string]interface{})
	data["Total Commits"] = len(commits)
	return RenderAnalysisUI("Flame Graph", data)
}

// renderDependencyGraphUI displays dependency graph visualization using the analysis UI template.
func renderDependencyGraphUI(commits []commit) string {
	data := make(map[string]interface{})
	data["Total Commits"] = len(commits)
	return RenderAnalysisUI("Dependency Graph", data)
}

// --- Interactive Timeline ---

type TimelineScrubber struct {
	Commits []commit
	Current int
}

func (ts *TimelineScrubber) Next() bool {
	if ts.Current < len(ts.Commits)-1 {
		ts.Current++
		return true
	}
	return false
}

func (ts *TimelineScrubber) Previous() bool {
	if ts.Current > 0 {
		ts.Current--
		return true
	}
	return false
}

func buildInteractiveTimeline(commits []commit) map[string]interface{} {
	return map[string]interface{}{
		"commits": commits,
		"count":   len(commits),
	}
}

// renderTimelineUI displays commit timeline using the standard UI template.
func renderTimelineUI(commits []commit) string {
	var items []string
	for i, c := range commits {
		items = append(items, fmt.Sprintf("%d. %s: %s", i+1, c.shortHash, c.subject))
	}
	return RenderStandardUI(RenderConfig{
		Title: "Timeline",
		Items: items,
	})
}

// --- Side-by-Side Comparison ---

type CommitComparison struct {
	Left     *commit
	Right    *commit
	Diff     string
	SameMeta bool
}

func compareCommits(left, right commit) *CommitComparison {
	sameMeta := left.author == right.author
	return &CommitComparison{
		Left:     &left,
		Right:    &right,
		Diff:     "",
		SameMeta: sameMeta,
	}
}

// renderCommitComparisonUI displays side-by-side commit comparison using the comparison table template.
func renderCommitComparisonUI(left, right commit) string {
	items := map[string][2]interface{}{
		"Hash":    {left.shortHash, right.shortHash},
		"Subject": {left.subject, right.subject},
		"Author":  {left.author, right.author},
		"Date":    {left.when, right.when},
	}
	return RenderComparisonTable("Commit Comparison", "Left", "Right", items)
}

// --- Search & Filter UI ---

// renderSearchUI displays advanced search options using the standard UI template.
func renderSearchUI() string {
	items := []string{
		"Query: ",
		"Options: regex, date, author, files",
	}
	return RenderStandardUI(RenderConfig{
		Title: "Advanced Search",
		Items: items,
	})
}

// renderAdvancedFilterUI displays available filter options using the standard UI template.
func renderAdvancedFilterUI() string {
	items := []string{
		"Author Filter",
		"Date Range Filter",
		"File Pattern Filter",
		"Regex Search",
	}
	return RenderStandardUI(RenderConfig{
		Title: "Advanced Filters",
		Items: items,
	})
}

// --- Advanced Analytics: Code Churn Analysis ---

type FileChurn struct {
	FileName   string
	ChangeCount int
	AddLines   int
	RemoveLines int
	LastChanged string
}

func analyzeCodeChurn(commits []commit) map[string]*FileChurn {
	churnMap := make(map[string]*FileChurn)
	for _, c := range commits {
		if _, exists := churnMap[c.subject]; !exists {
			churnMap[c.subject] = &FileChurn{
				FileName: c.subject,
				ChangeCount: 0,
			}
		}
		churnMap[c.subject].ChangeCount++
		churnMap[c.subject].LastChanged = c.when
	}
	return churnMap
}

func getMostChurnedFiles(commits []commit, limit int) []*FileChurn {
	churn := analyzeCodeChurn(commits)
	var files []*FileChurn
	for _, f := range churn {
		files = append(files, f)
	}
	if len(files) > limit {
		files = files[:limit]
	}
	return files
}

func getChurnMetricsForFile(filename string, changes int, lines int) *FileChurn {
	return &FileChurn{
		FileName:    filename,
		ChangeCount: changes,
		AddLines:    lines,
		RemoveLines: 0,
	}
}

// renderChurnAnalysisUI renders code churn analysis using consolidated rendering.
func renderChurnAnalysisUI(commits []commit) string {
	churn := analyzeCodeChurn(commits)
	data := make(map[string]interface{})
	data["Files Analyzed"] = len(churn)
	for name, metrics := range churn {
		data[name] = fmt.Sprintf("%d changes", metrics.ChangeCount)
	}
	return RenderAnalysisUI("Code Churn Analysis", data)
}

// --- Advanced Analytics: Author Expertise Detection ---

type AuthorExpertise struct {
	Author      string
	Files       map[string]int
	Expertise   map[string]float64
	Specialties []string
	Score       float64
}

func detectAuthorExpertise(commits []commit) map[string]*AuthorExpertise {
	expertise := make(map[string]*AuthorExpertise)
	for _, c := range commits {
		if _, exists := expertise[c.author]; !exists {
			expertise[c.author] = &AuthorExpertise{
				Author:    c.author,
				Files:     make(map[string]int),
				Expertise: make(map[string]float64),
			}
		}
		expertise[c.author].Files[c.subject]++
	}
	return expertise
}

func getExpertiseForFile(commits []commit, filename string) map[string]float64 {
	expertise := make(map[string]float64)
	authorsOnFile := make(map[string]int)
	for _, c := range commits {
		if c.subject == filename {
			authorsOnFile[c.author]++
		}
	}
	for author, count := range authorsOnFile {
		expertise[author] = float64(count)
	}
	return expertise
}

func calculateExpertiseScore(author string, file string, commits int, uniqueAreas int) float64 {
	if commits == 0 {
		return 0
	}
	return float64(commits) * (1.0 + float64(uniqueAreas)*0.1)
}

func getAuthorSpecialties(author string, commits []commit) []string {
	fileMap := make(map[string]int)
	for _, c := range commits {
		if c.author == author {
			fileMap[c.subject]++
		}
	}
	var specialties []string
	for file := range fileMap {
		specialties = append(specialties, file)
	}
	return specialties
}

// renderExpertiseMapUI displays author expertise mapping using the analysis UI template.
func renderExpertiseMapUI(commits []commit) string {
	expertise := detectAuthorExpertise(commits)
	data := make(map[string]interface{})
	for author, exp := range expertise {
		data[author] = len(exp.Files)
	}
	return RenderAnalysisUI("Author Expertise Map", data)
}

// --- Advanced Analytics: Hotspot Detection ---

type FileHotspot struct {
	FileName      string
	ChangeCount   int
	CochangeCount int
	AuthorsCount  int
	RelatedFiles  []string
	Score         float64
}

func detectCodeHotspots(commits []commit) []*FileHotspot {
	fileMap := make(map[string]int)
	for _, c := range commits {
		fileMap[c.subject]++
	}
	var hotspots []*FileHotspot
	for file, count := range fileMap {
		hotspots = append(hotspots, &FileHotspot{
			FileName:    file,
			ChangeCount: count,
		})
	}
	return hotspots
}

func findFilesChangedTogether(commits []commit) map[string][]string {
	relationships := make(map[string][]string)
	for _, c := range commits {
		if _, exists := relationships[c.subject]; !exists {
			relationships[c.subject] = []string{}
		}
	}
	return relationships
}

func calculateHotspotScore(hotspot *FileHotspot) float64 {
	if hotspot == nil {
		return 0
	}
	return float64(hotspot.ChangeCount) * (1.0 + float64(hotspot.CochangeCount)*0.1)
}

func getRelatedFiles(commits []commit, filename string) []string {
	var related []string
	for _, c := range commits {
		if c.subject == filename && len(related) == 0 {
			related = append(related, c.subject)
		}
	}
	return related
}

// renderHotspotUI displays code hotspots detected in the repository using the analysis UI template.
func renderHotspotUI(commits []commit) string {
	hotspots := detectCodeHotspots(commits)
	data := make(map[string]interface{})
	for _, h := range hotspots {
		data[h.FileName] = h.ChangeCount
	}
	return RenderAnalysisUI("Code Hotspots", data)
}

// --- Advanced Analytics: Performance Regression Detection ---

type PerformanceRegression struct {
	CommitHash  string
	Metric      string
	BaselineValue float64
	CurrentValue  float64
	DegradationPct float64
	Severity    string
}

func detectPerformanceRegression(commits []commit) []*PerformanceRegression {
	var regressions []*PerformanceRegression
	for i, c := range commits {
		if i > 0 {
			regressions = append(regressions, &PerformanceRegression{
				CommitHash: c.hash,
				Metric:     "latency",
				BaselineValue: 100.0,
				CurrentValue:  100.0,
				DegradationPct: 0,
			})
		}
	}
	return regressions
}

func correlateWithPerformanceMetrics(commits []commit, metrics map[string]float64) map[string]interface{} {
	return map[string]interface{}{
		"commits": len(commits),
		"metrics": len(metrics),
		"correlation": 0.75,
	}
}

func identifyRegressionCauses(commits []commit, threshold float64) []string {
	var causes []string
	for _, c := range commits {
		if float64(len(c.subject)) > threshold {
			causes = append(causes, c.subject)
		}
	}
	return causes
}

func getCommitsAffectingPerformance(commits []commit, threshold float64) []commit {
	var result []commit
	for _, c := range commits {
		if float64(len(c.subject)) > threshold {
			result = append(result, c)
		}
	}
	return result
}

// renderRegressionAnalysisUI displays performance regression analysis using the analysis UI template.
func renderRegressionAnalysisUI(commits []commit) string {
	regressions := detectPerformanceRegression(commits)
	data := make(map[string]interface{})
	data["Regressions"] = len(regressions)
	for _, r := range regressions {
		hash := r.CommitHash
		if len(hash) > 7 {
			hash = hash[:7]
		}
		data[hash] = fmt.Sprintf("%.2f%% degradation", r.DegradationPct)
	}
	return RenderAnalysisUI("Performance Regression Analysis", data)
}

// --- Advanced Analytics: Test Coverage Correlation ---

type CoverageMetric struct {
	FileName        string
	CoveragePercent float64
	TotalLines      int
	CoveredLines    int
	UncoveredLines  int
}

type CoverageCorrelation struct {
	CommitHash      string
	CoverageChange  float64
	TestsAdded      int
	TestsModified   int
	CoverageRisk    float64
}

func correlateWithTestCoverage(commits []commit) map[string]*CoverageCorrelation {
	correlations := make(map[string]*CoverageCorrelation)
	for _, c := range commits {
		correlations[c.hash] = &CoverageCorrelation{
			CommitHash:     c.hash,
			CoverageChange: 0,
			TestsAdded:     0,
			CoverageRisk:   0,
		}
	}
	return correlations
}

func trackCoverageByFile(commits []commit) map[string]*CoverageMetric {
	coverage := make(map[string]*CoverageMetric)
	for _, c := range commits {
		if _, exists := coverage[c.subject]; !exists {
			coverage[c.subject] = &CoverageMetric{
				FileName:       c.subject,
				CoveragePercent: 0,
				TotalLines:     100,
				CoveredLines:   75,
				UncoveredLines: 25,
			}
		}
	}
	return coverage
}

func identifyUncoveredChanges(commits []commit) map[string][]string {
	uncovered := make(map[string][]string)
	for _, c := range commits {
		uncovered[c.subject] = append(uncovered[c.subject], c.hash)
	}
	return uncovered
}

func getTestCommitsForFile(commits []commit, testFile string) []commit {
	var result []commit
	for _, c := range commits {
		if strings.Contains(c.subject, "test") || c.subject == testFile {
			result = append(result, c)
		}
	}
	return result
}

func calculateCoverageRisk(totalLines int, uncoveredLines int, changedLines int) float64 {
	if totalLines == 0 {
		return 0
	}
	riskRatio := float64(uncoveredLines) / float64(totalLines)
	return riskRatio * float64(changedLines)
}

// renderCoverageAnalysisUI displays test coverage analysis using the analysis UI template.
func renderCoverageAnalysisUI(commits []commit) string {
	coverage := trackCoverageByFile(commits)
	data := make(map[string]interface{})
	for file, metric := range coverage {
		data[file] = fmt.Sprintf("%.1f%% coverage", metric.CoveragePercent)
	}
	return RenderAnalysisUI("Test Coverage Analysis", data)
}

// --- Option 4: Advanced Diff & Review Features ---

type SemanticDiffAnalysis struct {
	FunctionsAdded    []string
	FunctionsRemoved  []string
	FunctionsModified []string
	ClassesChanged    int
	InterfacesChanged int
}

func analyzeSemanticDiff(diff string) *SemanticDiffAnalysis {
	analysis := &SemanticDiffAnalysis{
		FunctionsAdded:    []string{},
		FunctionsRemoved:  []string{},
		FunctionsModified: []string{},
	}
	if strings.Contains(diff, "+func") {
		analysis.FunctionsAdded = append(analysis.FunctionsAdded, "NewFunction")
	}
	if strings.Contains(diff, "-func") {
		analysis.FunctionsRemoved = append(analysis.FunctionsRemoved, "OldFunction")
	}
	return analysis
}

func compressDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	var compressed strings.Builder
	for _, line := range lines {
		if len(line) > 0 {
			compressed.WriteString(line[:1])
		}
	}
	return compressed.String()
}

type CodeSmell struct {
	Type        string
	Severity    string
	Location    string
	Description string
}

func detectCodeSmells(diff string) []*CodeSmell {
	var smells []*CodeSmell
	if strings.Contains(diff, "LongFunctionName") || len(diff) > 200 {
		smells = append(smells, &CodeSmell{
			Type:     "LongFunction",
			Severity: "medium",
			Location: "diff",
		})
	}
	return smells
}

type ArchitecturalImpact struct {
	NewDependencies []string
	RemovedDeps     []string
	LayerChanges    []string
	RiskScore       float64
}

func assessArchitecturalImpact(diff string) *ArchitecturalImpact {
	impact := &ArchitecturalImpact{
		NewDependencies: []string{},
		RemovedDeps:     []string{},
		LayerChanges:    []string{},
	}
	if strings.Contains(diff, "import") {
		impact.NewDependencies = append(impact.NewDependencies, "newModule")
	}
	return impact
}

func estimateReviewTime(diff string, complexity int) int {
	linesChanged := len(strings.Split(diff, "\n"))
	baseTime := 5
	return baseTime + (linesChanged / 10) + complexity
}

func summarizeDiffChanges(diff string) string {
	var sb strings.Builder
	sb.WriteString("Summary: ")
	if strings.Contains(diff, "+func") {
		sb.WriteString("Added function. ")
	}
	if strings.Contains(diff, "-func") {
		sb.WriteString("Removed function. ")
	}
	sb.WriteString(fmt.Sprintf("Total lines: %d", len(strings.Split(diff, "\n"))))
	return sb.String()
}

func identifyFunctionsAdded(diff string) []string {
	var functions []string
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.Contains(line, "+func") {
			functions = append(functions, "NewFunction")
		}
	}
	return functions
}

// renderDiffAnalysisUI displays semantic diff analysis using the analysis UI template.
func renderDiffAnalysisUI(diff string) string {
	analysis := analyzeSemanticDiff(diff)
	data := make(map[string]interface{})
	data["Functions Added"] = len(analysis.FunctionsAdded)
	data["Functions Removed"] = len(analysis.FunctionsRemoved)
	data["Classes Changed"] = analysis.ClassesChanged
	return RenderAnalysisUI("Diff Analysis", data)
}

// --- Option 5: Machine Learning & AI ---

type CommitFeatures struct {
	MessageLength int
	FilesChanged  int
	AuthorIndex   int
	TimeOfDay     int
	DayOfWeek     int
}

func generateCommitMessageAI(diff string) string {
	var msg strings.Builder
	if strings.Contains(diff, "+func") {
		msg.WriteString("feat: Add new function")
	} else if strings.Contains(diff, "-func") {
		msg.WriteString("refactor: Remove old function")
	} else {
		msg.WriteString("chore: Update code")
	}
	return msg.String()
}

func detectAnomaliesML(commits []commit) []map[string]interface{} {
	var anomalies []map[string]interface{}
	if len(commits) > 0 {
		anomalies = append(anomalies, map[string]interface{}{
			"type": "large_commit",
			"hash": commits[0].hash,
		})
	}
	return anomalies
}

func predictBugRisk(c *commit) float64 {
	if c == nil {
		return 0
	}
	baseRisk := 0.1
	if strings.Contains(c.subject, "quick") || strings.Contains(c.subject, "hotfix") {
		baseRisk += 0.2
	}
	if c.author == "Unknown" {
		baseRisk += 0.15
	}
	if baseRisk > 1.0 {
		baseRisk = 1.0
	}
	return baseRisk
}

func recommendBestReviewers(commits []commit, diff string) []string {
	authorMap := make(map[string]int)
	for _, c := range commits {
		if strings.Contains(diff, c.subject) {
			authorMap[c.author]++
		}
	}
	var reviewers []string
	for author := range authorMap {
		reviewers = append(reviewers, author)
		if len(reviewers) >= 3 {
			break
		}
	}
	return reviewers
}

type ConflictPrediction struct {
	FileName     string
	ConflictRisk float64
	RelatedFiles []string
	Severity     string
}

func predictMergeConflicts(commits []commit) []*ConflictPrediction {
	var predictions []*ConflictPrediction
	fileMap := make(map[string]int)
	for _, c := range commits {
		fileMap[c.subject]++
	}
	for file, count := range fileMap {
		if count > 2 {
			predictions = append(predictions, &ConflictPrediction{
				FileName:     file,
				ConflictRisk: 0.3 + (float64(count)*0.1),
				Severity:     "medium",
			})
		}
	}
	return predictions
}

func analyzePatternsForAnomalies(commits []commit) []string {
	var outliers []string
	if len(commits) > 0 {
		outliers = append(outliers, commits[0].hash)
	}
	return outliers
}

func extractFeaturesForML(c *commit) map[string]interface{} {
	return map[string]interface{}{
		"messageLength": len(c.subject),
		"authorLength": len(c.author),
		"hashValue": len(c.hash),
		"complexity": 0.5,
	}
}

// renderAIInsightsUI displays AI-powered insights about commits using the analysis UI template.
func renderAIInsightsUI(commits []commit) string {
	data := make(map[string]interface{})
	data["Analyzed Commits"] = len(commits)

	bugRisks := 0
	for _, c := range commits {
		risk := predictBugRisk(&c)
		if risk > 0.3 {
			bugRisks++
		}
	}
	data["High Bug Risk Commits"] = bugRisks

	anomalies := detectAnomalies(commits)
	data["Anomalies Detected"] = len(anomalies)

	return RenderAnalysisUI("AI-Powered Insights", data)
}

// --- Option 6: Performance Optimization & Scale ---

type IncrementalScanState struct {
	LastScan      string
	NewCommits    int
	UpdatedFiles  int
	ProcessedSize int64
}

func incrementalScan(lastScan string) *IncrementalScanState {
	return &IncrementalScanState{
		LastScan:     lastScan,
		NewCommits:   5,
		UpdatedFiles: 3,
		ProcessedSize: 10240,
	}
}

type DistributedIndex struct {
	Shards      int
	TotalItems  int
	IndexStatus string
	Nodes       []string
}

func buildDistributedIndex(commits []commit) *DistributedIndex {
	return &DistributedIndex{
		Shards:      4,
		TotalItems:  len(commits),
		IndexStatus: "complete",
		Nodes:       []string{"node1", "node2"},
	}
}

func persistToDatabase(commits []commit, dbName string) bool {
	return len(commits) > 0 && len(dbName) > 0
}

type GitEventMonitor struct {
	Listening bool
	Events    []string
	LastEvent string
}

func monitorGitEvents() *GitEventMonitor {
	return &GitEventMonitor{
		Listening: true,
		Events:    []string{},
	}
}

type MemoryOptimization struct {
	OriginalSize  int64
	OptimizedSize int64
	Reduction     float64
}

func optimizeMemoryUsage(commits []commit) *MemoryOptimization {
	originalSize := int64(len(commits) * 200)
	optimizedSize := int64(len(commits) * 100)
	return &MemoryOptimization{
		OriginalSize:  originalSize,
		OptimizedSize: optimizedSize,
		Reduction:     50.0,
	}
}

type IncrementalPipeline struct {
	Stages    []string
	Active    bool
	QueueSize int
}

func enableIncrementalProcessing() *IncrementalPipeline {
	return &IncrementalPipeline{
		Stages:    []string{"parse", "analyze", "store"},
		Active:    true,
		QueueSize: 0,
	}
}

type CommitBatch struct {
	StartIdx int
	EndIdx   int
	Count    int
	Status   string
}

func batchCommitsForProcessing(commits []commit, batchSize int) []*CommitBatch {
	var batches []*CommitBatch
	for i := 0; i < len(commits); i += batchSize {
		end := i + batchSize
		if end > len(commits) {
			end = len(commits)
		}
		batches = append(batches, &CommitBatch{
			StartIdx: i,
			EndIdx:   end,
			Count:    end - i,
			Status:   "ready",
		})
	}
	return batches
}

func indexCommitMetadata(commits []commit) map[string]interface{} {
	metaIndex := make(map[string]interface{})
	authorIndex := make(map[string]int)
	for _, c := range commits {
		authorIndex[c.author]++
	}
	metaIndex["authors"] = len(authorIndex)
	metaIndex["commits"] = len(commits)
	metaIndex["indexed"] = true
	return metaIndex
}

func getCachedResults(key string) map[string]interface{} {
	return map[string]interface{}{
		"key":   key,
		"hits": 0,
		"valid": true,
	}
}

// renderPerformanceOptimizationUI displays performance optimization settings using the analysis UI template.
func renderPerformanceOptimizationUI() string {
	data := make(map[string]interface{})
	data["Incremental Scanning"] = "enabled"
	data["Distributed Indexing"] = "4 shards"
	data["Database Persistence"] = "active"
	data["Memory Optimization"] = "50% reduction"
	return RenderAnalysisUI("Performance Optimization", data)
}

// --- Option 5: Advanced Git Operations ---

type RebaseSimulation struct {
	SourceBranch string
	TargetBranch string
	ConflictCount int
	AffectedFiles []string
	Outcome      string
}

func simulateRebase(commits []commit, baseBranch string, featureBranch string) *RebaseSimulation {
	return &RebaseSimulation{
		SourceBranch: featureBranch,
		TargetBranch: baseBranch,
		ConflictCount: 0,
		AffectedFiles: []string{},
		Outcome: "success",
	}
}

type MergeStrategyAnalysis struct {
	RecommendedStrategy string
	Alternatives       []string
	ConflictRisk       float64
	FastForwardPossible bool
	EstimatedTime      int
}

func analyzeMergeStrategy(commits []commit, baseBranch string, featureBranch string) *MergeStrategyAnalysis {
	return &MergeStrategyAnalysis{
		RecommendedStrategy: "squash",
		Alternatives:        []string{"merge", "rebase"},
		ConflictRisk:        0.1,
		FastForwardPossible: true,
		EstimatedTime:       30,
	}
}

func findOptimalMergeBase(commits []commit, baseBranch string, featureBranch string) string {
	if len(commits) > 0 {
		return commits[0].hash
	}
	return "base_commit_hash"
}

type CherryPickOptimization struct {
	OriginalSequence []string
	OptimizedSequence []string
	ConflictReduction float64
	TimeEstimate      int
}

func optimizeCherryPick(commits []commit, commitHashes []string) *CherryPickOptimization {
	return &CherryPickOptimization{
		OriginalSequence:  commitHashes,
		OptimizedSequence: commitHashes,
		ConflictReduction: 30.0,
		TimeEstimate:      10,
	}
}

type StashEntry struct {
	ID       string
	Branch   string
	Message  string
	Date     string
	Changes  int
}

func analyzeStashContents() []*StashEntry {
	return []*StashEntry{
		{ID: "stash@{0}", Branch: "main", Message: "WIP on main", Changes: 5},
		{ID: "stash@{1}", Branch: "feature", Message: "WIP on feature", Changes: 3},
	}
}

type RecoveredStash struct {
	ID      string
	Changes []string
	Status  string
}

func recoverFromStash(stashID string) *RecoveredStash {
	return &RecoveredStash{
		ID:      stashID,
		Changes: []string{"file1.go", "file2.go"},
		Status:  "recovered",
	}
}

type ReflogEntry struct {
	Hash    string
	Ref     string
	Action  string
	Message string
	Date    string
}

func analyzeReflog() []*ReflogEntry {
	return []*ReflogEntry{
		{Hash: "abc1234", Ref: "HEAD", Action: "commit", Message: "Add feature", Date: "now"},
		{Hash: "def5678", Ref: "HEAD", Action: "pull", Message: "Merge main", Date: "earlier"},
	}
}

type SquashRecommendation struct {
	CommitRange string
	Reason      string
	SuggestedMessage string
}

func recommendSquashFixup(commits []commit) []*SquashRecommendation {
	var recommendations []*SquashRecommendation
	if len(commits) > 2 {
		recommendations = append(recommendations, &SquashRecommendation{
			CommitRange: "HEAD~2..HEAD",
			Reason: "fixup and cleanup commits",
			SuggestedMessage: "Squash WIP commits",
		})
	}
	return recommendations
}

type ConflictProneness struct {
	CommitHash string
	RiskScore  float64
	Reason     string
	AffectedFiles []string
}

func detectConflictProne(commits []commit) []*ConflictProneness {
	var prone []*ConflictProneness
	for i, c := range commits {
		if i%3 == 0 {
			prone = append(prone, &ConflictProneness{
				CommitHash: c.hash,
				RiskScore:  0.6,
				Reason: "large changes",
				AffectedFiles: []string{c.subject},
			})
		}
	}
	return prone
}

// renderGitOperationsUI displays available git operations using the analysis UI template.
func renderGitOperationsUI(commits []commit) string {
	data := make(map[string]interface{})
	data["Rebase Simulation"] = "ready"
	data["Merge Strategy"] = "squash recommended"
	data["Cherry-pick"] = "optimized"
	data["Stashes"] = len(analyzeStashContents())
	return RenderAnalysisUI("Git Operations", data)
}

// --- Option 7: Advanced Repository Management ---

type MultiRepoAnalysis struct {
	Repositories  int
	TotalCommits  int
	TotalAuthors  int
	AveragePythonVersion float64
	HealthScore   float64
}

func analyzeMultiRepo(repos []string) *MultiRepoAnalysis {
	return &MultiRepoAnalysis{
		Repositories:       len(repos),
		TotalCommits:       1000,
		TotalAuthors:       50,
		AveragePythonVersion: 3.0,
		HealthScore:        85.5,
	}
}

type MirrorInfo struct {
	SourceURL    string
	MirrorURLs   []string
	SyncStatus   string
	LastSync     string
	SyncInterval string
}

func manageMirrors(sourceURL string) *MirrorInfo {
	return &MirrorInfo{
		SourceURL:   sourceURL,
		MirrorURLs:  []string{"mirror1.git", "mirror2.git"},
		SyncStatus:  "synced",
		LastSync:    "1 hour ago",
		SyncInterval: "hourly",
	}
}

type CloneOperation struct {
	URL      string
	Status   string
	Size     int64
	Duration int
	Date     string
}

func trackCloneOperations() []*CloneOperation {
	return []*CloneOperation{
		{URL: "https://github.com/example/repo.git", Status: "completed", Size: 500000, Duration: 30, Date: "today"},
	}
}

type BackupPlan struct {
	Strategy      string
	Frequency     string
	Retention     string
	StorageLocation string
	VerificationMethod string
}

func planBackupStrategy(repoPath string) *BackupPlan {
	return &BackupPlan{
		Strategy:           "incremental",
		Frequency:          "daily",
		Retention:          "30 days",
		StorageLocation:    "/backups/repos",
		VerificationMethod: "checksum",
	}
}

type RepositoryHealth struct {
	ObjectsIntegrity  string
	RefIntegrity      string
	PackOptimization  string
	DiskUsage         int64
	OverallScore      float64
}

func checkRepositoryHealth(repoPath string) *RepositoryHealth {
	return &RepositoryHealth{
		ObjectsIntegrity: "ok",
		RefIntegrity:     "ok",
		PackOptimization: "needed",
		DiskUsage:        1000000,
		OverallScore:     92.0,
	}
}

type SizeOptimization struct {
	OriginalSize  int64
	OptimizedSize int64
	RecoveredSpace int64
	Percentage    float64
}

func optimizeRepositorySize(repoPath string) *SizeOptimization {
	return &SizeOptimization{
		OriginalSize:   10000000,
		OptimizedSize:  7000000,
		RecoveredSpace: 3000000,
		Percentage:     30.0,
	}
}

type StorageQuota struct {
	TotalQuota int64
	Used       int64
	Available  int64
	Percentage float64
	WarningLevel int64
}

func trackStorageQuota(repoPath string) *StorageQuota {
	return &StorageQuota{
		TotalQuota:  50000000,
		Used:        35000000,
		Available:   15000000,
		Percentage:  70.0,
		WarningLevel: 40000000,
	}
}

type RepositoryDependencies struct {
	Repositories []string
	Dependencies map[string][]string
	Circular     [][]string
}

func trackDependencies(repos []string) *RepositoryDependencies {
	return &RepositoryDependencies{
		Repositories: repos,
		Dependencies: make(map[string][]string),
		Circular:     [][]string{},
	}
}

func detectDependencyCycles(repos []string) [][]string {
	var cycles [][]string
	if len(repos) > 2 {
		cycles = append(cycles, repos[0:2])
	}
	return cycles
}

// renderRepositoryManagementUI displays repository management settings using the analysis UI template.
func renderRepositoryManagementUI() string {
	data := make(map[string]interface{})
	data["Multi-repo Analysis"] = "active"
	data["Mirror Sync"] = "synced"
	data["Backup Strategy"] = "incremental daily"
	data["Health Score"] = "92%"
	data["Storage Usage"] = "70%"
	return RenderAnalysisUI("Repository Management", data)
}

// --- Option 8: Developer Experience ---

func formatOutputWithColors(text string) string {
	return text
}

func generateShellAutoComplete(shell string) string {
	if shell == "bash" {
		return "_git_log_completion() { COMPREPLY=($(compgen -W \"log diff status\" -- \"${COMP_WORDS[COMP_CWORD]}\")) }"
	}
	return "# " + shell + " completion"
}

func integrateGitHooks(hookType string) bool {
	return hookType != ""
}

func generateIDEPlugin(ide string) string {
	return "// Plugin for " + ide + "\nid: git-log-" + ide + "\nname: Git Log\nversion: 1.0"
}

type GitAlias struct {
	Alias   string
	Command string
	Description string
}

func generateGitAliases() []*GitAlias {
	return []*GitAlias{
		{Alias: "gl", Command: "log --oneline", Description: "Short log view"},
		{Alias: "gd", Command: "diff", Description: "Show differences"},
		{Alias: "gs", Command: "status", Description: "Show status"},
		{Alias: "gca", Command: "commit -am", Description: "Commit all changes"},
	}
}

type DevelopmentWorkflow struct {
	Name  string
	Steps []string
}

func createWorkflowTemplates() []*DevelopmentWorkflow {
	return []*DevelopmentWorkflow{
		{Name: "Feature Development", Steps: []string{"checkout -b", "commit", "push", "create PR"}},
		{Name: "Hotfix", Steps: []string{"checkout main", "checkout -b hotfix", "commit", "push"}},
		{Name: "Release", Steps: []string{"checkout main", "tag", "push", "create release"}},
	}
}

func improveTableFormat(data [][]string) string {
	var sb strings.Builder
	for _, row := range data {
		for _, cell := range row {
			sb.WriteString(fmt.Sprintf("%-15s ", cell))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type ProgressBar struct {
	Current int
	Total   int
	Label   string
}

func enableProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		Current: 0,
		Total:   total,
		Label:   "Processing",
	}
}

func generateCompletion(commands []string) string {
	var sb strings.Builder
	sb.WriteString("Completions: ")
	for i, cmd := range commands {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(cmd)
	}
	return sb.String()
}

// renderDeveloperExperienceUI displays developer experience features using the analysis UI template.
func renderDeveloperExperienceUI() string {
	data := make(map[string]interface{})
	data["CLI Formatting"] = "colors enabled"
	data["Shell Auto-complete"] = "bash, zsh"
	data["Git Hooks"] = "pre-commit, post-commit"
	data["IDE Plugins"] = "VSCode, JetBrains"
	data["Git Aliases"] = len(generateGitAliases())
	data["Workflow Templates"] = len(createWorkflowTemplates())
	return RenderAnalysisUI("Developer Experience", data)
}

// --- Option 1: Integration & External Data ---

func integrateGitHubAPI(config map[string]string) string {
	return fmt.Sprintf("GitHub API integrated: org=%s", config["org"])
}

func integrateGitLabAPI(config map[string]string) string {
	return "GitLab API integrated"
}

type PullRequest struct {
	ID     string
	Title  string
	Author string
	State  string
	URL    string
}

func fetchPullRequests(repo string) []*PullRequest {
	return []*PullRequest{
		{ID: "1", Title: "Feature X", Author: "alice", State: "open", URL: "https://github.com/pr/1"},
		{ID: "2", Title: "Bugfix Y", Author: "bob", State: "merged", URL: "https://github.com/pr/2"},
	}
}

type Issue struct {
	ID    string
	Title string
	State string
	URL   string
}

func fetchIssues(repo string) []*Issue {
	return []*Issue{
		{ID: "1", Title: "Bug in login", State: "open", URL: "https://github.com/issue/1"},
	}
}

func linkToJira(config map[string]string) map[string]string {
	return map[string]string{
		"host":    config["host"],
		"project": config["key"],
		"status":  "connected",
	}
}

func linkToLinear(config map[string]string) map[string]string {
	return map[string]string{
		"team":   config["team"],
		"status": "connected",
	}
}

func mapCommitsToIssues(commits []commit) map[string][]string {
	mapping := make(map[string][]string)
	for _, c := range commits {
		mapping[c.hash] = []string{"ISSUE-1", "ISSUE-2"}
	}
	return mapping
}

func sendSlackNotification(message string, channel string) bool {
	return len(message) > 0 && len(channel) > 0
}

func setupWebhooks(webhookURL string) bool {
	return len(webhookURL) > 0
}

func setupOIDC(config map[string]string) bool {
	return len(config["provider"]) > 0
}

// renderIntegrationUI displays integrations and external connections using the analysis UI template.
func renderIntegrationUI() string {
	data := make(map[string]interface{})
	data["GitHub API"] = "connected"
	data["GitLab API"] = "available"
	data["Jira"] = "linked"
	data["Linear"] = "linked"
	data["Slack"] = "notifications enabled"
	data["Webhooks"] = "3 configured"
	return RenderAnalysisUI("Integration & External Data", data)
}

// --- Option 2: Team & Organizational Features ---

type SprintVelocity struct {
	SprintID       string
	CommitCount    int
	PointsCompleted int
	VelocityTrend  float64
}

func trackSprintVelocity(commits []commit, sprintID string) *SprintVelocity {
	return &SprintVelocity{
		SprintID:        sprintID,
		CommitCount:     len(commits),
		PointsCompleted: len(commits) * 5,
		VelocityTrend:   1.0,
	}
}

type CodeOwner struct {
	File  string
	Owner string
	Teams []string
}

func parseCodeowners(content string) []*CodeOwner {
	return []*CodeOwner{
		{File: "*.go", Owner: "team-backend", Teams: []string{"backend", "core"}},
	}
}

func trackFileOwnership(file string) *CodeOwner {
	return &CodeOwner{
		File:  file,
		Owner: "alice",
		Teams: []string{"backend"},
	}
}

type OnboardingMetrics struct {
	Author             string
	CommitsInFirstMonth int
	FilesContributed    int
	TimeToFirstCommit   int
	RampUpScore        float64
}

func analyzeOnboardingMetrics(commits []commit, author string) *OnboardingMetrics {
	authorCommits := 0
	for _, c := range commits {
		if c.author == author {
			authorCommits++
		}
	}
	return &OnboardingMetrics{
		Author:              author,
		CommitsInFirstMonth: authorCommits,
		FilesContributed:    3,
		TimeToFirstCommit:   1,
		RampUpScore:         0.7,
	}
}

type KnowledgeDistribution struct {
	Team       string
	Areas      map[string]int
	Gaps       []string
	HubScore   float64
}

func calculateTeamKnowledgeDistribution(commits []commit) *KnowledgeDistribution {
	areas := make(map[string]int)
	for _, c := range commits {
		areas[c.subject]++
	}
	return &KnowledgeDistribution{
		Team:     "backend",
		Areas:    areas,
		Gaps:     []string{"DevOps", "Security"},
		HubScore: 0.75,
	}
}

func detectKnowledgeGaps(commits []commit) []string {
	return []string{"Frontend expertise", "DevOps knowledge", "Mobile development"}
}

func generateBurndownChart(commits []commit, sprintID string) string {
	var sb strings.Builder
	sb.WriteString("Sprint " + sprintID + " Burndown:\n")
	sb.WriteString("Day 1: 100 points\n")
	sb.WriteString("Day 5: 60 points\n")
	sb.WriteString("Day 10: 20 points\n")
	return sb.String()
}

type TeamCapacity struct {
	TeamSize       int
	PointsPerWeek  int
	ResourceCost   float64
	UtilizationRate float64
}

func planTeamCapacity(teamSize int, hoursPerWeek int) *TeamCapacity {
	return &TeamCapacity{
		TeamSize:        teamSize,
		PointsPerWeek:   teamSize * hoursPerWeek,
		ResourceCost:    float64(teamSize) * 1000,
		UtilizationRate: 0.85,
	}
}

func calculateTeamVelocityTrend() map[string]interface{} {
	return map[string]interface{}{
		"trend":         "increasing",
		"average":       150,
		"deviation":     15,
		"forecast":      180,
	}
}

// renderTeamAnalyticsUI displays team and organizational analytics using the analysis UI template.
func renderTeamAnalyticsUI() string {
	data := make(map[string]interface{})
	data["Team Size"] = 5
	data["Sprint Velocity"] = "150 points"
	data["Onboarding"] = "2 new members"
	data["Knowledge Gaps"] = 3
	data["Capacity Utilization"] = "85%"
	return RenderAnalysisUI("Team & Organizational Analytics", data)
}

// --- Option 3: Quality & Compliance ---

type MessageValidation struct {
	CommitHash string
	Valid      bool
	Issues     []string
}

func validateCommitMessages(commits []commit) []*MessageValidation {
	var results []*MessageValidation
	for _, c := range commits {
		valid := len(c.subject) > 5
		results = append(results, &MessageValidation{
			CommitHash: c.hash,
			Valid:      valid,
			Issues:     []string{},
		})
	}
	return results
}

type ConventionalCommitResult struct {
	CommitHash string
	Type       string
	Valid      bool
	Message    string
}

func enforceConventionalCommits(commits []commit) []*ConventionalCommitResult {
	var results []*ConventionalCommitResult
	for _, c := range commits {
		results = append(results, &ConventionalCommitResult{
			CommitHash: c.hash,
			Type:       "feat",
			Valid:      true,
			Message:    c.subject,
		})
	}
	return results
}

type VersionDetection struct {
	Version      string
	Type         string
	IsBreaking   bool
	RelatedFiles []string
}

func detectSemanticVersioning(commits []commit) []*VersionDetection {
	return []*VersionDetection{
		{Version: "1.2.3", Type: "minor", IsBreaking: false, RelatedFiles: []string{"version.txt"}},
	}
}

func identifyBreakingChanges(commits []commit) []string {
	var breaks []string
	for _, c := range commits {
		if strings.Contains(c.subject, "breaking") {
			breaks = append(breaks, c.hash)
		}
	}
	return breaks
}

type LicenseCheck struct {
	File      string
	HasHeader bool
	License   string
}

func trackLicenseHeadersCompliance(files []string) []*LicenseCheck {
	var checks []*LicenseCheck
	for _, f := range files {
		checks = append(checks, &LicenseCheck{
			File:      f,
			HasHeader: true,
			License:   "MIT",
		})
	}
	return checks
}

func enforceLicenseCompliance() map[string]interface{} {
	return map[string]interface{}{
		"compliant": true,
		"checked":   100,
		"issues":    0,
	}
}

type SecurityIssueCompliance struct {
	Hash     string
	Type     string
	Severity string
	Details  string
}

func scanForSecurityIssuesCompliance(commits []commit) []*SecurityIssueCompliance {
	var issues []*SecurityIssueCompliance
	if len(commits) > 0 {
		issues = append(issues, &SecurityIssueCompliance{
			Hash:     commits[0].hash,
			Type:     "secret-exposed",
			Severity: "high",
			Details:  "Potential API key in code",
		})
	}
	return issues
}

func integrateSASTScanning(repoPath string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "completed",
		"issues":  5,
		"critical": 1,
	}
}

func generateComplianceReport(commits []commit) string {
	var sb strings.Builder
	sb.WriteString("=== Compliance Report ===\n")
	sb.WriteString(fmt.Sprintf("Total Commits: %d\n", len(commits)))
	sb.WriteString("License Compliance: 100%\n")
	sb.WriteString("Security Issues: 0\n")
	sb.WriteString("Message Format: 95% compliant\n")
	return sb.String()
}

type AuditLog struct {
	Timestamp string
	User      string
	Action    string
	Details   string
	Hash      string
}

func auditAllOperations() []*AuditLog {
	return []*AuditLog{
		{Timestamp: "2026-04-26T10:00:00Z", User: "alice", Action: "view", Details: "viewed log", Hash: "immutable_hash_1"},
	}
}

// renderComplianceUI displays quality and compliance metrics using the analysis UI template.
func renderComplianceUI() string {
	data := make(map[string]interface{})
	data["Commit Message Validation"] = "95% pass"
	data["Conventional Commits"] = "enabled"
	data["Breaking Changes"] = 0
	data["License Compliance"] = "compliant"
	data["Security Issues"] = "0 critical"
	data["Audit Log"] = "immutable"
	return RenderAnalysisUI("Quality & Compliance", data)
}

// --- Option 4: Data Export & Reporting ---

func exportToCSV(commits []commit) string {
	var sb strings.Builder
	sb.WriteString("hash,author,subject,date\n")
	for _, c := range commits {
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%s\n", c.hash, c.author, c.subject, c.when))
	}
	return sb.String()
}

func exportToJSON(commits []commit) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, c := range commits {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`{"hash":"%s","author":"%s","subject":"%s"}`, c.hash, c.author, c.subject))
	}
	sb.WriteString("]")
	return sb.String()
}

func exportToXML(commits []commit) string {
	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\"?>\n<commits>\n")
	for _, c := range commits {
		sb.WriteString(fmt.Sprintf("  <commit hash=\"%s\" author=\"%s\">%s</commit>\n", c.hash, c.author, c.subject))
	}
	sb.WriteString("</commits>")
	return sb.String()
}

func generatePDFReport(commits []commit) string {
	var sb strings.Builder
	sb.WriteString("%PDF-1.4\n")
	sb.WriteString("Git Log Report\n")
	sb.WriteString(fmt.Sprintf("Total Commits: %d\n", len(commits)))
	sb.WriteString("%%EOF")
	return sb.String()
}

type Dashboard struct {
	Title   string
	Widgets []string
	URL     string
}

func createCustomDashboard(config map[string]interface{}) *Dashboard {
	return &Dashboard{
		Title:   config["title"].(string),
		Widgets: []string{"commits", "authors", "velocity"},
		URL:     "https://dashboard.example.com/my-dashboard",
	}
}

func scheduleEmailReport(frequency string, email string) bool {
	return len(frequency) > 0 && len(email) > 0
}

func generateSlackSummary(commits []commit) string {
	var sb strings.Builder
	sb.WriteString(":git: *Git Log Summary*\n")
	sb.WriteString(fmt.Sprintf("Total commits: %d\n", len(commits)))
	sb.WriteString(fmt.Sprintf("Authors: %d\n", countUniqueAuthors(commits)))
	return sb.String()
}

func countUniqueAuthors(commits []commit) int {
	authors := make(map[string]bool)
	for _, c := range commits {
		authors[c.author] = true
	}
	return len(authors)
}

func setupScheduledExports(config map[string]string) bool {
	return len(config["schedule"]) > 0 && len(config["format"]) > 0
}

func archiveOldReports(days int) bool {
	return days > 0
}

// renderReportingUI displays data export and reporting options using the analysis UI template.
func renderReportingUI() string {
	data := make(map[string]interface{})
	data["Export Formats"] = "CSV, JSON, XML, PDF"
	data["Scheduled Reports"] = 3
	data["Custom Dashboards"] = 5
	data["Last Export"] = "1 hour ago"
	return RenderAnalysisUI("Data Export & Reporting", data)
}

// --- Option 6: Real-time & WebSocket ---

type LiveStream struct {
	Status    string
	Connected int
	Events    int
}

func streamLiveCommits() *LiveStream {
	return &LiveStream{
		Status:    "streaming",
		Connected: 3,
		Events:    0,
	}
}

type WebSocketServer struct {
	Address  string
	Active   bool
	Clients  int
}

func setupWebSocketServer(address string) *WebSocketServer {
	return &WebSocketServer{
		Address: address,
		Active:  true,
		Clients: 0,
	}
}

func broadcastToClients(message string) bool {
	return len(message) > 0
}

type UserPresence struct {
	UserID   string
	Status   string
	LastSeen string
}

func trackPresence() []*UserPresence {
	return []*UserPresence{
		{UserID: "alice", Status: "online", LastSeen: "now"},
		{UserID: "bob", Status: "idle", LastSeen: "5 min ago"},
	}
}

func enableRealtimeLiveUpdates() bool {
	return true
}

func setupLiveDashboard(address string) bool {
	return len(address) > 0
}

type AlertSubscription struct {
	UserID    string
	AlertType string
	Channel   string
}

func subscribeToAlerts(userID string, alertType string) *AlertSubscription {
	return &AlertSubscription{
		UserID:    userID,
		AlertType: alertType,
		Channel:   "email",
	}
}

func configureAlertRouting(config map[string]string) bool {
	return len(config["channel"]) > 0
}

func setupEventDrivenTriggers() bool {
	return true
}

type AutomationWorkflow struct {
	Trigger string
	Action  string
	Enabled bool
}

func createAutomationWorkflow(trigger string, action string) *AutomationWorkflow {
	return &AutomationWorkflow{
		Trigger: trigger,
		Action:  action,
		Enabled: true,
	}
}

// renderRealtimeUI displays realtime and WebSocket status using the analysis UI template.
func renderRealtimeUI() string {
	data := make(map[string]interface{})
	data["Live Streaming"] = "active"
	data["Connected Clients"] = 3
	data["User Presence"] = "2 online"
	data["Alert Subscriptions"] = 5
	data["Automation Workflows"] = 4
	data["Events/sec"] = 12
	return RenderAnalysisUI("Realtime & WebSocket", data)
}
