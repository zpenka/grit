# Team and AI Module (`engine_team_ai.go`)

## Purpose

The team and AI module combines team collaboration analytics with AI-powered insights. It includes team velocity tracking, reviewer suggestions, commit classification, anomaly detection, security scanning, and compliance monitoring. Designed for team leads, security auditors, and release managers.

## Key Functions

### Team Analytics
- `calculateTeamStats(commits []commit) []teamStats`
  - Commits per team member
  - Average commit size, frequency
  - Activity over time periods
  
- `suggestReviewers(m model, file string) []reviewerSuggestion`
  - ML-based reviewer selection
  - Ranks by expertise (file familiarity)
  - Scores by availability
  
- `detectPairProgramming(commits []commit) []pairProgrammingData`
  - Find commits by paired authors
  - Co-change frequency metrics
  - Suggests collaboration patterns

### AI-Powered Insights
- `classifyCommit(subject string, hash string) commitClassification`
  - Categorize: feature, fix, refactor, docs, test
  - Keyword-based classification
  - Confidence scores
  
- `autoCompleteMessage(prefix string, commits []commit) []messageCompletion`
  - Suggest commit messages based on history
  - Learns from similar commits
  - Grammar/style consistency
  
- `detectAnomalies(commits []commit) []anomalyData`
  - Find unusual commits (large changes, verbose messages)
  - Potential issues: risky patterns
  - Severity scoring
  
- `findSimilarCommits(commits []commit, targetHash string) []similarCommit`
  - Find related commits by message/content
  - Useful for finding previous solutions

### Compliance & Security
- `checkSigningCompliance(commits []commit, enforced bool) map[string]signingStatus`
  - Verify GPG signatures
  - Enforcement status
  - Compliance reports
  
- `trackLicenseHeaders(hash string) []licenseHeader`
  - Check files for license headers
  - MIT, Apache, GPL detection
  - Missing header identification
  
- `scanForSecurityIssues(hash string, content string) []securityIssue`
  - Search for hardcoded secrets
  - Suspicious patterns (password, key, token)
  - Severity levels
  
- `detectSecrets(hash string, content string) []secretDetection`
  - Deep secret scanning
  - API keys, credentials, tokens
  - False positive reduction

### Release Management
- `detectSemver(commits []commit) []semverData`
  - Identify version bump commits
  - Semantic versioning patterns
  - Release markers
  
- `generateChangelog(commits []commit, version string) *changelogEntry`
  - Auto-generate changelog from commits
  - Categorize: features, bugfixes, breaking changes
  - Format for release notes
  
- `buildReleaseNotes(version string, commits []string) releaseNote`
  - Create release notes document
  - Highlights and contributors
  - Publication-ready format
  
- `trackVersionBumps(commits []commit) []versionBump`
  - Show version history
  - What changed between versions
  - Links commits to versions

## Design Decisions

### Non-Destructive Scanning
All security and compliance checks are read-only. No automatic corrections, only flagging and reporting.

### Confidence Scores
Anomaly detection and classification include confidence metrics. Users can adjust sensitivity.

### Team Privacy
Aggregated metrics respect individual privacy. No storing personal data beyond commit metadata.

### Multi-Level Severity
Findings categorized by severity (info, warning, critical) for prioritization.

## Dependencies

- **Standard library**: `fmt`, `strings`, `regexp`, `time`
- **Internal**: analytics, parsing modules
- **Optional**: External security scanning tools

## Testing

Tests cover:
- Team stats aggregation with varying team sizes
- Reviewer suggestion ranking accuracy
- Commit classification across different message formats
- Anomaly detection sensitivity
- Compliance checking edge cases
- Secret detection accuracy (avoiding false positives)

## Examples

```go
// Get team statistics
stats := calculateTeamStats(m.commits)
// Returns: [{author: "Alice", commits: 150, avgSize: 200, velocity: 15}, ...]

// Suggest reviewers for a file
suggestions := suggestReviewers(m, "auth/login.go")
// Returns: [{reviewer: "bob", expertise: 0.95, availability: 0.8, score: 0.76}, ...]

// Classify commit
class := classifyCommit("Fix: auth bypass in login", "abc123")
// Returns: {category: "fix", confidence: 0.95, reason: "Keyword 'Fix:' detected"}

// Detect anomalies
anomalies := detectAnomalies(m.commits)
// Returns: [{hash: "xyz789", type: "large", severity: 5, description: "10k+ lines changed"}]

// Check for hardcoded secrets
issues := scanForSecurityIssues("def456", content)
// Returns: [{hash: "def456", type: "hardcoded-secret", severity: "critical"}]

// Generate changelog
changelog := generateChangelog(m.commits, "v1.2.0")
// Returns: {version: "v1.2.0", features: [...], bugfixes: [...], breaking: [...]}
```

## Team Metrics

### Velocity Metrics
- Commits per author per week
- Lines of code added/removed
- Average commit size trends

### Collaboration Metrics
- Pair programming frequency
- Code review participation
- Merge conflict rates

### Quality Indicators
- Anomaly frequency
- Security issues found
- License compliance status

## Security Scanning

### Patterns Detected
- API keys (AWS, GitHub, etc.)
- Database credentials
- OAuth tokens
- Private keys (RSA, DSA)
- Connection strings

### False Positive Reduction
- Whitelisting common false positives
- Context analysis
- Entropy checking

## Compliance Features

### Standards Supported
- GPG signature verification
- License header tracking (MIT, Apache, GPL)
- GDPR data deletion tracking
- Audit trail generation

### Reporting
- Compliance dashboards
- Audit reports (exportable)
- Timeline of changes

## AI/ML Capabilities

Current implementation uses heuristics. Future versions may include:
- Actual ML models for classification
- Neural network anomaly detection
- NLP for message analysis

## Integration Points

- Team lead dashboard: Team stats and metrics
- Security audit: Secret and compliance scanning
- Release management: Version tracking and changelog
- Code review: Reviewer suggestions and commit classification

## Rendering

Results displayed using consolidated templates:
- `RenderDataGrid()`: For team stats tables
- `RenderAnalysisUI()`: For security findings
- `RenderStandardUI()`: For compliance status

## Performance Considerations

- Team stats: O(n) in commit count
- Reviewer suggestions: O(n × m) with file matching
- Secret scanning: O(n × content_size) but cached
- Compliance checks: O(n) with parallel processing

## Audit Trail

All compliance and security findings logged for:
- Compliance audits
- Security incident investigation
- Team metrics tracking

## Future Extensions

- Real ML models for commit classification
- Advanced anomaly detection (isolation forests)
- GitHub/GitLab native integration
- Automated issue creation for security findings
- Team workload balancing suggestions
- Code quality metrics (cyclomatic complexity)
