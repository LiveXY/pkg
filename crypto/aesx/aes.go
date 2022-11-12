package aesx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/livexy/pkg/bytex"
	"github.com/livexy/pkg/strx"
)

// AES加密
func Encrypt(key string, iv string, data string) string {
	if len(data) == 0 {
		return ""
	}
	key2, _ := base64.StdEncoding.DecodeString(key)
	iv2, _ := base64.StdEncoding.DecodeString(iv)
	block, _ := aes.NewCipher(key2)
	bs := block.BlockSize()
	origin := pkcs5Padding(strx.ToBytes(data), bs)
	cipher.NewCBCEncrypter(block, iv2).CryptBlocks(origin, origin)
	data = base64.StdEncoding.EncodeToString(origin)
	return data
}

// AES解密
func Decrypt(key string, iv string, data string) string {
	if len(data) == 0 {
		return ""
	}
	key2, _ := base64.StdEncoding.DecodeString(key)
	iv2, _ := base64.StdEncoding.DecodeString(iv)
	block, _ := aes.NewCipher(key2)
	origin, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	cipher.NewCBCDecrypter(block, iv2).CryptBlocks(origin, origin)
	data = bytex.ToStr(pkcs5UnPadding(origin))
	return data
}

func pkcs5Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origdata []byte) []byte {
	length := len(origdata)
	unpadding := int(origdata[length-1])
	return origdata[:(length - unpadding)]
}
