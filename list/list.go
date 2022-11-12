package list

import "sync"

type list[V any] struct {
	v     V
	items []V
	size  int
	lock  sync.RWMutex
}

func New[V any](size ...int) *list[V] {
	s := 1024
	if len(size) > 0 {
		s = size[0]
	}
	return &list[V]{items: make([]V, 0, s), size: s}
}

func (l *list[V]) Append(items ...V) *list[V] {
	l.lock.Lock()
	l.items = append(l.items, items...)
	l.lock.Unlock()
	return l
}

/*
	func (l *list[V]) Append2(items ...V) {
		l.lock.Lock()
		defer l.lock.Unlock()
		ll := len(items)
		if ll > 0 {
			nodes := make([]V, 0, len(l.items)+ll)
			copy(nodes, l.items[:])
			copy(nodes[len(l.items):], items)
			l.items = nodes
		}
	}
*/
func (l *list[V]) Clear(size ...int) *list[V] {
	l.lock.Lock()
	if len(size) > 0 {
		l.size = size[0]
	} else {
		l.size = 1024
	}
	l.items = make([]V, 0, l.size)
	l.lock.Unlock()
	return l
}

func (l *list[V]) Count() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.items)
}

func (l *list[V]) Empty() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.items) == 0
}

func (l *list[V]) List() []V {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.items
}

// 入队列
func (l *list[V]) Put(v V) *list[V] {
	return l.Append(v)
}

// 出队列
func (l *list[V]) Get() (V, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(l.items) > 0 {
		item := l.items[0]
		l.items = l.items[1:len(l.items)]
		return item, true
	}
	return l.v, false
}
