# Developer Guide

This guide explains how to develop, extend, and contribute to grit. It covers architecture deep-dives, testing strategies, and step-by-step contribution patterns.

## Getting Started

### Prerequisites
- Go 1.21+
- A terminal with 24-bit color support (recommended)
- Git

### Build & Run
```bash
go build -o grit .
./grit
```

### Testing
```bash
go test ./...              # Run all tests
go test -cover ./...       # With coverage
go test -run TestName      # Specific test
go test -v ./...           # Verbose output
```

All 914 tests should pass. Tests run automatically on PRs via GitHub Actions.

## Quick Troubleshooting

| Issue | Solution |
|-------|----------|
| Build fails with missing package | Run `go mod download` |
| Tests fail intermittently | Run again; some tests depend on timing |
| Terminal looks weird | Enable 24-bit color in your terminal settings |
| Changes don't appear | Rebuild with `go build -o grit .` |

## Architecture Overview

grit is organized into **13 focused modules**:

### Core Infrastructure
These provide foundational functionality:
- **engine_types.go** - All type definitions (150+)
- **engine_parsing.go** - Parse git data (commits, diffs, branches)
- **engine_navigation.go** - Handle UI navigation (cursor, panels, scrolling)
- **engine_filtering.go** - Filter and search commits
- **engine_cache.go** - Cache mechanisms (LRU, memoization)
- **engine_optimization.go** - Performance utilities

### Rendering
- **engine_render_consolidation.go** - Consolidated render templates
- **engine_rendering.go** - Main UI rendering engine (2700+ lines)

### Features (Modular)
Each module handles a logical feature group:
- **engine_analytics.go** - Analytics, bisect, code ownership
- **engine_git_ops.go** - Git operations (rebase, cherry-pick, reset, amend)
- **engine_workflows.go** - Worktrees, stash, tag management
- **engine_visualization.go** - Graphs, timelines, heatmaps
- **engine_integration.go** - GitHub, Jira, exports
- **engine_team_ai.go** - Team analytics, AI, compliance, versioning

## Architecture Deep-Dives

### Core Data Flow

```
User Input (Keybinding)
    ↓
grit.go: handleKeyBinding()
    ↓
Update() method processes event
    ↓
Model state changes (commits, cursor, filters, etc.)
    ↓
View() renders based on state
    ↓
Terminal displays output
```

### The Model Struct

The `model` struct is the **single source of truth** for application state:

```go
type model struct {
    // Navigation
    cursor          int          // Current commit position
    diffOffset      int          // Scroll position in diff panel
    focus           panel        // Which panel has focus (commits/diff)
    
    // Data
    commits         []commit     // All loaded commits
    diffLines       []diffLine   // Current diff
    branches        []string     // Available branches
    
    // Filters
    query           string       // Search query
    authorFilter    string       // Author-based filter
    sinceFilter     string       // Date-based filter
    
    // Features
    showAnalytics   bool         // Toggle analytics menu
    analyticsData   []...        // Feature computation results
    showGitOps      bool         // Toggle git ops menu
    // ... 50+ more feature flags and data fields
    
    // Caches
    dcache          *DiffCache   // LRU diff cache
    scache          map[...]     // Statistics cache
    recache         map[...]     // Regex pattern cache
}
```

**Key principle:** All state modifications happen on the model, making the codebase predictable and testable.

### Module Responsibilities

#### engine_parsing.go (Commit/Diff Parsing)
**Purpose:** Convert git command output into structured data

**Key Functions:**
- `parseCommits(output string) []commit` - Parse `git log` output
- `parseDiff(output string) []diffLine` - Parse `git show` output
- `parseFileItems(diff []diffLine) []fileItem` - Extract changed files

**Testing Approach:** Feed mock git output, verify parsed structures
```go
input := "hash\x00shortHash\x00Author\x001h ago\x00Subject\n"
commits := parseCommits(input)
assert(commits[0].author == "Author")
```

#### engine_filtering.go (Search & Filter)
**Purpose:** Reduce commit list based on user criteria

**Key Functions:**
- `filterCommits(commits []commit, query string) []commit` - Search by subject/author/hash
- `filterByAuthor(commits []commit, author string) []commit` - Author filter
- `isWithinDays(when string, days int) bool` - Date range filter

**Performance:** Caches compiled regex patterns to avoid recompiling on each keystroke

#### engine_navigation.go (Cursor & Panel Movement)
**Purpose:** Move through commits and toggle panels

**Key Functions:**
- `moveCursorUp(m model) model` - Move cursor up (wrap around)
- `moveCursorDown(m model) model` - Move cursor down (wrap around)
- `scrollDiffUp/Down(m model) model` - Scroll diff panel
- `switchPanel(m model) model` - Toggle between commit list and diff

**Boundary Handling:** All functions clamp cursor to valid range [0, len(commits)-1]

#### engine_cache.go (Performance Optimization)
**Purpose:** Store expensive computation results

**Three Cache Types:**
1. **DiffCache** (LRU) - Stores parsed diffs by commit hash
   - Bounded size to prevent memory issues
   - Tracks hit/miss rate for monitoring
   
2. **Statistics Cache** - Stores computed statistics
   - File change counts, author stats
   - Invalidated when commit list changes
   
3. **Regex Cache** - Stores compiled regex patterns
   - Patterns compiled once, used many times
   - Critical for search performance

**Example Usage:**
```go
// Without cache: recompile regex every keystroke
re := regexp.MustCompile(userQuery)

// With cache: compile once, reuse
re, _ := compileRegex(m.recache, userQuery)
```

