package esqueue

import (
	"testing"
	"time"
)

// TestEsQueue 测试高性能无锁队列的单次及批量入队、出队操作
func TestEsQueue(t *testing.T) {
	q := NewQueue[int](8, 1*time.Millisecond)
	if q.Capacity() != 8 {
		t.Errorf("期望容量为 8，实际结果为 %d", q.Capacity())
	}

	// Test Put
	ok, quantity := q.Put(1)
	if !ok || quantity != 1 {
		t.Errorf("Put 失败：状态为 %v，数量为 %d", ok, quantity)
	}

	// Test Get
	val, ok, quantity := q.Get()
	if !ok || quantity != 0 {
		t.Errorf("Get 失败：状态为 %v，数量为 %d", ok, quantity)
	}
	if *val.(*int) != 1 {
		t.Errorf("期望值为 1，实际结果为 %v", *val.(*int))
	}

	// Test Puts/Gets
	q.Puts([]int{1, 2, 3})
	if q.Quantity() != 3 {
		t.Errorf("期望数量为 3，实际结果为 %d", q.Quantity())
	}

	results := make([]any, 2)
	gets, rem := q.Gets(results)
	if gets != 2 || rem != 1 {
		t.Errorf("Gets 失败：获取数量为 %d，剩余数量为 %d", gets, rem)
	}
}
