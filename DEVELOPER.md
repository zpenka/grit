# Developer Guide

This guide explains how to develop, extend, and contribute to grit.

## Getting Started

### Prerequisites
- Go 1.21+
- A terminal with 24-bit color support (recommended)
- Git

### Build & Run
```bash
go build -o grit .
./grit
```

### Testing
```bash
go test ./...              # Run all tests
go test -cover ./...       # With coverage
go test -run TestName      # Specific test
go test -v ./...           # Verbose output
```

All 908 tests should pass. Tests run automatically on PRs via GitHub Actions.

## Architecture Overview

grit is organized into **14 focused modules**:

### Core Infrastructure
These provide foundational functionality:
- **engine_types.go** - All type definitions (150+)
- **engine_parsing.go** - Parse git data (commits, diffs, branches)
- **engine_navigation.go** - Handle UI navigation (cursor, panels, scrolling)
- **engine_filtering.go** - Filter and search commits
- **engine_cache.go** - Cache mechanisms (LRU, memoization)
- **engine_optimization.go** - Performance utilities

### Rendering
- **engine_render_consolidation.go** - Consolidated render templates
- **engine_rendering.go** - Main UI rendering engine (2700+ lines)

### Features (Modular)
Each module handles a logical feature group:
- **engine_analytics.go** - Analytics, bisect, code ownership
- **engine_git_ops.go** - Git operations (rebase, cherry-pick, reset, amend)
- **engine_workflows.go** - Worktrees, stash, tag management
- **engine_visualization.go** - Graphs, timelines, heatmaps
- **engine_integration.go** - GitHub, Jira, exports
- **engine_team_ai.go** - Team analytics, AI, compliance, versioning

## Adding a New Feature

### 1. Identify the Right Module
Choose the module that best fits your feature:
- Analytical calculation? → `engine_analytics.go`
- Git operation? → `engine_git_ops.go`
- Data visualization? → `engine_visualization.go`
- External integration? → `engine_integration.go`

### 2. Add the Feature Function
```go
// Example: Add to engine_analytics.go
// analyzeNewFeature calculates X from commits and returns Y.
// This function demonstrates the analysis pattern.
func analyzeNewFeature(commits []commit) []featureData {
    // Implementation
    return results
}
```

### 3. Add the Render Function
```go
// renderNewFeatureUI displays feature results using standard UI template.
func renderNewFeatureUI(m model, width int) string {
    data := make(map[string]interface{})
    // Populate data from m.featureData
    return RenderAnalysisUI("Feature Name", data)
}
```

### 4. Wire into Model & Keybindings
Update `engine_types.go`:
```go
type model struct {
    // ... existing fields ...
    featureData []featureData
    showFeature bool
}
```

Update `grit.go` keybinding handler to trigger feature.

### 5. Write Tests
Create `feature_test.go`:
```go
func TestAnalyzeNewFeature(t *testing.T) {
    commits := makeCommits(5)
    result := analyzeNewFeature(commits)
    
    if len(result) == 0 {
        t.Error("expected results, got empty")
    }
}
```

### 6. Update Documentation
- Add godoc comments to exported functions
- Update relevant module README
- Update CLAUDE.md if architectural impact

## Code Style

### Naming
- Functions: camelCase (unexported) or PascalCase (exported)
- Types: PascalCase
- Constants: UPPER_SNAKE_CASE
- Variables: camelCase

### Patterns
- Use `model` receiver for UI state mutations
- Use `[]commit` for immutable commit lists
- Cache expensive operations with `dcache`, `scache`, `recache`
- Lazy-load complex features

### Comments
- Add godoc comments to all exported functions
- Document WHY, not WHAT (code shows what it does)
- Include examples for complex functions

## Performance Considerations

### Caching
- `dcache`: LRU cache for diffs (large objects)
- `scache`: Statistics cache (frequently computed)
- `recache`: Regex pattern cache (compiled patterns)

### Lazy Loading
- Compute features only when user requests them
- Use `show*` flags to gate feature computation
- Leverage goroutines for parallel analysis (see `engine_optimization.go`)

### Large Repositories
- Use incremental loading for 10k+ commits
- Implement pagination for result display
- Profile with `pprof` before optimizing

## Testing Strategy

### Test Organization
Tests are in separate `*_test.go` files organized by module:
- `engine_analytics_test.go` - Analytics tests
- `navigation_test.go` - Navigation tests
- `filtering_test.go` - Filtering tests
- etc.

### Test Helpers
Located in `engine_test_helpers.go`:
- `makeCommits(n)` - Create test commits
- `makeDiffLines(n)` - Create test diff lines
- `makeNamedCommits()` - Create commits with specific authors

### Writing Tests
```go
func TestFeatureName(t *testing.T) {
    // Setup
    commits := makeCommits(10)
    
    // Execute
    result := featureFunction(commits)
    
    // Assert
    if result == nil {
        t.Fatal("expected result, got nil")
    }
    if len(result) != expected {
        t.Errorf("got %d, want %d", len(result), expected)
    }
}
```

## CI/CD

### GitHub Actions
Tests run automatically on:
- All pull requests to `main`
- All pushes to `main`

Workflow: `.github/workflows/tests.yml`

### Making Tests Required
To block merges on test failures:
1. Go to Settings → Branches
2. Add rule for `main` branch
3. Enable "Require status checks to pass before merging"
4. Select "test" as required check

## Common Tasks

### Run Specific Test
```bash
go test -run TestNavigationCursor ./...
```

### Check Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Format Code
```bash
go fmt ./...
```

### Run Linter
```bash
go vet ./...
```

### Debug with Print
```go
fmt.Printf("Debug: value=%v\n", value)
```

## Architecture Decisions

### Why Modular?
- Easier to navigate large codebase
- Clear separation of concerns
- Faster to find related code
- Simpler to test individual features

### Why Consolidated Rendering?
- Reduces duplication (50+ similar render functions)
- Consistent styling across features
- Easier to change UI globally

### Why LRU Caching?
- Frequently accessed diffs stay hot
- Bounded memory usage (prevents OOM)
- Simple to reason about behavior

## References

- [CLAUDE.md](CLAUDE.md) - Project overview and architecture
- [README.md](README.md) - User-facing documentation
- [ARCHITECTURE.md](ARCHITECTURE.md) - Architectural overview
- Go standard library: https://golang.org/pkg/
- Bubble Tea: https://github.com/charmbracelet/bubbletea
- Lipgloss: https://github.com/charmbracelet/lipgloss

## Getting Help

1. Check existing tests for usage examples
2. Look at similar features for patterns
3. Read module-specific READMEs
4. Ask in pull request discussions

## Contributing

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Write tests first (TDD)
3. Implement the feature
4. Ensure all tests pass: `go test ./...`
5. Create a PR with clear description
6. Address review feedback

Thank you for contributing! 🚀
