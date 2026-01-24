package syncx

import (
	"fmt"
	"sync"

	"errors"
)

// Group 并发控制组，确保相同 key 的任务只执行一次
type Group struct {
	m  map[string]*call
	mu sync.Mutex
}

type call struct {
	err  error
	wg   sync.WaitGroup
	dups int
}

// Shared 执行任务并返回是否为共享结果
func (g *Group) Shared(key string, fn func() error) (error, bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()
		return c.err, true
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.err, false
}

// Do 执行任务，确保相同 key 的并发调用只执行一次
func (g *Group) Do(key string, fn func() error) error {
	err, _ := g.Shared(key, fn)
	return err
}

func (g *Group) doCall(c *call, key string, fn func() error) {
	defer func() {
		c.wg.Done()
		g.mu.Lock()
		defer g.mu.Unlock()
		delete(g.m, key)
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				switch v := err.(type) {
				case error:
					c.err = fmt.Errorf("%w", v)
				case string:
					c.err = errors.New(v)
				default:
					c.err = fmt.Errorf("%v", v)
				}
			}
		}()
		c.err = fn()
	}()
}
