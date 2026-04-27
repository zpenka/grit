# Installation Guide

Multiple ways to install grit depending on your preference.

## Quick Start

### From Source (Recommended for Development)

**Requirements**: Go 1.21 or later

```bash
# Clone the repository
git clone https://github.com/zpenka/grit.git
cd grit

# Build the binary
go build -o grit .

# Run grit
./grit
```

### Using go install (Easiest for Users)

```bash
go install github.com/zpenka/grit@latest

# Then use it
grit
```

### macOS - Homebrew

```bash
brew install zpenka/tap/grit

# Then use it
grit
```

### Linux - Pre-built Binary

Download the latest release from [GitHub Releases](https://github.com/zpenka/grit/releases):

```bash
# Download (choose your OS/architecture)
curl -L https://github.com/zpenka/grit/releases/download/v1.0.0/grit-linux-amd64 -o grit

# Make executable
chmod +x grit

# Move to PATH
sudo mv grit /usr/local/bin/

# Verify
grit --version
```

### Windows - Binary

1. Download `grit-windows-amd64.exe` from [GitHub Releases](https://github.com/zpenka/grit/releases)
2. Add directory to PATH or run directly: `.\grit.exe`

## Build from Source

### Prerequisites

- Go 1.21 or later
- Git
- Terminal with 24-bit color support (recommended)

### Clone and Build

```bash
# Clone repository
git clone https://github.com/zpenka/grit.git
cd grit

# Build binary
go build -o grit .

# Verify
./grit --version
```

### Build for Distribution

Build for multiple platforms:

```bash
# Use included build script
./scripts/build-release.sh

# Or manually:
GOOS=linux GOARCH=amd64 go build -o grit-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o grit-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o grit-windows-amd64.exe
```

## Verification

### Test Installation

```bash
# Show version
grit --version

# Show help
grit --help

# Run in current directory (any git repo)
cd /path/to/git/repo
grit
```

### First Run

1. Navigate to any git repository
2. Run `grit`
3. Use arrow keys to navigate commits
4. Press `?` for help
5. Try `/` to search
6. Press `q` to quit

## Requirements

### System Requirements

- **OS**: Linux, macOS, Windows
- **Memory**: 512 MB RAM minimum, 2 GB recommended
- **Disk**: ~10 MB for binary
- **Terminal**: 80×24 minimum, 120×30 recommended
- **Color support**: 24-bit color recommended (256 color fallback supported)

### Git Requirements

- Git 2.0 or later installed on system
- Access to git repositories locally

### Go Requirements (For Building)

- Go 1.21 or later
- Standard Go build tools

## Troubleshooting

### "grit: command not found"

**Issue**: grit is installed but not in PATH

**Solution**:
```bash
# Check where grit was installed
which grit

# If not in PATH, add to .bashrc or .zshrc
export PATH=$PATH:/path/to/grit

# Or move to a directory in PATH
sudo mv grit /usr/local/bin/
```

### "command not found: go"

**Issue**: Go is not installed

**Solution**:
1. Install Go from [golang.org](https://golang.org/dl/)
2. Add to PATH: `export PATH=$PATH:/usr/local/go/bin`

### Terminal shows wrong colors

**Issue**: 24-bit color not supported

**Solution**:
1. Update terminal emulator
2. Set `TERM` environment variable:
   ```bash
   export TERM=xterm-256color
   ```
3. Use `-theme` flag if available

### Out of memory with large repo

**Issue**: Repository has 100k+ commits

**Solution**:
1. Use filters to reduce visible commits
2. Increase system RAM
3. Close other applications

### Slow performance

**Issue**: Repository is very large

**Solution**:
1. First load caches incrementally
2. Use filters (`/`) to focus on subset
3. Close and reopen grit if sluggish

## Uninstall

### If installed with go install

```bash
# Remove binary
rm $(which grit)

# Or
go clean -i github.com/zpenka/grit
```

### If installed with Homebrew

```bash
brew uninstall grit
```

### If installed manually

```bash
# Remove binary
rm /usr/local/bin/grit
```

## Update

### From go install

```bash
go install github.com/zpenka/grit@latest
```

### From Homebrew

```bash
brew upgrade grit
```

### From source

```bash
cd /path/to/grit
git pull origin main
go build -o grit .
```

## Environment Variables

Optional configuration:

```bash
# Diff cache size (default: 50)
export GRIT_CACHE_SIZE=100

# Default repository path
export GRIT_REPO_PATH=/path/to/repo

# Color theme
export GRIT_COLOR_THEME=dark
```

## Getting Help

- **[README.md](README.md)** - Project overview
- **[HELP.md](HELP.md)** - Command reference
- **[docs/QUICKSTART.md](docs/QUICKSTART.md)** - Interactive tutorial
- **[docs/EXAMPLES.md](docs/EXAMPLES.md)** - Example workflows
- **GitHub Issues** - Report bugs or request features

## Support

For issues, questions, or suggestions:
1. Check [HELP.md](HELP.md) for common commands
2. Review [docs/EXAMPLES.md](docs/EXAMPLES.md) for workflows
3. Open an issue on [GitHub](https://github.com/zpenka/grit/issues)

---

**Happy exploring!** 🚀
