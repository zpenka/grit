# Code Quality Standards

grit maintains high code quality standards across all metrics.

## Testing Requirements

### Coverage Standards
- **Target**: 80%+ code coverage
- **Critical paths**: 100% coverage
- **Acceptable minimum**: 70%

### Test Organization
- One `*_test.go` file per module
- Table-driven tests for multiple cases
- Clear setup-execute-assert pattern

### Test Execution
```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Benchmark Standards
- Benchmarks for performance-critical paths
- Baselines established for each release
- Regressions tracked and investigated

```bash
# Run benchmarks
go test -bench=. -benchtime=10s ./...

# Memory benchmarks
go test -bench=. -benchmem ./...
```

## Code Quality Metrics

### Cyclomatic Complexity
- **Target**: < 10 per function
- **Maximum**: 15 (refactor needed)
- **Tool**: `gocyclo`

### Lines per Function
- **Target**: < 50 lines
- **Maximum**: 100 lines (refactor needed)
- **Exceptions**: Complex algorithms documented

### Code Duplication
- **Target**: < 3% duplication
- **Tool**: `go-fuzz`, code review

## Performance Benchmarks

### Parse Operations (Commits, Diffs)
- Parse 1000 commits: < 100ms
- Parse diff: < 50ms (first), instant (cached)
- Filter 1000 commits: < 10ms

### UI Rendering
- Initial render: < 500ms
- Refresh render: < 50ms
- Scroll: < 20ms

### Memory Usage
- Baseline: < 50MB empty
- With 1000 commits: < 200MB
- Peak usage: < 500MB

## Standards Across Metrics

| Metric | Target | Minimum | Tool |
|--------|--------|---------|------|
| **Coverage** | 85% | 70% | `go test -cover` |
| **Benchmarks** | Tracked | < 5% regression | `go test -bench` |
| **Complexity** | < 10 | < 15 | Code review |
| **Duplication** | < 3% | < 5% | Code review |
| **Format** | 100% | 100% | `go fmt` |
| **Lint** | 0 issues | 0 issues | `go vet` |

## Quality Gates

### Before Commit
- [ ] Tests pass: `go test ./...`
- [ ] Coverage acceptable: `go test -cover ./...`
- [ ] Code formatted: `go fmt ./...`
- [ ] No lint issues: `go vet ./...`

### Before PR
- [ ] All tests pass on CI
- [ ] Coverage >= 80%
- [ ] No new linting issues
- [ ] Code review approved

### Before Release
- [ ] All quality metrics met
- [ ] Benchmarks within 5% of previous
- [ ] No performance regressions
- [ ] Documentation updated

## Continuous Integration

### GitHub Actions Workflow
```yaml
- Run tests: go test -v ./...
- Check coverage: go test -cover ./...
- Verify formatting: go fmt ./...
- Run linter: go vet ./...
- Compare benchmarks: baseline check
```

### Coverage Reporting
- Codecov integration enabled
- Coverage > 80% required for merge
- PR comments show coverage changes

## Performance Profiling

### CPU Profiling
```bash
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### Memory Profiling
```bash
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

### Use Cases
- Identify hotspots
- Track memory leaks
- Verify optimizations

## Code Review Standards

### Every PR Must Have
- [ ] Clear description of changes
- [ ] Tests included
- [ ] Coverage maintained/improved
- [ ] No performance regressions
- [ ] Documentation updated

### Code Review Checklist
- [ ] Logic is correct
- [ ] Edge cases handled
- [ ] Error handling present
- [ ] Performance acceptable
- [ ] Comments explain why
- [ ] Tests are comprehensive

## Refactoring Triggers

Refactor when:
- Cyclomatic complexity > 15
- Function > 100 lines
- Code duplication > 5%
- Performance regression > 5%
- Test coverage < 70%
- Code review feedback suggests it

## Tools & Commands

### Local Development
```bash
# Format code
go fmt ./...

# Lint
go vet ./...

# Test with coverage
go test -cover ./...

# Benchmark
go test -bench=. -benchtime=10s ./...

# Profile CPU usage
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof

# View code coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### CI/CD Integration
All tools run automatically on:
- Pull requests (required)
- Commits to main (required)
- Release builds (required)

## Metrics Dashboard

Current metrics (updated after each release):

| Metric | Current | Target | Status |
|--------|---------|--------|--------|
| **Code Coverage** | 80%+ | 85% | ✅ Good |
| **Test Count** | 383+ | - | ✅ Comprehensive |
| **Performance** | Baseline | < 5% regression | ✅ Stable |
| **Duplication** | < 3% | < 3% | ✅ Low |
| **Complexity** | < 10 avg | < 10 | ✅ Healthy |

## Improvement Plan

### Next Quarter
- [ ] Increase coverage to 90%
- [ ] Add performance benchmarks to CI
- [ ] Implement codecov integration
- [ ] Document optimization opportunities

### Next Year
- [ ] > 95% code coverage
- [ ] Performance benchmarks tracked
- [ ] Automated refactoring suggestions
- [ ] Zero critical issues

## References

- [CONTRIBUTING.md](CONTRIBUTING.md) - Code style guidelines
- [DEVELOPER.md](DEVELOPER.md) - Development guide
- [docs/RELEASE_CHECKLIST.md](docs/RELEASE_CHECKLIST.md) - Release process
- [benchmarks_test.go](benchmarks_test.go) - Performance tests

---

**Quality is everyone's responsibility!** 

Follow these standards to keep grit excellent.
