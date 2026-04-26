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
	var out []commit
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
func visibleCommits(m model) []commit {
	result := m.commits
	// Apply time filter first
	if m.sinceFilter > 0 {
		result = filterCommitsSince(result, m.sinceFilter)
	}
	// Apply author filter
	if m.authorFilter != "" {
		result = filterCommitsByAuthor(result, m.authorFilter)
	}
	// Apply search query filter last
	result = filterCommits(result, m.query)
	return result
}

// filterCommitsByAuthor returns commits whose author exactly matches the given author
// (case-insensitive).
func filterCommitsByAuthor(commits []commit, author string) []commit {
	if author == "" {
		return commits
	}
	var out []commit
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
	var out []commit
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
