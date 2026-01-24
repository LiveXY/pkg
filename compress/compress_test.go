package compress

import (
	"reflect"
	"testing"
)

// TestEncodeDecodeBytes 测试字节切片的 zstd 压缩与解压功能
func TestEncodeDecodeBytes(t *testing.T) {
	original := []byte("hello world, this is a zstd compression test. hello world!")
	compressed := EncodeBytes(original)
	if len(compressed) == 0 {
		t.Errorf("压缩后的字节流不能为空")
	}

	decompressed, err := DecodeBytes(compressed)
	if err != nil {
		t.Fatalf("DecodeBytes 错误：%v", err)
	}

	if !reflect.DeepEqual(original, decompressed) {
		t.Errorf("解压内容不匹配")
	}
}

// TestEncodeBase64DecodeBase64 测试字节切片压缩后转 Base64 及逆向还原的功能
func TestEncodeBase64DecodeBase64(t *testing.T) {
	original := "hello zstd"
	compressed64 := EncodeBase64([]byte(original))

	decompressed, err := DecodeBase64(compressed64)
	if err != nil {
		t.Fatalf("DecodeBase64 错误：%v", err)
	}

	if decompressed != original {
		t.Errorf("解压内容不匹配：实际结果为 %v，期望为 %v", decompressed, original)
	}
}
