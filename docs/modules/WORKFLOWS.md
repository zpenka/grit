# Workflows Module (`engine_workflows.go`)

## Purpose

The workflows module provides advanced git workflows beyond basic operations: worktree management, stash handling, tag operations, and reflog browsing. These are multi-step operations that require careful state tracking and user guidance.

## Key Functions

### Worktree Management
- `loadWorktrees(output string) []worktreeInfo`
  - Parse `git worktree list` output
  - Show all checked-out branches and work directories
  
- `createWorktree(m model, path string, branch string) model`
  - Create new worktree for parallel work
  - Doesn't affect main working directory
  
- `removeWorktree(m model, path string) model`
  - Clean up completed worktree
  - Requires clean state

### Stash Management
- `saveToStash(m model, description string) model`
  - Temporarily save work in progress
  - Full snapshot of working directory
  
- `applyStash(m model, stashID string) model`
  - Restore stashed changes
  - Handles conflicts if working directory changed
  
- `listStashes() []stashEntry`
  - Show all saved stashes with timestamps
  - Browse stash history

### Tag Operations
- `createTag(m model, name string, message string) model`
  - Create lightweight or annotated tags
  - Tags mark release versions
  
- `deleteTag(m model, name string) model`
  - Remove tag locally and optionally from remote
  
- `listTags() []tagInfo`
  - Show all tags with dates and commit references
  - Filter by pattern

### Reflog Browsing
- `browseReflog(m model) model`
  - Navigate git reflog (command history)
  - Find "lost" commits after resets
  
- `recoverCommit(m model, reflogEntry string) model`
  - Check out or cherry-pick from reflog
  - Safety net for accidental resets

## Design Decisions

### State-Preserving Operations
Worktrees and stashes preserve working state. Operations can be undone or continued.

### Multi-Step Guidance
UI guides users through complex workflows. Clear status messages and confirmation prompts.

### Recovery Tools
Reflog browsing and recovery tools help undo accidental destructive operations.

### Independent Operations
Each workflow (worktrees, stash, tags) is independent. Can use combinations.

## Dependencies

- **Internal**: git_utils for command execution
- **Standard library**: `fmt`, `strings`, `time`
- **External**: git command-line tool

## Testing

Tests cover:
- Worktree creation, listing, removal
- Stash save, apply, with conflict detection
- Tag creation (lightweight and annotated)
- Reflog parsing and recovery
- Edge cases: no stashes, no worktrees, empty reflog

## Examples

```go
// Create worktree for feature branch
m = createWorktree(m, "/tmp/grit-feature", "feature/new-analytics")
// Now working in /tmp/grit-feature while main repo unchanged

// Save work in progress
m = saveToStash(m, "WIP: Half-done refactoring")

// Later, restore stashed work
m = applyStash(m, "stash@{0}")

// Create release tag
m = createTag(m, "v1.2.0", "Release version 1.2.0")

// Find lost commit after accidental reset
reflog := browseReflog(m)
// Search for commit hash, then:
m = recoverCommit(m, "HEAD@{5}")
```

## Worktree Benefits

- Parallel work on multiple branches
- No stash/checkout juggling
- Separate build directories
- Team collaboration (different people in different worktrees)

## Stash Use Cases

- Switch branches mid-work without committing
- Experiment with changes and discard
- Save multiple work-in-progress snapshots
- Clean working directory for operations

## Tag Use Cases

- Release version marking
- Milestone tracking
- Reproducible builds (tag specific commits)
- Branch point documentation

## Recovery Operations

Reflog shows all HEAD movements:
```
abc123 HEAD@{0}: checkout: moving to main
def456 HEAD@{1}: commit: Fix bug
ghi789 HEAD@{2}: reset: moving to HEAD~1
```

Users can recover any previous state without loss.

## Integration Points

- Navigation: Browse worktrees, stashes, tags as lists
- Rendering: Show status in feature panel
- Keybinding handler: Initiate workflow operations

## Performance Considerations

- Worktree operations: Quick (just git commands)
- Stash save/apply: Depends on working tree size
- Reflog browsing: Fast (local operation)

## Future Extensions

- Interactive stash management UI
- Worktree templates for common branches
- Scheduled stash cleanup (old stashes)
- Tag push confirmation
- Bulk operations (delete multiple tags)
- Automatic reflog pruning
