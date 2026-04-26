package grit

import (
	"regexp"
)

// newDiffCache creates a new diff cache with LRU eviction.
func newDiffCache(size int) *diffCache {
	return &diffCache{
		data:    make(map[string][]diffLine),
		order:   []string{},
		maxSize: size,
	}
}

// set stores a diff in the cache with LRU eviction.
func (dc *diffCache) set(key string, lines []diffLine) {
	if _, exists := dc.data[key]; !exists {
		dc.order = append(dc.order, key)
		if len(dc.order) > dc.maxSize {
			oldest := dc.order[0]
			dc.order = dc.order[1:]
			delete(dc.data, oldest)
		}
	}
	dc.data[key] = lines
}

// get retrieves a diff from the cache and tracks hit count.
func (dc *diffCache) get(key string) ([]diffLine, bool) {
	lines, ok := dc.data[key]
	if ok {
		dc.hitCount++
	}
	return lines, ok
}

// getHitCount returns the number of cache hits.
func (dc *diffCache) getHitCount() int {
	return dc.hitCount
}

// newStatCache creates a new statistics cache.
func newStatCache(size int) *statCache {
	return &statCache{
		data:    make(map[string]commitStatistics),
		order:   []string{},
		maxSize: size,
	}
}

// getOrCompute gets cached stats or computes and caches them.
func (sc *statCache) getOrCompute(key string, lines []diffLine) commitStatistics {
	if stats, ok := sc.data[key]; ok {
		sc.hitCount++
		return stats
	}
	stats := commitStats(lines)
	sc.order = append(sc.order, key)
	if len(sc.order) > sc.maxSize {
		oldest := sc.order[0]
		sc.order = sc.order[1:]
		delete(sc.data, oldest)
	}
	sc.data[key] = stats
	return stats
}

// getHitCount returns the number of cache hits.
func (sc *statCache) getHitCount() int {
	return sc.hitCount
}

// newRegexCache creates a new regex pattern cache.
func newRegexCache(size int) *regexCache {
	return &regexCache{
		data:    make(map[string]*regexp.Regexp),
		maxSize: size,
	}
}

// compile compiles a regex pattern or returns cached version.
func (rc *regexCache) compile(pattern string) (*regexp.Regexp, error) {
	if re, ok := rc.data[pattern]; ok {
		rc.hitCount++
		return re, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if len(rc.data) < rc.maxSize {
		rc.data[pattern] = re
	}
	return re, nil
}

// compileRegex compiles a regex pattern for search.
func compileRegex(pattern string) (*regexp.Regexp, error) {
	return regexp.Compile(pattern)
}

// lazyLoadDiff loads diff asynchronously if not already loaded.
func lazyLoadDiff(m model) model {
	if m.cursor < len(m.commits) && len(m.diffLines) == 0 {
		m.loading = true
	}
	return m
}

// lazyLoadGraph builds commit graph on demand.
func lazyLoadGraph(m model) model {
	if len(m.commitGraph) == 0 && len(m.commits) > 0 {
		m.commitGraph = parseCommitGraph(m.commits)
	}
	return m
}

// lazyLoadStats computes commit statistics on demand.
func lazyLoadStats(m model) commitStatistics {
	return commitStats(m.diffLines)
}

// safeIsFileModified safely checks file modification without panicking.
func safeIsFileModified(hash, file string) bool {
	if hash == "" || file == "" {
		return false
	}
	return isFileModifiedInCommit(hash, file)
}

// safeParseCommitGraph safely parses graph, returning empty slice on error.
func safeParseCommitGraph(commits []commit) []graphNode {
	if commits == nil {
		return []graphNode{}
	}
	return parseCommitGraph(commits)
}
