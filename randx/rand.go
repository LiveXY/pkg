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
	letterIdxBits = 6                    // 6 bits 可以表示 64 个索引 (2^6)
	letterIdxMask = 1<<letterIdxBits - 1 // 掩码 00111111
	symbolIdxBits = 7
	symbolIdxMask = 1<<symbolIdxBits - 1 // 01111111 (7个1)
)

// 随机字符串
func Str(n int) string {
	if n < 1 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(n)
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	for i := 0; i < n; i++ {
		// 将随机字节映射到 letterBytes 范围
		idx := int(b[i] & letterIdxMask)
		if idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
		} else {
			i--
		}
	}
	return sb.String()
}

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

// 随机UUID
func UUID() string { return uuid.New().String() }

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
	rangeSize := int64(max - min + 1)
	n, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return min
	}
	return int(n.Int64()) + min
}

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

func Int63() int64 {
	max := new(big.Int).Lsh(big.NewInt(1), 63)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return n.Int64()
}

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