#### engine_rendering.go (UI Output)
**Purpose:** Convert model state to terminal display

**Key Functions:**
- `renderCommitList(m model) string` - Format commit list panel
- `renderDiffPanel(m model) string` - Format diff/feature panel
- `renderFooter(m model) string` - Show keyboard hints

**Strategy:** Use Lipgloss for styling, arrange panels with fixed widths based on terminal size

#### engine_analytics.go through engine_team_ai.go (Feature Modules)
**Purpose:** Compute insights from commit history

**Pattern for Each Feature:**
1. **Data function** - `analyze*() []resultType`
2. **Render function** - `render*UI() string`
3. **Toggle in model** - `show*` boolean flag
4. **Menu entry** - Add to feature menu

**Example:**
```go
// 1. Analyze
func analyzeCodeOwnership(commits []commit) []ownershipData { ... }

// 2. Render
func renderCodeOwnershipUI(m model, width int) string {
    return RenderAnalysisUI("Code Ownership", m.ownershipData)
}

// 3. Wire in model
m.showCodeOwnership = !m.showCodeOwnership

// 4. Menu
"Code Ownership" -> dispatchAnalyticsFeature(m, 0)
```

### State Management Patterns

**Immutability Pattern:**
```go
// Don't: mutate in place
func updateModel(m *model) {
    m.cursor++
}

// Do: return new model
func moveCursorDown(m model) model {
    m.cursor++
    return m
}
```

**Model Threading:**
All operations take `model` receiver and return modified `model`, enabling:
- Functional programming style
- Undo/redo capability (if needed)
- Predictable state transitions

**Lazy Computation:**
```go
// Only compute when user requests
if m.showAnalytics && len(m.analyticsData) == 0 {
    m.analyticsData = analyzeCodeOwnership(m.commits)
}
```

### Update Cycle (Bubble Tea Integration)

```
Update() receives user input
    ↓
Match key to handler (handleKeyBinding)
    ↓
Modify model based on action
    ↓
Return new model + Cmd (if any)
    ↓
View() called with new model
    ↓
Terminal renders output
    ↓
(repeat on next input)
```

**Key Files:**
- `grit.go:Update()` - Main event dispatcher
- `grit.go:View()` - Main render orchestrator
- `grit.go:handleKeyBinding()` - Key handling logic

## Testing Strategies

### TDD (Test-Driven Development) Workflow

grit follows **red-green-refactor** for new features:

1. **Red:** Write failing test for desired behavior
   ```go
   func TestNewFeature_Works(t *testing.T) {
       result := newFeature(testData)
       if !isCorrect(result) {
           t.Error("feature should work")
       }
   }
   ```

2. **Green:** Implement minimal code to pass test
   ```go
   func newFeature(data []commit) bool {
       return true  // Just enough to pass
   }
   ```

3. **Refactor:** Improve implementation while keeping tests green
   ```go
   func newFeature(data []commit) bool {
       // Real implementation
       return analyzeData(data)
   }
   ```

### Test Categories

#### Unit Tests (Most Common)
Test individual functions in isolation:

```go
// Test edge cases
func TestFilterCommits_Empty(t *testing.T) {
    result := filterCommits([]commit{}, "query")
    if len(result) != 0 {
        t.Error("expected empty result")
    }
}

// Test normal case
func TestFilterCommits_MatchesQuery(t *testing.T) {
    commits := []commit{
        {subject: "fix bug"},
        {subject: "add feature"},
    }
    result := filterCommits(commits, "fix")
    if len(result) != 1 {
        t.Error("should match one commit")
    }
}

// Test boundary
func TestFilterCommits_SpecialCharacters(t *testing.T) {
    commits := []commit{{subject: "fix: [bug]"}}
    result := filterCommits(commits, "[bug]")
    if len(result) != 1 {
        t.Error("should handle special chars")
    }
}
```

#### Integration Tests
Test multiple components working together:

```go
func TestNavigationWithFiltering(t *testing.T) {
    m := newModel(".")
    m.commits = makeCommits(100)
    
    // Filter commits
    m.query = "fix"
    filtered := visibleCommits(m)
    
    // Navigate filtered list
    m = moveCursorDown(m)
    
    // Verify state is consistent
    if m.cursor >= len(filtered) {
        t.Error("cursor out of bounds after navigation")
    }
}
```

#### Regression Tests
Ensure known bugs stay fixed:

```go
// Issue #42: Cursor not wrapping at end of list
func TestCursorWrapsAtEnd(t *testing.T) {
    m := newModel(".")
    m.commits = makeCommits(5)
    m.cursor = 4  // Last commit
    
    m = moveCursorDown(m)
    
    if m.cursor != 0 {
        t.Error("cursor should wrap to beginning")
    }
}
```

### Test Helpers (engine_test_helpers.go)

**Available Utilities:**
```go
makeCommits(n int) []commit
    // Creates n test commits with realistic data
    // IDs: "aaa1111", "bbb2222", etc.

makeDiffLines(n int) []diffLine
    // Creates n diff lines with mixed kinds
    // Useful for testing diff processing

makeNamedCommits(authors []string) []commit
    // Creates commits by specified authors
    // Good for testing author-based features
```

**Example Usage:**
```go
func TestAnalyticsWithAuthors(t *testing.T) {
    commits := makeNamedCommits([]string{"Alice", "Bob", "Charlie"})
    stats := analyzeAuthors(commits)
    if len(stats) != 3 {
        t.Error("should have stats for 3 authors")
    }
}
```

