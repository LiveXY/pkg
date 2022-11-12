package compress

import (
	"encoding/base64"

	"github.com/livexy/pkg/bytex"
	"github.com/livexy/pkg/logx"
	"github.com/livexy/pkg/strx"

	"github.com/klauspost/compress/zstd"
	"go.uber.org/zap"
)

var encoder, _ = zstd.NewWriter(nil)
var decoder, _ = zstd.NewReader(nil)

// 压缩数据 bytes
func EncodeBytes(d []byte) []byte {
	return encoder.EncodeAll(d, make([]byte, 0, len(d)))
}

// 解压缩数据 bytes
func DecodeBytes(s []byte) ([]byte, error) {
	d, err := decoder.DecodeAll(s, nil)
	if err != nil {
		logx.Error.Error("zstd解压错误", zap.Error(err))
	}
	return d, err
}

// 压缩数据 base64
func EncodeBase64(d []byte) string {
	return base64.StdEncoding.EncodeToString(encoder.EncodeAll(d, make([]byte, 0, len(d))))
}

// 解压缩数据 string
func DecodeString(s string) ([]byte, error) {
	d, err := decoder.DecodeAll(strx.ToBytes(s), nil)
	if err != nil {
		logx.Error.Error("zstd解压错误", zap.Error(err))
	}
	return d, err
}

// 解压缩数据 base64
func DecodeBase64(s string) (string, error) {
	o, _ := base64.StdEncoding.DecodeString(s)
	d, err := decoder.DecodeAll(o, nil)
	if err != nil {
		logx.Error.Error("zstd解压错误", zap.Error(err))
		return "", err
	}
	return bytex.ToStr(d), nil
}
