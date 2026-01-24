package bytex

import (
	"reflect"
	"testing"
)

// TestToStr 测试字节切片转换为字符串的功能（零拷贝）
func TestToStr(t *testing.T) {
	b := []byte("hello")
	got := ToStr(b)
	if got != "hello" {
		t.Errorf("ToStr 结果为 %v，期望为 hello", got)
	}
}

// TestCopy 测试字节切片的深度复制功能
func TestCopy(t *testing.T) {
	b := []byte("hello")
	got := Copy(b)
	if !reflect.DeepEqual(got, b) {
		t.Errorf("Copy 结果为 %v，期望为 %v", got, b)
	}
	// 验证是否为深拷贝
	got[0] = 'y'
	if b[0] == 'y' {
		t.Errorf("Copy 应该是深拷贝")
	}
}
