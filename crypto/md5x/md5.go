package md5x

import (
	"crypto/md5" // #nosec G501
	"fmt"
	"io"
)

// MD5加密
func MD5(sign string) string {
	h := md5.New() // #nosec G401
	_, _ = io.WriteString(h, sign)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}
