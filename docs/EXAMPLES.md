# Example Workflows

Real-world workflows showing how to use grit for common git tasks.

## Basic Workflows

### Workflow 1: Code Ownership Analysis

**Goal**: Understand who changed what parts of the codebase

```bash
grit                          # Start grit
a                             # Open code ownership analysis
```

What you'll see:
- List of files and authors
- Commit count per author
- Percentage of codebase owned by each person

**Use case**: Onboarding new developers, finding code experts

---

### Workflow 2: Find Code Hotspots

**Goal**: Identify frequently modified files (risk zones)

```bash
grit                          # Start grit
s                             # Open statistics & hotspots
```

What you'll see:
- Files sorted by change frequency
- Commit count per file
- Last modified dates

**Use case**: Focus code review efforts, plan refactoring

---

### Workflow 3: Bisect a Regression

**Goal**: Find the commit that introduced a bug

```bash
grit                          # Start grit
b                             # Start bisect mode
[navigate to commit]
y                             # Mark as good (no bug)
[navigate to another commit]
n                             # Mark as bad (has bug)
[bisect continues...]         # Eventually finds culprit
```

**Use case**: Identify bugs, understand root causes

---

## Search Workflows

### Workflow 4: Search by Keyword

**Goal**: Find all commits related to authentication

```bash
grit                          # Start grit
/auth                         # Search for "auth"
n                             # Find next match
N                             # Find previous match
```

**Use case**: Track feature history, find related changes

---

### Workflow 5: Find Commits by Author

**Goal**: See all work from a specific developer

```bash
grit                          # Start grit
?                             # Open filter options
[select: Author filter]
[type: john.doe@example.com]
Enter                         # Apply filter
```

**Use case**: Code review, team analytics, knowledge transfer

---

### Workflow 6: Analyze Changes in Date Range

**Goal**: See what changed last week

```bash
grit                          # Start grit
?                             # Open filter options
[select: Date range]
[start: 7 days ago]
[end: today]
Enter                         # Apply filter
```

**Use case**: Sprint reviews, weekly summaries

---

## Analysis Workflows

### Workflow 7: Team Velocity Tracking

**Goal**: Measure how fast your team delivers

```bash
grit                          # Start grit
t                             # Open team velocity metrics
```

What you'll see:
- Commits per author
- Average commit size
- Activity patterns

**Use case**: Sprint planning, capacity planning

---

### Workflow 8: Visualize Commit Distribution

**Goal**: See activity patterns over time

```bash
grit                          # Start grit
v                             # Open visualizations
[select: Activity heatmap]
```

What you'll see:
- Commits per day/week/month
- Team timezone patterns
- Peak activity times

**Use case**: Understand team work patterns, identify bottlenecks

---

### Workflow 9: Export Commit Data

**Goal**: Analyze commits in external tools (Excel, Python, etc)

```bash
grit                          # Start grit
[apply any filters you want]
e                             # Export
[select format: CSV/JSON/XML]
[data saved to file]
```

Then analyze in Excel, Python, R, etc.

**Use case**: Custom analysis, business intelligence

---

## Integration Workflows

### Workflow 10: Link GitHub Issues

**Goal**: See which commits are related to specific issues

```bash
grit                          # Start grit
[navigate to commit]
i                             # Show GitHub issues
```

You'll see:
- Linked GitHub PR numbers
- Issue status
- Related discussions

**Use case**: Issue tracking, release notes generation

---

### Workflow 11: Team Metrics Report

**Goal**: Create a team report for management

```bash
grit                          # Start grit
t                             # Team velocity
[note: commits, authors, trends]
e                             # Export
[format: JSON]
```

Then create report with the exported data.

**Use case**: Management reporting, team performance reviews

---

## Advanced Workflows

### Workflow 12: Investigate Breaking Change

**Goal**: Find when a feature stopped working

```bash
grit                          # Start grit
/breaking                     # Search for breaking changes
a                             # Analyze code ownership
[see who made changes]
b                             # Bisect for exact commit
```

**Use case**: Incident investigation, root cause analysis

---

### Workflow 13: Code Review Preparation

**Goal**: Prepare for reviewing commits from a colleague

```bash
grit                          # Start grit
/author:colleague             # Filter by author
s                             # See statistics
[review each file]
[make notes on findings]
```

**Use case**: Efficient code reviews, knowledge sharing

---

### Workflow 14: Release Preparation

**Goal**: Prepare release notes

```bash
grit                          # Start grit
[start from last release tag]
[navigate to head]
a                             # Analyze code ownership
s                             # Hotspots and stats
e                             # Export for changelog
```

Then combine into release notes.

**Use case**: Release management, changelog generation

---

## Tips for Workflows

1. **Combine filters** - Use `/keyword` + author filter for powerful searches
2. **Use bookmarks** - Press `m` on important commits, revisit with `b`
3. **Export often** - Save analysis results as CSV for external tools
4. **Scroll slowly** - Take time to read commit messages and authors
5. **Use features** - Each feature reveals different insights about code

---

## Workflow Checklists

### Before Deployment
- [ ] Review last week's commits (`?` → date range)
- [ ] Check team velocity (`t`)
- [ ] Verify code ownership (`a`)
- [ ] Export commit summary (`e` → JSON)

### Bug Investigation
- [ ] Search for related commits (`/bug`)
- [ ] Bisect to find regression (`b`)
- [ ] Analyze author changes (`a`)
- [ ] Export for report (`e`)

### Sprint Planning
- [ ] Check team velocity (`t`)
- [ ] Review code hotspots (`s`)
- [ ] Analyze by author (`?` → author filter)
- [ ] Create capacity plan

### Onboarding New Developer
- [ ] Show code ownership (`a`)
- [ ] Review hotspots (`s`)
- [ ] Share recent changes (`/recent`)
- [ ] Provide export for learning (`e`)

---

## Additional Resources

- [HELP.md](../HELP.md) - Complete command reference
- [docs/QUICKSTART.md](QUICKSTART.md) - Interactive tutorial
- [DEVELOPER.md](../DEVELOPER.md) - Technical details