### Common Testing Patterns

#### Testing Cache Behavior
```go
func TestDiffCacheHitsOnSecondAccess(t *testing.T) {
    cache := &DiffCache{}
    
    // First access: miss
    result1 := getDiff(cache, "abc123")
    
    // Second access: hit
    result2 := getDiff(cache, "abc123")
    
    if !cacheEqual(result1, result2) {
        t.Error("cache should return same result")
    }
    
    // Verify cache size is bounded
    if cache.len() > maxCacheSize {
        t.Error("cache exceeded max size")
    }
}
```

#### Testing State Mutations
```go
func TestModelImmutability(t *testing.T) {
    m1 := newModel(".")
    m1.cursor = 5
    
    // Operate on copy
    m2 := moveCursorDown(m1)
    
    // Original should be unchanged
    if m1.cursor != 5 {
        t.Error("original model was modified")
    }
    
    // New model should be different
    if m2.cursor != 6 {
        t.Error("new model should have updated cursor")
    }
}
```

#### Testing Error Conditions
```go
func TestParsingInvalidData(t *testing.T) {
    badInputs := []string{
        "",                    // Empty
        "no-delimiters",      // Missing fields
        "\x00\x00\x00",       // Only delimiters
    }
    
    for _, input := range badInputs {
        commits := parseCommits(input)
        // Should handle gracefully (nil or empty)
        if commits != nil && len(commits) > 0 {
            t.Errorf("should reject bad input: %q", input)
        }
    }
}
```

### Test Coverage Goals

**Current Coverage:** 71.7% overall, 83.4% in core package

**Target by Module:**
- **Core logic** (parsing, filtering, navigation): 85%+
- **Rendering**: 75%+ (harder to test, many permutations)
- **Git integration**: 40%+ (system-level, integration-tested)
- **Features**: 70%+ (each should have happy path + edge cases)

**Coverage Check:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View in browser
```

### Debugging Tests

**Print Debugging:**
```go
func TestFeature(t *testing.T) {
    result := featureFunc()
    t.Logf("Result: %#v", result)  // Use -v flag to see output
}
```

**Run with Verbose Output:**
```bash
go test -v -run TestFeatureName
```

**Run Single Test:**
```bash
go test -run TestFeatureName -v
```

**Run with Timeout:**
```bash
go test -timeout 5s ./...  # Fail if test hangs > 5s
```

## Adding a New Feature: Complete Walkthrough

This section demonstrates the complete process with a real example: "File Risk Score" feature that identifies risky files based on change frequency and test coverage.

### Step 1: Plan the Feature

Ask yourself:
- **What problem does it solve?** Find files that change frequently but have low test coverage
- **Where does it belong?** Analytics features → `engine_analytics.go`
- **What's the input?** List of commits
- **What's the output?** List of files with risk scores
- **How will users access it?** Via analytics menu (`a` key)

### Step 2: Define Data Types (engine_types.go)

Add the data structure to hold results:

```go
type fileRiskScore struct {
    path          string
    changeCount   int
    testCoverage  float64
    riskScore     float64  // 0-100
    lastModified  string
}
```

Add model fields:

```go
type model struct {
    // ... existing fields ...
    showFileRiskUI    bool
    fileRiskScores    []fileRiskScore
}
```

### Step 3: Write Tests First (TDD)

Create `engine_analytics_test.go` entries:

```go
func TestAnalyzeFileRisk_Empty(t *testing.T) {
    result := analyzeFileRisk([]commit{})
    if result != nil && len(result) > 0 {
        t.Error("empty commits should return empty or nil")
    }
}

func TestAnalyzeFileRisk_SingleFile(t *testing.T) {
    commits := []commit{
        {hash: "a", subject: "change file.go", files: []string{"file.go"}},
        {hash: "b", subject: "change file.go", files: []string{"file.go"}},
    }
    result := analyzeFileRisk(commits)
    if len(result) != 1 {
        t.Error("should detect 1 unique file")
    }
    if result[0].changeCount != 2 {
        t.Error("should count 2 changes")
    }
}

func TestAnalyzeFileRisk_CalculatesRiskScore(t *testing.T) {
    commits := makeCommits(10)
    result := analyzeFileRisk(commits)
    
    for _, risk := range result {
        if risk.riskScore < 0 || risk.riskScore > 100 {
            t.Errorf("risk score out of bounds: %f", risk.riskScore)
        }
    }
}

func TestAnalyzeFileRisk_RanksByRisk(t *testing.T) {
    // Verify results are sorted by risk (descending)
    result := analyzeFileRisk(makeCommits(20))
    for i := 0; i < len(result)-1; i++ {
        if result[i].riskScore < result[i+1].riskScore {
            t.Error("results should be sorted by risk descending")
        }
    }
}
```

Run tests (they should fail):
```bash
go test -run TestAnalyzeFileRisk -v
```

### Step 4: Implement the Feature

In `engine_analytics.go`:

```go
// analyzeFileRisk calculates risk score for each file based on
// change frequency and test coverage.
func analyzeFileRisk(commits []commit) []fileRiskScore {
    if len(commits) == 0 {
        return nil
    }
    
    // Count changes per file
    fileStats := make(map[string]*fileRiskScore)
    for _, c := range commits {
        for _, file := range c.files {
            if _, exists := fileStats[file]; !exists {
                fileStats[file] = &fileRiskScore{path: file}
            }
            fileStats[file].changeCount++
            fileStats[file].lastModified = c.when
        }
    }
    
    // Convert to slice and calculate risk scores
    var results []fileRiskScore
    for _, risk := range fileStats {
        risk.riskScore = calculateRiskScore(*risk, len(commits))
        results = append(results, *risk)
    }
    
    // Sort by risk (descending)
    sortByRisk(results)
    
    return results
}

