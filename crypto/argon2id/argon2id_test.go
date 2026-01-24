package argon2id

import (
	"testing"
)

// TestHashVerify 测试 Argon2id 密码哈希生成及其与原始密码的验证匹配功能
func TestHashVerify(t *testing.T) {
	pass := "secret123"
	hash := Hash(pass)
	if hash == "" {
		t.Fatalf("Hash 生成失败")
	}

	if !Verify(pass, hash) {
		t.Errorf("正确密码验证失败")
	}

	if Verify("wrongpass", hash) {
		t.Errorf("错误密码验证不应该通过")
	}
}
