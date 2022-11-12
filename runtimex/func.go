package runtimex

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var namemap map[uintptr]string
var lock sync.Mutex
func GetCurrFuncName() string {
	pcs := make([]uintptr, 1)
	runtime.Callers(2, pcs)
	pc := pcs[0]
	if name, ok := namemap[pc]; ok {
		return name
	} else {
		lock.Lock()
		if namemap == nil {
			namemap = make(map[uintptr]string)
		}
		name := runtime.FuncForPC(pc).Name()
		if index := strings.LastIndex(name, "."); index > 0 {
			name = name[index+1:]
		}
		namemap[pc] = name
		lock.Unlock()
		return name
	}
}

func GetCallFuncName(fn any) string {
	pc := reflect.ValueOf(fn).Pointer()
	if name, ok := namemap[pc]; ok {
		return name
	} else {
		lock.Lock()
		if namemap == nil {
			namemap = make(map[uintptr]string)
		}
		name := runtime.FuncForPC(pc).Name()
		if index := strings.LastIndex(name, "."); index > 0 {
			name = name[index+1:]
		}
		namemap[pc] = name
		lock.Unlock()
		return name
	}
}
