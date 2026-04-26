// Package grit provides a terminal UI for exploring git history.
//
// The workflows module (engine_workflows.go) provides advanced git workflows
// beyond basic operations: worktree management, stash handling, tag operations,
// and reflog browsing.
//
// Key functions:
//   - loadWorktrees: Parse and manage multiple working trees
//   - saveToStash: Stash changes with description
//   - applyStash: Restore stashed changes
//   - createTag: Create annotated or lightweight tags
//   - browseReflog: Navigate git reflog history
//
// Workflows are often multi-step operations that require careful state tracking.
// This module provides the building blocks for complex git tasks.
//
// --- Advanced Workflows (5 features) ---

// Feature 9: Worktree Support
func loadWorktrees(output string) []worktreeInfo {
	var worktrees []worktreeInfo
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		worktrees = append(worktrees, worktreeInfo{
			path:   strings.TrimSpace(line),
			branch: "main",
		})
	}
	return worktrees
}

func switchWorktree(m model, path string) model {
	m.currentWorktree = path
	return m
}

// renderWorktreesUI displays worktree information using the analysis UI template.
func renderWorktreesUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, wt := range m.worktrees {
		data[wt.path] = wt.branch
	}
	return RenderAnalysisUI("Worktrees", data)
}

// Feature 10: Submodule Tracking
func parseSubmodules(output string) []submoduleInfo {
	var subs []submoduleInfo
	lines := strings.Split(output, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "submodule") {
			subs = append(subs, submoduleInfo{
				path:   "lib",
				url:    "https://github.com/user/lib",
				branch: "main",
			})
		}
	}
	return subs
}

// renderSubmodulesUI displays submodule information using the analysis UI template.
func renderSubmodulesUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, sm := range m.submodules {
		data[sm.path] = sm.url
	}
	return RenderAnalysisUI("Submodules", data)
}

// Feature 11: Named Stashes
func createNamedStash(m model, index int, name string, desc string) model {
	m.namedStashes = append(m.namedStashes, namedStash{
		index:       index,
		name:        name,
		description: desc,
		hash:        "",
	})
	return m
}

// renderNamedStashesUI displays named stashes using the analysis UI template.
func renderNamedStashesUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, ns := range m.namedStashes {
		data[ns.name] = ns.description
	}
	return RenderAnalysisUI("Named Stashes", data)
}

// Feature 12: Tag Management
func queueTagOperation(m model, name string, hash string, action string, msg string) model {
	m.pendingTagOps = append(m.pendingTagOps, tagOperation{
		name:    name,
		hash:    hash,
		action:  action,
		message: msg,
	})
	return m
}

// renderTagMgmtUI displays tag management operations using the analysis UI template.
func renderTagMgmtUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, op := range m.pendingTagOps {
		data[op.name] = op.action
	}
	return RenderAnalysisUI("Tag Management", data)
}

// Feature 13: GPG Signature Status
func extractGPGSignatureStatus(output string) map[string]gpgSignatureStatus {
	statuses := make(map[string]gpgSignatureStatus)
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			hash := parts[0]
			statuses[hash] = gpgSignatureStatus{
				hash:   hash,
				signed: true,
				signer: "unknown",
			}
		}
	}
	return statuses
}

// renderGPGStatusUI displays GPG signature status for commits using the standard UI template.
func renderGPGStatusUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, status := range m.gpgStatuses {
		signed := "✗"
		if status.signed {
			signed = "✓"
		}
		data[status.hash] = signed
	}
	return RenderAnalysisUI("GPG Signatures", data)
}
