# Architecture Overview

**grit** has been refactored into focused, modular components. See [CLAUDE.md](CLAUDE.md) for the detailed package structure.

## Module Organization (Post-Phase-1)

### Core Infrastructure (Utilities & Foundation)
- **engine_types.go** - All 150+ type definitions
- **engine_parsing.go** - Git data parsing
- **engine_navigation.go** - UI navigation
- **engine_filtering.go** - Commit filtering
- **engine_cache.go** - Caching mechanisms
- **engine_optimization.go** - Performance optimization

### Rendering Layer
- **engine_render_consolidation.go** - Consolidated render templates
- **engine_rendering.go** - Main UI rendering engine

### Feature Modules
- **engine_analytics.go** - Analytics, bisect, code ownership
- **engine_git_ops.go** - Git operations (rebase, cherry-pick, etc.)
- **engine_workflows.go** - Worktrees, stash, tags
- **engine_visualization.go** - Graphs, timelines, heatmaps
- **engine_integration.go** - GitHub, Jira, exports
- **engine_team_ai.go** - Team analytics, AI, compliance

## Data Flow

### Commit Loading
```
git log → parseCommits() → model.commits → Render → UI
```

### Diff Processing
```
git show → parseDiff() → dcache (LRU) → renderDiffPanel() → UI
```

### Feature Execution
```
model.commits → Feature function (e.g., analyzeCodeOwnership)
→ model.*Data field → render*UI() → UI
```

## Key Design Patterns

1. **Modular Organization** - Clear separation of concerns across 14 focused modules
2. **Lazy Loading** - Features computed on-demand
3. **Caching Strategy** - LRU caches for diffs, stats, and regex patterns
4. **Consolidated Rendering** - Template-based UI rendering to reduce duplication
5. **Centralized State** - Model struct holds all application state

## Testing

- **371+ tests** covering all modules
- Tests organized by feature/module
- 100% pass rate maintained

## Performance Considerations

- Diff cache with configurable max size
- Statistics cache to avoid redundant computation
- Lazy loading for large repositories
- Incremental loading with UI updates
- Parallel processing for analysis operations

See [CLAUDE.md](CLAUDE.md) for more details on development and architecture decisions.
