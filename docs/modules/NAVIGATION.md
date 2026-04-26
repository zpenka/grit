# Navigation Module (`engine_navigation.go`)

## Purpose

The navigation module manages cursor position, panel switching, and scrolling within the UI. It maintains and updates all navigation state in response to user input. All functions are pure—they take a model and return an updated model.

## Key Functions

### `moveCursorUp(m model) model` / `moveCursorDown(m model) model`
Navigate up/down through the commit list
- Respects list boundaries (prevents over-scrolling)
- Resets diff offset when changing commits
- Works with filtered commit lists

### `switchPanel(m model, panel string) model`
Change the active UI panel
- Supports: "commits", "diff", "files", "feature"
- Updates focus state for keyboard routing
- Resets scroll position for new panel

### `scrollPanelUp(m model) model` / `scrollPanelDown(m model) model`
Scroll content within the active panel
- Tracks offset for both diff and file trees
- Respects content boundaries
- Smooth scrolling without jumping

### `setBookmark(m model, name string) model` / `goToBookmark(m model, name string) model`
Save and restore cursor positions
- Useful for navigating complex histories
- Up to 10 bookmarks supported
- Includes commit hash and scroll state

## Design Decisions

### Immutable Navigation
All functions take and return model rather than mutating in place. This makes navigation testable and composable with Bubble Tea's update cycle.

### Cursor + Offset Pattern
Maintains both absolute cursor position in full list and offset in filtered list. Prevents UI confusion when filters change.

### Panel-Specific State
Each panel (diff, file tree) has its own scroll offset. Switching panels preserves view position.

## Dependencies

- No external dependencies
- Uses standard model struct

## Testing

Tests cover:
- Boundary conditions (first/last commit)
- Panel switching and focus management
- Scrolling with various content sizes
- Bookmark save/restore cycles

## Examples

```go
// Navigate down through commits
m = moveCursorDown(m)

// Switch to diff panel
m = switchPanel(m, "diff")

// Scroll down in current panel
m = scrollPanelDown(m)

// Save current position as bookmark
m = setBookmark(m, "interesting-commit")

// Later, jump back to bookmark
m = goToBookmark(m, "interesting-commit")
```

## Performance Considerations

- All operations are O(1) in commit count
- Cursor position is lightweight to track
- Bookmarks stored in fixed array (no allocation)

## Integration with Bubble Tea

Navigation functions are called from Bubble Tea's `Update()` method in response to key presses. The model is threaded through each update step, maintaining consistent state.

## Future Extensions

- Incremental scrolling history (undo/redo navigation)
- Smart scroll to keep cursor centered
- Keyboard shortcuts for bookmarks
