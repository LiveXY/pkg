package runner

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// CallGroup 调度任务组，按指定间隔批量执行任务
type CallGroup struct {
	ticker *time.Ticker
	fnMap  sync.Map
	ctx    context.Context
	cancel context.CancelFunc
}

type taskNode struct {
	fn     any
	params []any
	sum    atomic.Int64
}

var globalCallers sync.Map

// TickerCall 注册一个定时调用的任务
func TickerCall(second int, name string, fn any, params ...any) {
	c := getOrNewCaller(second)
	key := fmt.Sprintf("%s:%v", name, params)
	c.fnMap.LoadOrStore(key, &taskNode{fn: fn, params: params})
}

// TickerSum 注册一个带累加参数的定时调用任务
func TickerSum(second int, name string, fn any, val int64, params ...any) {
	c := getOrNewCaller(second)
	key := fmt.Sprintf("%s:%v", name, params)
	actual, _ := c.fnMap.LoadOrStore(key, &taskNode{fn: fn, params: params})
	node := actual.(*taskNode)
	node.sum.Add(val)
}

func getOrNewCaller(second int) *CallGroup {
	if obj, ok := globalCallers.Load(second); ok {
		return obj.(*CallGroup)
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &CallGroup{
		ticker: time.NewTicker(time.Duration(second) * time.Second),
		ctx:    ctx,
		cancel: cancel,
	}
	actual, loaded := globalCallers.LoadOrStore(second, c)
	if !loaded {
		go c.run()
	} else {
		cancel()
		return actual.(*CallGroup)
	}
	return c
}

func (c *CallGroup) run() {
	defer c.ticker.Stop()
	for {
		select {
		case <-c.ticker.C:
			c.execute()
		case <-c.ctx.Done():
			c.execute()
			return
		}
	}
}

func (c *CallGroup) execute() {
	c.fnMap.Range(func(k, v any) bool {
		c.fnMap.Delete(k)
		node := v.(*taskNode)
		s := node.sum.Load()
		params := node.params
		if s > 0 {
			params = append(params, s)
		}
		go safeReflectCall(node.fn, params...)
		return true
	})
}

func (c *CallGroup) stop() {
	c.cancel()
}

// TickerStop 停止指定间隔或全部定时任务
func TickerStop(seconds ...int) {
	if len(seconds) == 0 {
		globalCallers.Range(func(k, v any) bool {
			v.(*CallGroup).stop()
			globalCallers.Delete(k)
			return true
		})
	} else {
		for _, s := range seconds {
			if obj, ok := globalCallers.Load(s); ok {
				obj.(*CallGroup).stop()
				globalCallers.Delete(s)
			}
		}
	}
}

func safeReflectCall(f any, params ...any) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("runner reflect call panic: %v\n", r)
		}
	}()
	v := reflect.ValueOf(f)
	in := make([]reflect.Value, len(params))
	for i, p := range params {
		in[i] = reflect.ValueOf(p)
	}
	v.Call(in)
}
