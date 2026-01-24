package compress

import (
	"encoding/base64"

	"github.com/livexy/pkg/bytex"
	"github.com/livexy/pkg/logx"
	"github.com/livexy/pkg/strx"

	"sync"

	"github.com/klauspost/compress/zstd"
	"go.uber.org/zap"
)

var encoderPool = sync.Pool{
	New: func() any {
		e, _ := zstd.NewWriter(nil)
		return e
	},
}

var decoderPool = sync.Pool{
	New: func() any {
		d, _ := zstd.NewReader(nil)
		return d
	},
}

// EncodeBytes 使用 zstd 压缩字节切片
func EncodeBytes(d []byte) []byte {
	enc := encoderPool.Get().(*zstd.Encoder)
	defer encoderPool.Put(enc)
	return enc.EncodeAll(d, make([]byte, 0, len(d)))
}

// DecodeBytes 使用 zstd 解压字节切片
func DecodeBytes(s []byte) ([]byte, error) {
	dec := decoderPool.Get().(*zstd.Decoder)
	defer decoderPool.Put(dec)
	d, err := dec.DecodeAll(s, nil)
	if err != nil {
		logx.Error.Error("zstd解压错误", zap.Error(err))
	}
	return d, err
}

// EncodeBase64 压缩字节切片并返回 Base64 字符串
func EncodeBase64(d []byte) string {
	enc := encoderPool.Get().(*zstd.Encoder)
	defer encoderPool.Put(enc)
	return base64.StdEncoding.EncodeToString(enc.EncodeAll(d, make([]byte, 0, len(d))))
}

// DecodeString 使用 zstd 解压字符串
func DecodeString(s string) ([]byte, error) {
	dec := decoderPool.Get().(*zstd.Decoder)
	defer decoderPool.Put(dec)
	d, err := dec.DecodeAll(strx.ToBytes(s), nil)
	if err != nil {
		logx.Error.Error("zstd解压错误", zap.Error(err))
	}
	return d, err
}

// DecodeBase64 解码 Base64 字符串并使用 zstd 解压
func DecodeBase64(s string) (string, error) {
	o, _ := base64.StdEncoding.DecodeString(s)
	dec := decoderPool.Get().(*zstd.Decoder)
	defer decoderPool.Put(dec)
	d, err := dec.DecodeAll(o, nil)
	if err != nil {
		logx.Error.Error("zstd解压错误", zap.Error(err))
		return "", err
	}
	return bytex.ToStr(d), nil
}
