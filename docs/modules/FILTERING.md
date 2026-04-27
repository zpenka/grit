# Filtering Module (`engine_filtering.go`)

## Purpose

The filtering module provides search and filtering capabilities for commits. It supports full-text search, regex patterns, author filtering, date range filtering, and extension-based filtering for quick navigation in large repositories.

## Key Functions

### `filterCommits(commits []commit, query string) []commit`
Text-based search across subject, author, and short hash
- Case-insensitive matching
- Empty query returns all commits
- Returns filtered slice preserving original order

### `visibleCommits(m model) []commit`
Apply all active filters to full commit list
- Combines query, author, date range, and extension filters
- Used for rendering and cursor validation
- Fast computation with early termination

### `filterByAuthor(commits []commit, author string) []commit`
Filter commits by exact author name
- Case-sensitive matching
- Useful for team analysis

### `filterByDateRange(commits []commit, since, until int64) []commit`
Filter commits by timestamp range
- Unix timestamps for boundaries
- Inclusive range
- Useful for sprint-based analysis

### `filterByExtension(commits []commit, ext string) []commit`
Filter commits that modified files with given extension
- Requires parsing diff to identify changed files
- Cached to avoid redundant parsing

## Design Decisions

### Pure Functions
All filtering functions take input and return filtered results without mutating state. Makes them easily testable and composable.

### Combinable Filters
Filters can be applied sequentially. The `visibleCommits()` function chains them together in a predictable order.

### Case Handling
Text search is case-insensitive; author search is case-sensitive to match git's behavior.

## Dependencies

- Standard library: `regexp`, `strings`
- `dcache` from optimization module for diff caching

## Testing

Tests cover:
- Empty commit lists
- Empty query (return all)
- Partial matches and full matches
- Case sensitivity
- Regex patterns
- Multiple simultaneous filters

## Examples

```go
// Search for "fix" in commit subjects
results := filterCommits(commits, "fix")

// Find all commits by John Doe
authors := filterByAuthor(commits, "John Doe")

// Commits in last month
now := time.Now().Unix()
month := now - (30 * 24 * 60 * 60)
recent := filterByDateRange(commits, month, now)

// All visible commits with active filters
visible := visibleCommits(m)
```

## Performance Considerations

- Search is O(n) in commit count
- Early termination when query is empty
- Regex patterns compiled once and cached
- Extension filtering requires diff parsing (cached)

## Integration Points

- Navigation module uses `visibleCommits()` for cursor validation
- Rendering module uses results for display
- Keybinding handler updates filter state in model

## Future Extensions

- Regex search with syntax validation
- Filter history/recency
- Complex boolean filters (AND/OR/NOT)
- Filter save/load for named searches
