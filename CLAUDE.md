# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**grit** is a terminal UI for exploring git history with 113+ integrated features. It enables interactive browsing of commits, diffs, blame information, analytics, and more—all within the terminal using Go and the Bubble Tea TUI framework.

**Key Stats:**
- 8,200+ lines of Go code
- 505 tests (100% pass rate)
- 150+ type definitions
- 200+ functions
- Zero external dependencies (besides Bubble Tea & Lipgloss for UI)

## Common Development Commands

### Building
```bash
go build -o grit .      # Build binary
go build -o grit ./cmd/grit  # If using cmd package layout
```

### Testing
```bash
go test ./...                    # Run all tests
go test -cover ./...             # With coverage report
go test -run TestFeatureName     # Run single test suite
go test ./core -v                # Verbose output for core package
```

### Running
```bash
./grit                   # Run the built binary from repo root
go run . 2>&1            # Run directly without building
```

### Linting & Quality
```bash
go fmt ./...             # Format all Go files
go vet ./...             # Run Go linter
```

## High-Level Architecture

### Layered Structure

1. **Core Types** (`core/types.go`, engine.go)
   - `commit`: git commit metadata
   - `model`: Central state holder (200+ fields)
   - `diffLine`, `lineKind`: Diff representation
   - `panel`: UI panel management

2. **Core Functions** (engine.go)
   - Parsing: `parseCommits()`, `parseDiff()`, `parseFileItems()`
   - Navigation: `moveCursorUp/Down()`, `switchPanel()`, scroll functions
   - Filtering: `filterCommits()`, `visibleCommits()`

3. **Feature Functions** (engine.go, ~50 functions)
   - **Bisect & Recovery**: bisection workflows
   - **Code Quality**: ownership analysis, hotspots, linting
   - **Analysis**: semantic search, activity heatmaps, merge analysis
   - **Team & Collaboration**: team stats, reviewer suggestions
   - **AI Insights**: auto-complete, classification, anomaly detection
   - **Compliance & Security**: GPG verification, secret scanning
   - **Release & Versioning**: semver detection, changelog generation
   - **Performance**: incremental loading, parallel processing

4. **UI Rendering** (engine.go, engine_render_consolidation.go)
   - Consolidated rendering templates: `RenderStandardUI()`, `RenderAnalysisUI()`, `RenderDataGrid()`
   - Feature-specific renderers: `render*UI()` functions (30+ total)
   - Styled output via Lipgloss

5. **Optimization** (engine_optimization.go)
   - **Caching**: Diff cache (LRU), statistics cache, regex pattern cache
   - **Lazy Loading**: Diffs & features computed on-demand
   - **Incremental Loading**: Progressive repo loading for large histories
   - **Parallel Processing**: Concurrent diff & analysis operations
   - **Memory Management**: Fixed-size caches, object pooling, circular buffers

### Package Structure (Post-Phase-1 Refactoring)
```
grit/
├── grit.go                           # Main entry, Bubble Tea integration
├── engine.go                         # Utilities, keybinding dispatch
│
├─ CORE INFRASTRUCTURE
├── engine_types.go                   # 150+ type definitions
├── engine_parsing.go                 # Git data parsing (commits, diffs, branches)
├── engine_navigation.go              # Cursor & panel navigation
├── engine_filtering.go               # Commit filtering logic
├── engine_cache.go                   # Caching mechanisms
├── engine_optimization.go            # Performance utilities
│
├─ RENDERING & UI
├── engine_render_consolidation.go    # Consolidated render templates
├── engine_rendering.go               # Main UI rendering (2700+ lines)
│
├─ FEATURES (MODULAR)
├── engine_analytics.go               # Analytics, bisect, code ownership
├── engine_git_ops.go                 # Rebase, cherry-pick, reset, amend
├── engine_workflows.go               # Worktrees, stash, tags
├── engine_visualization.go           # Graphs, timelines, heatmaps
├── engine_integration.go             # GitHub, Jira, exports
├── engine_team_ai.go                 # Team analytics, AI, compliance
│
├─ TESTING
├── engine_test_helpers.go            # Test utilities
├── engine_test.go                    # Core test helpers
├── *_test.go                         # 371 feature tests (organized by topic)
├── core/
│   ├── types.go         # Core data structures
│   ├── parser.go        # Commit/diff parsing
│   ├── filter.go        # Filtering operations
│   ├── utils.go         # Helper utilities
│   └── *_test.go        # Package-level tests
└── cmd/grit/            # Optional command package
```