func calculateRiskScore(risk fileRiskScore, totalCommits int) float64 {
    changeFrequency := float64(risk.changeCount) / float64(totalCommits) * 100
    // Assume ~70% test coverage for now (would integrate real data)
    testCoverage := 70.0
    
    // High frequency + low coverage = high risk
    risk := changeFrequency * (100 - testCoverage) / 100
    return math.Min(100, risk)
}

func sortByRisk(risks []fileRiskScore) {
    for i := 0; i < len(risks)-1; i++ {
        for j := i + 1; j < len(risks); j++ {
            if risks[j].riskScore > risks[i].riskScore {
                risks[i], risks[j] = risks[j], risks[i]
            }
        }
    }
}
```

Run tests again (they should pass):
```bash
go test -run TestAnalyzeFileRisk -v
```

### Step 5: Add Render Function

In `engine_analytics.go`:

```go
// renderFileRiskUI displays file risk analysis using the standard template.
func renderFileRiskUI(m model, width int) string {
    var items []string
    statusMap := make(map[string]string)
    
    for _, risk := range m.fileRiskScores {
        status := "low"
        if risk.riskScore > 70 {
            status = "high"
        } else if risk.riskScore > 40 {
            status = "medium"
        }
        
        item := fmt.Sprintf("%s (%d changes, %.0f risk)", 
            risk.path, risk.changeCount, risk.riskScore)
        items = append(items, item)
        statusMap[item] = status
    }
    
    return RenderStandardUI(RenderConfig{
        Title:     "File Risk Analysis",
        Items:     items,
        HasStatus: true,
        StatusMap: statusMap,
    })
}
```

### Step 6: Wire into Model & Menu

Update `engine_analytics.go` to add to menu:

```go
// In the analytics menu items (around line 50)
const analyticsMenuLen = 8  // Increase count

var analyticsMenuItems = []string{
    "Code Ownership",
    "Hotspots",
    "Bisect",
    "Complexity",
    "Statistics",
    "Heatmap",
    "Linting",
    "File Risk",  // New item
}

// In dispatchAnalyticsFeature (around line 100)
case 7:  // File Risk
    m.showFileRiskUI = !m.showFileRiskUI
    if m.showFileRiskUI && len(m.fileRiskScores) == 0 {
        m.fileRiskScores = analyzeFileRisk(m.commits)
    }
```

### Step 7: Add to Main View

In `grit.go`, add rendering in the `View()` method:

```go
// Around line 700, in View() where other features render
if m.showFileRiskUI {
    return renderFileRiskUI(m, m.width)
}
```

### Step 8: Update Model Type

In `engine_types.go`, add fields (done in Step 2):

```go
type model struct {
    // ... existing fields ...
    showFileRiskUI    bool
    fileRiskScores    []fileRiskScore
}
```

### Step 9: Run All Tests

```bash
go test ./...
```

All 914+ tests should pass (including your new ones).

### Step 10: Documentation

Add godoc comments:

```go
// analyzeFileRisk calculates risk score for each file based on
// change frequency and test coverage. Files that change frequently
// but have low test coverage score higher on the risk scale.
// Results are sorted by risk descending.
func analyzeFileRisk(commits []commit) []fileRiskScore { ... }
```

Update README.md:

```markdown
| `a` menu | File Risk | Shows files with highest change frequency vs. test coverage |
```

Update USAGE.md if needed:

```markdown
## Finding Risky Code
1. Press `a` → Select "File Risk"
2. High-risk files appear at top
3. Focus refactoring efforts on these files
```

### Step 11: Create PR

```bash
git checkout -b feature/file-risk-analysis
git add .
git commit -m "Add file risk analysis feature

Analyzes commit history to identify files with high change frequency
and low test coverage. Helps teams prioritize refactoring efforts.

- Added fileRiskScore data type
- Implemented analyzeFileRisk with sorting and scoring
- Added renderFileRiskUI for display
- Integrated with analytics menu
- Added 4 test cases covering empty, single file, score bounds, sorting"
```

Push and create PR for review.

---

**Summary:** Feature complete in 11 steps, all TDD with tests passing, fully integrated and documented.

## Common Pitfalls & Gotchas

### Pitfall 1: Mutating Model Instead of Returning
```go
// ❌ Wrong: Modifies in place
func moveCursor(m *model) {
    m.cursor++
}

// ✅ Right: Returns new model
func moveCursor(m model) model {
    m.cursor++
    return m
}
```

### Pitfall 2: Forgetting Cache Invalidation
```go
// ❌ Wrong: Cache becomes stale when commits change
m.commits = newCommits
// Cache still has old data!

// ✅ Right: Invalidate affected caches
m.commits = newCommits
m.scache = make(map[string]interface{})  // Clear stats
m.analyticsData = nil  // Clear feature data
```

### Pitfall 3: Unbounded Loops in Large Repos
```go
// ❌ Wrong: O(n²) on large commit lists
for _, c1 := range commits {
    for _, c2 := range commits {
        // Process pair
    }
}

