package list

import (
	"reflect"
	"testing"
)

// TestList 测试线程安全列表容器的 Append, Get, Count, Clear 等基础操作
func TestList(t *testing.T) {
	l := New[int](10)

	if !l.Empty() {
		t.Errorf("Empty 应该返回 true")
	}

	l.Append(1, 2, 3)
	if l.Count() != 3 {
		t.Errorf("Count 应该返回 3")
	}

	val, ok := l.Get()
	if !ok || val != 1 {
		t.Errorf("Get 结果不匹配：实际结果为 %v，状态为 %v", val, ok)
	}

	if l.Count() != 2 {
		t.Errorf("Get 后 Count 应该返回 2")
	}

	snapshot := l.List()
	wantSnapshot := []int{2, 3}
	if !reflect.DeepEqual(snapshot, wantSnapshot) {
		t.Errorf("List 结果不匹配：实际结果为 %v，期望为 %v", snapshot, wantSnapshot)
	}

	l.Clear()
	if !l.Empty() {
		t.Errorf("Clear 后 Empty 应该返回 true")
	}
}
