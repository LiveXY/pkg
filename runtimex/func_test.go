package runtimex

import (
	"strings"
	"testing"
)

// TestGetCurrFuncName 测试获取当前正在执行的函数名称的功能
func TestGetCurrFuncName(t *testing.T) {
	name := GetCurrFuncName()
	// Depending on how test is run, it might be TestGetCurrFuncName or TestGetCurrFuncName.func1
	if name != "TestGetCurrFuncName" && !strings.Contains(name, "TestGetCurrFuncName") {
		t.Errorf("期望函数名为 TestGetCurrFuncName，实际结果为 %v", name)
	}
}

func myTestFunc() {}

// TestGetCallFuncName 测试通过函数变量获取其原始定义函数名称的功能
func TestGetCallFuncName(t *testing.T) {
	name := GetCallFuncName(myTestFunc)
	if name != "myTestFunc" {
		t.Errorf("期望函数名为 myTestFunc，实际结果为 %v", name)
	}

	if got := GetCallFuncName(nil); got != "nil" {
		t.Errorf("期望结果为 nil，实际结果为 %v", got)
	}

	if got := GetCallFuncName(123); got != "not_func" {
		t.Errorf("期望结果为 not_func，实际结果为 %v", got)
	}
}
