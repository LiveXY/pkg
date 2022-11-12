package filex

import (
	"bufio"
	"io"
	"os"
)

const match = "/Page\x00"

// 获取PDF文件的页数
func pdfPages(reader io.ByteReader) (pages int) {
	i := 0
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return
		}
	check:
		switch match[i] {
		case 0:
			if !(b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z') {
				pages++
			}
			i = 0
			goto check
		case b:
			i++
		default:
			i = 0
		}
	}
}

// 获取PDF文件的页数
func PDFPages(path string) (pages int) {
	if reader, err := os.Open(path); err == nil {
		reader.Chdir()
		pages = pdfPages(bufio.NewReader(reader))
		reader.Close()
	}
	return
}