// ✅ Right: Use early exit or better algorithm
for i, c1 := range commits {
    for _, c2 := range commits[i+1:] {  // Only future commits
        // Process pair
    }
}
```

### Pitfall 4: UI Crashes on Empty Data
```go
// ❌ Wrong: Panics if no data
func renderFeature(m model) string {
    return fmt.Sprintf("%s", m.featureData[0].name)  // Index 0 doesn't exist
}

// ✅ Right: Check before accessing
func renderFeature(m model) string {
    if len(m.featureData) == 0 {
        return "No data available"
    }
    return fmt.Sprintf("%s", m.featureData[0].name)
}
```

### Pitfall 5: Regex Compilation in Hot Path
```go
// ❌ Wrong: Recompiles regex on every keystroke
func handleKey(m model) model {
    re := regexp.MustCompile(m.query)  // Slow!
    m.filtered = filterWithRegex(m.commits, re)
    return m
}

// ✅ Right: Cache compiled regex
func handleKey(m model) model {
    re, _ := compileRegex(m.recache, m.query)
    m.filtered = filterWithRegex(m.commits, re)
    return m
}
```

### Pitfall 6: Not Testing Edge Cases
```go
// ❌ Wrong: Only tests happy path
func TestFilter(t *testing.T) {
    result := filter([]commit{{subject: "fix"}}, "fix")
    if len(result) != 1 { t.Error("failed") }
}

// ✅ Right: Test edges too
func TestFilter(t *testing.T) {
    // Happy path
    result := filter([]commit{{subject: "fix"}}, "fix")
    assert(len(result) == 1)
    
    // Empty input
    result = filter([]commit{}, "fix")
    assert(len(result) == 0)
    
    // No matches
    result = filter([]commit{{subject: "feature"}}, "fix")
    assert(len(result) == 0)
    
    // Special characters
    result = filter([]commit{{subject: "fix [bug]"}}, "[bug]")
    assert(len(result) == 1)
}
```

### Pitfall 7: Terminal Size Not Respected
```go
// ❌ Wrong: Assumes 80-column terminal
items := items[:80]

// ✅ Right: Use actual terminal width
items := items[:m.width]
```

### Pitfall 8: Race Conditions in Goroutines
```go
// ❌ Wrong: Goroutine modifies shared map
go func() {
    m.cache[key] = value  // Race condition!
}()

// ✅ Right: Use channels or mutexes
results := make(chan cacheEntry)
go func() {
    results <- cacheEntry{key: key, value: value}
}()
// Process results serially
```

## Performance Profiling & Optimization

### Profiling grit

**Generate CPU Profile:**
```bash
go test -cpuprofile=cpu.prof -bench . ./...
go tool pprof cpu.prof
```

**Analyze Memory Usage:**
```bash
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

**Common pprof Commands:**
```
(pprof) top           # Show top 10 hotspots
(pprof) list funcName # Show function source with profiling data
(pprof) graph         # Generate call graph
```

### Optimization Checklist

1. **Profile First** - Find actual bottleneck before optimizing
2. **Cache Hits** - Are expensive operations cached?
3. **Algorithm** - Is the algorithm efficient? (O(n) vs O(n²)?)
4. **String Operations** - Avoid repeated string allocations
5. **Regex** - Are patterns compiled once?
6. **Goroutines** - Use for I/O-bound work, not CPU-bound
7. **Memory** - Check for unbounded allocations (maps, slices without capacity)

### Examples of Optimizations Done

**Example 1: Regex Caching**
```go
// Before: Recompiled on every search keystroke
re := regexp.MustCompile(query)

// After: Compiled once, cached
re, _ := compileRegex(m.recache, query)
// Result: 10x faster search response
```

**Example 2: Diff Caching**
```go
// Before: Reparsed diff for every render
diff := parseDiff(gitShow(commit))

// After: LRU cache with miss check
diff := getDiffFromCache(m.dcache, commit.hash)
// Result: Instant diff panel switching for viewed commits
```

**Example 3: Lazy Loading**
```go
// Before: Computed all features on load
m.flameData = buildFlame(commits)
m.heatmap = buildHeatmap(commits)
// Result: 10 second startup delay

// After: Compute on-demand
if m.showFlamegraph && m.flameData == nil {
    m.flameData = buildFlame(commits)
}
// Result: Instant startup
```

## Code Style

### Naming Conventions
```go
// Functions
func analyzeCodeOwnership(...) []...  // exported: PascalCase
func filterCommits(...) []commit      // exported: PascalCase
func parseCommits(...) []commit       // exported: PascalCase

func helpText() string                // unexported: camelCase
func cacheKey(hash string) string     // unexported: camelCase

// Types
type Model struct { }                 // PascalCase
type model struct { }                 // unexported: camelCase (central state)
type fileRiskScore struct { }         // Compound words: camelCase

// Constants
const MAX_CACHE_SIZE = 1000
const DEFAULT_WIDTH = 80
const CACHE_TTL_SECONDS = 300

// Variables
var globalCache *DiffCache            // Avoid globals when possible
m := newModel(".")                    // Model variable always 'm'
```

### Function Patterns

**Feature Functions** follow a consistent pattern:
```go
// Calculate feature: takes data, returns results
func analyzeFeature(commits []commit) []featureResult { ... }

// Render feature: takes model + width, returns string
func renderFeatureUI(m model, width int) string { ... }

// Dispatch from menu: updates model, returns it
func dispatchFeature(m model, idx int) model { ... }
```

