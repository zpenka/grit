package grit

import (
	"testing"
)

func TestCopyToClipboard_Valid(t *testing.T) {
	_ = copyToClipboard("test content")
}

func TestOpenInEditor_Valid(t *testing.T) {
	cmd := openInEditor(".", "abc123")
	if cmd == nil {
		t.Error("openInEditor should return non-nil Cmd")
	}
}

func TestFetchCommits_Valid(t *testing.T) {
	cmd := fetchCommits(".", "main")
	if cmd == nil {
		t.Error("fetchCommits should return non-nil Cmd")
	}
}

func TestFetchDiff_Valid(t *testing.T) {
	cmd := fetchDiff(".", "abc123")
	if cmd == nil {
		t.Error("fetchDiff should return non-nil Cmd")
	}
}

func TestFetchBranches_Valid(t *testing.T) {
	cmd := fetchBranches(".")
	if cmd == nil {
		t.Error("fetchBranches should return non-nil Cmd")
	}
}

func TestFetchBlame_Valid(t *testing.T) {
	cmd := fetchBlame(".", "abc123", "file.go")
	if cmd == nil {
		t.Error("fetchBlame should return non-nil Cmd")
	}
}

func TestFlashCmd_Valid(t *testing.T) {
	cmd := flashCmd()
	if cmd == nil {
		t.Error("flashCmd should return non-nil Cmd")
	}
}

func TestRenderCommitRow_Valid(t *testing.T) {
	m := newModel(".")
	commits := []commit{{hash: "abc123"}}
	_ = renderCommitRow(m, commits, 0, 80)
}

func TestRenderFileRow_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderFileRow(m, 0, 80)
}

func TestRenderBranchRow_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderBranchRow(m, 0, 80)
}

func TestRenderBlameRow_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderBlameRow(m, 0, 80)
}

func TestRenderDiffRow_Valid(t *testing.T) {
	m := newModel(".")
	_ = renderDiffRow(m, 0, 80)
}

func TestFirstVisibleHash_Valid(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{{hash: "abc123"}}
	_ = firstVisibleHash(m)
}
