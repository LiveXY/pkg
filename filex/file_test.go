package filex

import (
	"path/filepath"
	"reflect"
	"testing"
)

// TestFileInfo 测试获取文件路径、文件名和后缀的功能
func TestFileInfo(t *testing.T) {
	file := "/path/to/test.txt"
	dir, name, ext := FileInfo(file)
	if dir != "/path/to" {
		t.Errorf("期望目录为 /path/to，实际结果为 %v", dir)
	}
	if name != "test" {
		t.Errorf("期望名称为 test，实际结果为 %v", name)
	}
	if ext != ".txt" {
		t.Errorf("期望后缀为 .txt，实际结果为 %v", ext)
	}
}

// TestFileOperation 集成测试文件读写、存在检查、MD5计算、大小获取及父路径解析控制
func TestFileOperation(t *testing.T) {
	tmpDir := t.TempDir()
	fpath := filepath.Join(tmpDir, "test.txt")
	content := []byte("hello world\nline2")

	// Test WriteLine
	err := WriteLine(fpath, content)
	if err != nil {
		t.Fatalf("WriteLine 错误：%v", err)
	}

	// Test FileExist
	if !FileExist(fpath) {
		t.Errorf("FileExist 应该返回 true")
	}

	// Test ReadText
	text := ReadText(fpath)
	if text != "hello world\nline2" {
		t.Errorf("ReadText 结果为 %v", text)
	}

	// Test ReadLine
	lines := ReadLine(fpath)
	wantLines := []string{"hello world", "line2"}
	if !reflect.DeepEqual(lines, wantLines) {
		t.Errorf("ReadLine 结果为 %v，期望为 %v", lines, wantLines)
	}

	// Test FileMD5
	md5 := FileMD5(fpath)
	if md5 == "" {
		t.Errorf("FileMD5 结果不应为空")
	}

	// Test FileSize
	size := FileSize(fpath)
	if size != int64(len(content)) {
		t.Errorf("FileSize 结果为 %v，期望为 %d", size, len(content))
	}

	// Test GetParentPath
	parent := GetParentPath(fpath)
	if parent != tmpDir {
		t.Errorf("GetParentPath 结果为 %v，期望为 %v", parent, tmpDir)
	}
}

// TestGBKUTF8 测试 GBK 与 UTF-8 编码之间的相互转换功能
func TestGBKUTF8(t *testing.T) {
	utf8Str := "你好"
	gbk, err := UTF8ToGBK([]byte(utf8Str))
	if err != nil {
		t.Fatalf("UTF8ToGBK 错误：%v", err)
	}

	gotUTF8, err := GBKToUTF8(gbk)
	if err != nil {
		t.Fatalf("GBKToUTF8 错误：%v", err)
	}

	if string(gotUTF8) != utf8Str {
		t.Errorf("GBK/UTF8 转换失败：实际结果为 %v，期望为 %v", string(gotUTF8), utf8Str)
	}
}
