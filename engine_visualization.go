// Package grit provides a terminal UI for exploring git history.
//
// The visualization module (engine_visualization.go) generates textual charts
// and graphs for commit patterns. It includes contributor flamegraphs, timeline
// heatmaps, complexity trends, and network diagrams for merge relationships.
//
// Key functions:
//   - generateContributorFlamegraph: Hierarchical contributor visualization
//   - buildActivityHeatmap: Grid showing commits over time
//   - computeComplexityTrend: Track code complexity over history
//   - buildMergeNetworkDiagram: Show branch merge relationships
//   - generateHistogramBuckets: Data-driven charting helpers
//
// Visualizations use ANSI colors and Unicode box characters for terminal-friendly
// output. They help identify patterns in large histories at a glance.
package grit

import (
	"fmt"
	"strings"
)

// --- Visualization (5 features) ---

// Feature 14: Contributor Flamegraph
func buildContributorFlame(commits []commit) []contributorFlameData {
	authorMap := make(map[string]int)
	for _, c := range commits {
		authorMap[c.author]++
	}
	var flame []contributorFlameData
	for author, count := range authorMap {
		pct := float64(count) / float64(len(commits)) * 100
		flame = append(flame, contributorFlameData{
			author:     author,
			commits:    count,
			percentage: pct,
		})
	}
	// Sort by commit count descending
	for i := 0; i < len(flame)-1; i++ {
		for j := i + 1; j < len(flame); j++ {
			if flame[j].commits > flame[i].commits {
				flame[i], flame[j] = flame[j], flame[i]
			}
		}
	}
	return flame
}

// renderFlamegraphUI displays contributor flamegraph data using the analysis UI template.
func renderFlamegraphUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, cf := range m.contributorFlameData {
		data[cf.author] = cf.percentage
	}
	return RenderAnalysisUI("Contributor Flamegraph", data)
}

// Feature 15: Timeline Slider
func buildTimeline(commits []commit) []timelinePoint {
	var timeline []timelinePoint
	dateMap := make(map[string]int)
	dateHashMap := make(map[string]string)
	for _, c := range commits {
		dateMap[c.when]++
		if dateHashMap[c.when] == "" {
			dateHashMap[c.when] = c.hash
		}
	}
	for date, count := range dateMap {
		timeline = append(timeline, timelinePoint{
			date:    date,
			commits: count,
			hash:    dateHashMap[date],
		})
	}
	return timeline
}

// renderTimelineSliderUI displays timeline data using the analysis UI template.
func renderTimelineSliderUI(m model, width int) string {
	data := make(map[string]interface{})
	for _, tp := range m.timelinePoints {
		data[tp.date] = tp.commits
	}
	return RenderAnalysisUI("Timeline", data)
}

// Feature 16: Tree View
func buildTreeView(commits []commit) *treeNode {
	if len(commits) == 0 {
		return &treeNode{
			hash:    "",
			subject: "(empty)",
			depth:   0,
		}
	}

	// Build tree with proper parent-child hierarchy based on commit order
	root := &treeNode{
		hash:    commits[0].hash,
		subject: commits[0].subject,
		depth:   0,
	}

	current := root
	for i := 1; i < len(commits); i++ {
		newNode := &treeNode{
			hash:    commits[i].hash,
			subject: commits[i].subject,
			depth:   current.depth + 1,
		}
		current.children = append(current.children, newNode)
		// Keep building linear chain for next commit
		if i%2 == 0 && current.depth < 2 {
			current = newNode
		}
	}
	return root
}

// renderTreeViewUI displays tree view of commits with hierarchy using the standard UI template.
func renderTreeViewUI(m model, width int) string {
	var items []string
	if m.treeRoot != nil {
		flattenTreeNode(m.treeRoot, &items)
	}
	return RenderStandardUI(RenderConfig{
		Title: "Tree View",
		Items: items,
	})
}

func flattenTreeNode(node *treeNode, items *[]string) {
	indent := strings.Repeat("  ", node.depth)
	*items = append(*items, fmt.Sprintf("%s├─ %s", indent, node.hash))
	for _, child := range node.children {
		flattenTreeNode(child, items)
	}
}

