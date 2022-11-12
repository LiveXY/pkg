package esqueue

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

type esCache[V any] struct {
	value *V
	putno uint64
	getno uint64
}

type EsQueue[V any] struct {
	cache     []esCache[V]
	sleep time.Duration
	capacity  uint64
	capmod    uint64
	putpos    uint64
	getpos    uint64
}

func NewQueue[V any](capacity uint64, sleep time.Duration) *EsQueue[V] {
	q := new(EsQueue[V])
	q.capacity = minQuantity(capacity)
	q.capmod = q.capacity - 1
	q.putpos = 0
	q.getpos = 0
	q.sleep = sleep
	q.cache = make([]esCache[V], q.capacity)
	for i := range q.cache {
		cache := &q.cache[i]
		cache.getno = uint64(i)
		cache.putno = uint64(i)
	}
	cache := &q.cache[0]
	cache.getno = q.capacity
	cache.putno = q.capacity
	return q
}

func (q *EsQueue[V]) String() string {
	getpos := atomic.LoadUint64(&q.getpos)
	putpos := atomic.LoadUint64(&q.putpos)
	return fmt.Sprintf("Queue{capacity: %v, capmod: %v, putpos: %v, getpos: %v}", q.capacity, q.capmod, putpos, getpos)
}

func (q *EsQueue[V]) Capacity() uint64 {
	return q.capacity
}

func (q *EsQueue[V]) Quantity() uint64 {
	var putpos, getpos uint64
	var quantity uint64
	getpos = atomic.LoadUint64(&q.getpos)
	putpos = atomic.LoadUint64(&q.putpos)
	if putpos >= getpos {
		quantity = putpos - getpos
	} else {
		quantity = 0
	}
	return quantity
}

func (q *EsQueue[V]) Put(val V) (ok bool, quantity uint64) {
	var putpos, putposnew, getpos, posCnt uint64
	var cache *esCache[V]
	getpos = atomic.LoadUint64(&q.getpos)
	putpos = atomic.LoadUint64(&q.putpos)
	if putpos >= getpos {
		posCnt = putpos - getpos
	} else {
		posCnt = 0
	}
	if posCnt >= q.capacity {
		time.Sleep(q.sleep)
		return false, posCnt
	}
	putposnew = putpos + 1
	if !atomic.CompareAndSwapUint64(&q.putpos, putpos, putposnew) {
		runtime.Gosched()
		return false, posCnt
	}
	cache = &q.cache[putposnew&q.capmod]
	for {
		getno := atomic.LoadUint64(&cache.getno)
		putno := atomic.LoadUint64(&cache.putno)
		if putposnew == putno && getno == putno {
			cache.value = &val
			atomic.AddUint64(&cache.putno, q.capacity)
			return true, posCnt + 1
		} else {
			runtime.Gosched()
		}
	}
}

func (q *EsQueue[V]) Puts(values []V) (puts, quantity uint64) {
	var putpos, putposnew, getpos, posCnt, putCnt uint64
	getpos = atomic.LoadUint64(&q.getpos)
	putpos = atomic.LoadUint64(&q.putpos)
	if putpos >= getpos {
		posCnt = putpos - getpos
	} else {
		posCnt = 0
	}
	if posCnt >= q.capacity {
		time.Sleep(q.sleep)
		return 0, posCnt
	}
	if capputs, size := q.capacity-posCnt, uint64(len(values)); capputs >= size {
		putCnt = size
	} else {
		putCnt = capputs
	}
	putposnew = putpos + putCnt
	if !atomic.CompareAndSwapUint64(&q.putpos, putpos, putposnew) {
		runtime.Gosched()
		return 0, posCnt
	}
	for posnew, v := putpos+1, uint64(0); v < putCnt; posnew, v = posnew+1, v+1 {
		var cache = &q.cache[posnew&q.capmod]
		for {
			getno := atomic.LoadUint64(&cache.getno)
			putno := atomic.LoadUint64(&cache.putno)
			if posnew == putno && getno == putno {
				cache.value = &values[v]
				atomic.AddUint64(&cache.putno, q.capacity)
				break
			} else {
				runtime.Gosched()
			}
		}
	}
	return putCnt, posCnt + putCnt
}

func (q *EsQueue[V]) Get() (val any, ok bool, quantity uint64) {
	var putpos, getpos, getposnew, posCnt uint64
	var cache *esCache[V]
	putpos = atomic.LoadUint64(&q.putpos)
	getpos = atomic.LoadUint64(&q.getpos)
	if putpos >= getpos {
		posCnt = putpos - getpos
	} else {
		posCnt = 0
	}
	if posCnt < 1 {
		time.Sleep(q.sleep)
		return nil, false, posCnt
	}
	getposnew = getpos + 1
	if !atomic.CompareAndSwapUint64(&q.getpos, getpos, getposnew) {
		runtime.Gosched()
		return nil, false, posCnt
	}
	cache = &q.cache[getposnew&q.capmod]
	for {
		getno := atomic.LoadUint64(&cache.getno)
		putno := atomic.LoadUint64(&cache.putno)
		if getposnew == getno && getno == putno-q.capacity {
			val = cache.value
			cache.value = nil
			atomic.AddUint64(&cache.getno, q.capacity)
			return val, true, posCnt - 1
		} else {
			runtime.Gosched()
		}
	}
}

func (q *EsQueue[V]) Gets(values []any) (gets, quantity uint64) {
	var putpos, getpos, getposnew, posCnt, getCnt uint64
	putpos = atomic.LoadUint64(&q.putpos)
	getpos = atomic.LoadUint64(&q.getpos)
	if putpos >= getpos {
		posCnt = putpos - getpos
	} else {
		posCnt = 0
	}
	if posCnt < 1 {
		time.Sleep(q.sleep)
		return 0, posCnt
	}
	if size := uint64(len(values)); posCnt >= size {
		getCnt = size
	} else {
		getCnt = posCnt
	}
	getposnew = getpos + getCnt
	if !atomic.CompareAndSwapUint64(&q.getpos, getpos, getposnew) {
		runtime.Gosched()
		return 0, posCnt
	}
	for posnew, v := getpos+1, uint64(0); v < getCnt; posnew, v = posnew+1, v+1 {
		var cache = &q.cache[posnew&q.capmod]
		for {
			getno := atomic.LoadUint64(&cache.getno)
			putno := atomic.LoadUint64(&cache.putno)
			if posnew == getno && getno == putno-q.capacity {
				values[v] = cache.value
				cache.value = nil
				atomic.AddUint64(&cache.getno, q.capacity)
				break
			} else {
				runtime.Gosched()
			}
		}
	}
	return getCnt, posCnt - getCnt
}

// round 到最近的2的倍数
func minQuantity(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