### Data Flow

**Commit Loading:**
```
git log → parseCommits() → model.commits → Render → UI
```

**Diff Processing:**
```
git show → parseDiff() → dcache (LRU) → renderDiffPanel() → UI
```

**Feature Execution:**
```
model.commits → Feature function (e.g., analyzeCodeOwnership) 
→ model.*Data field → render*UI() → UI
```

## State Management

The `model` struct is the central state holder with 200+ fields:
- **Navigation**: `cursor` (position), `focus` (active panel), `diffOffset` (scroll)
- **Filters**: `query`, `authorFilter`, `sinceFilter`, `extensionFilters`
- **Features**: `show*` boolean flags (50+), `*Data` fields for feature results, `*History` for trends
- **Caches**: `dcache`, `scache`, `recache`

All state changes flow through the model during the Bubble Tea update/render cycle.

## Testing Strategy

Tests are organized into 5 categories:
1. **Core Tests** (30+): Parsing, navigation, filtering
2. **Feature Tests** (300+): One per feature implementation
3. **Integration Tests**: Multi-feature interactions
4. **Performance Tests**: Cache hit rates, optimization verification
5. **Regression Tests**: Known bug prevention

Test pattern: Setup fixture → Execute feature → Assert result

## Adding a New Feature

Follow this 6-step pattern:
1. Define a new type in the feature category (e.g., `type hotspotData struct {...}`)
2. Add model fields: `showFeature bool` and `*Data` field
3. Create feature function with TDD tests
4. Create `render*UI()` function following existing patterns
5. Add keybinding handler in `handleKeyBinding()`
6. Update README.md with feature documentation

Example: See ARCHITECTURE.md for a detailed walkthrough of adding "Code Hotspots."

## Performance Considerations

- **Diff Cache**: LRU cache with configurable max size, tracks hit rate
- **Statistics Cache**: Memoizes commit stats to avoid redundant computation
- **Lazy Loading**: Diffs and features computed only when activated
- **Incremental Loading**: Large repos load progressively with UI updates
- **Parallel Processing**: Multiple commits analyzed concurrently

Check cache metrics via model fields to identify optimization opportunities.

## Dependencies

- **Go 1.21+**
- **github.com/charmbracelet/bubbletea** (TUI framework)
- **github.com/charmbracelet/lipgloss** (Styling)
- Standard library: fmt, os, regexp, strings, time, sync, etc.

## Key Patterns & Conventions

- All git operations use `git -C <repoPath>` for working directory safety
- Diff lines use `lineKind` enum for classification (context, added, removed, hunk, meta)
- Commits stored in `model.commits` with optional filtering via `model.visibleCommits()`
- UI updates happen through Bubble Tea's `Update()` method
- Feature data stored in model with `show*` toggle flags
- Cache keys use commit hashes or compiled regex patterns
- Tests use fixtures and helper assertions for clarity

## References

- **[ARCHITECTURE.md](ARCHITECTURE.md)**: Module structure, data flow, and design patterns
- **[DEVELOPER.md](DEVELOPER.md)**: Complete developer guide and contribution workflow
- **[PROGRESS.md](PROGRESS.md)**: Phase-by-phase progress tracking
- **[CHANGELOG.md](CHANGELOG.md)**: Detailed changelog of features and changes
- **[README.md](README.md)**: User-facing documentation and quick start
- **[docs/modules/](docs/modules/)**: Detailed documentation for all 14 modules
