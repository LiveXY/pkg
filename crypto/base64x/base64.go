package base64x

import (
	"encoding/base64"
	"etms/pkg/bytex"
	"etms/pkg/strx"
)

// 加base64
func Base64(data string) string {
	return base64.StdEncoding.EncodeToString(strx.ToBytes(data))
}

// 解base64
func From64(data string) string {
	byte, _ := base64.StdEncoding.DecodeString(data)
	return bytex.ToStr(byte)
}
