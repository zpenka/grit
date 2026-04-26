# Git Operations Module (`engine_git_ops.go`)

## Purpose

The git operations module provides high-level git workflow tools including interactive rebase, cherry-pick, reset, and commit amendment. It wraps git commands with UI enhancements like conflict detection and operation previewing to reduce the risk of unintended changes.

## Key Functions

### Interactive Rebase
- `parseRebaseSequence(targetHash string, ops []string) rebaseSequence`
  - Plan interactive rebase operations
  - Validates operation sequence before execution
  
- `previewRebaseOperations(ops []rebaseOp) rebasePreview`
  - Show predicted outcome without applying
  - Detect conflicts and issues
  - Educational: Shows what will happen

### Cherry-Pick
- `prepareCherryPick(commits []commit, targetBranch string) cherryPickPlan`
  - Plan cherry-pick with conflict detection
  - Show which commits would conflict
  
- `executeCherryPick(plan cherryPickPlan) (success bool, conflicts []string)`
  - Apply commits to current branch
  - Return conflict locations if resolution needed

### Reset Operations
- `resetCommit(m model, hash string, mode string) model`
  - Undo commits (soft/hard/mixed)
  - Soft: Keeps changes, unstages files
  - Hard: Discards all changes
  - Mixed: Unstages but keeps working tree
  
- `resetToCommit(m model, hash string) model`
  - Move HEAD to specific commit
  - Preserves working directory by default

### Commit Amendment
- `previewAmendCommit(original string, newMsg string) amendPreview`
  - Show what changes will be made
  
- `amendCommit(m model, newMsg string) model`
  - Modify most recent commit message
  - Optionally add/remove files

### Advanced Operations
- `rewordCommit(m model, hash string, newMsg string) model`
  - Change message of non-HEAD commit (interactive rebase)
  
- `squashCommits(m model, hashes []string) model`
  - Combine multiple commits into one
  - Merge messages appropriately

## Design Decisions

### Preview-Based Workflow
All operations support previewing before execution. Users see exactly what will happen before making changes permanent.

### Conflict Detection
Built-in conflict detection for operations that might have conflicts (cherry-pick, rebase). User is warned before executing.

### Safety by Default
- Soft resets by default (preserve work)
- Conflicts must be manually resolved
- No automatic force-push

### Immutable Operations
All functions take and return model. Operation history is tracked for potential undo.

## Dependencies

- **Internal**: git_utils for command execution
- **Standard library**: `fmt`, `strings`, `errors`
- **External**: git command-line tool

## Testing

Tests cover:
- Rebase sequence validation
- Conflict detection accuracy
- Cherry-pick with and without conflicts
- Reset mode behavior (soft/hard/mixed)
- Commit amendment message formatting

## Examples

```go
// Preview interactive rebase
ops := []string{"pick", "reword", "squash"}
preview := previewRebaseOperations(parseRebaseSequence(baseHash, ops))
if len(preview.conflicts) == 0 {
    // Safe to execute
    executeRebase(preview)
}

// Cherry-pick with conflict detection
plan := prepareCherryPick([]commit{commit1, commit2}, "feature-branch")
if len(plan.conflicts) > 0 {
    fmt.Println("Conflicts detected in:", plan.conflicts)
    // User manually resolves or aborts
}

// Soft reset (preserve work)
m = resetCommit(m, "HEAD~3", "soft")

// Amend most recent commit
m = amendCommit(m, "Fix: Better implementation")
```

## Git Safety

All operations:
- Use `git -C <repo>` for working directory safety
- Validate commits before operations
- Support undo via reflog
- Never force-push unless explicitly requested

## Integration Points

- Model tracks current operation state
- Rendering shows preview/conflict state
- Navigation guides user through resolution
- Keybinding handler initiates operations

## Conflict Resolution

When conflicts arise:
1. Operation pauses (not applied)
2. UI shows conflict locations
3. User navigates to conflicts
4. Manual resolution or abort
5. Continue or retry operation

## Performance Considerations

- Previewing is fast (no actual git execution)
- Conflict detection scans changed files
- Large rebases can take time (but progress shown)

## Future Extensions

- Interactive conflict resolution UI
- Undo/redo for git operations
- Fixup/autosquash sequences
- Advanced rebase patterns (autostash, rerere)
- Commit signing enforcement
