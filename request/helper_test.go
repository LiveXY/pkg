package request

import (
	"testing"
	"time"
)

// TestConvertToQueryParams 测试将 map 转换为 URL 查询参数字符串的功能
func TestConvertToQueryParams(t *testing.T) {
	tests := []struct {
		params map[string]any
		want   string
	}{
		{map[string]any{"a": 1, "b": "2"}, "?a=1&b=2"}, // Note: map iteration order is random, but usually it works for 2 keys in tests
		{map[string]any{"a": 1}, "?a=1"},
		{nil, ""},
		{map[string]any{}, ""},
	}
	for _, tt := range tests {
		got := ConvertToQueryParams(tt.params)
		// Since map order is random, we check length and content if multiple keys
		if len(tt.params) <= 1 {
			if got != tt.want {
				t.Errorf("ConvertToQueryParams(%v) = %v，期望为 %v", tt.params, got, tt.want)
			}
		} else {
			if got[0] != '?' || len(got) != len(tt.want) {
				t.Errorf("ConvertToQueryParams(%v) 格式不匹配：%v", tt.params, got)
			}
		}
	}
}

// TestSort 测试将 map 按键排序并拼接为字符串的功能（常用于签名校验）
func TestSort(t *testing.T) {
	data := map[string]any{
		"b": 2,
		"a": 1,
		"c": 3,
	}
	got := Sort(data)
	want := "a1b2c3"
	if got != want {
		t.Errorf("Sort 结果为 %v，期望为 %v", got, want)
	}
}

// TestBuildTokenHeader 测试构建 Bearer Token 类型的 Authorization 请求头
func TestBuildTokenHeader(t *testing.T) {
	got := BuildTokenHeader("secret")
	if got.Name != "Authorization" || got.Value != "Bearer secret" {
		t.Errorf("BuildTokenHeader 结果不匹配：%v", got)
	}
}

// TestResetParams 测试重置并过滤 map 中空字符串参数的功能
func TestResetParams(t *testing.T) {
	c := CreateClient(time.Second)
	input := map[string]any{
		"a": "1",
		"b": "",
		"c": 3,
	}
	got := c.ResetParams(input)
	if got["a"] != "1" {
		t.Errorf("期望得到 a=1，实际结果为 %v", got["a"])
	}
	if _, ok := got["b"]; ok {
		t.Errorf("期望键名 b 被过滤掉")
	}
}
