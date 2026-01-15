package shax

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

// SHA256加密
func SHA256(sign string) string {
	hash := sha256.Sum256([]byte(sign))
	return fmt.Sprintf("%x", hash)
}

// SHA512加密
func SHA512(sign string) string {
	hash := sha512.Sum512([]byte(sign))
	return fmt.Sprintf("%x", hash)
}

// HMAC256加密
func HMAC256(message, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	return fmt.Sprintf("%x", h.Sum(nil))
}
