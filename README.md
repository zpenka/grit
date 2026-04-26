# grit

A terminal UI for exploring your git history. Browse commits, diffs, blame, analytics, and more—all in the terminal.

**312+ features** including:
- Interactive commit browsing
- Diff viewing with syntax highlighting
- Author and time-based filtering
- Git operations (rebase, cherry-pick, revert)
- Blame view for code ownership
- Analytics (hotspots, code ownership, complexity)
- Team collaboration metrics
- Export to markdown, patches, JIRA

## Installation

```bash
go install github.com/zpenka/grit@latest
```

Or build from source:

```bash
git clone github.com/zpenka/grit
cd grit
go build -o grit .
./grit
```

## Usage

```bash
grit
```

Run from within any git repository. Requires Go 1.21+.

## Requirements

- Go 1.21+
- A git repository
- Terminal with 24-bit color support (recommended)

## Testing

```bash
go test ./...           # Run all tests
go test -cover ./...    # Coverage report
```

## License

MIT
