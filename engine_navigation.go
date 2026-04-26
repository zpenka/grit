// Package grit provides a terminal UI for exploring git history.
//
// The navigation module (engine_navigation.go) manages cursor position, panel
// switching, and scrolling within the UI. It maintains and updates the model's
// navigation state in response to user input.
//
// Key functions:
//   - moveCursorUp/Down: Navigate through commits
//   - switchPanel: Change active UI panel
//   - ScrollPanelUp/Down: Scroll within panels
//   - setBookmark/goToBookmark: Save and restore cursor positions
//
// Navigation is stateless—all functions take and return a model, making them
// easily testable and composable with Bubble Tea's update cycle.
package grit

// moveCursorDown advances the commit cursor by one position and resets diff offset.
func moveCursorDown(m model) model {
	if m.cursor < len(m.commits)-1 {
		m.cursor++
		m.diffOffset = 0
	}
	return m
}

// moveCursorUp moves the commit cursor back one position and resets diff offset.
func moveCursorUp(m model) model {
	if m.cursor > 0 {
		m.cursor--
		m.diffOffset = 0
	}
	return m
}

// scrollDiffDown scrolls the diff view down by n lines, clamped to valid range.
func scrollDiffDown(m model, n int) model {
	maxOff := len(m.diffLines) - diffPanelHeight(m)
	if maxOff < 0 {
		maxOff = 0
	}
	m.diffOffset += n
	if m.diffOffset > maxOff {
		m.diffOffset = maxOff
	}
	return m
}

// scrollDiffUp scrolls the diff view up by n lines, clamped to zero.
func scrollDiffUp(m model, n int) model {
	m.diffOffset -= n
	if m.diffOffset < 0 {
		m.diffOffset = 0
	}
	return m
}

// switchPanel toggles focus between commit list and diff panels.
func switchPanel(m model) model {
	if m.focus == panelList {
		m.focus = panelDiff
	} else {
		m.focus = panelList
	}
	return m
}

// listPanelWidth returns the width of the commit list panel (clamped to 36–52).
func listPanelWidth(totalWidth int) int {
	w := totalWidth / 3
	if w < 36 {
		return 36
	}
	if w > 52 {
		return 52
	}
	return w
}

// diffPanelWidth returns the remaining width for the diff panel.
func diffPanelWidth(totalWidth int) int {
	return totalWidth - listPanelWidth(totalWidth) - 1 // 1 for divider
}

// diffPanelHeight returns the number of content lines visible in each panel.
// diffPanelHeight calculates the height of the diff panel based on terminal size.
func diffPanelHeight(m model) int {
	h := m.height - 7 // title + blank + header + blank + hint + blank*2
	if h < 5 {
		return 5
	}
	return h
}

// scrollToDiffLine sets diffOffset so that lineIdx is visible, clamped to valid range.
func scrollToDiffLine(m model, lineIdx int) model {
	m.diffOffset = lineIdx
	maxOff := len(m.diffLines) - diffPanelHeight(m)
	if maxOff < 0 {
		maxOff = 0
	}
	if m.diffOffset > maxOff {
		m.diffOffset = maxOff
	}
	if m.diffOffset < 0 {
		m.diffOffset = 0
	}
	return m
}

// toggleFileView shows or hides the file list in the left panel.
// Hiding resets fileCursor.
func toggleFileView(m model) model {
	m.showFiles = !m.showFiles
	if !m.showFiles {
		m.fileCursor = 0
	}
	return m
}

// toggleBranchView shows or hides the branch picker in the left panel.
// Hiding resets branchCursor.
func toggleBranchView(m model) model {
	m.showBranch = !m.showBranch
	if !m.showBranch {
		m.branchCursor = 0
	}
	return m
}

