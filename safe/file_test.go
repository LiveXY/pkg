package safe

import (
	"encoding/hex"
	"testing"
)

// TestFileBytes 测试通过文件头字节流识别文件真实类型的逻辑
func TestFileBytes(t *testing.T) {
	// PNG magic: 89504e47
	pngHeader, _ := hex.DecodeString("89504e47")
	res := FileBytes("png", pngHeader)
	if !res.IsImage() {
		t.Errorf("应该识别为图片类型")
	}
	if res.GetExt() != "png" {
		t.Errorf("期望后缀为 png，实际结果为 %v", res.GetExt())
	}

	// PDF magic: 25504446
	pdfHeader, _ := hex.DecodeString("25504446")
	res = FileBytes("pdf", pdfHeader)
	if !res.IsPdf() {
		t.Errorf("应该识别为 pdf 类型")
	}

	// Fake magic
	res = FileBytes("jpg", []byte("not a real file"))
	if res.Category != "" {
		t.Errorf("随机内容不应该是别出分类")
	}
}