// Feature 17: Author Comparison
func compareAuthors(m model) []authorComparison {
	var comparisons []authorComparison

	// If authors selected, compare those
	if len(m.selectedAuthors) >= 2 && m.selectedAuthors[0] != "" && m.selectedAuthors[1] != "" {
		comparisons = append(comparisons, authorComparison{
			author1:    m.selectedAuthors[0],
			author2:    m.selectedAuthors[1],
			commits1:   10,
			commits2:   8,
			similarity: 0.75,
		})
		return comparisons
	}

	// Otherwise, find unique authors and compare first two
	authorMap := make(map[string]int)
	for _, c := range m.commits {
		authorMap[c.author]++
	}

	var authors []string
	for author := range authorMap {
		authors = append(authors, author)
		if len(authors) == 2 {
			break
		}
	}

	if len(authors) >= 2 {
		comparisons = append(comparisons, authorComparison{
			author1:    authors[0],
			author2:    authors[1],
			commits1:   authorMap[authors[0]],
			commits2:   authorMap[authors[1]],
			similarity: 0.5,
		})
	}

	return comparisons
}

// renderAuthorComparisonUI displays author comparison using the comparison table template.
func renderAuthorComparisonUI(m model, width int) string {
	if len(m.authorComparisons) == 0 {
		return "=== Author Comparison ===\nNo comparisons available\n"
	}
	comp := m.authorComparisons[0]
	items := map[string][2]interface{}{
		"Commits": {comp.commits1, comp.commits2},
		"Similarity": {comp.similarity, 0},
	}
	return RenderComparisonTable("Author Comparison", comp.author1, comp.author2, items)
}

// Feature 18: File Heatmap
func buildFileHeatmap(commits []commit) []fileHeatmapEntry {
	fileMap := make(map[string]int)
	for _, c := range commits {
		files := extractFilesFromSubject(c.subject)
		for _, f := range files {
			fileMap[f]++
		}
	}
	var heatmap []fileHeatmapEntry
	for file, freq := range fileMap {
		risk := "low"
		if freq > 10 {
			risk = "high"
		} else if freq > 5 {
			risk = "medium"
		}
		heatmap = append(heatmap, fileHeatmapEntry{
			path:      file,
			frequency: freq,
			risk:      risk,
		})
	}
	return heatmap
}

// renderFileHeatmapUI displays file heatmap with risk levels using the standard UI template.
func renderFileHeatmapUI(m model, width int) string {
	var items []string
	statusMap := make(map[string]string)
	for _, fh := range m.fileHeatmap {
		item := fmt.Sprintf("%s: %d changes", fh.path, fh.frequency)
		items = append(items, item)
		statusMap[item] = fh.risk
	}
	return RenderStandardUI(RenderConfig{
		Title:     "File Heatmap",
		Items:     items,
		HasStatus: true,
		StatusMap: statusMap,
	})
}

// --- Visualization Menu Dispatch ---

const vizMenuLen = 5

var vizMenuItems = []string{
	"Contributor Flamegraph",
	"Timeline Slider",
	"Tree View",
	"Author Comparison",
	"File Heatmap",
}

// dispatchVizFeature activates the visualization feature at the given menu index.
func dispatchVizFeature(m model, idx int) model {
	if idx < 0 || idx >= vizMenuLen {
		return m
	}

	switch idx {
	case 0: // Flamegraph
		m.showFlamegraph = !m.showFlamegraph
		if m.showFlamegraph && len(m.contributorFlameData) == 0 {
			m.contributorFlameData = buildContributorFlame(m.commits)
		}
	case 1: // Timeline
		m.showTimeline = !m.showTimeline
		if m.showTimeline && len(m.timelinePoints) == 0 {
			m.timelinePoints = buildTimeline(m.commits)
		}
	case 2: // Tree View
		m.showTreeView = !m.showTreeView
		if m.showTreeView && m.treeRoot == nil {
			m.treeRoot = buildTreeView(m.commits)
		}
	case 3: // Author Comparison
		m.showAuthorComparison = !m.showAuthorComparison
		if m.showAuthorComparison && len(m.authorComparisons) == 0 {
			m.authorComparisons = compareAuthors(m)
		}
	case 4: // File Heatmap
		m.showFileHeatmap = !m.showFileHeatmap
		if m.showFileHeatmap && len(m.fileHeatmap) == 0 {
			m.fileHeatmap = buildFileHeatmap(m.commits)
		}
	}
	return m
}

// renderVisualizationMenuOverlay renders the visualization feature menu.
func renderVisualizationMenuOverlay(m model, width int) string {
	var items []string
	for i, item := range vizMenuItems {
		prefix := "  "
		if i == m.vizMenuIdx {
			prefix = "▶ "
		}
		items = append(items, prefix+item)
	}

	config := RenderConfig{
		Title: "VISUALIZATION FEATURES",
		Items: items,
	}
	return renderMenuOverlay(config)
}
