package grit

import (
	"testing"
)

func TestNewDiffCache_Valid(t *testing.T) {
	cache := newDiffCache(10)
	if cache == nil {
		t.Error("newDiffCache should return non-nil cache")
	}
}

func TestDiffCacheSet_Valid(t *testing.T) {
	cache := newDiffCache(10)
	lines := []diffLine{{kind: lineContext, text: "test"}}
	cache.set("key1", lines)
}

func TestDiffCacheGet_Valid(t *testing.T) {
	cache := newDiffCache(10)
	lines := []diffLine{{kind: lineContext, text: "test"}}
	cache.set("key1", lines)
	retrieved, ok := cache.get("key1")
	if !ok {
		t.Error("cache.get should find existing key")
	}
	if len(retrieved) != 1 {
		t.Errorf("expected 1 line, got %d", len(retrieved))
	}
}

func TestDiffCacheGetMissing_Valid(t *testing.T) {
	cache := newDiffCache(10)
	_, ok := cache.get("nonexistent")
	if ok {
		t.Error("cache.get should return false for missing key")
	}
}

func TestDiffCacheGetHitCount_Valid(t *testing.T) {
	cache := newDiffCache(10)
	count := cache.getHitCount()
	if count < 0 {
		t.Errorf("hit count should be non-negative, got %d", count)
	}
}

func TestNewStatCache_Valid(t *testing.T) {
	cache := newStatCache(10)
	if cache == nil {
		t.Error("newStatCache should return non-nil cache")
	}
}

func TestStatCacheGetOrCompute_Valid(t *testing.T) {
	cache := newStatCache(10)
	lines := []diffLine{{kind: lineAdded, text: "+ test"}}
	stats := cache.getOrCompute("key1", lines)
	if stats.insertions != 1 {
		t.Errorf("expected 1 insertion, got %d", stats.insertions)
	}
}

func TestStatCacheGetHitCount_Valid(t *testing.T) {
	cache := newStatCache(10)
	count := cache.getHitCount()
	if count < 0 {
		t.Errorf("hit count should be non-negative, got %d", count)
	}
}

func TestNewRegexCache_Valid(t *testing.T) {
	cache := newRegexCache(10)
	if cache == nil {
		t.Error("newRegexCache should return non-nil cache")
	}
}

func TestRegexCacheCompile_Valid(t *testing.T) {
	cache := newRegexCache(10)
	regex, err := cache.compile("test.*")
	if err != nil {
		t.Errorf("compile should not error: %v", err)
	}
	if regex == nil {
		t.Error("compile should return non-nil regex")
	}
}

func TestRegexCacheCompileInvalid_Valid(t *testing.T) {
	cache := newRegexCache(10)
	_, err := cache.compile("[invalid(")
	if err == nil {
		t.Error("compile should error for invalid regex")
	}
}

func TestCompileRegex_Valid(t *testing.T) {
	regex, err := compileRegex("test.*")
	if err != nil {
		t.Errorf("compileRegex should not error: %v", err)
	}
	if regex == nil {
		t.Error("compileRegex should return non-nil regex")
	}
}

func TestCompileRegexInvalid_Valid(t *testing.T) {
	_, err := compileRegex("[invalid(")
	if err == nil {
		t.Error("compileRegex should error for invalid regex")
	}
}

func TestLazyLoadDiff_Valid(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{{hash: "abc123", subject: "test"}}
	_ = lazyLoadDiff(m)
}

func TestLazyLoadGraph_Valid(t *testing.T) {
	m := newModel(".")
	m.commits = []commit{{hash: "abc123", subject: "test"}}
	_ = lazyLoadGraph(m)
}

func TestLazyLoadStats_Valid(t *testing.T) {
	m := newModel(".")
	stats := lazyLoadStats(m)
	if stats.insertions < 0 {
		t.Errorf("insertions should be non-negative, got %d", stats.insertions)
	}
}

func TestSafeIsFileModified_Valid(t *testing.T) {
	_ = safeIsFileModified("abc123", "file.go")
}

func TestSafeParseCommitGraph_Valid(t *testing.T) {
	commits := []commit{{hash: "abc123", subject: "test"}}
	nodes := safeParseCommitGraph(commits)
	if nodes == nil {
		t.Error("safeParseCommitGraph should return non-nil slice")
	}
}
