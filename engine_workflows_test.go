package grit

import (
	"strings"
	"testing"
)

// Tests for engine_workflows.go feature functions

func TestLoadWorktrees_ParsesOutput(t *testing.T) {
	output := `main	    /path/to/main
feature	  /path/to/feature
bugfix	   /path/to/bugfix`

	worktrees := loadWorktrees(output)

	if len(worktrees) == 0 {
		t.Error("loadWorktrees should parse non-empty output")
	}

	if len(worktrees) != 3 {
		t.Errorf("loadWorktrees should find 3 worktrees, found %d", len(worktrees))
	}

	for _, wt := range worktrees {
		if wt.path == "" || wt.branch == "" {
			t.Error("loadWorktrees should set path and branch")
		}
	}
}

func TestLoadWorktrees_HandlesEmpty(t *testing.T) {
	output := ""

	worktrees := loadWorktrees(output)

	if len(worktrees) != 0 {
		t.Errorf("loadWorktrees should return empty for empty input, got %d", len(worktrees))
	}
}

func TestSwitchWorktree_UpdatesModel(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)

	m = switchWorktree(m, "/path/to/worktree")

	if m.currentWorktree != "/path/to/worktree" {
		t.Errorf("switchWorktree should update currentWorktree")
	}
}

func TestParseSubmodules_ParsesOutput(t *testing.T) {
	output := ` aaa1111 submodule1 (v1.0.0)
 bbb2222 submodule2 (v2.0.0)
 ccc3333 submodule3 (detached)`

	submodules := parseSubmodules(output)

	if len(submodules) == 0 {
		t.Error("parseSubmodules should parse non-empty output")
	}

	for _, sm := range submodules {
		if sm.path == "" {
			t.Error("parseSubmodules should set submodule path")
		}
	}
}

func TestCreateNamedStash_AddsToModel(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)

	m = createNamedStash(m, 0, "wip-feature", "Work in progress")

	if len(m.namedStashes) == 0 {
		t.Error("createNamedStash should add stash to model")
	}

	lastStash := m.namedStashes[len(m.namedStashes)-1]
	if lastStash.name != "wip-feature" {
		t.Errorf("createNamedStash should set stash name, got %q", lastStash.name)
	}

	if lastStash.description != "Work in progress" {
		t.Errorf("createNamedStash should set description")
	}
}

func TestQueueTagOperation_AddsOperation(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)

	m = queueTagOperation(m, "v1.0.0", "aaa1111111111111111111111111111111111", "create", "Release v1.0.0")

	if len(m.pendingTagOps) == 0 {
		t.Error("queueTagOperation should add operation to model")
	}

	lastOp := m.pendingTagOps[len(m.pendingTagOps)-1]
	if lastOp.name != "v1.0.0" {
		t.Errorf("queueTagOperation should set tag name")
	}

	if lastOp.action != "create" {
		t.Errorf("queueTagOperation should set action")
	}
}

func TestExtractGPGSignatureStatus_ParsesStatus(t *testing.T) {
	output := `aaa1111 Good signature from Alice
bbb2222 Bad signature, key missing
ccc3333 No signature`

	status := extractGPGSignatureStatus(output)

	if len(status) == 0 {
		t.Error("extractGPGSignatureStatus should parse output")
	}

	for hash, sig := range status {
		if hash == "" {
			t.Error("extractGPGSignatureStatus should set hash key")
		}
		if sig.signer == "" && !sig.signed {
			t.Error("extractGPGSignatureStatus should set signer or signed status")
		}
	}
}

func TestRenderWorktreesUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showWorktrees = true
	m.worktrees = []worktreeInfo{
		{path: "/path/main", branch: "main", hash: "aaa1111111111111111111111111111111111"},
		{path: "/path/feature", branch: "feature", hash: "bbb2222222222222222222222222222222222"},
	}

	output := renderWorktreesUI(m, 80)

	if output == "" {
		t.Error("renderWorktreesUI should produce output")
	}

	if !strings.Contains(output, "main") && !strings.Contains(output, "feature") {
		t.Error("renderWorktreesUI output should include worktree names")
	}
}

func TestRenderSubmodulesUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showSubmodules = true
	m.submodules = []submoduleInfo{
		{path: "deps/dep1", url: "https://github.com/user/dep1.git", hash: "aaa1111111111111111111111111111111111", branch: "main"},
	}

	output := renderSubmodulesUI(m, 80)

	if output == "" {
		t.Error("renderSubmodulesUI should produce output")
	}
}

func TestRenderNamedStashesUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showNamedStashes = true
	m.namedStashes = []namedStash{
		{name: "wip", description: "Work in progress", hash: "aaa1111111111111111111111111111111111"},
	}

	output := renderNamedStashesUI(m, 80)

	if output == "" {
		t.Error("renderNamedStashesUI should produce output")
	}
}

func TestRenderTagMgmtUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showTagMgmt = true
	m.pendingTagOps = []tagOperation{
		{name: "v1.0.0", hash: "aaa1111111111111111111111111111111111", action: "create"},
	}

	output := renderTagMgmtUI(m, 80)

	if output == "" {
		t.Error("renderTagMgmtUI should produce output")
	}
}

func TestRenderGPGStatusUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(3)
	m.showGPGStatus = true
	m.gpgStatuses = map[string]gpgSignatureStatus{
		"aaa1111111111111111111111111111111111": {signed: true, signer: "Alice"},
	}

	output := renderGPGStatusUI(m, 80)

	if output == "" {
		t.Error("renderGPGStatusUI should produce output")
	}
}
