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

### Movement
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up / Previous commit |
| `↓` / `j` | Move down / Next commit |
| `←` / `h` | Scroll left |
| `→` / `l` | Scroll right |
| `Page Up` | Jump up (10 commits) |
| `Page Down` | Jump down (10 commits) |
| `Home` | First commit |
| `End` | Last commit |

### Panel Navigation
| Key | Action |
|-----|--------|
| `Tab` | Switch panel (commit → diff → files) |
| `Shift+Tab` | Previous panel |
| `1` | Commit list panel |
| `2` | Diff panel |
| `3` | File tree panel |
| `4` | Feature results panel |

### Search & Filter
| Key | Action |
|-----|--------|
| `/` | Search commits (text, author, hash) |
| `?` | Advanced filter options |
| `n` | Find next match |
| `N` | Find previous match |
| `Esc` | Clear search |

### Features
| Key | Action |
|-----|--------|
| `a` | Code ownership analysis |
| `s` | Code statistics & hotspots |
| `b` | Bisect regression |
| `t` | Team velocity metrics |
| `v` | Visualizations |
| `e` | Export (CSV/JSON/XML) |
| `c` | Copy commit hash |

### General
| Key | Action |
|-----|--------|
| `?` | Show help |
| `q` / `Ctrl+C` | Quit |
| `r` | Refresh |
| `u` | Undo last action |

## Examples

### Find commits by keyword
```bash
grit                    # Start grit
/bug                    # Search for commits with "bug"
n                       # Find next match
```

### Analyze code ownership
```bash
grit                    # Start grit
a                       # View code ownership analysis
```

### Explore with bisect
```bash
grit                    # Start grit
b                       # Start bisect
[mark commits as good/bad]
```

### Export results
```bash
grit                    # Start grit
e                       # Export
[select format: CSV/JSON/XML]
```

## Tips & Tricks

1. **Search is case-insensitive** - `/FIX` finds "fix", "Fix", "FIX"
2. **Author filter** - Use `/author:john` to filter by author
3. **Date filter** - Use `?` to access advanced date range filters
4. **Bookmarks** - Save interesting commits with `m` to revisit
5. **Diff caching** - First view of diff is slower, subsequent views are cached
6. **Large repos** - grit handles 10k+ commits with lazy loading

## Environment Variables

```bash
GRIT_CACHE_SIZE=100       # Diff cache entries (default: 50)
GRIT_REPO_PATH=/path      # Default repository
GRIT_COLOR_THEME=dark     # Color theme
```

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
