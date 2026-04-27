# Render Consolidation Module (`engine_render_consolidation.go`)

## Purpose

The render consolidation module provides unified UI templates and rendering patterns shared across all features. It eliminates duplication by centralizing layout logic while allowing features to customize content and styling. This ensures consistent styling across all 30+ features and makes global UI changes easy to implement.

## Key Functions

### Standard UI Template
- `RenderStandardUI(config RenderConfig) string`
  - Generic two-column layout
  - Left: Item list, Right: Status or detail
  - Used by: filters, searches, bookmarks, most lists
  
**Config Options:**
```go
type RenderConfig struct {
    Title       string              // Panel title
    Items       []string            // List items
    Selected    int                 // Currently selected item
    HasStatus   bool                // Show status column?
    StatusMap   map[string]string   // Item -> status mapping
    Details     string              // Right panel content
}
```

### Analysis UI Template
- `RenderAnalysisUI(title string, data map[string]interface{}) string`
  - Data grid layout with headers
  - Used by: analytics results, team stats, metrics
  - Supports: multi-column display, sorting hints
  
**Typical Data Map:**
```go
data := map[string]interface{}{
    "authors":     []string{"alice", "bob"},
    "commits":     []int{150, 75},
    "percentages": []float64{66.7, 33.3},
}
```

### Data Grid Template
- `RenderDataGrid(headers []string, rows [][]string) string`
  - Tabular data display
  - Column alignment and width management
  - Used by: detailed reports, exports, queries
  
**Example:**
```
Author      Commits  Avg Size  Last Activity
alice       150      250 lines 2 hours ago
bob         75       180 lines 1 day ago
eve         45       320 lines 3 days ago
```

### Panel Rendering
- `RenderPanel(title, content string, width, height int) string`
  - Generic panel with border and title
  - Used as building block for larger layouts
  - Supports: truncation, scrolling hints

### Supporting Functions
- `formatTable(data [][]string, colWidths []int) string`
  - Align columns in table
  - Handle long content with ellipsis
  
- `colorCode(text string, color string) string`
  - Apply color to text (added/removed/context)
  - Supports: green, red, gray, blue, yellow
  
- `truncate(text string, maxLen int) string`
  - Shorten text with ellipsis
  - Respects word boundaries

## Design Patterns

### Configuration-Based Rendering
Templates accept config objects rather than hard-coded values. Enables flexibility without code duplication.

### Semantic Coloring
- Green: Added content, success, positive metrics
- Red: Removed content, errors, negative metrics
- Yellow: Warnings, current selection, highlights
- Blue: Metadata, headings, structural info
- Gray: Context, secondary info, timestamps

### Responsive Layout
Templates adapt to terminal width:
- Minimum: 80 columns (required)
- Optimal: 120+ columns (preferred)
- Scaling: Column widths adjust proportionally

## Dependencies

- **Lipgloss**: Terminal styling, colors, borders
- **Standard library**: `fmt`, `strings`
- No external dependencies beyond styling

## Testing

Tests cover:
- Config rendering with various options
- Column alignment accuracy
- Color application
- Truncation behavior
- Width constraints

## Examples

```go
// Render list with status
config := RenderConfig{
    Title:     "Bookmarks",
    Items:     []string{"v1.0 Release", "Bug Fix", "Feature Branch"},
    Selected:  0,
    HasStatus: true,
    StatusMap: map[string]string{
        "v1.0 Release": "saved",
        "Bug Fix": "visiting",
        "Feature Branch": "saved",
    },
}
ui := RenderStandardUI(config)

// Render analysis results
data := map[string]interface{}{
    "Top Contributors": []string{"Alice (60%)", "Bob (30%)", "Eve (10%)"},
    "Activity": "High",
    "Stability": "Good",
}
ui := RenderAnalysisUI("Code Ownership", data)

// Render data grid
headers := []string{"File", "Author", "Changes", "Last Edit"}
rows := [][]string{
    {"main.go", "alice", "150", "2 hours ago"},
    {"auth.go", "bob", "75", "1 day ago"},
}
ui := RenderDataGrid(headers, rows)
```

## Feature Integration

All 30+ feature renderers follow this pattern:

```go
// Feature-specific render function
func renderFeatureUI(m model, width int) string {
    config := RenderConfig{
        Title: "Feature Name",
        Items: extractItems(m.featureData),
        // ... populate config
    }
    return RenderStandardUI(config)
}
```

## Styling Guidelines

### Consistent Spacing
- Panel borders: 1 character
- Padding: 1 space inside borders
- Column separation: 2 spaces
- Row spacing: Normal (no blank rows)

### Text Formatting
- Titles: Bold or uppercase
- Status: Color-coded
- Numbers: Right-aligned
- Text: Left-aligned

### Color Usage
Limited palette for terminal compatibility:
- ANSI 16 colors (basic terminal support)
- ANSI 256 colors (enhanced terminal support)
- RGB colors (24-bit truecolor terminals)

## Performance Considerations

Template rendering is O(n) in content size. Optimization strategies:

- Cache formatted strings where possible
- Avoid re-rendering unchanged panels
- Lazy format expensive columns
- Defer color application to final render

## Accessibility

- No color-only indicators (text + color)
- Sufficient contrast ratios
- Monospace font assumptions (terminal)
- Text-based formatting (works in screen readers)

## Extensibility

Adding new template:
1. Define config struct for parameters
2. Implement render function using Lipgloss
3. Add to consolidation module
4. Document usage examples
5. Add tests

## Troubleshooting

### Too Wide
- Terminal too narrow for content
- Solution: Truncate columns, use scrolling

### Colors Not Showing
- Terminal doesn't support colors
- Solution: App auto-detects, falls back to plain text

### Alignment Issues
- Variable-width fonts (not monospace)
- Solution: Display warning or adapt widths

## Future Extensions

- Dark/light theme support
- Colorblind-friendly modes
- Custom user themes
- Unicode symbol support
- Side-by-side panel layouts
- Responsive grid layouts
