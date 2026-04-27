# Documentation Map

Complete guide to all grit documentation. Start here to find what you're looking for.

## 🚀 Quick Start

**New to grit?**
- [README.md](README.md) - Overview and quick start
- [DEVELOPER.md](DEVELOPER.md) - Get up and running

**Want to contribute?**
- [CONTRIBUTING.md](CONTRIBUTING.md) - How to contribute
- [VERSION.md](VERSION.md) - Release process

**Need architecture details?**
- [ARCHITECTURE.md](ARCHITECTURE.md) - System design
- [docs/README.md](docs/README.md) - Module documentation index

---

## 📖 Documentation by Audience

### For Users

| Document | Purpose |
|----------|---------|
| [README.md](README.md) | Project overview, installation, features |
| [ARCHITECTURE.md](ARCHITECTURE.md) | How grit is organized |

### For Developers

| Document | Purpose |
|----------|---------|
| [DEVELOPER.md](DEVELOPER.md) | Setup, building, testing, common tasks |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Module structure and data flow |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Code style, testing, contribution workflow |
| [docs/README.md](docs/README.md) | Module-by-module reference |
| [docs/modules/](docs/modules/) | Deep dive into each module (13 files) |

### For Maintainers

| Document | Purpose |
|----------|---------|
| [VERSION.md](VERSION.md) | Semantic versioning and release process |
| [CHANGELOG.md](CHANGELOG.md) | What changed in each phase |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Code review guidelines |

### For AI Assistants (Claude, etc.)

| Document | Purpose |
|----------|---------|
| [CLAUDE.md](CLAUDE.md) | Architecture guidance for AI development |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System design and module overview |

---

## 📚 Core Documentation

### [README.md](README.md)
**User-facing project overview**
- Features and capabilities
- Installation instructions
- Quick start guide
- Statistics and architecture summary

### [DEVELOPER.md](DEVELOPER.md)
**Complete development guide**
- Building and running locally
- Testing and coverage
- Common development tasks
- Architecture overview
- Adding new features (6-step guide)
- Code style guidelines
- Contributing workflow

### [CONTRIBUTING.md](CONTRIBUTING.md)
**Contribution guidelines**
- Development setup
- Code style and patterns
- Writing tests
- Development workflow
- PR process and code review
- Documentation requirements
- Performance considerations
- Security best practices

### [ARCHITECTURE.md](ARCHITECTURE.md)
**System architecture overview**
- Module organization (14 modules)
- Data flow diagrams
- Design patterns
- Performance architecture

### [CLAUDE.md](CLAUDE.md)
**Guidance for AI code assistants**
- Project overview
- High-level architecture
- Common development commands
- Testing strategy
- Adding new features
- Performance considerations
- Code patterns and conventions

### [VERSION.md](VERSION.md)
**Release and versioning strategy**
- Semantic versioning (MAJOR.MINOR.PATCH)
- 7-step release process
- Pre-release versions (alpha, beta, rc)
- Branch strategy
- Backward compatibility
- Breaking change process

### [CHANGELOG.md](CHANGELOG.md)
**Detailed changelog by phase**
- Phase 1: Code organization
- Phase 2: Documentation
- Phase 3: Maintenance
- Phase 4: Documentation reorganization
- Statistics and metrics

---

## 🏗️ Module Documentation

All module documentation is in [docs/modules/](docs/modules/) and referenced in [docs/README.md](docs/README.md).

### Core Infrastructure
- [PARSING.md](docs/modules/PARSING.md) - Git data parsing and type conversion
- [NAVIGATION.md](docs/modules/NAVIGATION.md) - UI navigation (cursor, panels, scrolling)
- [FILTERING.md](docs/modules/FILTERING.md) - Search and commit filtering
- [CACHING.md](docs/modules/CACHING.md) - LRU and statistics caching
- [OPTIMIZATION.md](docs/modules/OPTIMIZATION.md) - Performance utilities and lazy loading

### Rendering & UI
- [RENDERING.md](docs/modules/RENDERING.md) - Main UI rendering engine (2,700+ lines)
- [RENDER_CONSOLIDATION.md](docs/modules/RENDER_CONSOLIDATION.md) - Unified rendering templates

### Features
- [ANALYTICS.md](docs/modules/ANALYTICS.md) - Code ownership, hotspots, bisection, velocity
- [GIT_OPS.md](docs/modules/GIT_OPS.md) - Rebase, cherry-pick, reset, amend
- [WORKFLOWS.md](docs/modules/WORKFLOWS.md) - Worktrees, stash, tags, reflog
- [VISUALIZATION.md](docs/modules/VISUALIZATION.md) - Flamegraphs, heatmaps, trends
- [INTEGRATION.md](docs/modules/INTEGRATION.md) - GitHub/Jira, CSV/JSON/XML export
- [TEAM_AI.md](docs/modules/TEAM_AI.md) - Team metrics, AI features, security, compliance

---

## 🔗 GitHub Templates

