package runtimex

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var funcNameCache sync.Map

func GetCurrFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	return getFuncNameByPC(pc)
}

func GetCallFuncName(fn any) string {
	if fn == nil {
		return "nil"
	}
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return "not_func"
	}
	pc := v.Pointer()
	return getFuncNameByPC(pc)
}

func getFuncNameByPC(pc uintptr) string {
	if name, ok := funcNameCache.Load(pc); ok {
		return name.(string)
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown"
	}
	fullName := f.Name()
	shortName := fullName
	if idx := strings.LastIndexByte(fullName, '.'); idx >= 0 {
		shortName = fullName[idx+1:]
	}
	if idx := strings.Index(shortName, ".func"); idx >= 0 {
		shortName = shortName[:idx]
	}
	funcNameCache.Store(pc, shortName)
	return shortName
}
