// 参考: https://www.garykessler.net/library/file_sigs.html
package safe

import (
	"bytes"
	"encoding/hex"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/livexy/linq"
)

var imgmap = map[string]string{
	"ffd8ff":   "jpg",
	"89504e47": "png",
	"474946":   "gif",
	"424d":     "bmp",
}
var zipmap = map[string]string{
	"504b0304": "zip",
}
var pdfmap = map[string]string{
	"25504446": "pdf",
}
var officemap = map[string]string{
	"504b030414000600":     "docx",
	"d0cf11e0a1b11ae1":     "doc",
	"504b03040a00":         "doc",
	"d0cf11e0a1b11ae10000": "wps",
	"cf11e0a1b11ae100":     "doc",
}
var videomap = map[string]string{
	"00000020667479706":            "mp4",
	"0000001":                      "mp4",
	"667479704d534e56":             "mp4",
	"6674797069736f6d":             "mp4",
	"667479706d703432":             "m4v",
	"667479704d345620":             "m4v",
	"6674797071742020":             "mov",
	"6d6f6f76":                     "mov",
	"464c5601":                     "flv",
	"4f67675300020000000000000000": "ogg",
	"1a45dfa3":                     "webm",
	"2e524d46":                     "rmvb",
	"465753":                       "swf",
	"5a5753":                       "swf",
	"445644":                       "dvd",
	"454e545259564344":             "vcd",
	"000001ba":                     "mpg",
	"000001b3":                     "mpg",
	"52494646":                     "mpg",
	"3026b2758e66cf11a6d9":         "wmv",
	"415649204c495354":             "avi",
}
var audiomap = map[string]string{
	"4d546864":         "mid",
	"494433":           "mp3",
	"57415645666d7420": "wav",
}

var imgexts = []string{"jpg", "jpeg", "png", "gif", "bmp"}
var zipexts = []string{"zip"}
var pdfexts = []string{"pdf"}
var officeexts = []string{"docx", "doc", "pdf", "wps", "xlsx", "xls", "pptx", "ppt", ".dotx", ".dot"}
var videoexts = []string{"mp4", "m4v", "mov", "flv", "ogg", "webm", "rmvb", "swf", "dvd", "vcd", "mpg", "mpeg", "wmv", "avi"}
var audioexts = []string{"mp3", "mid", "wav"}

func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// 获取文件类型
func getFileType(fsrc []byte, types map[string]string) (fileType string) {
	filecode := bytesToHexString(fsrc)
	for k, v := range types {
		if strings.HasPrefix(filecode, strings.ToLower(k)) || strings.HasPrefix(k, strings.ToLower(filecode)) {
			fileType = v
		}
	}
	return
}

func File(fpath string) *fileType {
	ext := strings.ToLower(strings.Trim(path.Ext(fpath), "."))
	fd, err := os.Open(fpath)
	if err == nil {
		var buffer [512]byte
		count, _ := fd.Read(buffer[:])
		fd.Close()
		return FileBytes(ext, buffer[:count])
	}
	return &fileType{}
}

func FileBytes(ext string, buffer []byte) *fileType {
	ftype := getFileType(buffer, imgmap)
	if len(ftype) > 0 && !linq.Contains(imgexts, ext) {
		ftype = ""
	}
	if len(ftype) == 0 {
		ftype = getFileType(buffer, pdfmap)
		if len(ftype) > 0 && !linq.Contains(pdfexts, ext) {
			ftype = ""
		}
	}
	if len(ftype) == 0 {
		ftype = getFileType(buffer, officemap)
		if len(ftype) > 0 && !linq.Contains(officeexts, ext) {
			ftype = ""
		}
	}
	if len(ftype) == 0 {
		ftype = getFileType(buffer, zipmap)
		if len(ftype) > 0 && !linq.Contains(zipexts, ext) {
			ftype = ""
		}
	}
	if len(ftype) == 0 {
		ftype = getFileType(buffer, audiomap)
		if len(ftype) > 0 && !linq.Contains(audioexts, ext) {
			ftype = ""
		}
	}
	if len(ftype) == 0 {
		ftype = getFileType(buffer, videomap)
		if len(ftype) > 0 && !linq.Contains(videoexts, ext) {
			ftype = ""
		}
	}
	return &fileType{fileType: ftype}
}

// 文件类型
type fileType struct{ fileType string }

// 是否图片文件
func (f *fileType) IsImageFile() bool { return linq.IndexOf(imgexts, f.fileType) != -1 }

// 是否视频文件
func (f *fileType) IsVideoFile() bool { return linq.IndexOf(videoexts, f.fileType) != -1 }

// 是否音频文件
func (f *fileType) IsAudioFile() bool { return linq.IndexOf(audioexts, f.fileType) != -1 }

// 是否压缩文件
func (f *fileType) IsZipFile() bool { return linq.IndexOf(zipexts, f.fileType) != -1 }

// 是否文档文件
func (f *fileType) IsOfficeFile() bool { return linq.IndexOf(officeexts, f.fileType) != -1 }
func (f *fileType) IsPdfFile() bool    { return linq.IndexOf(pdfexts, f.fileType) != -1 }

// 获取文件类型
func (f *fileType) GetType() string { return f.fileType }
