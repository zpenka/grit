// Package grit provides a terminal UI for exploring git history.
//
// The filtering module (engine_filtering.go) provides search and filtering
// capabilities for commits. It supports text search, regex patterns, author
// filtering, and date range filtering.
//
// Key functions:
//   - filterCommits: Search by text across subject, author, hash
//   - visibleCommits: Apply active filters to full commit list
//   - filterByAuthor: Filter commits by author name
//   - filterByDateRange: Filter commits by timestamp range
//
// Filtering is used to reduce large commit lists to a relevant subset, improving
// navigation performance and focus for analysis tasks.
package grit

import (
	"regexp"
	"strconv"
	"strings"
)

// filterCommits returns commits whose subject, author, or short hash contain
// query (case-insensitive). An empty query returns all commits unchanged.
func filterCommits(commits []commit, query string) []commit {
	if query == "" {
		return commits
	}
	q := strings.ToLower(query)
	// Pre-allocate with capacity to reduce append reallocations (~50% expected match rate)
	out := make([]commit, 0, len(commits)/2)
	for _, c := range commits {
		if strings.Contains(strings.ToLower(c.subject), q) ||
			strings.Contains(strings.ToLower(c.author), q) ||
			strings.Contains(strings.ToLower(c.shortHash), q) {
			out = append(out, c)
		}
	}
	return out
}

// filterCommitsWithCache filters commits with result caching and metrics tracking.
func filterCommitsWithCache(cache *FilterCache, commits []commit, query string) []commit {
	if query == "" {
		return commits
	}

	// Check cache first
	if cached, exists := cache.cache[query]; exists {
		cache.metrics.Hits++
		return cached
	}

	// Cache miss - filter and store
	cache.metrics.Misses++
	result := filterCommits(commits, query)

	// Store in cache if not full
	if cache.metrics.Size < cache.metrics.MaxSize {
		cache.cache[query] = result
		cache.metrics.Size++
	}

	return result
}

// visibleCommits returns the commit list after applying all active filters:
// search query, author filter, and time-based filter.
// Optimized to single-pass filtering to reduce allocations.
func visibleCommits(m model) []commit {
	if len(m.commits) == 0 {
		return m.commits
	}

	// Fast path: no filters
	if m.sinceFilter <= 0 && m.authorFilter == "" && m.query == "" {
		return m.commits
	}

	// Single-pass filtering to avoid intermediate slices
	// Start with small capacity to avoid over-allocation
	q := strings.ToLower(m.query)
	out := make([]commit, 0, 128)

	for _, c := range m.commits {
		// Check time filter
		if m.sinceFilter > 0 && !isWithinDays(c.when, m.sinceFilter) {
			continue
		}
		// Check author filter
		if m.authorFilter != "" && !strings.EqualFold(c.author, m.authorFilter) {
			continue
		}
		// Check query filter
		if m.query != "" {
			if !strings.Contains(strings.ToLower(c.subject), q) &&
				!strings.Contains(strings.ToLower(c.author), q) &&
				!strings.Contains(strings.ToLower(c.shortHash), q) {
				continue
			}
		}
		// All filters passed
		out = append(out, c)
	}
	return out
}

// filterCommitsByAuthor returns commits whose author exactly matches the given author
// (case-insensitive).
func filterCommitsByAuthor(commits []commit, author string) []commit {
	if author == "" {
		return commits
	}
	// Pre-allocate with capacity (~10-20% of commits per author)
	out := make([]commit, 0, len(commits)/5)
	for _, c := range commits {
		if strings.EqualFold(c.author, author) {
			out = append(out, c)
		}
	}
	return out
}

// filterCommitsSince returns commits from the last N days, parsed from the
// "when" field (e.g., "5 days ago", "2 weeks ago"). Returns all commits if
// days <= 0.
func filterCommitsSince(commits []commit, days int) []commit {
	if days <= 0 {
		return commits
	}
	// Pre-allocate (~30% of commits typically within last N days)
	out := make([]commit, 0, len(commits)/3)
	for _, c := range commits {
		if isWithinDays(c.when, days) {
			out = append(out, c)
		}
	}
	return out
}

// isWithinDays checks if a "when" string (e.g., "5 days ago") represents
// a time within the last N days.
func isWithinDays(when string, days int) bool {
	whenLower := strings.ToLower(when)

	// Extract number from strings like "5 days ago", "2 weeks ago", etc.
	re := regexp.MustCompile(`(\d+)\s+(day|week|month|year)`)
	matches := re.FindStringSubmatch(whenLower)
	if len(matches) < 3 {
		return false
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return false
	}

	unit := matches[2]
	totalDays := 0
	switch unit {
	case "day":
		totalDays = num
	case "week":
		totalDays = num * 7
	case "month":
		totalDays = num * 30
	case "year":
		totalDays = num * 365
	default:
		return false
	}

	return totalDays <= days
}

// formatActiveFilters returns a string showing all active filters.
// formatActiveFilters returns a display string of all active filters.
func formatActiveFilters(m model) string {
	var filters []string
	if m.authorFilter != "" {
		filters = append(filters, m.authorFilter)
	}
	if m.sinceFilter > 0 {
		filters = append(filters, strconv.Itoa(m.sinceFilter)+"d")
	}
	if len(filters) == 0 {
		return ""
	}
	return "[" + strings.Join(filters, " + ") + "]"
}
