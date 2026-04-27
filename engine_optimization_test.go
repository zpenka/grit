package grit

import (
	"testing"
)

// TestCacheMetrics_GetHitRate tests hit rate calculation
func TestCacheMetrics_GetHitRate_NoActivity(t *testing.T) {
	m := CacheMetrics{Hits: 0, Misses: 0}
	rate := m.GetHitRate()
	if rate != 0 {
		t.Errorf("expected 0 hit rate for no activity, got %f", rate)
	}
}

func TestCacheMetrics_GetHitRate_AllHits(t *testing.T) {
	m := CacheMetrics{Hits: 10, Misses: 0}
	rate := m.GetHitRate()
	if rate != 1.0 {
		t.Errorf("expected 1.0 hit rate for all hits, got %f", rate)
	}
}

func TestCacheMetrics_GetHitRate_AllMisses(t *testing.T) {
	m := CacheMetrics{Hits: 0, Misses: 10}
	rate := m.GetHitRate()
	if rate != 0 {
		t.Errorf("expected 0 hit rate for all misses, got %f", rate)
	}
}

func TestCacheMetrics_GetHitRate_Mixed(t *testing.T) {
	m := CacheMetrics{Hits: 7, Misses: 3}
	rate := m.GetHitRate()
	expected := 0.7
	if rate < expected-0.01 || rate > expected+0.01 {
		t.Errorf("expected ~0.7 hit rate, got %f", rate)
	}
}

// TestCacheMetrics_ShouldEvict tests eviction decision
func TestCacheMetrics_ShouldEvict_BelowCapacity(t *testing.T) {
	m := CacheMetrics{Size: 5, MaxSize: 10}
	if m.ShouldEvict() {
		t.Error("should not evict when below capacity")
	}
}

func TestCacheMetrics_ShouldEvict_AtCapacity(t *testing.T) {
	m := CacheMetrics{Size: 10, MaxSize: 10}
	if !m.ShouldEvict() {
		t.Error("should evict when at capacity")
	}
}

func TestCacheMetrics_ShouldEvict_OverCapacity(t *testing.T) {
	m := CacheMetrics{Size: 15, MaxSize: 10}
	if !m.ShouldEvict() {
		t.Error("should evict when over capacity")
	}
}

// TestLazyLoader_Load tests lazy initialization
func TestLazyLoader_LoadExecutes(t *testing.T) {
	called := false
	loader := NewLazyLoader(func() interface{} {
		called = true
		return "data"
	})

	if called {
		t.Error("loader function should not be called on creation")
	}

	result := loader.Load()
	if !called {
		t.Error("loader function should be called on Load()")
	}
	if result != "data" {
		t.Errorf("expected 'data', got %v", result)
	}
}

