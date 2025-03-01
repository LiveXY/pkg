package compress

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/livexy/pkg/bytex"
	"github.com/livexy/pkg/strx"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 解压文件
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
		fpath := path.Join(dir, decodeName)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, f.Mode())
			if err != nil {
				return list, err
			}
		} else {
			if strings.Contains(decodeName, string(os.PathSeparator)) {
				tpath := filepath.Dir(fpath)
				err := os.MkdirAll(tpath, os.ModePerm)
				if err != nil {
					return list, err
				}
			}
			fr, err := f.Open()
			if err != nil {
				return list, err
			}
			fw, err := os.OpenFile(filepath.Clean(fpath), os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode())
			if err != nil {
				_ = fr.Close()
				return list, err
			}
			_, err = io.Copy(fw, fr)
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

// 压缩文件
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
