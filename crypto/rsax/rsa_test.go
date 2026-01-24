package rsax

import (
	"testing"
)

// TestRSA 测试 RSA 密钥对生成、公钥加密及私钥解密的完整流程
func TestRSA(t *testing.T) {
	pub, pri, ok := RSAGenKey(1024)
	if !ok {
		t.Fatalf("RSAGenKey 密钥生成失败")
	}

	data := "hello rsa"
	encrypted := Encrypt(pub, data)
	if encrypted == "" {
		t.Errorf("RSA 加密失败")
	}

	decrypted := Decrypt(pri, encrypted)
	if decrypted != data {
		t.Errorf("RSA 解密结果不匹配：实际结果为 %v，期望为 %v", decrypted, data)
	}
}
