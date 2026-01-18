package list

import (
	"slices"
	"sync"
)

type List[V any] struct {
	items []V
	mu    sync.RWMutex
}

func New[V any](capacity ...int) *List[V] {
	capSize := 1024
	if len(capacity) > 0 && capacity[0] > 0 {
		capSize = capacity[0]
	}
	return &List[V]{
		items: make([]V, 0, capSize),
	}
}

func (l *List[V]) Append(items ...V) *List[V] {
	if len(items) == 0 {
		return l
	}
	l.mu.Lock()
	l.items = append(l.items, items...)
	l.mu.Unlock()
	return l
}

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

func (l *List[V]) List() []V {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if len(l.items) == 0 {
		return nil
	}
	return slices.Clone(l.items)
}

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

func (l *List[V]) Count() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.items)
}

func (l *List[V]) Empty() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.items) == 0
}
