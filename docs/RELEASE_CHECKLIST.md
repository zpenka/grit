# Release Checklist

Follow this checklist before each release.

## Pre-Release (1 week before)

- [ ] Create release branch: `git checkout -b release/v1.x.x`
- [ ] Update version number in code
- [ ] Run full test suite: `go test ./...`
- [ ] Check test coverage: `go test -cover ./...`
- [ ] Verify benchmarks: `go test -bench=.`

## Code Quality (1 week before)

- [ ] Run linter: `go vet ./...`
- [ ] Format code: `go fmt ./...`
- [ ] Review recent PRs
- [ ] Check for TODOs/FIXMEs in code
- [ ] Update CHANGELOG.md with features
- [ ] Update VERSION.md release notes

## Documentation (3 days before)

- [ ] Update README.md if needed
- [ ] Verify all module docs are current
- [ ] Check HELP.md for new features
- [ ] Update INSTALL.md if instructions changed
- [ ] Review docs/EXAMPLES.md for relevance
- [ ] Test all documentation links

## Testing (3 days before)

- [ ] Run tests on multiple systems: Linux, macOS, Windows
- [ ] Test large repository (5k+ commits)
- [ ] Test small repository
- [ ] Manual feature testing (each major feature)
- [ ] Performance benchmark
- [ ] Memory profile check

## Benchmarks (2 days before)

- [ ] Run benchmarks: `go test -bench=. -benchtime=10s ./...`
- [ ] Compare with previous release
- [ ] Document any performance changes
- [ ] Note any regressions

## Build Artifacts (1 day before)

- [ ] Build binaries for all platforms: `./scripts/build-release.sh v1.x.x`
- [ ] Verify all binaries run: `./grit-linux-amd64 --version`
- [ ] Create checksums: `sha256sum grit-*`
- [ ] Test installation methods:
  - [ ] go install method
  - [ ] Binary download
  - [ ] Homebrew (if applicable)

## Release Day

- [ ] Final test of main branch
- [ ] Create annotated tag: `git tag -a v1.x.x -m "Release v1.x.x"`
- [ ] Push tag: `git push origin v1.x.x`
- [ ] Create GitHub Release with:
  - [ ] Release notes (from CHANGELOG.md)
  - [ ] Binary downloads
  - [ ] Checksums
  - [ ] Installation instructions link

## Post-Release (1 day after)

- [ ] Verify GitHub release is visible
- [ ] Test `go install github.com/zpenka/grit@v1.x.x`
- [ ] Announce on social media (if applicable)
- [ ] Update project README with new version
- [ ] Merge release branch to main
- [ ] Create new development cycle

## Continuous Monitoring (After release)

- [ ] Monitor issues for bugs
- [ ] Watch test failures
- [ ] Track performance reports
- [ ] Collect user feedback
- [ ] Plan next release features

## Common Issues & Solutions

### Tests failing before release
- [ ] Debug locally
- [ ] Fix on release branch
- [ ] Re-run full test suite
- [ ] Do NOT release with failing tests

### Performance regression
- [ ] Identify changed code
- [ ] Profile with pprof
- [ ] Optimize or defer feature
- [ ] Re-benchmark
- [ ] Document in release notes

### Build artifacts won't run
- [ ] Check Go version
- [ ] Verify dependencies
- [ ] Test build script
- [ ] Try building manually
- [ ] Debug on target platform

### Documentation unclear
- [ ] Request review from team
- [ ] Update before release
- [ ] Include in release notes
- [ ] Link to new documentation

## Release Notes Template

```markdown
## Version 1.x.x - YYYY-MM-DD

### New Features
- Feature 1 description
- Feature 2 description

### Improvements
- Performance improvement description
- UI/UX improvement description

### Bug Fixes
- Bug 1 fixed
- Bug 2 fixed

### Documentation
- New guide/documentation added
- Module documentation updated

### Breaking Changes
(if any)
- Breaking change 1
- Migration path: ...

### Contributors
- @username1
- @username2

### Performance
- 15% faster diff loading
- 20% reduced memory usage

### Testing
- 383+ tests passing
- Code coverage: 80%+
```

## Sign-off

- [ ] Release manager sign-off
- [ ] Tech lead review
- [ ] Security review (if applicable)
- [ ] Documentation review

---

**Release completed successfully!** 🎉

Next: Plan features for next release.
