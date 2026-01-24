package shax

import (
	"testing"
)

// TestSHA256 测试 SHA256 哈希计算功能
func TestSHA256(t *testing.T) {
	got := SHA256("123456")
	want := "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"
	if got != want {
		t.Errorf("SHA256 结果不匹配")
	}
}

// TestSHA512 测试 SHA512 哈希计算功能
func TestSHA512(t *testing.T) {
	got := SHA512("123456")
	if len(got) != 128 {
		t.Errorf("SHA512 长度不匹配")
	}
}

// TestHMAC256 测试 HMAC-SHA256 消息认证码计算功能
func TestHMAC256(t *testing.T) {
	got := HMAC256("message", "key")
	if len(got) != 64 {
		t.Errorf("HMAC256 长度不匹配")
	}
}