**Navigation Functions** modify model:
```go
func moveCursorDown(m model) model { ... }
func scrollDiffUp(m model) model { ... }
func switchPanel(m model) model { ... }
```

**Parsing Functions** convert strings to types:
```go
func parseCommits(output string) []commit { ... }
func parseDiff(output string) []diffLine { ... }
func parseFileItems(diff []diffLine) []fileItem { ... }
```

### Comments & Documentation

**Good Comments** (why/context):
```go
// If navigation history isn't at the end, discard future history.
// This implements a git-like history model where going back then
// moving forward doesn't restore the "future" branch.
if m.navHistoryIdx < len(m.navHistory)-1 {
    m.navHistory = m.navHistory[:m.navHistoryIdx+1]
}
```

**Bad Comments** (obvious what code does):
```go
// ❌ Obvious - comment adds no value
m.cursor++  // Increment cursor
```

**Godoc Format** (exported functions):
```go
// analyzeCodeOwnership analyzes git commit history to determine
// which developers own each file based on change frequency.
//
// The analysis considers both number and recency of changes,
// giving more weight to recent modifications. Results are sorted
// by ownership percentage descending.
//
// Example:
//    commits := []commit{ ... }
//    ownership := analyzeCodeOwnership(commits)
//    fmt.Println(ownership[0].file)  // Most-owned file
func analyzeCodeOwnership(commits []commit) []ownershipData { ... }
```

**When to Comment:**
- Non-obvious algorithms or heuristics
- Workarounds for Go language quirks
- Performance-critical decisions
- Cross-module dependencies
- Known limitations or TODOs

**When NOT to Comment:**
- Self-explanatory code with good names
- Loop mechanics (for i := 0; ...)
- Type assertions (if val, ok := ...)
- Straightforward function bodies

## Module-Specific Guides

### engine_parsing.go - Git Data Parsing

**Responsibility:** Convert raw git output to structured Go types

**Key Challenges:**
- Git output varies by version
- Special characters in commit messages
- Large diffs (100k+ lines)

**Best Practices:**
```go
// Always check for sufficient fields
if len(fields) < 5 {
    return nil  // Malformed line
}

// Handle special characters
subject := strings.TrimSpace(subject)
subject = strings.Replace(subject, "\n", " ", -1)

// Stream process large diffs
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    // Process one line at a time, not whole file
}
```

### engine_filtering.go - Search & Filtering

**Responsibility:** Reduce commit list based on user queries

**Key Challenges:**
- Regex performance on repeated queries
- Case-insensitive matching
- Partial hash matching (first N chars)

**Best Practices:**
```go
// Cache compiled regexes
re, _ := compileRegex(m.recache, query)

// Combine filter criteria
visible := filterByAuthor(commits, author)
visible = filterByDate(visible, since)
visible = filterByQuery(visible, query)  // Chain filters

// Short-circuit empty results
if len(commits) == 0 {
    return nil
}
```

### engine_rendering.go - UI Generation

**Responsibility:** Convert model state to terminal display

**Key Challenges:**
- Terminal width/height constraints
- Unicode character width (emoji)
- Color ANSI codes

**Best Practices:**
```go
// Respect terminal dimensions
items := items[:min(len(items), m.height-2)]

// Use Lipgloss for styling
style := lipgloss.NewStyle().Foreground(lipgloss.Color("41"))
output := style.Render(text)

// Format for readability
// Don't: All text in one line
// Do: Break into logical sections (header, items, footer)
```

### engine_cache.go - Caching Patterns

**Responsibility:** Store expensive computation results

**Strategies:**
1. **LRU Cache** (diffs) - Fixed size, evict least recent
2. **Map Cache** (stats) - Unlimited, invalidate on change
3. **Pattern Cache** (regex) - Store compiled patterns

**Best Practices:**
```go
// Check cache first
if cached, ok := m.scache[key]; ok {
    return cached  // Hit
}

// Compute on miss
result := expensive_computation()
m.scache[key] = result
return result

// Invalidate on relevant changes
if commits changed {
    m.scache = make(map[string]interface{})
}
```

## Debugging & Troubleshooting

### Debugging Techniques

**Print Debugging (Logging to stderr):**
```go
func Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m := m.(model)
    fmt.Fprintf(os.Stderr, "Debug: cursor=%d, commits=%d\n", m.cursor, len(m.commits))
    // ... rest of Update
}
```

**View stderr output:**
```bash
./grit 2>debug.log  # Capture stderr
tail -f debug.log   # Watch debug output
```

**Conditional Debugging:**
```go
const DEBUG = true

func debugLog(format string, args ...interface{}) {
    if DEBUG {
        fmt.Fprintf(os.Stderr, format+"\n", args...)
    }
}
```

### Testing Locally

**Test against real repo:**
```bash
cd /path/to/my/repo
./grit                    # Test interactively
```

**Test with specific history:**
```bash
# Test with small repo
git clone --depth 100 <url> test-repo
cd test-repo
./grit
```

**Replay test scenarios:**
```bash
# Create test fixture repo
mkdir test-repo && cd test-repo
git init
echo "file1" > file.txt && git add . && git commit -m "Initial commit"
echo "file2" >> file.txt && git add . && git commit -m "Add content"
# ... create scenario

cd - && ./grit test-repo   # Test against fixture
```

### Common Issues & Solutions

**UI looks garbled:**
- Check terminal width is 80+ columns
- Check terminal supports 24-bit color
- Try different terminal emulator

