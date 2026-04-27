# Contributing to grit

Thank you for your interest in contributing to grit! This document provides guidelines for contributing code, documentation, and issues.

## Getting Started

### Prerequisites
- Go 1.21+
- Git
- A terminal with 24-bit color support (recommended)

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/grit.git`
3. Add upstream: `git remote add upstream https://github.com/zpenka/grit.git`
4. Create a feature branch: `git checkout -b feature/your-feature-name`
5. Build and test: `go build -o grit .` and `go test ./...`

## Code Style Guidelines

### Naming Conventions
- Functions: camelCase (unexported) or PascalCase (exported)
- Types: PascalCase
- Constants: UPPER_SNAKE_CASE
- Variables: camelCase

### Patterns
- Use `model` receiver for UI state mutations
- Keep functions under 100 lines when possible
- Cache expensive operations with `dcache`, `scache`, or `recache`
- Lazy-load features for performance

### Comments
- Add godoc comments to all exported functions
- Document WHY, not WHAT (code shows what it does)
- Keep comments concise (one line preferred)

## Writing Tests

### Test Organization
- Create `*_test.go` files alongside code files
- Use table-driven tests for multiple cases
- Follow the pattern: Setup → Execute → Assert

### Test Example
```go
func TestAnalyzeCodeOwnership(t *testing.T) {
    commits := makeCommits(10)
    result := analyzeCodeOwnership(commits)
    
    if len(result) == 0 {
        t.Error("expected results, got empty")
    }
}
```

### Testing Best Practices
- Write tests first (TDD preferred)
- Aim for high coverage but not 100%
- Use helper functions from `engine_test_helpers.go`
- Test edge cases (empty lists, nil, etc.)

## Development Workflow

### Branch Naming
- Feature: `feature/descriptive-name`
- Bugfix: `fix/issue-description`
- Documentation: `docs/doc-topic`
- Refactor: `refactor/component-name`

### Branch Cleanup
After your PR is merged:
1. Delete your feature branch locally: `git branch -d feature-name`
2. Delete from remote: `git push origin --delete feature-name`
3. Keep main updated: `git checkout main && git pull upstream main`

### Commit Messages
- Use imperative mood ("Add feature" not "Added feature")
- First line: 50 characters or less
- Reference issues: "Fixes #123" or "Related to #456"
- Separate subject from body with blank line

## Pull Request Process

1. **Before submitting:**
   - Run all tests: `go test ./...`
   - Format code: `go fmt ./...`
   - Run linter: `go vet ./...`

2. **When creating PR:**
   - Use the PR template (auto-filled)
   - Reference related issues
   - Describe what and why, not just what
   - Add links to relevant documentation

3. **During review:**
   - Respond to feedback promptly
   - Push additional commits (don't force-push)
   - Ask questions if guidance is unclear
   - Be respectful and collaborative

4. **After merge:**
   - Delete your feature branch
   - Pull latest main
   - Close related issues

## Code Review

All code changes require review before merging. Reviewers look for:

- **Correctness**: Does the code do what it claims?
- **Style**: Does it follow the guidelines?
- **Tests**: Are there tests? Do they pass?
- **Performance**: Could this impact performance negatively?
- **Documentation**: Are new features documented?

### Being a Good Code Reviewer
- Be constructive and kind
- Suggest improvements, don't demand changes
- Approve when satisfied
- Highlight excellent code and practices

## Documentation

### When to Document
- New features must include documentation
- Public functions need godoc comments
- Complex logic needs inline comments explaining WHY
- Architecture changes need ARCHITECTURE.md updates

### Documentation Files
- **DEVELOPER.md**: Development guide
- **CONTRIBUTING.md**: This file
- **VERSION.md**: Versioning strategy
- **docs/modules/**: Module-specific READMEs
- **CHANGELOG.md**: Release notes

## Reporting Issues

### Bug Reports
Use the bug report template:
- Describe the bug clearly
- Provide steps to reproduce
- Show expected vs actual behavior
- Include environment (OS, Go version, terminal)

### Feature Requests
Use the feature request template:
- Describe the desired behavior
- Explain why it's useful
- Provide examples if possible
- Suggest implementation if you have ideas

## Performance Considerations

When adding features:
- Consider impact on startup time
- Use caching for expensive operations
- Lazy-load features that aren't always needed
- Profile with `pprof` if unsure

## Security

- Never commit secrets or credentials
- Report security issues privately (don't create public issues)
- Be cautious with user input
- Validate at system boundaries

## Getting Help

- Read [DEVELOPER.md](DEVELOPER.md) for development info
- Check [CLAUDE.md](CLAUDE.md) for architecture
- Look at existing code for patterns
- Ask in discussions or issues

## Code of Conduct

Be respectful, inclusive, and professional. We welcome contributors of all backgrounds and experience levels.

## Recognition

Contributors are recognized in:
- Commit history
- CHANGELOG.md
- Release notes
- Project README (for significant contributions)

---

**Thank you for contributing to grit!** Your work helps make git history exploration better for everyone.
