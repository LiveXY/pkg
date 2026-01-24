package num

import (
	"reflect"
	"testing"
	"time"
)

// TestInt64ToStr 测试 int64 转字符串功能
func TestInt64ToStr(t *testing.T) {
	if got := Int64ToStr(123); got != "123" {
		t.Errorf("Int64ToStr(123) = %v，期望为 123", got)
	}
}

// TestUInt64ToStr 测试 uint64 转字符串功能
func TestUInt64ToStr(t *testing.T) {
	if got := UInt64ToStr(123); got != "123" {
		t.Errorf("UInt64ToStr(123) = %v，期望为 123", got)
	}
}

// TestUIntToStr 测试 uint 转字符串功能
func TestUIntToStr(t *testing.T) {
	if got := UIntToStr(123); got != "123" {
		t.Errorf("UIntToStr(123) = %v，期望为 123", got)
	}
}

// TestIntToStr 测试 int 转字符串功能
func TestIntToStr(t *testing.T) {
	if got := IntToStr(123); got != "123" {
		t.Errorf("IntToStr(123) = %v，期望为 123", got)
	}
}

// TestFloatToStr 测试 float32 转字符串功能（指定精度）
func TestFloatToStr(t *testing.T) {
	if got := FloatToStr(1.234, 2); got != "1.23" {
		t.Errorf("FloatToStr(1.234, 2) = %v，期望为 1.23", got)
	}
}

// TestFloat64ToStr 测试 float64 转字符串功能（指定精度）
func TestFloat64ToStr(t *testing.T) {
	if got := Float64ToStr(1.234, 2); got != "1.23" {
		t.Errorf("Float64ToStr(1.234, 2) = %v，期望为 1.23", got)
	}
}

// TestInt8ToStr 测试 int8 转字符串功能
func TestInt8ToStr(t *testing.T) {
	if got := Int8ToStr(int8(12)); got != "12" {
		t.Errorf("Int8ToStr(12) = %v，期望为 12", got)
	}
}

// TestInt8ArrayToStr 测试 int8 数组转分隔符连接字符串功能
func TestInt8ArrayToStr(t *testing.T) {
	if got := Int8ArrayToStr([]int8{1, 2, 3}, ","); got != "1,2,3" {
		t.Errorf("Int8ArrayToStr([1,2,3]) = %v，期望为 1,2,3", got)
	}
}

// TestInt8ToStrArray 测试 int8 数组转字符串数组功能（包含过滤逻辑）
func TestInt8ToStrArray(t *testing.T) {
	got := Int8ToStrArray([]int8{1, -1, 2, 1}, true)
	want := []string{"1", "2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Int8ToStrArray = %v，期望为 %v", got, want)
	}
}

// TestArrayAddInt8 测试 int8 数组元素累加功能
func TestArrayAddInt8(t *testing.T) {
	source := []int8{1, 2}
	got := ArrayAddInt8(source, 1)
	want := []int8{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ArrayAddInt8 = %v，期望为 %v", got, want)
	}
}

// TestFloatRound 测试浮点数四舍五入功能
func TestFloatRound(t *testing.T) {
	if got := FloatRound(1.234, 2); got != 1.23 {
		t.Errorf("FloatRound(1.234, 2) = %v，期望为 1.23", got)
	}
}

// TestToStr 测试任意类型转字符串的通用转换功能
func TestToStr(t *testing.T) {
	tests := []struct {
		input any
		want  string
	}{
		{"str", "str"},
		{123, "123"},
		{int64(456), "456"},
		{true, "true"},
		{[]byte("bytes"), "bytes"},
		{1.23, "1.23"},
		{int8(8), "8"},
		{uint(9), "9"},
		{nil, ""},
	}
	for _, tt := range tests {
		if got := ToStr(tt.input); got != tt.want {
			t.Errorf("ToStr(%v) = %v，期望为 %v", tt.input, got, tt.want)
		}
	}
	// 测试时间
	tm := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	if got := ToStr(tm); got != "2023-01-01 12:00:00" {
		t.Errorf("ToStr(时间转换) 结果为 %v", got)
	}
}

// TestGetZHNum 测试数字转中文大写/读法功能
func TestGetZHNum(t *testing.T) {
	tests := []struct {
		num  int
		want string
	}{
		{0, "零"},
		{1, "一"},
		{10, "十"},
		{11, "十一"},
		{20, "二十"},
		{100, "一百"},
		{101, "一百零一"},
		{1000, "一千"},
		{10000, "一万"},
	}
	for _, tt := range tests {
		if got := GetZHNum(tt.num); got != tt.want {
			t.Errorf("GetZHNum(%d) = %v，期望为 %v", tt.num, got, tt.want)
		}
	}
}
