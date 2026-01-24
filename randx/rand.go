package randx

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const symbolBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+,.?/:;{}[]`~"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	symbolIdxBits = 7
	symbolIdxMask = 1<<symbolIdxBits - 1
)

// Str 生成指定长度的随机字母字符串
func Str(n int) string {
	if n < 1 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(n)
	cache, remain := make([]byte, n), n
	for i := 0; i < n; {
		if remain == 0 {
			if _, err := rand.Read(cache); err != nil {
				panic("crypto/rand failed: " + err.Error())
			}
			remain = n
		}
		idx := int(cache[remain-1] & letterIdxMask)
		remain--
		if idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i++
		}
	}
	return sb.String()
}

// Symbol 生成指定长度的包含符号的随机字符串
func Symbol(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	cache, remain := make([]byte, n), n
	for i := 0; i < n; {
		if remain == 0 {
			if _, err := rand.Read(cache); err != nil {
				panic("crypto/rand failed: " + err.Error())
			}
			remain = n
		}
		idx := int(cache[remain-1] & symbolIdxMask)
		remain--
		if idx < len(symbolBytes) {
			sb.WriteByte(symbolBytes[idx])
			i++
		}
	}
	return sb.String()
}

// UUID 生成随机 UUID 字符串
func UUID() string { return uuid.New().String() }

// OrderID 生成基于当前时间和随机数的订单 ID
func OrderID() string {
	begin := time.Now().Format("20060102150405")
	end := CaptchaCode(7)
	return begin + end
}

// CaptchaCode 生成指定长度的数字验证码字符串
func CaptchaCode(len int) string {
	min := int64(math.Pow10(len - 1))
	max := int64(math.Pow10(len) - 1)
	code := Int64(min, max)
	return strconv.FormatInt(code, 10)
}

// Int 生成指定范围内的随机整数 [min, max]
func Int(min, max int) int {
	if max < min {
		return min
	}
	rangeSize := int64(max - min + 1)
	n, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return min
	}
	return int(n.Int64()) + min
}

// Int64 生成指定范围内的随机 int64 [min, max]
func Int64(min, max int64) int64 {
	if max < min {
		return min
	}
	delta := new(big.Int).SetInt64(max - min + 1)
	n, err := rand.Int(rand.Reader, delta)
	if err != nil {
		panic("crypto/rand 失败: " + err.Error())
	}
	return n.Int64() + min
}

// Int63 生成一个随机的非负 int64
func Int63() int64 {
	max := new(big.Int).Lsh(big.NewInt(1), 63)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return n.Int64()
}

// Int63n 生成 [0, n) 范围内的随机 int64
func Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}
	max := big.NewInt(n)
	result, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return result.Int64()
}
