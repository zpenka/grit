# grit - Command Line Help

Terminal UI for exploring git history interactively.

## Usage

```bash
grit [OPTIONS] [REPOSITORY]
```

## Commands & Options

### Starting grit

```bash
grit                    # Current directory
grit /path/to/repo      # Specific repository
grit -h                 # Help
grit -v                 # Version
```

## Keyboard Navigation

### Commit List Navigation
| Key | Action |
|-----|--------|
| `j` / `↓` | Move down to next commit |
| `k` / `↑` | Move up to previous commit |
| `5j` | Jump 5 commits down (count prefix works with j/k) |
| `g` | Jump to first commit |
| `G` | Jump to last commit |
| `/` | Search commits by subject, author, or hash |
| `Ctrl+/` | Clear search |

### Panel Navigation
| Key | Action |
|-----|--------|
| `Tab` / `l` | Switch to diff panel |
| `h` | Switch to commit list |
| `f` | Toggle file list panel |
| `b` | Show branch picker |
| `B` | Show blame for current file |

### Feature Panels
| Key | Action |
|-----|--------|
| `a` | Analytics submenu (ownership, hotspots, linting, bisect, heatmap, stats, complexity) |
| `v` | Visualization submenu (flamegraph, timeline, tree view, author comparison, file heatmap) |
| `t` | Team/AI submenu (team stats, reviewer suggestions, velocity, classification, security, changelog) |
| `r` | Interactive rebase preview |
| `c` | Cherry-pick mode (toggle commits to pick) |
| `x` | Cycle reset mode (soft → mixed → hard) |

### Diff Panel Navigation
| Key | Action |
|-----|--------|
| `j` / `↓` | Scroll down one line |
| `k` / `↑` | Scroll up one line |
| `d` / `space` | Scroll down half page |
| `u` | Scroll up half page |
| `g` | Jump to top |
| `G` | Jump to bottom |

### General
| Key | Action |
|-----|--------|
| `y` | Copy current commit hash |
| `e` | Open commit diff in `$EDITOR` |
| `?` | Show help overlay |
| `Esc` | Close any open panel or menu |
| `q` / `Ctrl+C` | Quit grit |

## Examples

### Find commits by keyword
```bash
grit                    # Start grit
/bug                    # Search for commits with "bug"
j                       # Move to next result
```

### Analyze code ownership
```bash
grit                    # Start grit
a                       # Open analytics submenu
[navigate with j/k and press Enter to select Code Ownership]
```

### Use bisect to find a regression
```bash
grit                    # Start grit
a                       # Open analytics submenu
[navigate to Bisect and press Enter]
[mark commits as good/bad in the bisect UI]
```

### View team metrics
```bash
grit                    # Start grit
t                       # Open team/AI submenu
[navigate with j/k and select Team Statistics]
```

### Explore visualizations
```bash
grit                    # Start grit
v                       # Open visualization submenu
[navigate with j/k and select Contributor Flamegraph]
```

## Tips & Tricks

1. **Quick navigation** - Use count prefixes: `5j` jumps 5 commits down
2. **Search is live** - Type `/` and characters appear as you type
3. **Fast navigation** - `g` goes to first commit, `G` goes to last
4. **Multiple panels** - Open file list with `f` to jump between files
5. **Diff caching** - First view of diff is slower, subsequent views are cached
6. **Large repos** - grit handles 10k+ commits with lazy loading
7. **Press ? for help** - The help overlay shows all keybindings at any time

## Performance Notes

- **Startup**: 1-2 seconds for 1k commits, 5-10s for 10k+ commits
- **Diff loading**: First diff takes 100-500ms, cached diffs are instant
- **Searches**: Real-time, incremental (no delay)
- **Large repos**: Use filters to reduce visible commits

## Troubleshooting

### Terminal too small
- Minimum: 80 columns × 24 rows
- Recommended: 120 columns × 30 rows

### Colors not showing
- Ensure terminal supports 24-bit color
- Check `TERM` environment variable

### Slow performance
- Use filters to reduce commits
- Check `grit -v` for version (newer is faster)

## See Also

- [README.md](README.md) - Project overview
- [DEVELOPER.md](DEVELOPER.md) - Development guide
- [docs/EXAMPLES.md](docs/EXAMPLES.md) - Example workflows
- [docs/QUICKSTART.md](docs/QUICKSTART.md) - Tutorial