Located in [.github/](./github/):
- [pull_request_template.md](.github/pull_request_template.md) - PR checklist
- [ISSUE_TEMPLATE/bug_report.md](.github/ISSUE_TEMPLATE/bug_report.md) - Bug report form
- [ISSUE_TEMPLATE/feature_request.md](.github/ISSUE_TEMPLATE/feature_request.md) - Feature request

---

## 📊 Statistics & Metadata

### Documentation Coverage
- **Total files**: 24 markdown files
- **Total lines**: 3,600+ documentation lines
- **Root docs**: 8 files
- **Module docs**: 13 files
- **GitHub templates**: 3 files

### Code Coverage
- **Total code**: 8,200+ lines of Go
- **Tests**: 371+ (100% passing)
- **Type definitions**: 150+
- **Functions**: 200+
- **Modules**: 14 focused modules

---

## 🎯 Finding What You Need

### "I want to..."

**...understand the codebase**
1. Start: [README.md](README.md) - Overview
2. Deep dive: [ARCHITECTURE.md](ARCHITECTURE.md) - Design
3. Learn modules: [docs/README.md](docs/README.md) - Index
4. Specific module: [docs/modules/](docs/modules/) - Details

**...set up development environment**
1. Read: [DEVELOPER.md](DEVELOPER.md) - Complete guide
2. Build: `go build -o grit .`
3. Test: `go test ./...`

**...contribute code**
1. Read: [CONTRIBUTING.md](CONTRIBUTING.md) - Guidelines
2. Create feature: Follow 6-step guide in [DEVELOPER.md](DEVELOPER.md)
3. Submit: Use PR template from [.github/pull_request_template.md](.github/pull_request_template.md)

**...understand a specific module**
1. Find module: [docs/README.md](docs/README.md) - Module list
2. Read docs: [docs/modules/](docs/modules/) - Detailed info
3. Check tests: `*_test.go` files - Usage examples

**...make a release**
1. Read: [VERSION.md](VERSION.md) - Release process
2. Follow: 7-step process outlined
3. Tag: Git tags with semantic versioning

**...see what changed**
1. Recent changes: [CHANGELOG.md](CHANGELOG.md) - By phase
2. Git history: `git log --oneline`
3. Commits: Browse on GitHub

---

## 📝 How Documentation is Organized

```
Root Level (User & Developer Docs)
├── README.md                 # Project overview
├── DEVELOPER.md              # Dev setup & workflow
├── CONTRIBUTING.md           # Contribution guidelines
├── ARCHITECTURE.md           # System design
├── VERSION.md                # Release strategy
├── CHANGELOG.md              # Change history
├── CLAUDE.md                 # AI assistant guidance
└── DOCS.md (this file)       # Documentation map

Module Documentation
├── docs/README.md            # Module index
└── docs/modules/             # 13 detailed module docs
    ├── PARSING.md
    ├── NAVIGATION.md
    ├── FILTERING.md
    ├── CACHING.md
    ├── OPTIMIZATION.md
    ├── RENDERING.md
    ├── RENDER_CONSOLIDATION.md
    ├── ANALYTICS.md
    ├── GIT_OPS.md
    ├── WORKFLOWS.md
    ├── VISUALIZATION.md
    ├── INTEGRATION.md
    └── TEAM_AI.md

GitHub Templates
└── .github/
    ├── pull_request_template.md
    └── ISSUE_TEMPLATE/
        ├── bug_report.md
        └── feature_request.md
```

---

## ✅ Documentation Quality Checklist

All documentation includes:
- ✅ Clear purpose statement
- ✅ Key functions/features listed
- ✅ Design decisions explained
- ✅ Code examples where applicable
- ✅ Performance considerations
- ✅ Testing strategy
- ✅ Integration points
- ✅ Future extensions

---

## 🚀 Phase Progress

| Phase | Completed | Files | Focus |
|-------|-----------|-------|-------|
| Phase 1 | ✅ 2026-04-26 | Code | Code organization into 14 modules |
| Phase 2 | ✅ 2026-04-26 | Docs | Godoc comments + module READMEs |
| Phase 3 | ✅ 2026-04-26 | Contrib | Maintenance, contribution, versioning |
| Phase 4 | ✅ 2026-04-27 | Docs | Documentation reorganization |

See [CHANGELOG.md](CHANGELOG.md) for detailed phase information.

---

## 💡 Tips for Using This Documentation

1. **Use DOCS.md (this file) as your entry point** - It's organized by audience and task
2. **README.md is for users** - Share it with people new to grit
3. **DEVELOPER.md is your daily reference** - Bookmark it
4. **Module docs are detailed references** - Read when working on that module
5. **CONTRIBUTING.md is required reading** - Before submitting PRs
6. **VERSION.md is for release** - Only needed when making releases

---

**Questions?** Check [CONTRIBUTING.md](CONTRIBUTING.md) "Getting Help" section.

**Last Updated**: April 27, 2026
