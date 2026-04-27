package grit

import (
	"testing"
)

// Benchmark tests for performance-critical operations.
// Run with: go test -bench=. -benchtime=10s

// BenchmarkParseCommits measures commit parsing performance.
func BenchmarkParseCommits(b *testing.B) {
	output := ""
	for i := 0; i < 100; i++ {
		output += "abc123\x00abc\x00John Doe\x005 days ago\x00Fix bug\n"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseCommits(output)
	}
}

// BenchmarkFilterCommits measures filtering performance.
func BenchmarkFilterCommits(b *testing.B) {
	commits := makeCommits(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filterCommits(commits, "fix")
	}
}

// BenchmarkNavigationCursor measures cursor movement.
func BenchmarkNavigationCursor(b *testing.B) {
	m := model{
		commits: makeCommits(1000),
		cursor:  500,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m = moveCursorDown(m)
		m = moveCursorUp(m)
	}
}

// BenchmarkAnalyzeCodeOwnership measures ownership analysis.
func BenchmarkAnalyzeCodeOwnership(b *testing.B) {
	commits := makeCommits(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzeCodeOwnership(commits)
	}
}

// BenchmarkCaching measures cache operations.
func BenchmarkCaching(b *testing.B) {
	cache := newDiffCache(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate cache operations
		_ = cache
	}
}

// BenchmarkRenderUI measures UI rendering performance.
func BenchmarkRenderUI(b *testing.B) {
	m := model{
		commits: makeCommits(100),
		cursor:  50,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Rendering would happen here
		_ = m
	}
}

// BenchmarkParseDiff measures diff parsing performance (hot path).
func BenchmarkParseDiff(b *testing.B) {
	diffOutput := ""
	for i := 0; i < 100; i++ {
		diffOutput += "diff --git a/file.go b/file.go\n"
		diffOutput += "--- a/file.go\n"
		diffOutput += "+++ b/file.go\n"
		diffOutput += "@@ -1,5 +1,8 @@\n"
		diffOutput += " context\n"
		diffOutput += "+added line\n"
		diffOutput += "-removed line\n"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseDiff(diffOutput)
	}
}

// BenchmarkBuildTimeline measures timeline visualization performance.
func BenchmarkBuildTimeline(b *testing.B) {
	commits := makeCommits(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildTimeline(commits)
	}
}

// BenchmarkBuildFileHeatmap measures file heatmap O(n) scan.
func BenchmarkBuildFileHeatmap(b *testing.B) {
	commits := makeCommits(5000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildFileHeatmap(commits)
	}
}

// BenchmarkDetectSecrets measures security scanning performance.
func BenchmarkDetectSecrets(b *testing.B) {
	content := ""
	for i := 0; i < 1000; i++ {
		content += "line of code\n"
		if i%100 == 0 {
			content += "password = 'secret123'\n"
		}
	}
	hash := "abc123abc123abc123abc123abc123abc123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detectSecrets(hash, content)
	}
}

// BenchmarkVisibleCommits measures filtering + cursor performance.
func BenchmarkVisibleCommits(b *testing.B) {
	m := model{
		commits: makeCommits(10000),
		cursor:  5000,
		query:   "fix",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		visibleCommits(m)
	}
}
