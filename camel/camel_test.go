package camel

import (
	"testing"
)

// TestBigCamel 测试下画线转大驼峰命名的功能
func TestBigCamel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user_id", "UserID"},
		{"user_name", "UserName"},
		{"api_url", "APIURL"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := BigCamel(tt.input); got != tt.want {
			t.Errorf("BigCamel(%v) = %v，期望为 %v", tt.input, got, tt.want)
		}
	}
}

// TestSmallCamel 测试其他命名转小驼峰命名的功能
func TestSmallCamel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"UserID", "userId"},
		{"UserName", "userName"},
		{"APIURL", "apiUrl"},
	}
	for _, tt := range tests {
		if got := SmallCamel(tt.input); got != tt.want {
			t.Errorf("SmallCamel(%v) = %v，期望为 %v", tt.input, got, tt.want)
		}
	}
}

// TestUnBigCamel 测试驼峰命名转回下画线命名的功能
func TestUnBigCamel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"UserID", "user_id"},
		{"UserName", "user_name"},
		{"APIURL", "api_url"},
	}
	for _, tt := range tests {
		if got := UnBigCamel(tt.input); got != tt.want {
			t.Errorf("UnBigCamel(%v) = %v，期望为 %v", tt.input, got, tt.want)
		}
	}
}
