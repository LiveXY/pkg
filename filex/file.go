package filex

import (
	"bufio"
	"bytes"
	"crypto/md5" // #nosec G501
	"encoding/hex"
	"image"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/livexy/pkg/bytex"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 去空行读文件
func ReadLine(fpath string) []string {
	var data []string
	f, _ := os.Open(filepath.Clean(fpath))
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := readLine(r)
		if len(line) > 0 {
			data = append(data, line)
		}
		if err != nil {
			break
		}
	}
	return data
}
func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return bytex.ToStr(line), err
}

// 读文本文件
func ReadText(fpath string) string {
	bytes, _ := os.ReadFile(filepath.Clean(fpath))
	return bytex.ToStr(bytes)
}

// 写文件数据
func WriteLine(fpath string, data []byte) error {
	file, err := os.OpenFile(filepath.Clean(fpath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// 目录是否存在
func DirExist(fpath string) bool {
	exist := false
	if obj, e := os.Stat(fpath); e == nil {
		exist = obj.IsDir()
	}
	return exist
}

// 获取文件信息
func FileInfo(file string) (string, string, string) {
	fpath := filepath.Dir(file)
	fname := filepath.Base(file)
	fext := strings.ToLower(path.Ext(file))
	return fpath, strings.TrimSuffix(fname, fext), fext
}

// 判断文件是否存在
func FileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 复制文件
func FileCopy(file, dst string) error {
	src, err := os.Open(filepath.Clean(file))
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(filepath.Clean(dst))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

// 创建目录
func MkDir(p string) error {
	return os.MkdirAll(p, 0750)
}

// 删除文件
func DeleteFile(fpath string) error {
	return os.Remove(fpath)
}

// 文件MD5
func FileMD5(fpath string) string {
	file, err := os.Open(filepath.Clean(fpath))
	if err != nil {
		return ""
	}
	defer file.Close()
	md5 := md5.New() // #nosec G401
	if _, err = io.Copy(md5, file); err != nil {
		return ""
	}
	md5str := hex.EncodeToString(md5.Sum(nil))
	return md5str
}

// 文件大小
func FileSize(fpath string) int64 {
	file, err := os.Open(filepath.Clean(fpath))
	if err != nil {
		return 0
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return 0
	}
	fsize := info.Size()
	return fsize
}

// 目录大小
func DirSize(path string) int64 {
	var size int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size
}

// 图片尺寸
func GetImageSize(fpath string) (int, int, error) {
	imgfile, err := os.Open(filepath.Clean(fpath))
	if err != nil {
		return 0, 0, err
	}
	defer imgfile.Close()
	img, _, err := image.DecodeConfig(imgfile)
	if err != nil {
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}

// 获取父路径
func GetParentPath(fpath string) string {
	fpath = strings.TrimRight(fpath, "/\\")
	return filepath.Dir(fpath)
}

// GBK 转 UTF-8
func GBKToUTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func UTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
