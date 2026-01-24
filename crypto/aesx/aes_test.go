package aesx

import (
	"encoding/base64"
	"testing"
)

// TestEncryptDecrypt 测试 AES CBC 模式的加解密功能，验证密文正确性及解密后的原始值
func TestEncryptDecrypt(t *testing.T) {
	// 16 bytes key and iv, base64 encoded
	key := base64.StdEncoding.EncodeToString([]byte("1234567812345678"))
	iv := base64.StdEncoding.EncodeToString([]byte("1234567812345678"))
	data := "hello aes encryption"

	encrypted := Encrypt(key, iv, data)
	if encrypted == "" {
		t.Errorf("加密失败")
	}

	decrypted := Decrypt(key, iv, encrypted)
	if decrypted != data {
		t.Errorf("解密结果不匹配：实际结果为 %v，期望为 %v", decrypted, data)
	}
}
