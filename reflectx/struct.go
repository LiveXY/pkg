package reflectx

import (
	"reflect"
	"strings"
	"sync"
)

// TableStruct 结构体元数据
type TableStruct struct {
	Name   string
	Fields int
	PKs    []string
}

var tableCache sync.Map

// GetTableStruct 获取结构体的表元数据信息
func GetTableStruct[T any]() TableStruct {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if val, ok := tableCache.Load(typ); ok {
		return val.(TableStruct)
	}
	table := parseTableStruct(typ)
	tableCache.Store(typ, table)
	return table
}

func parseTableStruct(typ reflect.Type) TableStruct {
	numField := typ.NumField()
	table := TableStruct{
		Name:   typ.Name(),
		Fields: numField,
	}
	for i := 0; i < numField; i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("gorm")
		if tag == "" {
			continue
		}
		if strings.Contains(tag, "primaryKey") {
			columnName := ""
			parts := strings.Split(tag, ";")
			for _, part := range parts {
				if strings.HasPrefix(part, "column:") {
					columnName = strings.TrimPrefix(part, "column:")
					break
				}
			}
			if columnName == "" {
				columnName = field.Name
			}
			table.PKs = append(table.PKs, columnName)
		}
	}
	return table
}
