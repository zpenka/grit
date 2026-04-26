package grit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// newModel creates a fresh model instance for the given repository path.
func newModel(repoPath string) model {
	return model{
		repoPath: repoPath,
		focus:    panelList,
	}
}


// truncate cuts s to at most max visible runes, appending "…" if shortened.
// truncate shortens a string to max runes and appends "…" if truncated.
func truncate(s string, max int) string {
	r := []rune(s)
	if max <= 0 {
		return ""
	}
	if len(r) <= max {
		return s
	}
	if max == 1 {
		return "…"
	}
	return string(r[:max-1]) + "…"
}

// firstWord returns the first space-delimited word of s.
func firstWord(s string) string {
	if i := strings.Index(s, " "); i >= 0 {
		return s[:i]
	}
	return s
}




// parseBranches parses the output of "git branch -a", stripping the current-branch
// marker (*) and skipping ref-pointer lines (e.g. "origin/HEAD -> origin/main").
func parseBranches(output string) []string {
	var branches []string
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimPrefix(line, "* ")
		if strings.Contains(line, " -> ") {
			continue
		}
		branches = append(branches, line)
	}
	return branches
}

// parseCurrentBranch returns the name of the currently checked-out branch from
// "git branch -a" output (the line prefixed with "* ").
func parseCurrentBranch(output string) string {
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if strings.HasPrefix(line, "* ") {
			return strings.TrimPrefix(line, "* ")
		}
	}
	return ""
}

// parseBlameLine parses one line of "git blame --date=short" output.
// Format: "hash (Author Name   YYYY-MM-DD  linenum) content"
func parseBlameLine(line string) (blameLine, bool) {
	paren := strings.Index(line, "(")
	close := strings.Index(line, ")")
	if paren < 0 || close < 0 || close < paren {
		return blameLine{}, false
	}
	hash := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line[:paren]), "^"))
	if len(hash) > 7 {
		hash = hash[:7]
	}
	meta := strings.Fields(line[paren+1 : close])
	if len(meta) < 3 {
		return blameLine{}, false
	}
	// last field: line number; second-to-last: date; rest: author
	lineNum, err := strconv.Atoi(meta[len(meta)-1])
	if err != nil {
		return blameLine{}, false
	}
	date := meta[len(meta)-2]
	author := strings.Join(meta[:len(meta)-2], " ")
	var text string
	if close+2 <= len(line) {
		text = line[close+2:]
	} else if close+1 < len(line) {
		text = line[close+1:]
	}
	return blameLine{
		shortHash: hash,
		author:    author,
		date:      date,
		lineNum:   lineNum,
		text:      text,
	}, true
}

// parseBlame parses all lines from "git blame --date=short" output,
// skipping any lines that don't match the expected format.
func parseBlame(output string) []blameLine {
	var lines []blameLine
	for _, line := range strings.Split(output, "\n") {
		if bl, ok := parseBlameLine(line); ok {
			lines = append(lines, bl)
		}
	}
	return lines
}

// currentFile returns the path of the file whose diff section is currently
// visible, based on the last fileItem whose diffIdx is <= diffOffset.
// currentFile returns the path of the currently selected file in file view mode.
func currentFile(m model) string {
	if len(m.fileItems) == 0 {
		return ""
	}
	cur := m.fileItems[0].path
	for _, fi := range m.fileItems {
		if fi.diffIdx <= m.diffOffset {
			cur = fi.path
		}
	}
	return cur
}

// parseCount converts a digit string to a navigation count.
// Empty string or zero returns 1; values above 200 are capped.
func parseCount(s string) int {
	if s == "" {
		return 1
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 1
	}
	if n > 200 {
		return 200
	}
	return n
}



// commitStats calculates statistics for a commit's diff.
// commitStats calculates insertion, deletion, and file change counts from diff lines.
func commitStats(lines []diffLine) commitStatistics {
	var stats commitStatistics
	for _, line := range lines {
		if line.kind == lineMeta && strings.HasPrefix(line.text, "diff --git") {
			stats.filesChanged++
		}
		if line.kind == lineAdded {
			stats.insertions++
		}
		if line.kind == lineRemoved {
			stats.deletions++
		}
	}
	return stats
}

