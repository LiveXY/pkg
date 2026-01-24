package util

import (
	"strings"
	"testing"
)

// TestFormatFileSize 测试文件大小格式化功能，覆盖从 B 到 EB 的各种单位
func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size uint64
		want string
	}{
		{0, "0.00B"},
		{100, "100.00B"},
		{1024, "1.00KB"},
		{1024 * 1024, "1.00MB"},
		{1024 * 1024 * 1024, "1.00GB"},
		{1024 * 1024 * 1024 * 1024, "1.00TB"},
		{1024 * 1024 * 1024 * 1024 * 1024, "1.00EB"},
		{1536, "1.50KB"},
	}
	for _, tt := range tests {
		if got := FormatFileSize(tt.size); got != tt.want {
			t.Errorf("FormatFileSize(%d) = %v，期望为 %v", tt.size, got, tt.want)
		}
	}
}

// TestFormatUseTime 测试英文/简写格式的耗时格式化功能
func TestFormatUseTime(t *testing.T) {
	tests := []struct {
		time int64
		want string
	}{
		{0, ""},
		{30, "30s"},
		{60, "1m"},
		{90, "1m30s"},
		{3600, "1h"},
		{3660, "1h1m"},
		{216000, "1d"}, // 根据源码逻辑 60*60*60 是分界线
	}
	for _, tt := range tests {
		if got := FormatUseTime(tt.time); got != tt.want {
			t.Errorf("FormatUseTime(%d) = %v，期望为 %v", tt.time, got, tt.want)
		}
	}
}

// TestFormatCNUseTime 测试中文格式的耗时格式化功能
func TestFormatCNUseTime(t *testing.T) {
	tests := []struct {
		time int64
		want string
	}{
		{0, ""},
		{30, "30秒"},
		{60, "1分钟"},
		{90, "1分30秒"},
		{3600, "1小时"},
		{3660, "1时1分"},
		{216000, "1天"},
	}
	for _, tt := range tests {
		if got := FormatCNUseTime(tt.time); got != tt.want {
			t.Errorf("FormatCNUseTime(%d) = %v，期望为 %v", tt.time, got, tt.want)
		}
	}
}

// TestUseTime 测试耗时统计工具函数的输出格式
func TestUseTime(t *testing.T) {
	got := UseTime("TestTask", func(i int) {}, 10)
	if !strings.HasPrefix(got, "TestTask执行10次用时：") {
		t.Errorf("UseTime() 输出格式不正确：%v", got)
	}
}

// TestToStr 测试多类型参数转换为下划线连接字符串的功能
func TestToStr(t *testing.T) {
	tests := []struct {
		params []any
		want   string
	}{
		{[]any{"a", "b"}, "a_b"},
		{[]any{"a", int(1), int8(2), int32(3), int64(4)}, "a_1_2_3_4"},
		{[]any{uint(5), uint8(6), uint32(7), uint64(8)}, "5_6_7_8"},
		{[]any{float32(1.2), float64(3.4)}, "1.20_3.40"},
		{[]any{true, false}, "True_False"},
		{[]any{"mixed", 100, true}, "mixed_100_True"},
	}
	for _, tt := range tests {
		if got := ToStr(tt.params...); got != tt.want {
			t.Errorf("ToStr(%v) = %v，期望为 %v", tt.params, got, tt.want)
		}
	}
}
