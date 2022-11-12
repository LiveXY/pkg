package reflectx

import (
	"reflect"
	"sync"
)

type structcache struct {
	name string
	fields int
}

var namemap map[any]structcache
var lock sync.Mutex
func GetStructName[T any]() (string, int) {
	var t T
	if cache, ok := namemap[t]; ok {
		return cache.name, cache.fields
	} else {
		lock.Lock()
		if namemap == nil {
			namemap = make(map[any]structcache)
		}
		tof := reflect.TypeOf(t)
		vk := tof.Kind()
		if vk == reflect.Ptr {
			tof = tof.Elem()
		}
		name := tof.Name()
		fields := tof.NumField()
		namemap[t] = structcache{name: name, fields: fields}
		lock.Unlock()
		return name, fields
	}
}
