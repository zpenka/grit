# Visualization Module (`engine_visualization.go`)

## Purpose

The visualization module generates textual charts and graphs for commit patterns. It includes contributor flamegraphs, timeline heatmaps, complexity trends, and network diagrams for merge relationships. These visualizations help identify patterns in large histories at a glance.

## Key Functions

### Flamegraph Visualization
- `generateContributorFlamegraph(commits []commit) string`
  - Hierarchical view of contributions
  - Width represents commit count
  - Shows temporal distribution
  - Useful for spotting dominant contributors

### Timeline Heatmap
- `buildActivityHeatmap(commits []commit, granularity string) [][]string`
  - Grid showing commits over time periods
  - Granularity: hourly, daily, weekly, monthly
  - Colors intensity: more commits = darker color
  - Shows team activity patterns

### Complexity Trends
- `computeComplexityTrend(commits []commit) []complexityData`
  - Estimated code complexity over time
  - Lines added/removed trends
  - Cyclomatic complexity approximations
  - Shows codebase health trajectory

### Merge Network Diagram
- `buildMergeNetworkDiagram(commits []commit) string`
  - Shows branch merge relationships
  - ASCII art diagram
  - Node = branch, edge = merge
  - Reveals team branching patterns

### Supporting Functions
- `generateHistogramBuckets(commits []commit, bucketSize int) []histogramBucket`
  - Group commits into time buckets
  - Used by multiple visualizations
  - Flexible bucketing (hourly, daily, weekly)
  
- `colorizeIntensity(value, max float64) string`
  - Map numeric value to color intensity
  - Used in heatmaps
  - Supports ANSI 256 colors

## Visualization Types

### ASCII Art Diagrams
```
Merge Network:
    main ─────┬──────┐
              ├─ feature/auth
              └─ feature/api

Activity Heatmap:
Mon  ░░░░░░░▒▒▒▒▒░░░░░░░  (Low activity)
Tue  ░░░░░▒▒▒▒▒▒▒▒░░░░░░  (Medium activity)
Wed  ░▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░  (High activity)
```

### Flamegraph
Shows contribution hierarchy with proportional widths.

## Design Decisions

### Terminal-Friendly Format
All visualizations use ANSI colors and Unicode box characters compatible with terminal emulators.

### Data Aggregation
Visualizations aggregate commits over time periods for clarity. Configurable granularity (hourly to monthly).

### Complementary Views
Each visualization answers different questions:
- Flamegraph: "Who contributed most?"
- Heatmap: "When was activity highest?"
- Complexity: "Is code getting better or worse?"
- Network: "How do branches interact?"

## Dependencies

- **Standard library**: `fmt`, `strings`, `unicode`
- **Internal**: analytics module for underlying data

## Testing

Tests cover:
- Chart generation with various commit counts
- Time bucketing accuracy
- Merge network detection
- Complexity calculation correctness
- Color mapping for intensity

## Examples

```go
// Generate contributor flamegraph
graph := generateContributorFlamegraph(m.commits)
// Output:
// Alice ████████████████ (60% of commits)
// Bob   ██████ (20%)
// Eve   ████ (13%)
// Charlie ██ (7%)

// Build activity heatmap (daily granularity)
heatmap := buildActivityHeatmap(m.commits, "daily")
// Shows commits per day for past week/month

// Compute complexity over time
trend := computeComplexityTrend(m.commits)
// Returns: [{date: "2024-01-01", complexity: 2.1, loc: 5000}, ...]

// Show merge relationships
diagram := buildMergeNetworkDiagram(m.commits)
// ASCII art showing how branches merge together
```

## Use Cases

### Team Leads
- Identify load imbalance among contributors
- Spot activity patterns and availability
- Track code quality trends

### Project Managers
- Monitor project velocity (commits over time)
- Identify busy periods and bottlenecks
- Plan team capacity

### Code Reviewers
- Find high-complexity areas (risk zones)
- Identify code ownership patterns
- Spot unusual commit patterns (potential problems)

## Performance Considerations

- Flamegraph: O(n) in commit count
- Heatmap: O(n) plus bucketing
- Complexity: O(n) with parsing (cached)
- Network: O(n) for merge detection

All visualizations lazy-load and cache results.

## Rendering

Visualizations displayed in feature panel using consolidated templates:
- `RenderAnalysisUI()`: For structured data
- Direct string rendering: For ASCII art diagrams

## Integration Points

- Analytics module: Used for data extraction
- Rendering module: Displayed in feature panel
- Keybinding handler: User requests visualizations

## Color Scheme

Heatmap intensity (low to high):
- ░ (light) → ▒ (medium) → ▓ (dark) → █ (very dark)

Branch diagram:
- Different colors for different branches
- Merge edges highlighted

## Future Extensions

- Pie charts for contribution distribution
- Line graphs for complexity/velocity trends
- Scatter plots for commit size vs. time
- Gantt charts for team schedules
- Interactive zoom in time heatmaps
