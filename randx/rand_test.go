package randx

import (
	"strconv"
	"testing"
)

// TestStr 测试随机字母字符串生成功能
func TestStr(t *testing.T) {
	n := 10
	got := Str(n)
	if len(got) != n {
		t.Errorf("Str 长度为 %v，期望为 %v", len(got), n)
	}
}

// TestSymbol 测试包含特殊符号的随机字符串生成功能
func TestSymbol(t *testing.T) {
	n := 20
	got := Symbol(n)
	if len(got) != n {
		t.Errorf("Symbol 长度为 %v，期望为 %v", len(got), n)
	}
}

// TestUUID 测试 UUID 生成功能
func TestUUID(t *testing.T) {
	got := UUID()
	if len(got) != 36 {
		t.Errorf("UUID 长度为 %v，期望为 36", len(got))
	}
}

// TestOrderID 测试基于时间戳的订单 ID 生成功能
func TestOrderID(t *testing.T) {
	got := OrderID()
	// 20060102150405 is 14 chars + 7 chars captcha = 21 chars
	if len(got) != 21 {
		t.Errorf("OrderID 长度为 %v，期望为 21", len(got))
	}
}

// TestCaptchaCode 测试纯数字验证码生成功能
func TestCaptchaCode(t *testing.T) {
	n := 6
	got := CaptchaCode(n)
	if len(got) != n {
		t.Errorf("CaptchaCode 长度为 %v，期望为 %v", len(got), n)
	}
	if _, err := strconv.Atoi(got); err != nil {
		t.Errorf("CaptchaCode 应该全位数字：%v", got)
	}
}

// TestInt 测试指定范围内的随机整数生成功能
func TestInt(t *testing.T) {
	min, max := 10, 20
	for i := 0; i < 100; i++ {
		got := Int(min, max)
		if got < min || got > max {
			t.Errorf("Int 结果为 %v，超出了范围 [%d, %d]", got, min, max)
		}
	}
}
