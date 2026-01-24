package list

import (
	"slices"
	"sync"
)

// List 线程安全的切片容器，支持并发 Append 和 Get (FIFO)
type List[V any] struct {
	items []V
	mu    sync.RWMutex
}

// New 创建一个新的 List 实例，可指定初始容量
func New[V any](capacity ...int) *List[V] {
	capSize := 1024
	if len(capacity) > 0 && capacity[0] > 0 {
		capSize = capacity[0]
	}
	return &List[V]{
		items: make([]V, 0, capSize),
	}
}

// Append 添加元素到列表末尾
func (l *List[V]) Append(items ...V) *List[V] {
	if len(items) == 0 {
		return l
	}
	l.mu.Lock()
	l.items = append(l.items, items...)
	l.mu.Unlock()
	return l
}

// Get 从列表头部取出一个元素并移除它
func (l *List[V]) Get() (V, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.items) == 0 {
		var zero V
		return zero, false
	}
	item := l.items[0]
	var zero V
	l.items[0] = zero
	l.items = l.items[1:]
	if len(l.items) == 0 && cap(l.items) > 2048 {
		l.items = make([]V, 0, 1024)
	}
	return item, true
}

// List 返回当前所有元素的副本
func (l *List[V]) List() []V {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if len(l.items) == 0 {
		return nil
	}
	return slices.Clone(l.items)
}

// Clear 清空列表
func (l *List[V]) Clear(capacity ...int) *List[V] {
	capSize := 1024
	if len(capacity) > 0 {
		capSize = capacity[0]
	}
	l.mu.Lock()
	clear(l.items)
	l.items = make([]V, 0, capSize)
	l.mu.Unlock()
	return l
}

// Count 获取当前元素数量
func (l *List[V]) Count() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.items)
}

// Empty 检查列表是否为空
func (l *List[V]) Empty() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.items) == 0
}
