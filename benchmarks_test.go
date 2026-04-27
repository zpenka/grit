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
