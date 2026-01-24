package template

import (
	"testing"
)

// TestFastTemplate 测试简易、高性能的模板变量替换功能
func TestFastTemplate(t *testing.T) {
	text := "Hello {{.Name}}"
	param := map[string]any{"Name": "World"}
	got, err := FastTemplate(text, param)
	if err != nil {
		t.Errorf("FastTemplate 错误：%v", err)
	}
	if got != "Hello World" {
		t.Errorf("FastTemplate 结果为 %v，期望为 Hello World", got)
	}
}

// TestTextTemplate 测试基于标准库 text/template 的文本模板处理功能
func TestTextTemplate(t *testing.T) {
	text := "Hello {{.Name}}"
	param := map[string]any{"Name": "World"}
	got, err := TextTemplate(text, param)
	if err != nil {
		t.Errorf("TextTemplate 错误：%v", err)
	}
	if got != "Hello World" {
		t.Errorf("TextTemplate 结果为 %v，期望为 Hello World", got)
	}
}

// TestHTMLTemplate 测试基于标准库 html/template 的安全 HTML 模板处理功能，验证其自动转义特性
func TestHTMLTemplate(t *testing.T) {
	text := "<div>Hello {{.Name}}</div>"
	param := map[string]any{"Name": "<b>World</b>"}
	got, err := HTMLTemplate(text, param)
	if err != nil {
		t.Errorf("HTMLTemplate 错误：%v", err)
	}
	// HTML 模板应该转义标签
	if got != "<div>Hello &lt;b&gt;World&lt;/b&gt;</div>" {
		t.Errorf("HTMLTemplate 结果为 %v", got)
	}
}
