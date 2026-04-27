# Analytics Module (`engine_analytics.go`)

## Purpose

The analytics module provides insights into commit patterns, code ownership, activity trends, and quality metrics. It includes author statistics, code hotspots, collaboration analysis, and bisection support. This module transforms raw commit data into actionable insights.

## Key Functions

### Code Ownership Analysis
- `analyzeAuthorStats(commits []commit) []authorStats`
  - Commits per author, average commit size
  - Identifies primary maintainers
  
- `analyzeCodeOwnership(commits []commit) []fileOwnership`
  - Who changed which files
  - Line counts per author
  - Useful for onboarding and team planning

### Activity Patterns
- `analyzeActivityHeatmap(commits []commit) [][]int`
  - Commits per hour/day/week
  - Identifies team timezone and activity patterns
  - Used for timeline visualization
  
- `analyzeHotspots(commits []commit) []fileHotspot`
  - Most frequently modified files
  - Churn metrics
  - Identifies stability risks

### Collaboration Analysis
- `analyzePairProgramming(commits []commit) []pairData`
  - Commits by paired authors
  - Co-change frequency
  
- `analyzeMergePatterns(commits []commit) []mergeMetric`
  - Merge frequency, conflict patterns
  - Branch management insights

### Advanced Analytics
- `analyzeCodeComplexity(commits []commit) []complexityTrend`
  - Estimated complexity over time
  - File complexity changes
  
- `detectMicroservicePatterns(commits []commit) []microserviceHint`
  - Suggests service boundaries
  - Based on file path patterns and change isolation

## Feature Complexity: Bisection

- `initiateBisect(m model, badHash, goodHash string) model`
  - Start binary search for regression
  - Guides user through commit space
  
- `updateBisectState(m model, result string) model`
  - Record test results (good/bad)
  - Narrows search space
  - Can skip untestable commits

## Design Decisions

### On-Demand Computation
Analytics are computed only when user requests them. This keeps startup fast and memory usage low.

### Caching Results
Results cached in `scache` statistics cache. Multiple accesses don't recompute.

### Pure Functions
All analytics functions take commits and return analysis—no side effects. Makes them testable and composable.

### Comprehensive Author Handling
Handles variations: same author different emails, co-authored commits, email format changes.

## Dependencies

- Standard library: `fmt`, `strings`, `regexp`, `time`, `sort`
- **Internal**: caching, optimization modules

## Testing

Tests cover:
- Single author repositories
- Multiple authors with varying contribution levels
- Large commit histories (1000+ commits)
- Hotspot identification accuracy
- Bisection logic (correct narrowing of search space)

## Examples

```go
// Analyze author statistics
stats := analyzeAuthorStats(m.commits)
// Returns: [{author: "Alice", commits: 150, avgSize: 200}, ...]

// Find frequently modified files
hotspots := analyzeHotspots(m.commits)
// Returns: [{file: "main.go", changeCount: 45, lastChanged: 2 days ago}, ...]

// Start binary search for regression
m = initiateBisect(m, "badHash", "goodHash")

// Record result and narrow search
m = updateBisectState(m, "bad")  // Current commit has bug
```

## Performance Considerations

- Author stats: O(n) in commit count
- Hotspots: O(n × m) where m = avg files per commit
- Complexity analysis: O(n) with caching
- Bisection: O(log n) iterations to find regression

All expensive operations lazy-load and cache results.

## Rendering

Results displayed using consolidated templates:
- `RenderAnalysisUI()`: Used for stats, ownership, hotspots
- `RenderDataGrid()`: Used for detailed data tables
- `RenderStandardUI()`: Used for lists with status

## Integration Points

- Navigation: Bisection guides cursor movement
- Rendering: Results displayed in feature panel
- Keybinding handler: User initiates analytics and bisection

## Future Extensions

- Machine learning classification (feature/fix/refactor)
- Commit impact scoring
- Bus factor analysis
- Code review metrics
- Performance regression detection
- Commit quality scoring
