# Documentation Index

Welcome to grit's comprehensive documentation. This directory contains detailed information about each module and system component.

## Quick Navigation

- **[User Guide](../README.md)** - Getting started for end users
- **[Developer Guide](../DEVELOPER.md)** - Development setup and workflow
- **[Contributing Guide](../CONTRIBUTING.md)** - How to contribute
- **[Architecture Overview](../ARCHITECTURE.md)** - System design
- **[All Documentation](../DOCS.md)** - Complete documentation map

---

## Table of Contents

### Core Infrastructure

These modules provide foundational functionality for the entire system.

| Module | Purpose | Key Functions |
|--------|---------|----------------|
| **[PARSING.md](modules/PARSING.md)** | Git data parsing | Parse commits, diffs, file trees |
| **[NAVIGATION.md](modules/NAVIGATION.md)** | UI navigation | Cursor movement, panel switching, bookmarks |
| **[FILTERING.md](modules/FILTERING.md)** | Search & filtering | Text search, regex, author/date filters |
| **[CACHING.md](modules/CACHING.md)** | Performance caching | LRU cache, statistics, regex compilation |
| **[OPTIMIZATION.md](modules/OPTIMIZATION.md)** | Performance utilities | Lazy loading, parallel processing, metrics |

### Rendering & UI

Handles all terminal display and user interface rendering.

| Module | Purpose | Key Functions |
|--------|---------|----------------|
| **[RENDERING.md](modules/RENDERING.md)** | Main UI rendering | Commit list, diff panel, file tree, feature panels |
| **[RENDER_CONSOLIDATION.md](modules/RENDER_CONSOLIDATION.md)** | Unified templates | Standard UI, analysis UI, data grid patterns |

### Feature Modules

Each module implements a logical feature group with multiple integrated features.

| Module | Purpose | Features |
|--------|---------|----------|
| **[ANALYTICS.md](modules/ANALYTICS.md)** | Code analysis | Code ownership, hotspots, bisection, velocity |
| **[GIT_OPS.md](modules/GIT_OPS.md)** | Git workflows | Rebase, cherry-pick, reset, amend |
| **[WORKFLOWS.md](modules/WORKFLOWS.md)** | Advanced workflows | Worktrees, stash, tags, reflog |
| **[VISUALIZATION.md](modules/VISUALIZATION.md)** | Data visualization | Flamegraphs, heatmaps, complexity trends |
| **[INTEGRATION.md](modules/INTEGRATION.md)** | External systems | GitHub, Jira, CSV/JSON/XML export |
| **[TEAM_AI.md](modules/TEAM_AI.md)** | Team & AI features | Team velocity, commit classification, security |

---

## Architecture Overview

### Design Patterns

- **Modular Organization** - 14 focused modules with clear responsibilities
- **Lazy Loading** - Features computed on-demand for performance
- **Caching Strategy** - LRU caches for diffs, statistics, regex patterns
- **Consolidated Rendering** - Unified templates reduce code duplication
- **Centralized State** - Single `model` struct holds all application state

### Data Flow

```
Git Commands → Parsing → Caching → Features → Rendering → UI
```

### Module Dependencies

```
Core Infrastructure (types, parsing, navigation, filtering, cache, optimization)
        ↓
Rendering (consolidation templates + main rendering)
        ↓
Features (analytics, git_ops, workflows, visualization, integration, team_ai)
```

---

## Getting Started with Development

1. **Read first**: [DEVELOPER.md](../DEVELOPER.md)
2. **Understand architecture**: [ARCHITECTURE.md](../ARCHITECTURE.md)
3. **Explore modules**: Start with [PARSING.md](modules/PARSING.md) and work outward
4. **Add features**: Follow [CONTRIBUTING.md](../CONTRIBUTING.md) guidelines

---

## Testing

Each module has comprehensive tests:

```bash
go test ./...                    # Run all tests (371+)
go test -run TestModule ./...    # Run specific module tests
go test -cover ./...             # With coverage
```

Tests are organized in `*_test.go` files alongside implementation.

---

## Performance Considerations

- Large repositories (10k+ commits): See [OPTIMIZATION.md](modules/OPTIMIZATION.md)
- Diff caching strategy: See [CACHING.md](modules/CACHING.md)
- Lazy loading patterns: See [OPTIMIZATION.md](modules/OPTIMIZATION.md)

---

## Code Quality

### Style Guidelines

See [CONTRIBUTING.md](../CONTRIBUTING.md) for:
- Naming conventions
- Comment guidelines
- Testing requirements
- PR process

### Key Patterns

- Use `model` receiver for state mutations
- Cache expensive operations
- Lazy-load features for performance
- Document WHY, not WHAT

---

## Release & Versioning

See [VERSION.md](../VERSION.md) for:
- Semantic versioning rules
- Release process (7 steps)
- Tagging strategy
- Backward compatibility

---

## Project Statistics

- **Code**: 8,200+ lines of Go
- **Tests**: 371+ (100% pass rate)
- **Types**: 150+ definitions
- **Functions**: 200+
- **Modules**: 14 focused modules
- **Documentation**: 3,600+ lines (24 files)

---

## Phase Progress

- ✅ **Phase 1**: Code organization into 14 modules
- ✅ **Phase 2**: Comprehensive documentation + godoc
- ✅ **Phase 3**: Maintenance & contribution guidelines
- ✅ **Phase 4**: Documentation reorganization

See [CHANGELOG.md](../CHANGELOG.md) for detailed phase summaries.

---

## Quick Links

- **User Documentation**: [README.md](../README.md)
- **Architecture**: [ARCHITECTURE.md](../ARCHITECTURE.md)
- **Development**: [DEVELOPER.md](../DEVELOPER.md)
- **Contributing**: [CONTRIBUTING.md](../CONTRIBUTING.md)
- **Releases**: [VERSION.md](../VERSION.md)
- **Changes**: [CHANGELOG.md](../CHANGELOG.md)

---

**Last Updated**: April 27, 2026
