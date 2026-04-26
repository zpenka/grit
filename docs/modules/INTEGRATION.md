# Integration Module (`engine_integration.go`)

## Purpose

The integration module connects grit to external systems: GitHub for PR linking, Jira for issue tracking, and various export formats for data sharing. It enables workflows that combine git history with external issue tracking and enables sharing analysis results.

## Key Functions

### GitHub Integration
- `linkGitHubPR(commit commit, repo string) *githubPR`
  - Search GitHub for PR containing commit hash
  - Returns PR number, title, author, merge status
  - Requires GitHub API token
  
- `fetchPRDetails(prNum int, owner string, repo string) *prDetails`
  - Get PR metadata, reviewers, comments
  - Shows review status and timeline

### Jira Integration
- `linkJiraIssue(commit commit, project string) *jiraIssue`
  - Search Jira for issue referenced in commit message
  - Matches common patterns: PROJ-123, #123
  - Returns issue key, title, status
  
- `fetchIssueDetails(issueKey string) *issueDetails`
  - Get full issue metadata, assignee, resolution

### Data Export
- `exportToCSV(commits []commit) string`
  - Export commit list as CSV
  - Columns: hash, author, date, message, files changed, additions, deletions
  - Importable into Excel/Sheets
  
- `exportToJSON(data interface{}) string`
  - JSON export of commits or analysis results
  - Useful for external tools and dashboards
  
- `exportToXML(data interface{}) string`
  - XML export for system integration
  - Structured schema for compatibility

### Metrics Export
- `generateBuildReport(commits []commit, startDate, endDate int64) *buildReport`
  - Export metrics for CI/CD integration
  - Velocity, quality metrics, deployment info

## Design Decisions

### Read-Only by Default
Integration functions currently read from external systems. No writing back to GitHub/Jira prevents accidental mutations.

### API Token Security
Credentials stored outside git repo (environment variables or secure config).

### Format Flexibility
Multiple export formats support different use cases and downstream tools.

### Offline Fallback
If external services unavailable, shows local data without failure.

## Dependencies

- **External APIs**: GitHub, Jira (HTTP)
- **Standard library**: `fmt`, `strings`, `encoding/json`, `net/http`
- **Internal**: authentication module (credentials management)

## Testing

Tests cover:
- GitHub API error handling
- Jira issue matching (various formats)
- CSV export formatting
- JSON/XML structure validity
- Rate limiting gracefully

## Examples

```go
// Link commit to GitHub PR
pr := linkGitHubPR(commit, "zpenka/grit")
// Returns: {Number: 42, Title: "Fix buffer overflow", Author: "alice"}

// Find associated Jira issue
issue := linkJiraIssue(commit, "PROJ")
// Returns: {Key: "PROJ-123", Title: "Performance regression", Status: "Done"}

// Export commits as CSV
csv := exportToCSV(m.commits)
// Output:
// hash,author,date,message,files_changed,additions,deletions
// abc123,alice,2024-01-01,Fix bug,3,15,8
// def456,bob,2024-01-02,Add feature,5,120,20

// Export analysis results as JSON
json := exportToJSON(m.codeOwnership)
// Output: [{"file": "main.go", "author": "alice", "lines": 1500}, ...]
```

## Authentication

Environment variables:
```bash
export GITHUB_TOKEN=ghp_xxxxx
export JIRA_URL=https://jira.company.com
export JIRA_USER=user@company.com
export JIRA_TOKEN=xxxxx
```

## GitHub Features

- PR detection from commit hashes
- PR review count and status
- Merge status and date
- CI/CD status from PR checks
- Suggestions for reviewers based on PR history

## Jira Features

- Issue detection from commit messages
- Issue status and resolution
- Assignee and reporter info
- Sprint and epic associations
- Link velocity to issues (commits per issue)

## Export Use Cases

### CSV Export
- Team reports and presentations
- Statistical analysis in external tools
- Audit trails and compliance records

### JSON Export
- API integration with other tools
- Dashboard feeding
- Data warehouse import
- Team communication (shareable links)

### XML Export
- Enterprise system integration
- Legacy system compatibility
- Schema validation

## Data Privacy

- No personal data collected beyond commit info
- API tokens never logged or shared
- Results can be anonymized for export

## Integration Points

- Navigation: Browse linked PRs/issues
- Rendering: Show external links in feature panel
- Keybinding handler: Trigger GitHub/Jira lookups

## Performance Considerations

- GitHub/Jira API calls cached (avoid rate limiting)
- Async lookups prevent UI blocking
- Results persist in model for session duration

## Rate Limiting

GitHub API: 60 requests/hour (unauthenticated), 5000/hour (authenticated)
Jira API: Varies by instance

Caching helps stay within limits.

## Error Handling

- Missing tokens: Skip integration gracefully
- API errors: Show fallback local data
- Network issues: Queue for retry

## Future Extensions

- Write to GitHub (create issues/PRs)
- Bidirectional Jira sync
- Other platforms (GitLab, Azure DevOps)
- Custom webhooks
- Automated issue creation from anomalies
- Slack/Teams notifications for milestones
