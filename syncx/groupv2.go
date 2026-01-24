package syncx

import (
	"fmt"
	"runtime/debug"
	"sync"
)

// GroupResultV2 泛型结果结构体
type GroupResultV2[T any] struct {
	Val T
	Err error
}

type callv2[T any] struct {
	wg  sync.WaitGroup
	res GroupResultV2[T]
}

// GroupV2 泛型并发控制组
type GroupV2[T any] struct {
	mu sync.Mutex
	m  map[string]*callv2[T]
}

// Do 执行带泛型返回值的任务，确保相同 key 的并发调用只执行一次
func (g *GroupV2[T]) Do(key string, fn func() GroupResultV2[T]) GroupResultV2[T] {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*callv2[T])
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.res
	}

	c := new(callv2[T])
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.docallv2(c, key, fn)
	return c.res
}

func (g *GroupV2[T]) docallv2(c *callv2[T], key string, fn func() GroupResultV2[T]) {
	defer func() {
		g.mu.Lock()
		delete(g.m, key)
		g.mu.Unlock()
		c.wg.Done()
	}()
	func() {
		defer func() {
			if r := recover(); r != nil {
				c.res.Err = fmt.Errorf("panic: %v\nstack: %s", r, debug.Stack())
			}
		}()
		c.res = fn()
	}()
}