**Tests fail randomly:**
- Some tests depend on timing
- Run again: `go test ./... && go test ./...`
- Check for race conditions: `go test -race ./...`

**Binary won't run after changes:**
- Clear build cache: `go clean`
- Rebuild: `go build -o grit .`
- Check syntax: `go vet ./...`

**Memory usage increasing:**
- Check for cache leaks (clearing stale entries)
- Run: `go test -memprofile=mem.prof ./...`
- Analyze: `go tool pprof -alloc_space mem.prof`

## CI/CD & Quality Checks

### GitHub Actions Workflow

Tests run automatically:
- **On PR creation/update**: Pull request testing
- **On push to main**: Post-merge verification
- **Workflow location**: `.github/workflows/tests.yml`

**What runs:**
- `go test ./...` - All tests (914 tests)
- Code coverage tracking
- Test result reporting

### Required Status Checks

The `main` branch has these requirements:
- All tests must pass
- No broken builds
- Coverage maintained (71.7%+)

**To block merges on failure:**
1. Go to repository Settings → Branches
2. Click "main" branch rule
3. Enable "Require status checks to pass"
4. Select "test" as required check

### Pre-Commit Checks

Before pushing, run locally:
```bash
# Formatting
go fmt ./...

# Linting
go vet ./...

# Full test suite
go test -cover ./...

# All at once
go fmt ./... && go vet ./... && go test ./...
```

### GitHub PR Workflow

1. Create feature branch: `git checkout -b feature/name`
2. Make changes, commit with clear messages
3. Push: `git push -u origin feature/name`
4. Create PR on GitHub
5. Fix any failing tests
6. Request review
7. Merge when approved

**PR Title Format:**
```
verb: brief description (under 70 chars)

Examples:
- Add file risk analysis feature
- Fix cursor wrapping at end of list
- Optimize diff caching for large repos
```

**PR Description Format:**
```
## Summary
Brief paragraph explaining what and why

## Changes
- Bullet list of changes
- One line per significant change

## Test Plan
- How to manually test this
- Edge cases covered
- Performance impact (if any)
```

## Common Development Tasks

### Task: Add a New Test
```bash
# 1. Create test in appropriate file
# engine_analytics_test.go, engine_parsing_test.go, etc.

# 2. Run just that test
go test -run TestFeatureName -v

# 3. Check coverage for that module
go test -cover github.com/zpenka/grit
```

### Task: Find a Bug
```bash
# 1. Add regression test that fails
func TestBugName(t *testing.T) {
    // Reproduce the bug
    result := buggyFunction()
    if !correctBehavior(result) {
        t.Error("bug reproduction")
    }
}

# 2. Run the test (should fail)
go test -run TestBugName -v

# 3. Fix the bug
# (edit source code)

# 4. Verify test passes
go test -run TestBugName -v
```

### Task: Optimize Performance
```bash
# 1. Profile the code
go test -cpuprofile=cpu.prof ./...

# 2. Analyze results
go tool pprof cpu.prof
(pprof) top        # See hotspots

# 3. Implement optimization
# (edit source code)

# 4. Verify improvement
go test -cpuprofile=cpu2.prof ./...
go tool pprof -base=cpu.prof cpu2.prof  # Compare
```

### Task: Check Code Style
```bash
# Format all files
go fmt ./...

# Check for style issues
go vet ./...

# Run tests to catch errors
go test ./...
```

### Task: Review Code Changes
```bash
# See what changed
git diff

# See staged changes
git diff --staged

# See all changes in branch
git diff main..HEAD
```

### Task: Run Specific Tests
```bash
# Single test
go test -run TestFeatureName -v

# Multiple tests matching pattern
go test -run TestFeature -v

# All tests in a file
go test -run ^TestName package_test.go

# All tests excluding pattern
go test -run '^Test(?!Slow)' ./...

# With coverage
go test -run TestFeatureName -cover

# With verbose output
go test -run TestFeatureName -v
```

### Task: Debug a Failing Test
```bash
# Run with verbose output
go test -run TestName -v

# Run with print statements visible
go test -run TestName -v -count=1

# Run with race detector
go test -race -run TestName

# Run multiple times to catch flakiness
go test -run TestName -count=10
```

## Architecture Decision Log

### Decision 1: Modular Organization

**Question:** How should 113+ features be organized in code?

**Options Considered:**
- Single 50k-line file (rejected - unmaintainable)
- Feature folder per feature (rejected - too many files)
- Feature groups by category (chosen)

**Result:** 13 focused modules, each ~1000-3000 LOC. Easier navigation, clear separation of concerns, faster context switching.

### Decision 2: Consolidated Rendering Templates

**Question:** 50+ feature render functions with duplicate code - how to reduce?

**Options Considered:**
- Keep individual functions (rejected - maintenance nightmare)
- Create shared render templates (chosen)

**Result:** 3 template functions: `RenderStandardUI`, `RenderAnalysisUI`, `RenderDataGrid`. Reduced code by 500+ lines, consistent styling.

### Decision 3: LRU Cache for Diffs

**Question:** Diffs are expensive to parse - how to optimize?

**Options Considered:**
- Cache all diffs (rejected - unbounded memory)
- No caching (rejected - slow)
- LRU cache with fixed size (chosen)

**Result:** Bounded memory usage, automatic eviction of least-used diffs, huge speed improvement for browsing.

## Complete Contribution Workflow

This section walks you through the complete process of contributing a feature or fix.

### Phase 1: Planning (5-10 minutes)

