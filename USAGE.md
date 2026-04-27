# grit - Comprehensive Usage Guide

This guide covers everything from basic usage to advanced workflows for exploring git history with grit.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Core Features Walkthrough](#core-features-walkthrough)
3. [Advanced Workflows](#advanced-workflows)
4. [Troubleshooting](#troubleshooting)
5. [Performance Tips](#performance-tips)
6. [Keyboard Shortcuts Reference](#keyboard-shortcuts-reference)

## Getting Started

### Launch grit

```bash
# From any git repository
./grit

# Or install and run from anywhere
go install github.com/zpenka/grit@latest
grit
```

### Understanding the Interface

grit displays commits in a **two-panel layout**:
- **Left panel**: List of commits (navigable with `j`/`k`)
- **Right panel**: Diff details, features, or analysis results

The **footer bar** shows available keys and current state. Press `?` to see the full help overlay.

## Core Features Walkthrough

### Basic Navigation

#### Browsing Commits

```
j / k           - Move down / up in commit list
5j / 5k         - Jump 5 commits (works with any number)
g               - Go to first commit (oldest)
G               - Go to last commit (newest)
```

**Example workflow:**
1. Press `g` to jump to the oldest commit
2. Press `5j` repeatedly to scan through history in chunks
3. Press `k` to move back one at a time when you spot something interesting

#### Switching Panels

```
l / h           - Switch focus to diff / commit list
Tab             - Toggle between panels (faster)
Enter           - Expand/collapse current panel
```

**Example:** You're viewing a diff and want to see the next commit. Press `h` to focus the commit list, then `j` to move down.

### Viewing Commit Details

#### Diff Panel (`l` to focus)

- **-** red lines: Code removed
- **+** green lines: Code added
- **@@** symbols: Hunk headers showing line numbers
- Diff automatically appears for the selected commit

**Navigation in diff:**
```
j / k           - Scroll diff up / down
Page Up/Down    - Jump to previous/next hunk
Ctrl+Home/End   - Jump to start/end of diff
```

#### File List (`f` key)

Press `f` to toggle a file list showing all files changed in the current commit:

```
j / k           - Navigate files
Enter           - View that file's diff
l               - Jump to file in full diff
```

**Example:** Looking at a large commit with 50+ files? Use `f` to quickly find the specific file you want to examine.

### Searching and Filtering

#### Quick Search (`/` key)

```
/               - Start search mode
Type query      - Search commits by:
                  • Commit message/subject
                  • Author name
                  • Partial hash
Ctrl+/          - Clear current search
```

**Example searches:**
```
/fix bug        - Find commits with "fix bug" in subject
/Alice          - Find all commits by Alice
/abc123         - Find commits with hash containing "abc123"
```

Search is **case-insensitive** and **regex-enabled** for advanced patterns.

#### Author Filter (`a` submenu)

Shows analytics menu. Navigate to filtering options:
```
a               - Open Analytics menu
j/k             - Select "Author Filter" option
Enter           - Activate
```

Filter commits to show only those from specific authors, making large team histories navigable.

### Branches and Tags

#### Branch Picker (`b` key)

```
b               - Show branch selector
j/k             - Navigate branches
Enter           - Switch to selected branch
Esc             - Close menu
```

**Tip:** Use branch switching to explore parallel development on different branches without leaving grit.

### Code Analysis

#### Analytics Menu (`a` key)

Open the analytics submenu to access code-focused insights:

```
a               - Open Analytics menu
```

Available options:
- **Code Ownership** - Who owns each file/component?
- **Hotspots** - Which files change most frequently?
- **Bisect** - Binary search to find when bugs were introduced
- **Complexity** - How complex is the code getting?
- **Statistics** - Stats for the current commit
- **Heatmap** - Activity visualization over time
- **Code Linting** - Lint violations and patterns

**Example: Finding a hotspot**
1. Press `a` → Select "Hotspots"
2. Review the most-changed files
3. Focus on the 2-3 files with highest change frequency
4. These are good candidates for refactoring

#### Blame Mode (`B` key)

Shows who last modified each line of the current file:

```
B               - Toggle blame overlay
```

Displays author and commit hash for each line. Useful for:
- Understanding why code was written a certain way
- Finding the context of a bug
- Learning from senior developers' code patterns

### Visualization Menu (`v` key)

Access visual representations of commit history:

```
v               - Open Visualization menu
```

Available options:
- **Contributor Flamegraph** - Visual hierarchy of commits per author
- **Timeline Slider** - Commit activity over time
- **Tree View** - Hierarchical view of commits and branches
- **Author Comparison** - Side-by-side stats for two developers
- **File Heatmap** - Visualization of file change frequency

**Example: Using Timeline**
1. Press `v` → Select "Timeline Slider"
2. See commit distribution across dates
3. Identify periods of high activity vs. calm periods
4. Use this to understand project velocity

### Team & AI Menu (`t` key)

Collaboration and intelligence features:

```
t               - Open Team/AI menu
```

Available options:
- **Team Statistics** - Commits, lines changed, code ownership per person
- **Reviewer Suggestions** - Who should review changes?
- **Velocity** - Team commit velocity and trends
- **Classification** - Auto-classify commits (feature/fix/refactor)
- **Security** - Secret scanning and GPG verification
- **Compliance** - Semantic versioning, changelog generation

**Example: Reviewer suggestions**
1. Press `t` → Select "Reviewer Suggestions"
2. For any file, grit suggests reviewers based on code history
3. Helps find the best person to review your PR

### Integration Menu (`i` key)

Connect with external tools:

```
i               - Open Integration menu
```

Available options:
- **GitHub PR Links** - Link commits to GitHub PRs
- **Jira Ticket Links** - Link to Jira issues
- **Issue References** - Find referenced GitHub issues
- **Export** - Export commits/diffs to multiple formats
- **Conflict UI** - Handle merge conflicts

**Example: Exporting for a report**
1. Press `i` → Select "Export"
2. Choose format (markdown, CSV, JSON, XML, patch series)
3. Enter filename
4. Commits are exported with all metadata

### Git Operations Menu (`g` key)

Advanced git workflows:

```
g               - Open Git Operations menu
r (or g→r)      - Interactive rebase mode
c (or g→c)      - Cherry-pick multiple commits
x (or g→x)      - Reset to commit (soft/mixed/hard)
```

**Interactive Rebase Example:**
1. Press `r` to enter rebase mode
2. Preview which commits will be rebased
3. Mark commits for squash, reword, drop, etc.
4. Execute the rebase

**Cherry-Pick Example:**
1. Press `c` to toggle cherry-pick mode
2. Use `j`/`k` to navigate and Space/Enter to mark commits
3. Marked commits are highlighted
4. Execute cherry-pick operation on all marked commits

## Advanced Workflows

### Workflow 1: Finding When a Bug Was Introduced

**Scenario:** Code in production is broken, but you don't know when it broke.

**Solution: Use bisect**

1. Press `a` → Select "Bisect" mode
2. Mark commit as "good" (before bug) or "bad" (after bug)
3. grit binary searches to find the exact culprit commit
4. When found, examine the commit to understand what broke

This reduces searching through 1000s of commits to ~10 comparisons.

### Workflow 2: Understanding Code Ownership

**Scenario:** You need to understand who maintains different parts of the codebase.

1. Press `a` → Select "Code Ownership"
2. grit analyzes commit history to assign ownership scores
3. See breakdown per file/module
4. Press `b` → Select author to drill down into their contributions

### Workflow 3: Preparing a Release

**Scenario:** Need to generate release notes and identify breaking changes.

1. Press `t` → Select "Changelog Generation"
   - grit auto-generates changelog with categorized commits
2. Copy markdown output
3. Press `i` → Select "Export"
   - Export all commits since last tag
   - Share with team for review

### Workflow 4: Performance Analysis

**Scenario:** Code is getting slower. Find what changed.

1. Press `a` → Select "Complexity Analysis"
   - See how code complexity evolved over time
2. Jump to commits where complexity spiked
3. Examine diffs to understand what added complexity
4. Use `B` (blame) to see who made the change
5. Reach out to that person for context

### Workflow 5: Team Onboarding

**Scenario:** New team member needs to understand the codebase structure.

1. Press `v` → Select "Contributor Flamegraph"
   - Visual hierarchy of who works on what
2. Press `v` → Select "File Heatmap"
   - Shows which files are most active (and thus most important)
3. Press `a` → Select "Code Ownership"
   - Shows who's expert in each area
4. Have new person start with high-ownership areas

## Troubleshooting

### Issue: grit is slow with large repositories

**Symptoms:** Scrolling is sluggish, diffs take time to load.

**Solutions:**
1. **Limit history:** Start with a branch instead of all history
   - Press `b` → Select feature branch instead of main
2. **Use filters:** Narrow down commits with `/` search
   - This reduces what grit needs to process
3. **Increase cache:** Set environment variable (if supported)
   - `GRIT_CACHE_SIZE=5000 ./grit`

### Issue: Diffs don't appear or are cut off

**Symptoms:** Right panel is blank or shows partial diff.

**Solutions:**
1. **Verify terminal height:** Make sure your terminal window is tall enough
   - grit needs at least 10 lines height
2. **Try expanding panel:** Press `Enter` to expand the diff panel
3. **Switch focus:** Press `l` to ensure diff panel is focused
4. **Refresh:** Press `Ctrl+r` to force a refresh

### Issue: Search results seem inconsistent

**Symptoms:** Same search returns different results.

**Solutions:**
1. **Clear cache:** Press `Ctrl+/` to clear search and reset cache
2. **Use exact match:** Wrap search in quotes for literal matching
   - `/exact "phrase only"`
3. **Check case sensitivity:** Search is case-insensitive by default
   - For case-sensitive, use regex flag: `/(?-i)CaseSensitive`

### Issue: Cherry-pick or rebase fails

**Symptoms:** Operation shows error message.

**Solutions:**
1. **Check for conflicts:** Conflicts prevent operation
   - Press `i` → Select "Conflict UI" to resolve
2. **Verify clean working directory:** Ensure no uncommitted changes
   - `git status` before using git operations in grit
3. **Review target commits:** Make sure selected commits are valid
   - Some commits may not be cherry-pickable to current branch

### Issue: Can't find commits with specific properties

**Symptoms:** Search for author name returns no results.

**Solutions:**
1. **Check exact spelling:** Author names are case-insensitive but must match
   - If author is "Alice Johnson", search `/alice` works but `/alic` doesn't
2. **Use broader search:** Start with partial match
   - `/alice` then use j/k to navigate results
3. **Try different fields:** Search also matches subject
   - If you can't find by author, search by commit message words

## Performance Tips

### Tip 1: Use Filtering Wisely

Instead of browsing all 10,000 commits:

```
/refactor       - Find all "refactor" commits
/Alice          - Find all Alice's commits
/file.go        - Find commits touching specific file
```

Combine these approaches to narrow down to what matters.

### Tip 2: Branch Switching

Large monorepos slow down across all branches. Focus on one:

```
b               - Switch to relevant feature branch
                  (fewer commits = faster navigation)
```

### Tip 3: Lazy Load Features

Don't open features you're not using:
- Each feature (flamegraph, heatmap, etc.) processes all commits
- Open only the one you need right now
- Close with `Esc` when done to free resources

### Tip 4: Keyboard Over Mouse

Terminal navigation is faster than mouse:
- `5j` + `5k` (jumping chunks) is faster than scrolling
- `g` / `G` (first/last) is instant
- Muscle memory = speed

### Tip 5: Bookmark Important Commits

Mark commits you want to return to:
```
m               - Mark/bookmark current commit
j to jump between bookmarks using marked navigation
```

This is faster than searching multiple times.

## Keyboard Shortcuts Reference

### Navigation
| Key | Action |
|-----|--------|
| `j` / `k` | Move down / up |
| `5j` / `5k` | Jump 5 commits |
| `g` | Go to first commit |
| `G` | Go to last commit |
| `l` / `h` | Focus diff / commit list |
| `Tab` | Toggle focus |

### Viewing
| Key | Action |
|-----|--------|
| `f` | Toggle file list |
| `b` | Branch picker |
| `B` | Blame for current file |
| `/` | Search commits |
| `Ctrl+/` | Clear search |

### Submenus
| Key | Action |
|-----|--------|
| `a` | Analytics (ownership, hotspots, bisect, linting, complexity, stats, heatmap) |
| `v` | Visualizations (flamegraph, timeline, tree, comparison, heatmap) |
| `t` | Team/AI (stats, suggestions, velocity, classification, security, changelog) |
| `i` | Integration (GitHub, Jira, export, conflict UI) |
| `g` | Git Operations (rebase, cherry-pick, reset) |

### Git Operations
| Key | Action |
|-----|--------|
| `r` | Interactive rebase preview |
| `c` | Cherry-pick mode (toggle commits) |
| `x` | Cycle reset mode (soft → mixed → hard) |
| `y` | Copy current commit hash |
| `e` | Open in editor |

### UI Control
| Key | Action |
|-----|--------|
| `?` | Show help overlay |
| `Esc` | Close any panel/menu |
| `q` | Quit grit |

## Quick Examples by Use Case

### I want to understand why a file changed

```
/filename.go            # Search for commits changing this file
B                       # Blame to see who changed what lines
Select commit           # Look at the full diff context
```

### I want to know who's the expert on this code

```
a                       # Analytics menu
Select "Code Ownership" # See ownership breakdown
Find expert's name      # They'll be top contributor
Press / and search      # Review their key commits
```

### I want to prepare a release

```
t                       # Team menu
Select "Changelog"      # Auto-generate from commits
i                       # Integration menu
Select "Export"         # Export changelog + commits
Markdown format         # For release notes
```

### I want to find when something broke

```
a                       # Analytics menu
Select "Bisect"         # Binary search commits
Mark good/bad           # Tell grit which commits work/fail
Follow hints            # Bisect guides you to the culprit
```

### I want to move commits to another branch

```
c                       # Cherry-pick mode
Navigate and select     # Use Space/Enter to mark commits
j/k to navigate         # Mark all you want
Execute                 # Cherry-pick all marked commits
```

## Tips for Different Roles

### For Developers
- Use `/` search to find related commits quickly
- Use `B` blame to understand code changes
- Use `r` rebase to organize your commits before pushing
- Use `e` to edit commits directly in your editor

### For Code Reviewers
- Use `a` → "Hotspots" to find risky changes first
- Use `a` → "Code Ownership" to understand context
- Use `t` → "Reviewer Suggestions" to help other reviewers
- Use `v` → "Author Comparison" to understand developer patterns

### For Tech Leads
- Use `t` → "Team Statistics" for velocity tracking
- Use `v` → "Contributor Flamegraph" to see team distribution
- Use `a` → "Complexity Analysis" to identify refactoring targets
- Use `t` → "Security" for compliance tracking

### For DevOps/SRE
- Use `a` → "Heatmap" to spot problem periods
- Use `a` → "Bisect" to find when things broke
- Use `i` → "Export" to generate audit trails
- Use `t` → "Changelog" for automated release notes

---

**Need help?** Press `?` in grit to see keybindings. Check [README.md](README.md) for installation, [DEVELOPER.md](DEVELOPER.md) for hacking on grit itself.
