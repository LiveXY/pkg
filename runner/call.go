package runner

import (
	"reflect"
	"sync"
	"time"

	"github.com/livexy/pkg/util"
)

var callmap sync.Map

type call struct {
	ticker *time.Ticker
	fnmap  sync.Map
}

type fn struct {
	fn     any
	params []any
	sum    int64
}

func (c *call) start() {
	for {
		<-c.ticker.C
		c.do()
	}
}

func (c *call) do() {
	c.fnmap.Range(func(k, v any) bool {
		val := v.(*fn)
		if val.sum > 0 {
			param := val.params
			param = append(param, val.sum)
			reflectcall(val.fn, param...)
		} else {
			reflectcall(val.fn, val.params...)
		}
		c.fnmap.Delete(k)
		return true
	})
}

func (c *call) stop() {
	c.ticker.Stop()
	c.do()
}

func (c *call) call(name string, f any, ps ...any) {
	key := name + ":" + util.ToStr(ps...)
	if _, ok := c.fnmap.Load(key); !ok {
		c.fnmap.Store(key, &fn{fn: f, params: ps})
	}
}

func (c *call) sum(name string, f any, val int64, ps ...any) {
	key := name + ":" + util.ToStr(ps...)
	var ff *fn
	obj, ok := c.fnmap.Load(key)
	if !ok {
		ff = &fn{fn: f, params: ps, sum:0 }
	} else {
		ff = obj.(*fn)
	}
	ff.sum += val
	c.fnmap.Store(key, ff)
}

func newcall(second int) *call {
	c := &call{}
	c.ticker = time.NewTicker(time.Duration(second) * time.Second)
	go c.start()
	return c
}

func TickerCall(second int, name string, fn any, params ...any) {
	var c *call
	obj, ok := callmap.Load(second)
	if ok {
		c = obj.(*call)
	} else {
		c = newcall(second)
		callmap.Store(second, c)
	}
	c.call(name, fn, params...)
}

func TickerSum(second int, name string, fn any, val int64, params ...any) {
	var c *call
	obj, ok := callmap.Load(second)
	if ok {
		c = obj.(*call)
	} else {
		c = newcall(second)
		callmap.Store(second, c)
	}
	c.sum(name, fn, val, params...)
}

func TickerStop(seconds ...int) {
	if len(seconds) == 0 {
		callmap.Range(func(k, v any) bool {
			v.(*call).stop()
			callmap.Delete(k)
			return true
		})
	} else {
		for _, second := range seconds {
			if obj, ok := callmap.Load(second); ok {
				obj.(*call).stop()
				callmap.Delete(second)
			}
		}
	}
}

func reflectcall(f any, params ...any) {
	fun := reflect.ValueOf(f)
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	fun.Call(in)
}