func TestLazyLoader_LoadCaches(t *testing.T) {
	callCount := 0
	loader := NewLazyLoader(func() interface{} {
		callCount++
		return "data"
	})

	loader.Load()
	loader.Load()
	loader.Load()

	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestLazyLoader_IsLoaded_InitialState(t *testing.T) {
	loader := NewLazyLoader(func() interface{} {
		return "data"
	})

	if loader.IsLoaded() {
		t.Error("should not be loaded initially")
	}
}

func TestLazyLoader_IsLoaded_AfterLoad(t *testing.T) {
	loader := NewLazyLoader(func() interface{} {
		return "data"
	})
	loader.Load()

	if !loader.IsLoaded() {
		t.Error("should be loaded after Load()")
	}
}

func TestLazyLoader_Reset(t *testing.T) {
	loader := NewLazyLoader(func() interface{} {
		return "data"
	})
	loader.Load()

	if !loader.IsLoaded() {
		t.Error("should be loaded before reset")
	}

	loader.Reset()

	if loader.IsLoaded() {
		t.Error("should not be loaded after Reset()")
	}
}

func TestLazyLoader_NilFunction(t *testing.T) {
	loader := NewLazyLoader(nil)
	result := loader.Load()

	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

// TestMemoryPool_Get tests getting from pool
func TestMemoryPool_GetFromEmptyPool(t *testing.T) {
	pool := NewMemoryPool(1)
	result := pool.Get(func() interface{} {
		return "new_item"
	})

	if result != "new_item" {
		t.Errorf("expected 'new_item', got %v", result)
	}
}

func TestMemoryPool_GetFromPool(t *testing.T) {
	pool := NewMemoryPool(1)
	pool.Put("pooled_item")
	result := pool.Get(func() interface{} {
		return "new_item"
	})

	if result != "pooled_item" {
		t.Errorf("expected 'pooled_item' from pool, got %v", result)
	}
}

func TestMemoryPool_GetMultipleItems(t *testing.T) {
	pool := NewMemoryPool(2)
	pool.Put("item1")
	pool.Put("item2")

	item1 := pool.Get(func() interface{} { return "new" })
	item2 := pool.Get(func() interface{} { return "new" })

	if item1 != "item1" || item2 != "item2" {
		t.Error("should retrieve pooled items in order")
	}
}

// TestMemoryPool_Put tests putting into pool
func TestMemoryPool_PutIntoPool(t *testing.T) {
	pool := NewMemoryPool(1)
	pool.Put("item1")

	// Get should retrieve the item
	result := pool.Get(func() interface{} {
		return "new"
	})

	if result != "item1" {
		t.Errorf("expected pooled 'item1', got %v", result)
	}
}

func TestMemoryPool_PutFullPool(t *testing.T) {
	pool := NewMemoryPool(1)
	pool.Put("item1")

	// This should be discarded as pool is full
	pool.Put("item2")

	result := pool.Get(func() interface{} {
		return "new"
	})

	if result != "item1" {
		t.Errorf("expected 'item1', got %v", result)
	}
}

// TestBatchProcessor_Add tests adding items
func TestBatchProcessor_AddBelowSize(t *testing.T) {
	bp := NewBatchProcessor(3, 100)

	result := bp.Add("item1")
	if result {
		t.Error("should not be full after 1 item when size is 3")
	}
}

func TestBatchProcessor_AddAtSize(t *testing.T) {
	bp := NewBatchProcessor(3, 100)

	bp.Add("item1")
	bp.Add("item2")
	result := bp.Add("item3")

	if !result {
		t.Error("should be full when batch reaches size")
	}
}

func TestBatchProcessor_AddOverSize(t *testing.T) {
	bp := NewBatchProcessor(2, 100)

	bp.Add("item1")
	bp.Add("item2")
	result := bp.Add("item3")

	if !result {
		t.Error("should indicate full when exceeding size")
	}
}

// TestBatchProcessor_IsFull tests full check
func TestBatchProcessor_IsFullEmpty(t *testing.T) {
	bp := NewBatchProcessor(3, 100)

	if bp.IsFull() {
		t.Error("should not be full when empty")
	}
}

func TestBatchProcessor_IsFullPartial(t *testing.T) {
	bp := NewBatchProcessor(3, 100)

	bp.Add("item1")
	if bp.IsFull() {
		t.Error("should not be full when partially filled")
	}
}

func TestBatchProcessor_IsFullComplete(t *testing.T) {
	bp := NewBatchProcessor(2, 100)

	bp.Add("item1")
	bp.Add("item2")
	if !bp.IsFull() {
		t.Error("should be full when reaching batch size")
	}
}

// TestBatchProcessor_Get tests retrieving and resetting
func TestBatchProcessor_GetItems(t *testing.T) {
	bp := NewBatchProcessor(2, 100)

	bp.Add("item1")
	bp.Add("item2")
	items := bp.Get()

	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if items[0] != "item1" || items[1] != "item2" {
		t.Error("items should be in order")
	}
}

func TestBatchProcessor_GetResets(t *testing.T) {
	bp := NewBatchProcessor(2, 100)

	bp.Add("item1")
	bp.Add("item2")
	bp.Get()

	if bp.IsFull() {
		t.Error("should not be full after Get()")
	}
	if len(bp.Get()) != 0 {
		t.Error("should have no items after Get()")
	}
}

// TestCircularBuffer_Push tests adding items
func TestCircularBuffer_PushSingleItem(t *testing.T) {
	cb := NewCircularBuffer(3)
	cb.Push("item1")

	if cb.count != 1 {
		t.Errorf("expected count 1, got %d", cb.count)
	}
}

func TestCircularBuffer_PushMultiple(t *testing.T) {
	cb := NewCircularBuffer(3)
	cb.Push("item1")
	cb.Push("item2")
	cb.Push("item3")

	if cb.count != 3 {
		t.Errorf("expected count 3, got %d", cb.count)
	}
}

func TestCircularBuffer_PushOverflow(t *testing.T) {
	cb := NewCircularBuffer(2)
	cb.Push("item1")
	cb.Push("item2")
	cb.Push("item3") // Should overflow

	// Buffer should still report max size
	if cb.count > 2 {
		t.Errorf("buffer should not exceed size 2, got %d", cb.count)
	}
}

// TestCircularBuffer_GetAll tests retrieving items
func TestCircularBuffer_GetAllEmpty(t *testing.T) {
	cb := NewCircularBuffer(3)
	items := cb.GetAll()

	if len(items) != 0 {
		t.Errorf("expected empty slice, got %d items", len(items))
	}
}

func TestCircularBuffer_GetAllFilled(t *testing.T) {
	cb := NewCircularBuffer(3)
	cb.Push("item1")
	cb.Push("item2")
	cb.Push("item3")

	items := cb.GetAll()
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}

// TestCircularBuffer_Size tests size tracking
func TestCircularBuffer_SizeEmpty(t *testing.T) {
	cb := NewCircularBuffer(5)
	if cb.Size() != 0 {
		t.Errorf("expected size 0 for empty buffer, got %d", cb.Size())
	}
}

func TestCircularBuffer_SizePartial(t *testing.T) {
	cb := NewCircularBuffer(5)
	cb.Push("item1")
	cb.Push("item2")

	if cb.Size() != 2 {
		t.Errorf("expected size 2, got %d", cb.Size())
	}
}

func TestCircularBuffer_SizeFull(t *testing.T) {
	cb := NewCircularBuffer(3)
	cb.Push("item1")
	cb.Push("item2")
	cb.Push("item3")

	if cb.Size() != 3 {
		t.Errorf("expected size 3, got %d", cb.Size())
	}
}

// TestRateLimiter_Allow tests rate limiting
func TestRateLimiter_AllowInitial(t *testing.T) {
	rl := NewRateLimiter(1)

	if !rl.Allow() {
		t.Error("should allow first request")
	}
}

func TestRateLimiter_AllowExceeded(t *testing.T) {
	rl := NewRateLimiter(1)
	rl.Allow()

	// Second request immediately should be denied
	if rl.Allow() {
		t.Error("should not allow exceeding rate limit")
	}
}

// TestMetrics_RecordOperation tests recording operations
func TestMetrics_RecordSuccessfulOperation(t *testing.T) {
	m := NewMetrics()
	m.RecordOperation("test_op", true)

	rate := m.GetSuccessRate("test_op")
	if rate != 1.0 {
		t.Errorf("expected 1.0 success rate, got %f", rate)
	}
}

func TestMetrics_RecordFailedOperation(t *testing.T) {
	m := NewMetrics()
	m.RecordOperation("test_op", false)

	rate := m.GetSuccessRate("test_op")
	if rate != 0 {
		t.Errorf("expected 0 success rate, got %f", rate)
	}
}

func TestMetrics_RecordMixed(t *testing.T) {
	m := NewMetrics()
	m.RecordOperation("test_op", true)
	m.RecordOperation("test_op", true)
	m.RecordOperation("test_op", false)

	rate := m.GetSuccessRate("test_op")
	expected := 2.0 / 3.0
	if rate < expected-0.01 || rate > expected+0.01 {
		t.Errorf("expected ~0.67 success rate, got %f", rate)
	}
}

func TestMetrics_UnknownOperation(t *testing.T) {
	m := NewMetrics()

	rate := m.GetSuccessRate("unknown_op")
	if rate != 0 {
		t.Errorf("expected 0 for unknown operation, got %f", rate)
	}
}

func TestMetrics_MultipleOperations(t *testing.T) {
	m := NewMetrics()

	m.RecordOperation("op1", true)
	m.RecordOperation("op2", false)
	m.RecordOperation("op1", false)

	rate1 := m.GetSuccessRate("op1")
	rate2 := m.GetSuccessRate("op2")

	if rate1 != 0.5 {
		t.Errorf("expected 0.5 for op1, got %f", rate1)
	}
	if rate2 != 0 {
		t.Errorf("expected 0 for op2, got %f", rate2)
	}
}