// generateCommitMessage generates a suggested commit message based on the diff.
func generateCommitMessage(lines []diffLine, filename string) string {
	stats := commitStats(lines)
	isNew := false
	isDeleted := false
	isBreaking := false

	for _, line := range lines {
		if strings.Contains(line.text, "new file mode") {
			isNew = true
		}
		if strings.Contains(line.text, "deleted file mode") {
			isDeleted = true
		}
		// Simple breaking change detection: removed function/interface definitions
		if line.kind == lineRemoved &&
			(strings.Contains(line.text, "func ") || strings.Contains(line.text, "interface")) {
			isBreaking = true
		}
	}

	scope := filenameToScope(filename)
	var verb string

	switch {
	case isDeleted:
		verb = "remove"
	case isNew:
		verb = "add"
	case stats.deletions > stats.insertions*2:
		verb = "refactor"
	default:
		verb = "update"
	}

	msg := verb
	if scope != "" {
		msg += "(" + scope + ")"
	}
	if isBreaking {
		msg += "!"
	}
	msg += ": " + capitalizeFirst(verb) + " changes"

	return msg
}

// filenameToScope extracts a scope from a filename (e.g., "auth.go" -> "auth").
// filenameToScope extracts the directory scope from a file path for navigation.
func filenameToScope(filename string) string {
	if filename == "" {
		return ""
	}
	base := filename
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	if idx := strings.LastIndex(base, "."); idx > 0 {
		base = base[:idx]
	}
	return strings.ToLower(base)
}

// capitalizeFirst capitalizes the first letter of a string.
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}


// detectLanguage detects the programming language from a filename.
func detectLanguage(filename string) string {
	ext := filename
	if idx := strings.LastIndex(filename, "."); idx >= 0 {
		ext = filename[idx:]
	}

	langMap := map[string]string{
		".go":         "go",
		".py":         "python",
		".js":         "javascript",
		".ts":         "typescript",
		".rb":         "ruby",
		".java":       "java",
		".cpp":        "cpp",
		".c":          "c",
		".rs":         "rust",
		".sh":         "bash",
		".sql":        "sql",
		".html":       "html",
		".css":        "css",
		".json":       "json",
		".yaml":       "yaml",
		".yml":        "yaml",
		".xml":        "xml",
		".md":         "markdown",
		"Makefile":    "makefile",
		"Dockerfile":  "dockerfile",
		".gitignore":  "gitignore",
		".env":        "dotenv",
	}

	if lang, ok := langMap[ext]; ok {
		return lang
	}
	if lang, ok := langMap[filename]; ok {
		return lang
	}

	// Default to text
	return "text"
}

// miniMapPosition calculates the position of a scroll indicator (0-height).
// miniMapPosition calculates the position indicator for a visual minimap of commits.
func miniMapPosition(cursor, panelHeight, totalCommits int) int {
	if totalCommits <= 1 {
		return 0
	}
	if panelHeight <= 1 {
		return 0
	}

	// Map cursor position to panel height
	position := (cursor * (panelHeight - 1)) / (totalCommits - 1)
	if position > panelHeight-1 {
		position = panelHeight - 1
	}
	return position
}

// --- diffStatBadge ---

// diffStatBadge formats commit statistics as a compact badge (e.g., "3 files +10 -5").
func diffStatBadge(stats commitStatistics) string {
	var parts []string
	if stats.filesChanged > 0 {
		parts = append(parts, fmt.Sprintf("%d file%s", stats.filesChanged, pluralize(stats.filesChanged)))
	}
	if stats.insertions > 0 {
		parts = append(parts, fmt.Sprintf("+%d", stats.insertions))
	}
	if stats.deletions > 0 {
		parts = append(parts, fmt.Sprintf("-%d", stats.deletions))
	}
	if len(parts) == 0 {
		return "0 changes"
	}
	return strings.Join(parts, " ")
}

// pluralize returns "s" if count != 1, else "".
// pluralize returns "s" if count != 1, else empty string, for grammatical pluralization.
func pluralize(count int) string {
	if count != 1 {
		return "s"
	}
	return ""
}

// --- goToCommit ---


// --- copyAsPatch ---

// copyAsPatch generates a patch file format from a commit hash and diff lines.
// copyAsPatch formats diff lines as a git patch with headers.
func copyAsPatch(hash string, lines []diffLine) string {
	var sb strings.Builder
	sb.WriteString("From: " + hash + "\n")
	sb.WriteString("Subject: Commit patch\n")
	sb.WriteString("\n")
	for _, line := range lines {
		sb.WriteString(line.text)
		sb.WriteString("\n")
	}
	return sb.String()
}

