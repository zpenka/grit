package grit

import (
	"testing"
)

// Tests for engine_visualization.go feature functions

func TestBuildTreeView_HandlesEmptyCommits(t *testing.T) {
	commits := []commit{}
	tree := buildTreeView(commits)

	if tree == nil {
		t.Error("buildTreeView should return a tree node for empty commits")
	}
}

func TestBuildTreeView_BuildsHierarchy(t *testing.T) {
	// Create commits with parent relationships
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "Root commit"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Child of aaa"},
		{hash: "ccc3333333333333333333333333333333333", subject: "Another child of aaa"},
		{hash: "ddd4444444444444444444444444444444444", subject: "Grandchild of bbb"},
	}

	tree := buildTreeView(commits)

	if tree == nil {
		t.Error("buildTreeView should return non-nil tree")
	}

	// Tree should not be flat (depth-1) if it has proper structure
	if tree.children != nil && len(tree.children) > 0 {
		// Check that we have nested structure, not just depth-1
		hasNesting := false
		for _, child := range tree.children {
			if child.children != nil && len(child.children) > 0 {
				hasNesting = true
				break
			}
		}
		// A proper tree with these commits should have some nesting
		if len(commits) > 2 && !hasNesting {
			t.Error("buildTreeView should create nested structure, not just flat tree")
		}
	}
}

func TestBuildContributorFlame_HasValidData(t *testing.T) {
	commits := []commit{
		{author: "Alice", hash: "aaa1111111111111111111111111111111111"},
		{author: "Bob", hash: "bbb2222222222222222222222222222222222"},
		{author: "Alice", hash: "ccc3333333333333333333333333333333333"},
	}

	flame := buildContributorFlame(commits)

	if len(flame) == 0 {
		t.Error("buildContributorFlame should return data")
	}

	for _, entry := range flame {
		if entry.author == "" {
			t.Error("buildContributorFlame should set author name")
		}
		if entry.commits < 0 {
			t.Error("buildContributorFlame should have non-negative commit count")
		}
	}
}

func TestBuildTimeline_CreatesSequence(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", when: "1 day ago", subject: "First"},
		{hash: "bbb2222222222222222222222222222222222", when: "2 days ago", subject: "Second"},
		{hash: "ccc3333333333333333333333333333333333", when: "3 days ago", subject: "Third"},
	}

	timeline := buildTimeline(commits)

	if len(timeline) == 0 {
		t.Error("buildTimeline should return timeline data")
	}

	// Verify timeline entries have data
	for _, point := range timeline {
		if point.hash == "" {
			t.Error("buildTimeline should set commit hash")
		}
	}
}

func TestCompareAuthors_ReturnsPairings(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{
		{author: "Alice", hash: "aaa1111111111111111111111111111111111"},
		{author: "Bob", hash: "bbb2222222222222222222222222222222222"},
		{author: "Alice", hash: "ccc3333333333333333333333333333333333"},
		{author: "Bob", hash: "ddd4444444444444444444444444444444444"},
	}

	comparisons := compareAuthors(m)

	// Should have comparison data
	if len(comparisons) == 0 {
		t.Error("compareAuthors should return comparison data")
	}

	for _, comp := range comparisons {
		if comp.author1 == "" && comp.author2 == "" {
			t.Error("compareAuthors should set author names")
		}
	}
}

func TestBuildFileHeatmap_TracksFileActivity(t *testing.T) {
	commits := []commit{
		{hash: "aaa1111111111111111111111111111111111", subject: "Edit main.go"},
		{hash: "bbb2222222222222222222222222222222222", subject: "Edit main.go"},
		{hash: "ccc3333333333333333333333333333333333", subject: "Edit utils.go"},
	}

	heatmap := buildFileHeatmap(commits)

	if len(heatmap) == 0 {
		t.Error("buildFileHeatmap should return heatmap data")
	}

	for _, entry := range heatmap {
		if entry.path == "" {
			t.Error("buildFileHeatmap should set file path")
		}
		if entry.frequency < 0 {
			t.Error("buildFileHeatmap should have non-negative frequency")
		}
	}
}

func TestFlattenTreeNode_ProducesValidList(t *testing.T) {
	// Create a simple tree
	root := &treeNode{
		hash: "root",
		children: []*treeNode{
			{hash: "child1"},
			{hash: "child2"},
		},
	}

	var items []string
	flattenTreeNode(root, &items)

	if len(items) == 0 {
		t.Error("flattenTreeNode should produce items")
	}

	// All items should have content
	for _, item := range items {
		if item == "" {
			t.Error("flattenTreeNode should produce non-empty items")
		}
	}
}

func TestRenderFlamegraphUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.showFlamegraph = true

	output := renderFlamegraphUI(m, 80)

	if output == "" {
		t.Error("renderFlamegraphUI should produce output")
	}

	// Should contain some commit data
	if len(output) < 10 {
		t.Error("renderFlamegraphUI output too short")
	}
}

func TestRenderTimelineSliderUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.showTimeline = true

	output := renderTimelineSliderUI(m, 80)

	if output == "" {
		t.Error("renderTimelineSliderUI should produce output")
	}
}

func TestRenderTreeViewUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.showTreeView = true

	output := renderTreeViewUI(m, 80)

	if output == "" {
		t.Error("renderTreeViewUI should produce output")
	}
}

func TestRenderAuthorComparisonUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeNamedCommits()
	m.showAuthorComparison = true

	output := renderAuthorComparisonUI(m, 80)

	if output == "" {
		t.Error("renderAuthorComparisonUI should produce output")
	}
}

func TestRenderFileHeatmapUI_ProducesOutput(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(5)
	m.showFileHeatmap = true

	output := renderFileHeatmapUI(m, 80)

	if output == "" {
		t.Error("renderFileHeatmapUI should produce output")
	}
}
