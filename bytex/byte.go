package bytex

import "unsafe"

// ToStr 将字节切片转换为字符串（零拷贝）
func ToStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// Copy 复制字节切片
func Copy(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}
