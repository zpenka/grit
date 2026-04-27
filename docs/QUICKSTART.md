# Quick Start Tutorial

Get up and running with grit in 5 minutes.

## Step 1: Installation

Choose one method:

```bash
# Method A: Using go install (easiest)
go install github.com/zpenka/grit@latest

# Method B: Build from source
git clone https://github.com/zpenka/grit.git
cd grit
go build -o grit .

# Method C: Homebrew (macOS)
brew install zpenka/tap/grit
```

Verify installation:
```bash
grit --version
```

See [INSTALL.md](../INSTALL.md) for more installation options.

---

## Step 2: First Run

Open any git repository and start grit:

```bash
cd /path/to/your/git/repo
grit
```

You should see:
- **Left panel**: List of commits (most recent at bottom)
- **Right panel**: Details of selected commit
- **Bottom**: Status bar with help hints

---

## Step 3: Basic Navigation

Try these keyboard shortcuts:

| Key | What happens |
|-----|--------------|
| `↓` or `j` | Move down to next commit |
| `↑` or `k` | Move up to previous commit |
| `Page Down` | Jump 10 commits forward |
| `Page Up` | Jump 10 commits backward |
| `Home` | Jump to first commit |
| `End` | Jump to last commit |

**Try it**: Press `↓` a few times to navigate through commits.

---

## Step 4: View Commit Details

When you select a commit:
- **Subject line** - Commit message title
- **Author** - Who made the change
- **Date** - When it was committed
- **Diff** - Detailed changes (files, additions, deletions)

**Try it**: Press `→` to see the full diff of a commit.

---

## Step 5: Search & Filter

Find specific commits:

```bash
grit
/bug        # Search for "bug"
n           # Find next match
N           # Find previous match
Esc         # Clear search
```

Search works on:
- Commit messages
- Author names
- Commit hashes

**Try it**: Type `/` followed by a word from a commit message.

---

## Step 6: Features

Try grit's powerful features:

| Key | Feature |
|-----|---------|
| `a` | Code ownership (who changed what) |
| `s` | Statistics & hotspots (what changed most) |
| `t` | Team metrics (who commits most) |
| `v` | Visualizations (activity heatmap) |
| `e` | Export (CSV, JSON, XML) |

**Try it**: Press `a` to see code ownership analysis.

---

## Step 7: Get Help

Access the help menu:

```bash
grit
?           # Show help
```

Or read the help documentation:
- [HELP.md](../HELP.md) - Complete reference
- [docs/EXAMPLES.md](EXAMPLES.md) - Workflow examples

---

## Common Tasks

### Search by Author

```bash
grit
/author:john    # Show commits from john
```

### Find Recent Changes

```bash
grit
/              # Start search
recent         # Search for "recent" changes
```

### Analyze Team Activity

```bash
grit
t              # Team metrics
```

### Export for Analysis

```bash
grit
e              # Export
CSV            # Choose format
```

---

## Tips for Success

1. **Use keyboard shortcuts** - They're much faster than mouse
2. **Combine features** - Search first, then analyze with features
3. **Read commit messages** - They explain the "why" behind changes
4. **Use bookmarks** - Press `m` on important commits to mark them
5. **Check author info** - Learn who maintains which parts

---

## Next Steps

Explore these based on your needs:

**Want to understand code?**
→ [docs/EXAMPLES.md](EXAMPLES.md) workflow #1 (Code Ownership)

**Want to track bugs?**
→ [docs/EXAMPLES.md](EXAMPLES.md) workflow #3 (Bisect)

**Want team insights?**
→ [docs/EXAMPLES.md](EXAMPLES.md) workflow #7 (Team Velocity)

**Want to try advanced features?**
→ [HELP.md](../HELP.md) - Complete keyboard reference

---

## Troubleshooting

**Terminal looks weird?**
- Try pressing `r` to refresh
- Make sure terminal is at least 80 columns wide

**Colors are wrong?**
- Update your terminal emulator
- Set `export TERM=xterm-256color`

**grit is slow?**
- Use filters (`/`) to reduce visible commits
- Check system resources

**Not finding what you need?**
- See [HELP.md](../HELP.md) for complete reference
- Check [docs/EXAMPLES.md](EXAMPLES.md) for workflows
- Open an issue on [GitHub](https://github.com/zpenka/grit/issues)

---

## You're Ready! 🎉

You now know:
✅ How to install grit
✅ How to navigate commits
✅ How to search and filter
✅ How to analyze code
✅ How to export data

Start exploring your git history!

```bash
grit
```

---

**Questions?** See [HELP.md](../HELP.md) or open an issue on GitHub.
