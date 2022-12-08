package reflectx

import (
	"reflect"
	"strings"
	"sync"
)

type TableStruct struct {
	Name   string
	Fields int
	PKs    []string
}

var tablemap map[any]TableStruct
var lock sync.Mutex

func GetTableStruct[T any]() TableStruct {
	var t T
	if cache, ok := tablemap[t]; ok {
		return cache
	} else {
		lock.Lock()
		if tablemap == nil {
			tablemap = make(map[any]TableStruct)
		}
		tof := reflect.TypeOf(t)
		vk := tof.Kind()
		if vk == reflect.Ptr {
			tof = tof.Elem()
		}
		table := TableStruct{Name: tof.Name(), Fields: tof.NumField()}
		for i := 0; i < tof.NumField(); i++ {
			tag := tof.Field(i).Tag.Get("gorm")
			if strings.Contains(tag, "primaryKey") {
				in := strings.Index(tag, "column:")
				if in != -1 {
					table.PKs = append(table.PKs, tag[in+7:])
				}
			}
		}
		tablemap[t] = table
		lock.Unlock()
		return table
	}
}
