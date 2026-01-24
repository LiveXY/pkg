package safe

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	TypeImage  = "image"
	TypeVideo  = "video"
	TypeAudio  = "audio"
	TypeZip    = "zip"
	TypeOffice = "office"
	TypePdf    = "pdf"
)

type MagicConfig struct {
	Magic    string
	Ext      string
	Category string
}

var magicTable = []MagicConfig{
	{"ffd8ff", "jpg", TypeImage},
	{"89504e47", "png", TypeImage},
	{"474946", "gif", TypeImage},
	{"424d", "bmp", TypeImage},
	{"504b0304", "zip", TypeZip},
	{"25504446", "pdf", TypePdf},
	{"504b030414000600", "docx", TypeOffice},
	{"d0cf11e0a1b11ae1", "doc", TypeOffice},
	{"0000002066747970", "mp4", TypeVideo},
	{"464c5601", "flv", TypeVideo},
	{"1a45dfa3", "webm", TypeVideo},
	{"4d546864", "mid", TypeAudio},
	{"4 ID3", "mp3", TypeAudio},
	{"57415645666d7420", "wav", TypeAudio},
}

var extGroups = map[string][]string{
	TypeImage:  {"jpg", "jpeg", "png", "gif", "bmp"},
	TypeZip:    {"zip"},
	TypePdf:    {"pdf"},
	TypeOffice: {"docx", "doc", "pdf", "wps", "xlsx", "xls", "pptx", "ppt"},
	TypeVideo:  {"mp4", "m4v", "mov", "flv", "ogg", "webm", "rmvb", "swf", "mpg", "avi"},
	TypeAudio:  {"mp3", "mid", "wav"},
}

// FileTypeResult 文件类型检测结果
type FileTypeResult struct {
	Ext      string
	Category string
}

// File 检测指定路径文件的类型
func File(fpath string) *FileTypeResult {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fpath), "."))
	fd, err := os.Open(filepath.Clean(fpath))
	if err != nil {
		return &FileTypeResult{}
	}
	defer fd.Close()
	var buffer [32]byte
	n, _ := fd.Read(buffer[:])
	return FileBytes(ext, buffer[:n])
}

// FileBytes 根据字节流检测文件类型
func FileBytes(userExt string, buffer []byte) *FileTypeResult {
	if len(buffer) == 0 {
		return &FileTypeResult{}
	}
	fileCode := hex.EncodeToString(buffer)
	for _, conf := range magicTable {
		if strings.HasPrefix(fileCode, conf.Magic) {
			allowedExts := extGroups[conf.Category]
			if slices.Contains(allowedExts, userExt) {
				return &FileTypeResult{
					Ext:      conf.Ext,
					Category: conf.Category,
				}
			}
		}
	}
	return &FileTypeResult{}
}

// IsImage 是否为图片类型
func (r *FileTypeResult) IsImage() bool { return r.Category == TypeImage }

// IsVideo 是否为视频类型
func (r *FileTypeResult) IsVideo() bool { return r.Category == TypeVideo }

// IsAudio 是否为音频类型
func (r *FileTypeResult) IsAudio() bool { return r.Category == TypeAudio }

// IsZip 是否为压缩包类型
func (r *FileTypeResult) IsZip() bool { return r.Category == TypeZip }

// IsOffice 是否为 Office 文档类型
func (r *FileTypeResult) IsOffice() bool { return r.Category == TypeOffice }

// IsPdf 是否为 PDF 类型
func (r *FileTypeResult) IsPdf() bool { return r.Category == TypePdf }

// GetExt 获取文件扩展名
func (r *FileTypeResult) GetExt() string { return r.Ext }
