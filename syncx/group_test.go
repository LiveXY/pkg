package syncx

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// TestGroupDo 测试并发控制功能，确保相同 key 的任务在并发调用时只执行一次
func TestGroupDo(t *testing.T) {
	var g Group
	var calls int32
	fn := func() error {
		atomic.AddInt32(&calls, 1)
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	key := "test-key"
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			g.Do(key, fn)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("期望仅调用 1 次，实际调用次数为 %d", calls)
	}
}

// TestGroupShared 测试并发共享结果功能，验证第一个调用者和后续并发调用者的返回状态
func TestGroupShared(t *testing.T) {
	var g Group
	fn := func() error {
		time.Sleep(100 * time.Millisecond)
		return errors.New("error")
	}

	key := "test-key"
	err, shared := g.Shared(key, fn)
	if err == nil || shared {
		t.Errorf("首次调用错误：err=%v，shared=%v", err, shared)
	}
}
