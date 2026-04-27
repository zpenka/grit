# Main Rendering Module (`engine_rendering.go`)

## Purpose

The main rendering module is the heart of the UI, responsible for converting the model state into terminal output. It renders the four-panel interface: commit list, diff view, file tree, and feature panels. At 66KB+, this is the largest module and coordinates all visual output.

## Key Rendering Groups

### Commit List Panel
Displays filtered commits with metadata
- **Functions**: `renderCommitList()`, `renderCommitRow()`, `renderCommitBadge()`
- Shows: Subject, author, time ago, diff stats
- Styling: Color-coded by type (feature/fix/refactor), highlighting current commit
- Cached for performance

### Diff Panel
Shows unified diff with syntax highlighting and line numbers
- **Functions**: `renderDiffPanel()`, `renderDiffLine()`, `highlightDiffSyntax()`
- Features: Line numbers, context lines highlighted differently
- Caching: Diff parsing cached in `dcache`, rendering reused across frames
- Performance: Only visible portion rendered

### File Tree Panel
Shows changed files organized by directory structure
- **Functions**: `renderFileTree()`, `renderFileItem()`, `renderFileStats()`
- Shows: File path, change type, hunks, additions/deletions
- Navigation: Expandable directories
- Context: Shows blame info on hover

### Feature Panels
Results from analytics, bisect, team stats, visualizations
- **Functions**: One per feature (30+ total `render*UI()` functions)
- Uses: Consolidated templates from `engine_render_consolidation.go`
- Updates: Feature-specific rendering coordinated by keybinding handler

## Key Functions

### Core Rendering
- `renderFrame()`: Main entry point, assembles all four panels
- `layoutPanels()`: Divides terminal space among panels
- `applyTheme()`: Colors, borders, styling with Lipgloss

### Panel-Specific
- `renderCommitList()`: Commit browser with filtering
- `renderDiffPanel()`: Unified diff display
- `renderFileTree()`: File/directory tree
- `renderFeaturePanel()`: Current feature results

## Design Patterns

### Consolidated Templates
Uses `RenderStandardUI()`, `RenderAnalysisUI()`, `RenderDataGrid()` from engine_render_consolidation.go to eliminate duplication. All 30+ features reuse the same templates.

### Lazy Rendering
Only renders visible lines. Large diffs render only visible portion. Improves frame rate with large content.

### Caching
Diff parsing results cached in `dcache`. Rendered strings memoized to avoid redundant formatting.

## Dependencies

- **Lipgloss**: Terminal styling, colors, borders
- **Standard library**: `fmt`, `strings`, `regexp`
- **Internal**: Consolidation templates, caching module

## Performance Considerations

Large diffs (5000+ lines):
- Only visible ~50 lines rendered per frame
- Diff parsing cached to avoid re-parsing on scroll
- Line highlighting deferred until visible

Commit lists (10k+ commits):
- Only ~30 commits visible at once
- Badge computation deferred
- Author highlighting uses pre-computed colors

## Testing

Tests cover:
- Panel sizing and layout
- Line wrapping and truncation
- Syntax highlighting correctness
- Cache hit rates for typical usage
- Performance with large diffs/commits

## Examples

```go
// Main rendering entry point (called from Bubble Tea)
ui := renderFrame(m, terminalWidth, terminalHeight)

// Render single diff line with colors
line := renderDiffLine(diffLine{
    kind: "added",
    text: "+new feature",
    num:  42,
})

// Apply theme styling
styled := applyTheme(ui, "dark-mode")
```

## Integration with Bubble Tea

Called from Bubble Tea's `View()` method:
1. Model state is read (cursor, filters, feature toggles)
2. Panels are rendered based on state
3. Terminal output is returned to Bubble Tea
4. Terminal displays the result

## Color Scheme

- **Added lines**: Green (`+`)
- **Removed lines**: Red (`-`)
- **Context lines**: Gray
- **Hunk headers**: Blue
- **Current commit**: Yellow highlight
- **Authors**: Assigned colors for consistency

## Future Extensions

- Configurable themes (light/dark/colorblind)
- Syntax highlighting for more languages
- Diff stats visualization (sparklines)
- Blame information overlays
- Side-by-side diff view
