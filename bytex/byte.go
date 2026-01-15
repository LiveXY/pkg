package bytex

import "unsafe"

func ToStr(b []byte) string {
	// #nosec G103
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func ToStr2(b []byte) string {
	// #nosec G103
	return *(*string)(unsafe.Pointer(&b))
}

func Copy(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}
