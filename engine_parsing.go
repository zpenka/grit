// Package grit provides a terminal UI for exploring git history.
//
// The parsing module (engine_parsing.go) handles conversion of raw git command
// output into structured data types. It parses commit logs, diffs, and file
// listings from git commands, transforming text output into strongly-typed Go
// structs for use throughout the application.
//
// Key functions:
//   - parseCommits: Convert git log output to commit structs
//   - parseDiff: Parse unified diff format into diff line structs
//   - parseFileItems: Extract file paths from git ls-tree output
//
// The parsing module serves as the bridge between raw git data and the internal
// type system. All data entering the system flows through these functions.
package grit

import (
	"strings"
)

// parseCommits parses output of:
//
//	git log --format="%H%x00%h%x00%an%x00%ar%x00%s"
func parseCommits(output string) []commit {
	var commits []commit
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "\x00", 5)
		if len(parts) < 5 {
			continue
		}
		commits = append(commits, commit{
			hash:      parts[0],
			shortHash: parts[1],
			author:    parts[2],
			when:      parts[3],
			subject:   parts[4],
		})
	}
	return commits
}

// parseCommitsWithPool parses commits using memory pooling for efficiency.
func parseCommitsWithPool(output string) []commit {
	pool := NewMemoryPool(100)
	var commits []commit

	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "\x00", 5)
		if len(parts) < 5 {
			continue
		}

		// Get commit from pool or create new
		c := pool.Get(func() interface{} {
			return &commit{}
		}).(*commit)

		// Populate fields
		c.hash = parts[0]
		c.shortHash = parts[1]
		c.author = parts[2]
		c.when = parts[3]
		c.subject = parts[4]

		commits = append(commits, *c)

		// Return to pool for reuse
		pool.Put(c)
	}
	return commits
}

// parseDiff classifies each line of a unified diff by type.
func parseDiff(raw string) []diffLine {
	var lines []diffLine
	for _, text := range strings.Split(raw, "\n") {
		var kind lineKind
		switch {
		case strings.HasPrefix(text, "@@"):
			kind = lineHunk
		case strings.HasPrefix(text, "diff "),
			strings.HasPrefix(text, "index "),
			strings.HasPrefix(text, "--- "),
			strings.HasPrefix(text, "+++ "),
			strings.HasPrefix(text, "new file"),
			strings.HasPrefix(text, "deleted file"),
			strings.HasPrefix(text, "similarity"),
			strings.HasPrefix(text, "rename"):
			kind = lineMeta
		case strings.HasPrefix(text, "+"):
			kind = lineAdded
		case strings.HasPrefix(text, "-"):
			kind = lineRemoved
		default:
			kind = lineContext
		}
		lines = append(lines, diffLine{kind, text})
	}
	return lines
}

// processDiffBatch processes diff lines using batch processing for efficiency.
func processDiffBatch(processor *BatchProcessor, lines []diffLine) []diffLine {
	var results []diffLine

	for _, line := range lines {
		processor.Add(line)

		// Process batch when full
		if processor.IsFull() {
			batch := processor.Get()
			for _, item := range batch {
				if dl, ok := item.(diffLine); ok {
					results = append(results, dl)
				}
			}
		}
	}

	// Process remaining items
	remaining := processor.Get()
	for _, item := range remaining {
		if dl, ok := item.(diffLine); ok {
			results = append(results, dl)
		}
	}

	return results
}

// parseFileItemsFromDiff scans diffLines for "diff --git" boundaries and returns each
// file's path and the index of its boundary line in diffLines.
func parseFileItemsFromDiff(lines []diffLine) []fileItem {
	var items []fileItem
	for i, line := range lines {
		if line.kind != lineMeta || !strings.HasPrefix(line.text, "diff --git ") {
			continue
		}
		parts := strings.Fields(line.text)
		if len(parts) < 4 {
			continue
		}
		path := strings.TrimPrefix(parts[3], "b/")
		items = append(items, fileItem{path: path, diffIdx: i})
	}
	return items
}
