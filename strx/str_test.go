package strx

import (
	"testing"
)

// TestToInt 测试字符串转 int 功能
func TestToInt(t *testing.T) {
	if got := ToInt("123"); got != 123 {
		t.Errorf("ToInt(123) 结果为 %v，期望为 123", got)
	}
	if got := ToInt("abc"); got != 0 {
		t.Errorf("ToInt(abc) 结果为 %v，期望为 0", got)
	}
}

// TestToUInt 测试字符串转 uint 功能
func TestToUInt(t *testing.T) {
	if got := ToUInt("123"); got != 123 {
		t.Errorf("ToUInt(123) 结果为 %v，期望为 123", got)
	}
}

// TestToBool 测试字符串转 bool 功能
func TestToBool(t *testing.T) {
	if got := ToBool("true"); got != true {
		t.Errorf("ToBool(true) 结果为 %v，期望为 true", got)
	}
	if got := ToBool("1"); got != true {
		t.Errorf("ToBool(1) 结果为 %v，期望为 true", got)
	}
	if got := ToBool("0"); got != false {
		t.Errorf("ToBool(0) 结果为 %v，期望为 false", got)
	}
}

// TestFormat 测试占位符 {0}, {1} 格式化功能
func TestFormat(t *testing.T) {
	got := Format("hello {0}, {1}", "world", 123)
	want := "hello world, 123"
	if got != want {
		t.Errorf("Format 结果为 %v，期望为 %v", got, want)
	}
}

// TestToBytes 测试字符串转字节切片（零拷贝）功能
func TestToBytes(t *testing.T) {
	s := "hello"
	got := ToBytes(s)
	if string(got) != s {
		t.Errorf("ToBytes 结果为 %s，期望为 %s", got, s)
	}
}

// TestInt8Contains 测试分隔符字符串中是否包含指定 int8 的功能
func TestInt8Contains(t *testing.T) {
	if !Int8Contains("1,2,3", 2, ",") {
		t.Errorf("Int8Contains 应该返回 true")
	}
	if Int8Contains("1,2,3", 4, ",") {
		t.Errorf("Int8Contains 应该返回 false")
	}
}

// TestToTime 测试多种日期格式字符串转 time.Time 的功能
func TestToTime(t *testing.T) {
	_, err := ToTime("2023-01-01 12:00:00")
	if err != nil {
		t.Errorf("ToTime 及其格式化解析错误：%v", err)
	}
	_, err = ToTime("2023/01/01")
	if err != nil {
		t.Errorf("ToTime 斜杠格式解析错误：%v", err)
	}
}

// TestToMap 测试字符串数组转 map 功能
func TestToMap(t *testing.T) {
	got := ToMap([]string{"a", " b "})
	if _, ok := got["a"]; !ok {
		t.Errorf("ToMap 缺少键名 'a'")
	}
	if _, ok := got["b"]; !ok {
		t.Errorf("ToMap 缺少键名 'b'")
	}
}

// TestPad 测试字符串左右填充功能
func TestPad(t *testing.T) {
	if got := PadLeft("1", 3, '0'); got != "001" {
		t.Errorf("PadLeft 结果为 %v，期望为 001", got)
	}
	if got := PadRight("1", 3, '0'); got != "100" {
		t.Errorf("PadRight 结果为 %v，期望为 100", got)
	}
}

// TestRepeat 测试字符重复生成功能
func TestRepeat(t *testing.T) {
	if got := Repeat('a', 3); got != "aaa" {
		t.Errorf("Repeat 结果为 %v，期望为 aaa", got)
	}
}

// TestSub 测试字符串安全截取（支持中文字符）的功能
func TestSub(t *testing.T) {
	s := "hello世界"
	if got := Sub(s, 0, 5); got != "hello" {
		t.Errorf("Sub 截取结果为 %v，期望为 hello", got)
	}
	if got := Sub(s, 5, 7); got != "世界" {
		t.Errorf("Sub 截取结果为 %v，期望为 世界", got)
	}
}
