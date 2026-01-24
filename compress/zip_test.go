package compress

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestZipUnZip 集成测试文件的 Zip 压缩与解压还原功能
func TestZipUnZip(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some test files
	file1 := filepath.Join(tmpDir, "test1.txt")
	file2 := filepath.Join(tmpDir, "test2.txt")
	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)

	zipPath := filepath.Join(tmpDir, "test.zip")
	files := []string{file1, file2}

	// Test Zip
	err := Zip(zipPath, files)
	if err != nil {
		t.Fatalf("Zip 压缩错误：%v", err)
	}

	// Test UnZip
	unzipDir := filepath.Join(tmpDir, "unzip")
	os.Mkdir(unzipDir, 0755)

	extractedFiles, err := UnZip(zipPath, unzipDir)
	if err != nil {
		t.Fatalf("UnZip 解压错误：%v", err)
	}

	// Note: Zip implementation uses filepath.Base(f) as standard name in zip
	wantFiles := []string{"test1.txt", "test2.txt"}
	if !reflect.DeepEqual(extractedFiles, wantFiles) {
		t.Errorf("解压出的文件列表不匹配：实际结果为 %v，期望为 %v", extractedFiles, wantFiles)
	}

	// Verify content
	c1, _ := os.ReadFile(filepath.Join(unzipDir, "test1.txt"))
	if string(c1) != "content1" {
		t.Errorf("test1.txt 内容不匹配")
	}
}