// --- parseGitReferences ---

// parseGitReferences extracts issue/PR numbers from a commit message (e.g., #123, fixes #456).
func parseGitReferences(msg string) []string {
	re := regexp.MustCompile(`#(\d+)`)
	matches := re.FindAllStringSubmatch(msg, -1)
	var refs []string
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			refs = append(refs, match[1])
			seen[match[1]] = true
		}
	}
	return refs
}

// --- isMergeCommit ---

// isMergeCommit checks if a commit is a merge commit.
// isMergeCommit detects if a commit is a merge based on diff metadata.
func isMergeCommit(lines []diffLine) bool {
	for _, line := range lines {
		if strings.HasPrefix(line.text, "Merge:") {
			return true
		}
		if strings.Contains(strings.ToLower(line.text), "merge branch") {
			return true
		}
	}
	return false
}

// --- getMergeParents ---

// getMergeParents extracts parent hashes from a merge commit.
// getMergeParents extracts the parent commit hashes from a merge commit's diff.
func getMergeParents(lines []diffLine) []string {
	for _, line := range lines {
		if strings.HasPrefix(line.text, "Merge:") {
			parts := strings.Fields(strings.TrimPrefix(line.text, "Merge:"))
			if len(parts) >= 2 {
				return parts[:2]
			}
		}
	}
	return nil
}

// --- hunk structure and parseHunks ---

type hunk struct {
	startLine int
	endLine   int
	header    string
}

// parseHunks extracts all hunks from diff lines.
func parseHunks(lines []diffLine) []hunk {
	var hunks []hunk
	var lastHunk hunk
	hunkCount := 0

	for i, line := range lines {
		if line.kind == lineHunk {
			if hunkCount > 0 {
				hunks = append(hunks, lastHunk)
			}
			lastHunk = hunk{
				startLine: i,
				endLine:   i,
				header:    line.text,
			}
			hunkCount++
		} else if hunkCount > 0 {
			lastHunk.endLine = i
		}
	}
	if hunkCount > 0 {
		hunks = append(hunks, lastHunk)
	}
	return hunks
}

// --- toggleLineComment ---

// toggleLineComment adds or removes a comment on a specific diff line.
func toggleLineComment(m model, lineIdx int, comment string) model {
	if m.comments == nil {
		m.comments = make(map[int]string)
	}
	if comment == "" {
		delete(m.comments, lineIdx)
	} else {
		m.comments[lineIdx] = comment
	}
	return m
}

// --- compileRegex ---

// compileRegex compiles a regex pattern for search.

// --- parseDateRange ---

// parseDateRange parses a date or date range (e.g., "2024-01-15" or "2024-01-01..2024-01-31").
func parseDateRange(input string) (*time.Time, *time.Time, error) {
	if input == "" {
		return nil, nil, nil
	}

	// Check for range format
	if strings.Contains(input, "..") {
		parts := strings.Split(input, "..")
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid date range format")
		}
		start, err1 := time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
		end, err2 := time.Parse("2006-01-02", strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return nil, nil, fmt.Errorf("invalid date format")
		}
		return &start, &end, nil
	}

	// Single date
	date, err := time.Parse("2006-01-02", input)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid date format")
	}
	return &date, nil, nil
}

// --- filterCommitsByFile ---

// filterCommitsByFile returns commits that touched the specified file.
// Note: This is infrastructure-ready; actual filtering requires git queries.
func filterCommitsByFile(commits []commit, file string) []commit {
	if file == "" {
		return commits
	}
	// This would typically query git for commits touching the file
	// For now, return infrastructure (called from UI)
	return []commit{}
}

// --- parseTags ---

// parseTags parses git tag output (format: "hash tagname" per line).
func parseTags(output string) []string {
	var tags []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Handle both "hash tagname" and "tagname" formats
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			// Take the last non-hash part
			tag := parts[len(parts)-1]
			if !strings.HasPrefix(tag, "[") { // Skip metadata
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

// ===== OPTION 1: UI INTEGRATION =====

// renderStatsBadgeInList renders a compact stats badge for a commit list row.
