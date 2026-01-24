package syncx

import (
	"fmt"
	"runtime/debug"
	"sync"
)

// GroupResultV3 增强型泛型结果结构体，包含共享标记
type GroupResultV3[T any] struct {
	Val    T
	Err    error
	Shared bool
}

// GroupV3 增强型泛型并发控制组，支持 Channel 返回
type GroupV3[T any] struct {
	mu sync.Mutex
	m  map[string]*callv3[T]
}

type callv3[T any] struct {
	wg    sync.WaitGroup
	res   GroupResultV3[T]
	dups  int
	chans []chan<- GroupResultV3[T]
}

// Do 执行带泛型返回值的任务，并返回结果是否为共享
func (g *GroupV3[T]) Do(key string, fn func() GroupResultV3[T]) GroupResultV3[T] {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*callv3[T])
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()
		c.res.Shared = c.dups > 0
		return c.res
	}
	c := new(callv3[T])
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	g.docallv3(c, key, fn)
	c.res.Shared = c.dups > 0
	return c.res
}

// DoChan 执行任务并通过 Channel 返回结果
func (g *GroupV3[T]) DoChan(key string, fn func() GroupResultV3[T]) <-chan GroupResultV3[T] {
	ch := make(chan GroupResultV3[T], 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*callv3[T])
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &callv3[T]{chans: []chan<- GroupResultV3[T]{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	go g.docallv3(c, key, fn)
	return ch
}

func (g *GroupV3[T]) docallv3(c *callv3[T], key string, fn func() GroupResultV3[T]) {
	var (
		normalReturn = false
		recovered    = false
	)
	defer func() {
		if !normalReturn && !recovered {
			c.res.Err = fmt.Errorf("runtime.Goexit callv3ed")
		}
		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c {
			delete(g.m, key)
		}
		res := GroupResultV3[T]{c.res.Val, c.res.Err, c.dups > 0}
		for _, ch := range c.chans {
			ch <- res
		}
	}()
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = true
				c.res.Err = fmt.Errorf("panic: %v\nstack: %s", r, debug.Stack())
			}
		}()
		c.res = fn()
		normalReturn = true
	}()
}

// Forget 手动删除指定 key 的缓存任务
func (g *GroupV3[T]) Forget(key string) {
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}
