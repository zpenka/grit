# grit - Terminal UI for Git History

A powerful, feature-rich terminal UI for exploring git repositories with **113+ integrated features**. Built in Go with zero external dependencies (besides Bubble Tea for TUI).

## Key Features

- **Interactive Browsing** - Explore commits, diffs, blame information interactively
- **Code Analysis** - Code ownership, hotspots, team activity heatmaps
- **Git Operations** - Interactive rebase, cherry-pick, reset, amend workflows
- **Export & Integration** - GitHub/Jira linking, CSV/JSON/XML export
- **Team Metrics** - Velocity tracking, collaboration analytics, reviewer suggestions
- **Visualization** - Contributor flamegraphs, activity heatmaps, complexity analysis
- **Security** - Secret scanning, GPG signature verification, license header tracking
- **Compliance** - Semantic versioning detection, changelog auto-generation, audit trails

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
go test ./...                    # Run all tests
go test -cover ./...             # With coverage report
go test -run TestFeatureName     # Run specific test
```

All 840 tests pass with 100% success rate. Tests run automatically on PRs via GitHub Actions.

## Documentation

**Start here:** [**DOCS.md**](DOCS.md) - Complete documentation map organized by audience and task

- **[DEVELOPER.md](DEVELOPER.md)** - Getting started and development guide
- **[CLAUDE.md](CLAUDE.md)** - High-level overview for AI assistants
- **[CHANGELOG.md](CHANGELOG.md)** - Detailed changelog of changes by phase
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contribution guidelines and code standards
- **[VERSION.md](VERSION.md)** - Release process and semantic versioning
- **[docs/modules/](docs/modules/)** - Detailed documentation for all 14 modules:
  - Core infrastructure: Parsing, Navigation, Filtering, Caching, Optimization
  - Rendering: Main UI engine and consolidation templates
  - Features: Analytics, Git Ops, Workflows, Visualization, Integration, Team/AI

## Architecture Overview

grit is organized into **14 focused modules**:

### Core Infrastructure (6 modules)
Engine types, parsing, navigation, filtering, caching, and optimization

### Rendering & UI (2 modules)
Consolidated rendering templates and main UI rendering engine

### Features (6 modules)
- Analytics (code ownership, hotspots, bisection)
- Git Ops (rebase, cherry-pick, reset, amend)
- Workflows (worktrees, stash, tags, reflog)
- Visualization (flamegraphs, heatmaps, trends)
- Integration (GitHub, Jira, exports)
- Team/AI (metrics, classification, security, compliance)

For detailed architecture information, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Statistics

- **8,200+** lines of Go code
- **840** tests (100% pass rate)
- **150+** type definitions
- **200+** functions
- **14** focused modules
- **Zero** external dependencies (besides Bubble Tea & Lipgloss)

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines, and [DEVELOPER.md](DEVELOPER.md) for:
- Development workflow
- Code style guidelines
- Testing requirements
- How to add new features

## License

MIT
