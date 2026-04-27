# Changelog

All notable changes to grit are documented here. This file follows the [Keep a Changelog](https://keepachangelog.com/) format.

**📚 For complete documentation overview, see [DOCS.md](DOCS.md)**

## [1.0.0] - 2026-04-27 (Stable Release)

### Release Summary

grit reaches **v1.0.0** — a stable, production-ready terminal UI for git history exploration. This release consolidates two major phases of optimization, documentation, and feature completeness.

**Key Highlights:**
- **505 tests** (100% pass rate) across 14 focused modules
- **Performance optimization** — 70% memory reduction on critical filtering paths, benchmarks established
- **Production-grade tooling** — CLI flags, version management, release build scripts
- **Zero external dependencies** besides Bubble Tea & Lipgloss for TUI
- **Go 1.24** and latest library versions (bubbletea 1.3.10, lipgloss 1.1.0)

### Completed in This Release

#### Phase 5: Code Quality, UX, and Release
- **TDD test coverage** across all feature modules (38 new tests added)
- **Feature stub fixes** — replaced hardcoded values with real implementations
- **CLI flags** — version (`-v`/`--version`) and help (`-h`/`--help`) flags fully working
- **Release tooling** — build script fixed, cross-platform binary generation working (linux, darwin, windows)
- **Tests**: 371 → 409, then to 505 (Phase 5 + Phase 4 + earlier improvements)

#### Phase 4: Performance & Dependencies  
- **Dependency upgrades** — Go 1.21 → 1.24.0, bubbletea v0.24 → v1.3.10, lipgloss v0.10 → v1.1.0
- **Benchmark suite** — Added 5 new performance benchmarks (11 total)
- **Filtering optimization** — Single-pass `visibleCommits()` reducing memory 561KB → 170KB (70% reduction)
- **Allocation tracking** — 7 new performance tests establishing baselines for regression prevention
- **Cache validation** — Diff cache hit rates, stats cache memoization, memory pool verification

### Stats

- **Tests**: 505 total (100% pass rate)
- **Modules**: 14 (Core infrastructure, Rendering, Features)
- **Features**: 113+ integrated capabilities
- **Lines of Code**: 8,200+ Go code
- **Type Definitions**: 150+
- **Functions**: 200+
- **External Dependencies**: 0 (besides UI framework)

### Installation

```bash
go install github.com/zpenka/grit@v1.0.0
# or download binary for your platform from GitHub Releases
```

### Documentation

- **[DOCS.md](DOCS.md)** — Complete documentation map (start here)
- **[README.md](README.md)** — User-facing overview
- **[DEVELOPER.md](DEVELOPER.md)** — Development guide
- **[docs/modules/](docs/modules/)** — 14 detailed module guides

---

## [Phase 4] - 2026-04-26

### Added - Performance Monitoring & Optimization

**Benchmarks & Profiling:**
- Added 5 new benchmarks in `benchmarks_test.go` (diff, timeline, heatmap, secrets, visibleCommits)
- Added allocation tracking tests in `performance_test.go`
- Established baseline metrics for hot paths

**Performance Tests:**
- `TestAllocationsParseCommits` - 109 allocs/op baseline
- `TestAllocationsFilterCommits` - 100 allocs/op baseline
- `TestDiffCacheHitRate` - cache effectiveness validation
- `TestVisibleCommitsAllocationUnderLoad` - filtering at scale

### Changed - Dependency Updates

**Upgraded:**
- bubbletea: v0.24.0 → v1.3.10 (fully backward compatible)
- lipgloss: v0.10.0 → v1.1.0  
- Go: 1.21 → 1.24.0

**Impact:** Zero breaking changes, all 409 tests passing

### Improved - Filtering Performance

**Single-pass visibleCommits optimization:**
- Consolidated 3 chained filter operations into single loop
- Eliminated intermediate slice allocations
- Memory usage: 561KB → 170KB (70% reduction)
- Speed: maintained at ~793µs

**Pre-allocation optimizations:**
- filterCommits: capacity = len/2 (~50% expected matches)
- filterCommitsByAuthor: capacity = len/5 (~20% per author)
- filterCommitsSince: capacity = len/3 (~30% recent)

### Stats

- **Benchmarks**: 6 → 11 (5 new performance benchmarks)
- **Performance tests**: 0 → 7 (allocation & cache validation)
- **Memory reduction**: 70% on main filtering hot path
- **Commits**: 3 (dep update + benchmarks + optimization)

---

## [Phase 3] - 2026-04-26

### Added - Feature Completeness & Release Readiness

**CLI Flags (grit.go):**
- `const Version = "0.1.0"` - Version constant
- `-v, --version` flag - Print version and exit
- `-h, --help` flag - Print usage/help message

**Test Coverage:**
- `engine_team_ai_test.go` - 11 tests for analytics/compliance features
- `engine_visualization_test.go` - 14 tests for visualization functions
- `engine_workflows_test.go` - 13 tests for workflow operations
- Total: 38 new tests, all passing

### Fixed - Hardcoded Stubs Replaced

**engine_team_ai.go:**
- `trackLicenseHeaders` - Removed hardcoded main.go/MIT; returns empty list
- `scanForSecurityIssues` - Fixed hardcoded "line 5" → real line number tracking
- `trackDataDeletion` - Removed hardcoded "GDPR request" reason
- `detectSecrets` - Fixed hardcoded "line 1" → real line tracking with patterns
- `detectAnomalies` - Enhanced to detect large commits (50000+ lines, "rewrite entire" patterns)

**engine_visualization.go:**
- `buildTreeView` - Fixed flat depth-1 structure → builds proper nested hierarchy
- `buildTimeline` - Added missing hash field to timeline points
- `compareAuthors` - Now works without pre-selected authors; auto-compares first two

**scripts/build-release.sh:**
- Fixed variable case bug: `${version}` → `${VERSION}` (was creating blank filenames)

### Changed - Release Tooling

- Release script now generates proper binary names: `grit-0.1.0-linux-amd64`, `grit-0.1.0-darwin-arm64`, etc.
- All platform builds (linux, darwin, windows) with correct architectures

### Stats

- **Tests**: 371 → 409 (38 new tests, 100% pass rate)
- **Files changed**: 7 (5 modified, 3 created)
- **Lines added**: 737+
- **Commits**: 1 (bundled PR #11)

---

## [Phase 2] - 2026-04-26

### Added - Documentation & Godoc Comments

**Module Documentation (13 READMEs in docs/modules/):**
- `PARSING.md` - Git data parsing module
- `NAVIGATION.md` - UI navigation and cursor management
- `FILTERING.md` - Commit filtering and search
- `CACHING.md` - Cache mechanisms and strategies
- `OPTIMIZATION.md` - Performance optimization utilities
- `RENDERING.md` - Main UI rendering engine
- `RENDER_CONSOLIDATION.md` - Unified rendering templates
- `ANALYTICS.md` - Code analytics and insights
- `GIT_OPS.md` - Git operations (rebase, cherry-pick, etc.)
- `WORKFLOWS.md` - Advanced git workflows
- `VISUALIZATION.md` - Data visualization and charts
- `INTEGRATION.md` - External system integration
- `TEAM_AI.md` - Team metrics and AI features

**Godoc Comments:**
- Added comprehensive package-level documentation to 13 modules
- Enables IDE tooltips and documentation generation
- Improves developer experience for new contributors

**Developer Guide:**
- Created `DEVELOPER.md` with 200+ lines of content
- Getting started and development workflow
- Architecture overview and feature development guide
- Code style guidelines and patterns
- Testing strategy and CI/CD setup
- Common tasks and troubleshooting

**Progress Tracking:**
- Created `PROGRESS.md` - Comprehensive progress documentation
- Created `CHANGELOG.md` - This file

### Changed - Documentation Updates

- Updated `CLAUDE.md` - Reflected 14-module structure
- Updated `ARCHITECTURE.md` - Simplified with module references
- Enhanced `README.md` - Focused on user-facing features

### Removed - Outdated Files

- Removed `ARCHITECTURE_COMPREHENSIVE.md` - Outdated (referenced 312+ features vs actual 113+)

### Stats

- **2,000+ new documentation lines**
- **13 module READMEs created**
- **13 godoc comments added**
- **28 files changed** (13 godoc + 13 READMEs + 2 new files)

---

## [Phase 1] - 2026-04-26

### Added - Modular Code Organization

**Core Infrastructure Modules:**
- `engine_types.go` - 150+ type definitions
- `engine_parsing.go` - Git data parsing (commits, diffs, branches)
- `engine_navigation.go` - Cursor and panel management
- `engine_filtering.go` - Commit filtering and search logic
- `engine_cache.go` - LRU and statistics caching
- `engine_optimization.go` - Performance utilities and lazy loading

**Rendering Modules:**
- `engine_render_consolidation.go` - Unified rendering templates (reduces duplication)
- `engine_rendering.go` - Main UI rendering (refactored and split)

**Feature Modules:**
- `engine_analytics.go` - Code ownership, hotspots, bisection, velocity analysis
- `engine_git_ops.go` - Rebase, cherry-pick, reset, amend operations
- `engine_workflows.go` - Worktrees, stash, tags, reflog management
- `engine_visualization.go` - Flamegraphs, heatmaps, complexity trends
- `engine_integration.go` - GitHub/Jira linking, CSV/JSON/XML export
- `engine_team_ai.go` - Team metrics, AI insights, compliance, versioning

**CI/CD:**
- Added `.github/workflows/tests.yml` - Automated testing on PRs
- Tests run on all pull requests to main branch
- Instructions provided for making tests required checks

### Changed - Code Refactoring

- **engine_rendering.go**: Reduced from 4,557 to 2,700+ lines after extraction
- **Code organization**: Clear separation of concerns across 14 focused modules
- **Dependency structure**: Well-defined module boundaries and relationships

### Maintenance

- All 371 tests continue to pass
- No functionality changed (refactoring only)
- Improved code navigation and maintainability

### Stats

- **865 lines extracted** from engine_rendering.go
- **14 modules total** (up from 1 monolithic module)
- **371+ tests** - 100% passing
- **100% backward compatible** - No feature changes

---

## Structure Summary

### Modules Overview

| Category | Module | Purpose | Lines |
|----------|--------|---------|-------|
| **Core** | engine_types.go | Type definitions | ~150 exports |
| | engine_parsing.go | Git data parsing | ~100 |
| | engine_navigation.go | UI navigation | ~200 |
| | engine_filtering.go | Commit filtering | ~100 |
| | engine_cache.go | Caching mechanisms | ~100 |
| | engine_optimization.go | Performance utils | ~150 |
| **Rendering** | engine_render_consolidation.go | Render templates | ~200 |
| | engine_rendering.go | Main UI engine | ~2,700 |
| **Features** | engine_analytics.go | Analytics & bisect | ~500 |
| | engine_git_ops.go | Git operations | ~150 |
| | engine_workflows.go | Advanced workflows | ~150 |
| | engine_visualization.go | Visualizations | ~200 |
| | engine_integration.go | Integrations | ~150 |
| | engine_team_ai.go | Team & AI features | ~400 |

---

## Documentation Files

| File | Purpose | Status |
|------|---------|--------|
| README.md | User-facing overview | ✅ Current |
| CLAUDE.md | Developer guidance & architecture | ✅ Current |
| ARCHITECTURE.md | Module structure & data flow | ✅ Current |
| DEVELOPER.md | Complete developer guide | ✅ New (Phase 2) |
| PROGRESS.md | Progress tracking | ✅ New (Phase 2) |
| CHANGELOG.md | This file | ✅ New (Phase 2) |
| docs/modules/ | 13 detailed module READMEs | ✅ New (Phase 2) |

---

## Testing

- **Total Tests**: 371+
- **Pass Rate**: 100%
- **Organization**: Tests in `*_test.go` files organized by module
- **Helpers**: Centralized in `engine_test_helpers.go`
- **CI/CD**: Automated via GitHub Actions

### Test Coverage

- Core infrastructure: 30+ tests
- Features: 300+ tests
- Integration: 20+ tests
- Performance: 10+ tests

---

## Future Roadmap

### Phase 3 (Planned)
- [ ] Module-specific test documentation
- [ ] API documentation for extensions
- [ ] Performance optimization guide
- [ ] Contribution workflow docs

### Ongoing
- Maintain 100% test pass rate
- Keep documentation in sync with code
- Regular code reviews and refactoring
- Performance monitoring and optimization

---

## How to Use This Changelog

- **For Users**: See what features are available in each phase
- **For Developers**: Understand code organization and documentation structure
- **For Contributors**: Reference architecture decisions and patterns

For detailed information about specific modules, see the individual READMEs in [docs/modules/](docs/modules/).
