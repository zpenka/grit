# Parsing Module (`engine_parsing.go`)

## Purpose

The parsing module converts raw git command output into strongly-typed Go structs. It serves as the bridge between git and the internal type system, transforming text into data structures used throughout the application.

## Key Functions

### `parseCommits(output string) []commit`
Parses `git log` output in the format: `%H%x00%h%x00%an%x00%ar%x00%s`
- Full hash, short hash, author, relative time, subject
- Handles large commit histories efficiently
- Returns commit slice ready for filtering and display

### `parseDiff(output string) []diffLine`
Parses unified diff format from `git show` or `git diff`
- Classifies lines: context, added, removed, hunk header, metadata
- Extracts file names and change statistics
- Handles binary files gracefully

### `parseFileItems(output string) []fileItem`
Parses `git ls-tree` output for file trees
- Extracts file paths and permissions
- Supports recursive tree traversal
- Used by file browser panel

## Design Decisions

### Null-Byte Separators
Uses `%x00` (null byte) in git log format for reliable field separation. This avoids issues with colons or pipes in commit messages.

### No Validation
Input validation happens at system boundaries (git commands). Internal parsing functions trust their inputs based on git's guarantees.

### Streaming-Friendly
Functions are designed to process output incrementally, supporting large repositories with potentially thousands of commits.

## Dependencies

- Standard library: `strings`, `bytes`
- No external dependencies

## Testing

Tests cover:
- Large commit histories (100+ commits)
- Various commit message formats
- Diff formats (unified, context, etc.)
- Edge cases: empty diffs, binary files

## Examples

```go
// Parse commit log
logOutput := `abc123\x00abc\x00John Doe\x005 days ago\x00Fix bug`
commits := parseCommits(logOutput)

// Parse diff
diffOutput := `diff --git a/file.go b/file.go
index 1234567..89abcdef
--- a/file.go
+++ b/file.go
@@ -10,3 +10,4 @@
 context line
-old line
+new line`
lines := parseDiff(diffOutput)
```

## Performance Considerations

- Parsing is O(n) in output size
- No copying of input; uses string slices
- Memory usage proportional to number of commits
- Cache diff parsing results in `dcache`

## Future Extensions

- Support for custom git log formats
- Streaming parser for very large diffs
- Incremental parsing for large histories
