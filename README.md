# grit - Terminal UI for Git History

A powerful, feature-rich terminal UI for exploring git repositories with 113+ integrated features. Built in Go with zero external dependencies (besides Bubble Tea for TUI).

## Key Features

- **Interactive Browsing** - Explore commits, diffs, blame information
- **Analysis** - Code ownership, hotspots, team activity heatmaps
- **Git Operations** - Rebase, cherry-pick, reset, amend workflows
- **Export & Integration** - GitHub/Jira linking, CSV/JSON/XML export
- **Team Metrics** - Velocity tracking, collaboration analytics
- **Visualization** - Commit graphs, timelines, complexity analysis

## Quick Start

```bash
go build -o grit .
./grit
```

Run from any git repository. Requires Go 1.21+.

## Installation

```bash
go install github.com/zpenka/grit@latest
```

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