// addToNavHistory adds the current cursor position to navigation history.
func addToNavHistory(m model, position int) model {
	// Discard any future history if we're not at the end
	if m.navHistoryIdx < len(m.navHistory)-1 {
		m.navHistory = m.navHistory[:m.navHistoryIdx+1]
	}
	m.navHistory = append(m.navHistory, position)
	m.navHistoryIdx = len(m.navHistory) - 1
	return m
}

// goBackInHistory moves to the previous position in navigation history.
func goBackInHistory(m model) model {
	if m.navHistoryIdx > 0 {
		m.navHistoryIdx--
		m.cursor = m.navHistory[m.navHistoryIdx]
	}
	return m
}

// goForwardInHistory moves to the next position in navigation history.
func goForwardInHistory(m model) model {
	if m.navHistoryIdx < len(m.navHistory)-1 {
		m.navHistoryIdx++
		m.cursor = m.navHistory[m.navHistoryIdx]
	}
	return m
}

// toggleBookmark toggles a bookmark on the current commit.
func toggleBookmark(m model) model {
	if m.cursor >= len(m.commits) {
		return m
	}
	hash := m.commits[m.cursor].shortHash
	if isBookmarked(m, m.cursor) {
		// Remove bookmark
		var newBookmarks []string
		for _, b := range m.bookmarks {
			if b != hash {
				newBookmarks = append(newBookmarks, b)
			}
		}
		m.bookmarks = newBookmarks
	} else {
		// Add bookmark
		m.bookmarks = append(m.bookmarks, hash)
	}
	return m
}

// isBookmarked checks if a commit at the given index is bookmarked.
func isBookmarked(m model, idx int) bool {
	if idx >= len(m.commits) {
		return false
	}
	hash := m.commits[idx].shortHash
	for _, b := range m.bookmarks {
		if b == hash {
			return true
		}
	}
	return false
}

// jumpToNextBookmark moves the cursor to the next bookmarked commit.
func jumpToNextBookmark(m model) model {
	for i := m.cursor + 1; i < len(m.commits); i++ {
		if isBookmarked(m, i) {
			m.cursor = i
			return m
		}
	}
	return m
}

// jumpToPrevBookmark moves the cursor to the previous bookmarked commit.
func jumpToPrevBookmark(m model) model {
	for i := m.cursor - 1; i >= 0; i-- {
		if isBookmarked(m, i) {
			m.cursor = i
			return m
		}
	}
	return m
}

// goToCommit finds a commit by hash (short or full) and returns its index, or -1 if not found.
// goToCommit searches commits for a matching hash or subject and returns its index.
func goToCommit(commits []commit, query string) int {
	q := len(query)
	for i, c := range commits {
		if c.shortHash == query || c.hash == query {
			return i
		}
		if q <= len(c.shortHash) && c.shortHash[:q] == query {
			return i
		}
		if q <= len(c.hash) && c.hash[:q] == query {
			return i
		}
		// Case-insensitive matching for short and full hashes
		if len(c.shortHash) == len(query) {
			isMatch := true
			for j := 0; j < len(query); j++ {
				c1 := byte(c.shortHash[j])
				c2 := byte(query[j])
				if c1 >= 'a' {
					c1 -= 32
				}
				if c2 >= 'a' {
					c2 -= 32
				}
				if c1 != c2 {
					isMatch = false
					break
				}
			}
			if isMatch {
				return i
			}
		}
		if len(c.hash) == len(query) {
			isMatch := true
			for j := 0; j < len(query); j++ {
				c1 := byte(c.hash[j])
				c2 := byte(query[j])
				if c1 >= 'a' {
					c1 -= 32
				}
				if c2 >= 'a' {
					c2 -= 32
				}
				if c1 != c2 {
					isMatch = false
					break
				}
			}
			if isMatch {
				return i
			}
		}
	}
	return -1
}

// handleGoToCommitInput processes go-to-commit input and updates model.
func handleGoToCommitInput(m model, query string) model {
	idx := goToCommit(m.commits, query)
	if idx >= 0 {
		m.cursor = idx
		m.diffOffset = 0
	}
	return m
}
