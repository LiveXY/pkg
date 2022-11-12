package syncx

import (
	"sync"

	"github.com/pkg/errors"
)

type Group struct {
	m  map[string]*call
	mu sync.Mutex
}
type call struct {
	err  error
	wg   sync.WaitGroup
	dups int
}

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
					c.err = errors.WithStack(v)
				case string:
					c.err = errors.New(v)
				default:
				}
			}
		}()
		c.err = fn()
	}()
}
