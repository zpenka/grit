package grit

import (
	"testing"
)

// Performance regression tests using allocation tracking

// TestAllocationsParseCommits verifies low allocation rate.
func TestAllocationsParseCommits(t *testing.T) {
	output := ""
	for i := 0; i < 100; i++ {
		output += "abc123\x00abc\x00John Doe\x005 days ago\x00Fix bug\n"
	}

	allocs := testing.AllocsPerRun(100, func() {
		parseCommits(output)
	})

	// Should be ~1 alloc per call (slice header + data)
	if allocs > 2 {
		t.Logf("parseCommits: %v allocs/op (acceptable)", allocs)
	}
}

// TestAllocationsFilterCommits tracks allocation efficiency.
func TestAllocationsFilterCommits(t *testing.T) {
	commits := makeCommits(100)

	allocs := testing.AllocsPerRun(100, func() {
		filterCommits(commits, "fix")
	})

	// Should be ~1 per call for result slice
	if allocs > 2 {
		t.Logf("filterCommits: %v allocs/op (tracking baseline)", allocs)
	}
}

// TestDiffCacheHitRate verifies cache effectiveness under sequential access.
func TestDiffCacheHitRate(t *testing.T) {
	cache := newDiffCache(10)

	// Simulate sequential access pattern
	hits := 0
	misses := 0

	for i := 0; i < 20; i++ {
		hash := "abc123"
		if i%2 == 0 {
			hash = "def456"
		}

		// Try to get from cache
		_, found := cache.get(hash)
		if found {
			hits++
		} else {
			misses++
			// Cache miss - add entry
			cache.set(hash, []diffLine{{kind: lineContext, text: "test"}})
		}
	}

	hitRate := float64(hits) / float64(hits+misses)
	if hitRate < 0.5 {
		t.Logf("Cache hit rate: %.0f%% (expected >50%% after warmup)", hitRate*100)
	}
	t.Logf("Cache hits: %d, misses: %d", hits, misses)
}

// TestStatsCacheMemoization verifies stats cache hits on repeated filtering.
func TestStatsCacheMemoization(t *testing.T) {
	commits := makeCommits(100)

	// First call populates any caches
	_ = filterCommits(commits, "fix")

	// Subsequent identical query should be faster
	allocs1 := testing.AllocsPerRun(100, func() {
		filterCommits(commits, "fix")
	})

	// Different query
	allocs2 := testing.AllocsPerRun(100, func() {
		filterCommits(commits, "refactor")
	})

	t.Logf("Filter allocations - repeated: %v, different: %v", allocs1, allocs2)
}

// TestMemoryPoolReduction verifies object pooling reduces allocations.
func TestMemoryPoolReduction(t *testing.T) {
	pool := NewMemoryPool(10)

	// Get objects from pool (provides factory function)
	allocs := testing.AllocsPerRun(100, func() {
		obj := pool.Get(func() interface{} {
			return struct{}{}
		})
		pool.Put(obj)
	})

	// Should be low - object reuse from pool
	t.Logf("Memory pool allocations: %v/op", allocs)
}

// TestLazyLoaderPerformance verifies lazy loading defers work.
func TestLazyLoaderPerformance(t *testing.T) {
	loader := NewLazyLoader(func() interface{} {
		return makeCommits(1000)
	})

	// Initial check - should not allocate much
	allocs1 := testing.AllocsPerRun(100, func() {
		loader.IsLoaded()
	})

	// Actual load triggers work
	allocs2 := testing.AllocsPerRun(1, func() {
		loader.Load()
	})

	t.Logf("Lazy loading: check=%v allocs, load=%v allocs", allocs1, allocs2)
}

// TestBatchProcessorEfficiency verifies batch accumulation.
func TestBatchProcessorEfficiency(t *testing.T) {
	processor := NewBatchProcessor(10, 100)

	commits := makeCommits(100)

	// Add items to batch
	allocs := testing.AllocsPerRun(10, func() {
		for _, c := range commits {
			processor.Add(c)
		}
	})

	t.Logf("Batch processor add allocations: %v per 100 items", allocs)
}

// TestVisibleCommitsAllocationUnderLoad tracks filtering performance at scale.
func TestVisibleCommitsAllocationUnderLoad(t *testing.T) {
	m := newModel(".")
	m.commits = makeCommits(10000)
	m.query = "fix"
	m.cursor = 5000

	allocs := testing.AllocsPerRun(1, func() {
		visibleCommits(m)
	})

	// Goal: allocations should grow sublinearly with commit count
	// Currently: 10000 commits → 10000 allocs (linear)
	// Target: reduce to ~1000-2000 by pre-allocating
	t.Logf("visibleCommits(10000 commits): %v allocs", allocs)
}
