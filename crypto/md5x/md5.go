package md5x

import (
	"crypto/md5"
	"fmt"
	"io"
)

// MD5加密
func MD5(sign string) string {
	h := md5.New()
	io.WriteString(h, sign)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}
