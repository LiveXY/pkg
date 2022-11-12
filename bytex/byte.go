package bytex

import "unsafe"

func ToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func Copy(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}
