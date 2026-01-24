package compress

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/livexy/pkg/bytex"
	"github.com/livexy/pkg/strx"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// MaxDecompressSize 限制最大解压大小，防止压缩包炸弹
const MaxDecompressSize = 2 * 1024 * 1024 * 1024

// UnZip 解压指定的 zip 文件到目标目录，支持自动处理 GB18030 编码的文件名
func UnZip(zipfile, dir string) ([]string, error) {
	var list []string
	cf, err := zip.OpenReader(zipfile)
	if err != nil {
		return list, err
	}
	defer cf.Close()
	var decodeName string
	for _, f := range cf.File {
		if strings.HasPrefix(f.Name, "__MACOSX") || strings.HasSuffix(f.Name, ".DS_Store") {
			continue
		}
		if f.Flags == 0 {
			i := bytes.NewReader(strx.ToBytes(f.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := io.ReadAll(decoder)
			decodeName = bytex.ToStr(content)
		} else {
			decodeName = f.Name
		}
		fpath := filepath.Join(dir, decodeName)
		rel, err := filepath.Rel(dir, fpath)
		if err != nil || strings.HasPrefix(rel, "..") || strings.HasPrefix(rel, "/") {
			continue
		}

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, f.Mode())
			if err != nil {
				return list, err
			}
		} else {
			if strings.Contains(decodeName, string(os.PathSeparator)) {
				tpath := filepath.Dir(fpath)
				err := os.MkdirAll(tpath, 0750)
				if err != nil {
					return list, err
				}
			}
			fr, err := f.Open()
			if err != nil {
				return list, err
			}
			fw, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode())
			if err != nil {
				_ = fr.Close()
				return list, err
			}
			lr := io.LimitReader(fr, MaxDecompressSize)
			_, err = io.Copy(fw, lr)
			if err != nil {
				_ = fr.Close()
				_ = fw.Close()
				return list, err
			}
			list = append(list, decodeName)
			_ = fr.Close()
			_ = fw.Close()
		}
	}
	return list, nil
}

// Zip 将指定的文件列表压缩为 zip 包
func Zip(zipPath string, files []string) error {
	if len(files) == 0 {
		return errors.New("nil file")
	}
	fzip, err := os.Create(filepath.Clean(zipPath))
	if err != nil {
		return err
	}
	defer fzip.Close()
	zipfile := zip.NewWriter(fzip)
	defer zipfile.Close()
	for _, f := range files {
		filename := filepath.Base(f)
		fw, err := zipfile.Create(filename)
		if err != nil {
			return err
		}
		fc, err := os.ReadFile(filepath.Clean(f))
		if err != nil {
			return err
		}
		_, err = fw.Write(fc)
		if err != nil {
			return err
		}
	}
	return nil
}
