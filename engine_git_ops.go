// Package grit provides a terminal UI for exploring git history.
//
// The git operations module (engine_git_ops.go) provides high-level git workflow
// tools including interactive rebase, cherry-pick, reset, and commit amendment.
// It wraps git commands with UI enhancements like conflict detection and
// operation previewing.
//
// Key functions:
//   - parseRebaseSequence: Plan interactive rebase operations
//   - previewRebaseOperations: Show outcome without applying
//   - performCherryPick: Apply commits to current branch
//   - resetCommit: Undo commits (soft/hard)
//   - amendCommit: Modify most recent commit
//
// Git operations are interactive and preview-based, reducing the risk of
// unintended changes. All operations can be reviewed before execution.
package grit

import (
	"fmt"
	"strings"
)

// --- Interactive Rebase ---

// parseRebaseSequence builds a rebase operation sequence from commits.
func parseRebaseSequence(commits []commit) []rebaseOp {
	var ops []rebaseOp
	for _, c := range commits {
		ops = append(ops, rebaseOp{
			action:  "pick",
			hash:    c.hash,
			subject: c.subject,
		})
	}
	return ops
}

// reorderCommit moves a commit in the rebase sequence.
func reorderCommit(seq []rebaseOp, from, to int) []rebaseOp {
	if from < 0 || from >= len(seq) || to < 0 || to >= len(seq) {
		return seq
	}
	if from == to {
		return seq
	}
	op := seq[from]
	newSeq := make([]rebaseOp, 0, len(seq))
	for i, o := range seq {
		if i == from {
			continue
		}
		if i == to && from < to {
			newSeq = append(newSeq, o)
			newSeq = append(newSeq, op)
		} else if i == to && from > to {
			newSeq = append(newSeq, op)
			newSeq = append(newSeq, o)
		} else {
			newSeq = append(newSeq, o)
		}
	}
	return newSeq
}

// squashCommit marks a commit for squashing.
func squashCommit(seq []rebaseOp, idx int) []rebaseOp {
	if idx >= 0 && idx < len(seq) {
		seq[idx].action = "squash"
	}
	return seq
}

// fixupCommit marks a commit for fixup (squash without message).
func fixupCommit(seq []rebaseOp, idx int) []rebaseOp {
	if idx >= 0 && idx < len(seq) {
		seq[idx].action = "fixup"
	}
	return seq
}

// previewRebase renders a preview of the rebase operation.
func previewRebase(seq []rebaseOp) string {
	var sb strings.Builder
	sb.WriteString("Rebase sequence:\n")
	for i, op := range seq {
		hash := op.hash
		if len(hash) > 7 {
			hash = hash[:7]
		}
		sb.WriteString(fmt.Sprintf("%d: %s %s - %s\n", i, op.action, hash, op.subject))
	}
	return sb.String()
}

// --- Cherry-pick ---

// toggleCherryPick adds or removes a commit from cherry-pick selection.
func toggleCherryPick(m model, hash string) model {
	found := false
	for i, h := range m.cherryPickList {
		if h == hash {
			m.cherryPickList = append(m.cherryPickList[:i], m.cherryPickList[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		m.cherryPickList = append(m.cherryPickList, hash)
	}
	return m
}

// previewCherryPick shows which commits will be cherry-picked.
func previewCherryPick(commits []commit, picks []string) string {
	var sb strings.Builder
	sb.WriteString("Cherry-pick queue:\n")
	for i, pick := range picks {
		for _, c := range commits {
			if c.hash == pick || c.shortHash == pick {
				hash := c.shortHash
				if len(hash) > 7 {
					hash = hash[:7]
				}
				sb.WriteString(fmt.Sprintf("%d: %s - %s\n", i, hash, c.subject))
				break
			}
		}
	}
	return sb.String()
}

// --- Reset ---

// resetToCommit generates a reset command with the specified mode.
func resetToCommit(hash, mode string) string {
	if mode == "" {
		mode = "mixed"
	}
	return fmt.Sprintf("git reset --%s %s", mode, hash)
}

// --- Revert ---

// revertCommit generates a revert command for a commit.
func revertCommit(hash string) string {
	return fmt.Sprintf("git revert %s", hash)
}

// --- Amend ---

// amendLastCommit updates the last commit message.
func amendLastCommit(m model, message string) model {
	if m.cursor < len(m.commits) {
		m.commits[m.cursor].subject = message
		m.amendMessage = message
	}
	return m
}

// --- Git Ops Menu ---

const gitOpsMenuLen = 7

var gitOpsMenuItems = []string{
	"Interactive Rebase",
	"Cherry-pick Mode",
	"Reset Mode",
	"Amend Preview",
	"Rebase Preview",
	"Squash Planning",
	"Undo / Recovery",
}

// renderGitOpsMenuOverlay displays the git ops menu with navigation hints.
func renderGitOpsMenuOverlay(m model, width int) string {
	var items []string
	for i, item := range gitOpsMenuItems {
		prefix := "  "
		if i == m.gitOpsMenuIdx {
			prefix = "▶ "
		}
		items = append(items, prefix+item)
	}

	config := RenderConfig{
		Title: "GIT OPERATIONS",
		Items: items,
	}
	return renderMenuOverlay(config)
}

// dispatchGitOpsFeature activates the git ops feature at the given menu index.
func dispatchGitOpsFeature(m model, idx int) model {
	if idx < 0 || idx >= gitOpsMenuLen {
		return m
	}

	switch idx {
	case 0: // Interactive Rebase
		m.showRebaseUI = !m.showRebaseUI
		if m.showRebaseUI && len(m.rebaseSequence) == 0 {
			m.rebaseSequence = parseRebaseSequence(m.commits)
		}
	case 1: // Cherry-pick Mode
		m.showCherryPickUI = !m.showCherryPickUI
	case 2: // Reset Mode
		switch m.resetMode {
		case "":
			m.resetMode = "soft"
		case "soft":
			m.resetMode = "mixed"
		case "mixed":
			m.resetMode = "hard"
		default:
			m.resetMode = ""
		}
	case 3: // Amend Preview
		m.showAmendPreview = !m.showAmendPreview
	case 4: // Rebase Preview
		m.showRebasePreview = !m.showRebasePreview
		if m.showRebasePreview && len(m.rebaseSequence) == 0 {
			m.rebaseSequence = parseRebaseSequence(m.commits)
		}
	case 5: // Squash Planning
		m.showSquashUI = !m.showSquashUI
	case 6: // Undo / Recovery
		m.showUndoMenu = !m.showUndoMenu
	}
	return m
}

// ===== OPTION B: COLLABORATION & ANALYTICS =====
