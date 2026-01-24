package runner

import (
	"sync/atomic"
	"testing"
)

// TestSafeReflectCall 测试反射调用函数的功能，并验证其 Panic 捕获与恢复机制
func TestSafeReflectCall(t *testing.T) {
	var count int32
	fn := func(val int) {
		atomic.StoreInt32(&count, int32(val))
	}

	safeReflectCall(fn, 123)
	if atomic.LoadInt32(&count) != 123 {
		t.Errorf("safeReflectCall 更新计数器失败")
	}

	// Test panic recovery
	panicFn := func() {
		panic("test panic")
	}
	// Should not crash
	safeReflectCall(panicFn)
}