1. **Identify the problem/feature**
   - Bug fix: Find failing behavior
   - Feature: Define requirements
   
2. **Choose the module**
   ```
   Analytics → engine_analytics.go
   Git ops → engine_git_ops.go
   Visualization → engine_visualization.go
   etc.
   ```

3. **Design the implementation**
   - What data structures needed?
   - What functions needed?
   - How to wire into UI?

### Phase 2: Setup (2 minutes)

```bash
# Create feature branch
git checkout main
git pull
git checkout -b feature/your-feature

# Verify tests pass
go test ./...
```

### Phase 3: TDD Implementation (30-60 minutes)

1. **Write tests first** (Red phase)
   - Create `*_test.go` file
   - Write 5-10 test cases
   - Run: `go test -run TestName -v` (should fail)

2. **Implement feature** (Green phase)
   - Add core logic to appropriate module
   - Keep implementation minimal
   - Run: `go test -run TestName -v` (should pass)

3. **Refactor** (Refactor phase)
   - Improve code clarity
   - Add comments explaining why
   - Extract helper functions
   - Run: `go test ./...` (all tests pass)

4. **Wire into UI** (Integration)
   - Update `engine_types.go` model struct
   - Add menu entry or keybinding
   - Add render function
   - Update `grit.go` View() method

### Phase 4: Testing (10-20 minutes)

```bash
# Run all tests
go test ./...

# Check coverage
go test -cover ./...

# Test in real repo
go build -o grit .
cd /path/to/git/repo
/path/to/grit

# Manual testing:
# 1. Navigate to feature (key or menu)
# 2. Verify output looks correct
# 3. Test edge cases (empty data, large data, etc.)
# 4. Check terminal doesn't crash
```

### Phase 5: Documentation (5-10 minutes)

1. **Add godoc comments**
   ```go
   // featureName does something useful.
   // It analyzes commits to find patterns.
   // Results are sorted by frequency.
   func featureName(commits []commit) []result { ... }
   ```

2. **Update README.md**
   - Add feature to list if user-facing
   - Add keybinding to reference table

3. **Update USAGE.md**
   - Add example of using the feature
   - Add to relevant workflow section

4. **Update CLAUDE.md**
   - Add if architectural impact

### Phase 6: Code Quality (5-10 minutes)

```bash
# Format code
go fmt ./...

# Check for issues
go vet ./...

# Run tests one more time
go test ./...

# Check coverage
go test -cover ./...
```

### Phase 7: Commit & Push (5 minutes)

```bash
# Stage changes
git add .

# Commit with clear message
git commit -m "Add feature: descriptive title

Details about the change:
- What problem it solves
- How it works
- Testing approach

Include any breaking changes or special notes."

# Push to GitHub
git push -u origin feature/your-feature
```

### Phase 8: Create PR (2 minutes)

1. Go to https://github.com/zpenka/grit/pulls
2. Click "New Pull Request"
3. Select: `base: main` ← `compare: feature/your-feature`
4. Fill in PR template:
   ```
   ## Summary
   One paragraph explaining what and why
   
   ## Changes
   - Bullet list
   - Of changes made
   
   ## Test Plan
   - How to test manually
   - Edge cases covered
   ```
5. Click "Create Pull Request"

### Phase 9: Review & Feedback (Variable)

- GitHub CI runs tests automatically
- Address any failed tests
- Respond to code review feedback
- Discuss approach if needed
- Update code based on feedback

### Phase 10: Merge & Celebrate (1 minute)

1. Get approval from reviewer
2. Ensure CI passes
3. Click "Squash and merge" (for clean history)
4. Delete feature branch
5. Done! 🎉

---

**Estimated time:** 1-2 hours for a typical feature

## Getting Help

**For questions, check in this order:**

1. **Look at existing code** - Similar features in same module
   ```bash
   grep -n "func render" engine_analytics.go
   ```

2. **Check tests** - Tests show how functions should work
   ```bash
   grep -n "TestAnalytics" engine_analytics_test.go
   ```

3. **Read module docs** - Architecture and patterns
   ```bash
   cat docs/modules/analytics.md
   ```

4. **Search codebase** - Find similar implementations
   ```bash
   git log -p --grep="feature name"
   ```

5. **Ask in PR discussion** - Get feedback on approach

## References

**Project Documentation:**
- [CLAUDE.md](CLAUDE.md) - High-level overview for AI assistants
- [README.md](README.md) - User-facing documentation  
- [USAGE.md](USAGE.md) - User guide with examples
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
- [docs/modules/](docs/modules/) - Architecture and module details
- [CHANGELOG.md](CHANGELOG.md) - Detailed version history

**External Resources:**
- [Go 1.21+ Documentation](https://golang.org/pkg/)
- [Bubble Tea TUI Framework](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Styling Library](https://github.com/charmbracelet/lipgloss)
- [Git Documentation](https://git-scm.com/doc)

**Guides in This File:**
- [Architecture Overview](#architecture-overview) - System design
- [Architecture Deep-Dives](#architecture-deep-dives) - Detailed module explanations
- [Testing Strategies](#testing-strategies) - TDD and test patterns
- [Adding a New Feature](#adding-a-new-feature-complete-walkthrough) - Step-by-step example
- [Common Pitfalls](#common-pitfalls--gotchas) - Mistakes to avoid
- [Complete Contribution Workflow](#complete-contribution-workflow) - End-to-end process

---

**Thank you for contributing to grit!** 🚀

Questions? Check the references above or open an issue to discuss your approach before diving into code.
