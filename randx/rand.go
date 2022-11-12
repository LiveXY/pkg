package randx

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 随机UUID
func UUID() string { return uuid.New().String() }

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const symbolBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+,.?/:;{}[]`~"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())
var randx = rand.New(src)

// 随机字符串
func Str(n int) string {
	if n < 1 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
func Symbol(n int) string {
	if n < 1 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(symbolBytes) {
			sb.WriteByte(symbolBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
func OrderID() string {
	begin := time.Now().Format("20060102150405")
	end := CaptchaCode(7)
	return begin + end
}
func CaptchaCode(len int) string {
	min := int64(math.Pow10(len - 1))
	max := int64(math.Pow10(len) - 1)
	code := Int64(min, max)
	return strconv.FormatInt(code, 10)
}

// 随机数
func Int(min, max int) int {
	if max < min {
		return min
	}
	return randx.Intn(max-min+1) + min
}
func Int64(min, max int64) int64 {
	if max < min {
		return min
	}
	return randx.Int63n(max-min+1) + min
}
func Int63() int64 {
	return randx.Int63()
}
func Int63n(n int64) int64 {
	return randx.Int63n(n)
}
