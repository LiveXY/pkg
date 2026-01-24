package base64x

import (
	"testing"
)

// TestBase64 测试字符串的 Base64 编码功能
func TestBase64(t *testing.T) {
	got := Base64("hello")
	want := "aGVsbG8="
	if got != want {
		t.Errorf("Base64(hello) = %v，期望为 %v", got, want)
	}
}

// TestFrom64 测试 Base64 字符串的解码还原功能
func TestFrom64(t *testing.T) {
	got := From64("aGVsbG8=")
	want := "hello"
	if got != want {
		t.Errorf("From64(aGVsbG8=) = %v，期望为 %v", got, want)
	}
}
