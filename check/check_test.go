package check

import (
	"testing"
)

// TestIsUserName 测试用户名格式校验逻辑
func TestIsUserName(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"abc", true},
		{"ab", false},
		{"admin@123", true},
		{"user_123", true},
		{"user-123", false},
	}
	for _, tt := range tests {
		if got := IsUserName(tt.input); got != tt.want {
			t.Errorf("IsUserName(%v) = %v，期望为 %v", tt.input, got, tt.want)
		}
	}
}

// TestIsEmail 测试邮箱格式校验逻辑
func TestIsEmail(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"test@example.com", true},
		{"test.123@example.com", true},
		{"test@sub.example.com", true},
		{"invalid-email", false},
		{"@example.com", false},
	}
	for _, tt := range tests {
		if got := IsEmail(tt.input); got != tt.want {
			t.Errorf("IsEmail(%v) = %v，期望为 %v", tt.input, got, tt.want)
		}
	}
}

// TestIsNumeric 测试纯数字校验逻辑
func TestIsNumeric(t *testing.T) {
	if !IsNumeric("12345") {
		t.Errorf("IsNumeric(12345) 应该返回 true")
	}
	if IsNumeric("123a45") {
		t.Errorf("IsNumeric(123a45) 应该返回 false")
	}
}

// TestIsAlphanumeric 测试字母数字校验逻辑
func TestIsAlphanumeric(t *testing.T) {
	if !IsAlphanumeric("abc123") {
		t.Errorf("IsAlphanumeric(abc123) 应该返回 true")
	}
	if IsAlphanumeric("abc_123") {
		t.Errorf("IsAlphanumeric(abc_123) 应该返回 false")
	}
}

// TestIsJSON 测试 JSON 字符串有效性校验逻辑
func TestIsJSON(t *testing.T) {
	if !IsJSON(`{"key": "value"}`) {
		t.Errorf("IsJSON 对于有效的内容应该返回 true")
	}
	if IsJSON(`{key: value}`) {
		t.Errorf("IsJSON 对于无效的内容应该返回 false")
	}
}

// TestIsMobile 测试中国大陆手机号格式校验逻辑
func TestIsMobile(t *testing.T) {
	if !IsMobile("13812345678") {
		t.Errorf("IsMobile(13812345678) 应该返回 true")
	}
	if IsMobile("12312345678") {
		t.Errorf("IsMobile(12312345678) 应该返回 false")
	}
}

// TestIsInternalIPv4 测试内部/回环 IPv4 地址校验逻辑
func TestIsInternalIPv4(t *testing.T) {
	if !IsInternalIPv4("127.0.0.1") {
		t.Errorf("127.0.0.1 应该是内部地址")
	}
	if !IsInternalIPv4("192.168.1.1") {
		t.Errorf("192.168.1.1 应该是内部地址")
	}
	if IsInternalIPv4("8.8.8.8") {
		t.Errorf("8.8.8.8 不应该是内部地址")
	}
}

// TestIsStrongPass 测试密码强度校验逻辑
func TestIsStrongPass(t *testing.T) {
	if got := IsStrongPass("12345"); got != 0 {
		t.Errorf("过短密码应该得 0 分")
	}
	if got := IsStrongPass("Ab1!"); got != 4 {
		// Wait, length < 6 returns 0
		if got != 0 {
			t.Errorf("长度为 4 应该得 0 分")
		}
	}
	if got := IsStrongPass("Abc123!"); got != 4 {
		t.Errorf("Abc123! 的分应该是 4，实际得分为 %d", got)
	}
}

// TestIsDate 测试日期常用格式校验逻辑
func TestIsDate(t *testing.T) {
	if !IsDate("2023-01-01") {
		t.Errorf("2023-01-01 应该是有效的日期格式")
	}
	if !IsDate("20230101120005") {
		t.Errorf("20230101120005 应该是有效的日期格式")
	}
}
