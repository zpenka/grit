# Versioning Strategy

grit follows [Semantic Versioning](https://semver.org/) for releases and version management.

## Version Format

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
```

### Examples
- `1.0.0` - Initial release
- `1.2.0` - Minor release (new features)
- `1.2.1` - Patch release (bug fix)
- `2.0.0` - Major release (breaking changes)
- `1.0.0-alpha` - Alpha prerelease
- `1.0.0-rc.1` - Release candidate

## Semantic Versioning Rules

### MAJOR version
- Increment when making incompatible API changes
- May break existing workflows or configurations
- Requires migration guide for users
- Reset MINOR and PATCH to 0

### MINOR version
- Increment when adding functionality in a backward-compatible manner
- New features that don't break existing code
- New deprecations announced
- Reset PATCH to 0

### PATCH version
- Increment for backward-compatible bug fixes
- Performance improvements
- Documentation updates
- No functional changes to API

## Release Process

### Step 1: Prepare Release Branch
```bash
git checkout main
git pull origin main
git checkout -b release/v1.2.0
```

### Step 2: Update Version References
- Update `VERSION` constant in code (if applicable)
- Update version in documentation
- Update CHANGELOG.md with release notes

### Step 3: Create Release Commit
```bash
git commit -m "Release v1.2.0

- List major features
- List bug fixes
- List documentation updates
- Mention contributors"
```

### Step 4: Tag Release
```bash
git tag -a v1.2.0 -m "Release version 1.2.0"
git push origin release/v1.2.0
git push origin v1.2.0
```

### Step 5: Create Release PR
- Create PR from release/v1.2.0 to main
- Reference any related issues
- Link to CHANGELOG.md entry
- Request review from maintainers

### Step 6: Merge and Create GitHub Release
- Merge PR to main
- Go to GitHub Releases
- Create release from tag v1.2.0
- Copy CHANGELOG.md section as release notes
- Publish release

### Step 7: Cleanup
```bash
git checkout main
git pull origin main
git branch -d release/v1.2.0
```

## Development Versions

### Pre-release Versions
Used for testing before official release:
- **Alpha**: `v1.0.0-alpha` - Early development, may be unstable
- **Beta**: `v1.0.0-beta` - Feature-complete, testing phase
- **Release Candidate**: `v1.0.0-rc.1` - Final testing before release

### Build Metadata
Used for internal builds:
```
v1.0.0+build.20260427
v1.0.0+ci.github.abc123
```

## Branch Strategy

### Long-lived Branches
- **main**: Production-ready code, current release
- **develop**: Next release development (optional)

### Temporary Branches
- **feature/**: New features
- **fix/**: Bug fixes
- **release/**: Release preparation
- **hotfix/**: Critical production fixes

## Deployment Versioning

### Release Schedule
- Patch releases: As needed (usually monthly)
- Minor releases: Quarterly
- Major releases: Annually or as needed

### Version Planning
1. Plan features for next release
2. Estimate timeline
3. Announce version in issues/discussions
4. Create tracking issue for release
5. Update CHANGELOG.md as features complete
6. Release when ready

## Backward Compatibility

### Compatibility Guarantees
- Within MAJOR version: API/CLI compatibility maintained
- MINOR versions: New features won't break existing code
- PATCH versions: No new features, only fixes

### Deprecation Process
1. Announce deprecation in release notes
2. Add deprecation warnings in code
3. Maintain deprecated feature for 2+ releases
4. Remove in next MAJOR version
5. Document migration path

## Breaking Changes

When introducing breaking changes:
1. Only in MAJOR version releases
2. Announce early (multiple releases before)
3. Provide clear migration guide
4. Update all documentation
5. Consider compatibility layer if possible

## Version Numbering Examples

| Change | Version | Type |
|--------|---------|------|
| Bug fix | 1.0.1 | Patch |
| New feature | 1.1.0 | Minor |
| Breaking change | 2.0.0 | Major |
| Critical fix | 1.0.2 | Patch |
| Multiple features | 1.5.0 | Minor |
| Performance improvement | 1.0.2 | Patch |

## Git Tags

### Tag Format
- Format: `v1.2.3` (prefix with 'v')
- Annotated tags only: `git tag -a v1.2.3 -m "Release v1.2.3"`
- Include changelog in tag message

### Tag Cleanup
- Never delete released tags
- Use new tag if error found
- Document corrections in CHANGELOG.md

## Communicating Versions

### Announcements
1. GitHub Releases page
2. Changelog/Release notes
3. Project README (if major)
4. Social media for major releases

### Release Notes Contents
```
## v1.2.0 - 2026-04-27

### New Features
- Feature 1 description
- Feature 2 description

### Bug Fixes
- Fixed issue with X
- Fixed issue with Y

### Documentation
- Updated DEVELOPER.md
- Added new module documentation

### Contributors
- @username1
- @username2
```

## Future Considerations

As grit matures:
- Consider long-term support (LTS) versions
- Document upgrade paths for major versions
- Maintain multiple release branches if needed
- Consider semantic commit conventions

---

**Release Process Owner**: Project maintainers
**Last Updated**: 2026-04-27
