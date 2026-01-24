package md5x

import (
	"testing"
)

// TestMD5 测试 MD5 哈希计算功能
func TestMD5(t *testing.T) {
	got := MD5("123456")
	want := "e10adc3949ba59abbe56e057f20f883e"
	if got != want {
		t.Errorf("MD5(123456) = %v，期望为 %v", got, want)
	}
}
