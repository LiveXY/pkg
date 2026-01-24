package num

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Int64ToStr 将 int64 转换为字符串
func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

// UInt64ToStr 将 uint64 转换为字符串
func UInt64ToStr(num uint64) string {
	return strconv.FormatUint(num, 10)
}

// UIntToStr 将 uint 转换为字符串
func UIntToStr(num uint) string {
	return strconv.FormatUint(uint64(num), 10)
}

// IntToStr 将 int 转换为字符串
func IntToStr(num int) string {
	return strconv.Itoa(num)
}

// FloatToStr 将 float32 转换为字符串，可指定精度
func FloatToStr(num float32, precs ...int) string {
	prec := 2
	if len(precs) > 0 {
		prec = precs[0]
	}
	return strconv.FormatFloat(float64(num), 'f', prec, 64)
}

// Float64ToStr 将 float64 转换为字符串，可指定精度
func Float64ToStr(num float64, precs ...int) string {
	prec := 2
	if len(precs) > 0 {
		prec = precs[0]
	}
	return strconv.FormatFloat(num, 'f', prec, 64)
}

// Int8ToStr 将 int8 转换为字符串
func Int8ToStr(num int8) string {
	return strconv.Itoa(int(num))
}

// Int8ArrayToStr 将 int8 切片转换为带分隔符的字符串
func Int8ArrayToStr(list []int8, s string) string {
	var sb strings.Builder
	for i, v := range list {
		if i > 0 {
			sb.WriteString(s)
		}
		sb.WriteString(Int8ToStr(v))
	}
	return sb.String()
}

// Int8ToStrArray 将 int8 切片转换为字符串切片，可过滤非正数
func Int8ToStrArray(data []int8, gtzero bool) (list []string) {
	seen := make(map[string]struct{})
	for _, v := range data {
		if gtzero && v <= 0 {
			continue
		}
		s := Int8ToStr(v)
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			list = append(list, s)
		}
	}
	return
}

// ArrayAddInt8 为 int8 切片中的每个元素增加指定数值
func ArrayAddInt8(source []int8, val int8) []int8 {
	for i := range source {
		source[i] += val
	}
	return source
}

// FloatRound 对浮点数进行指定位数的四舍五入
func FloatRound(val float64, n int32) float64 {
	p := math.Pow(10, float64(n))
	return math.Round(val*p) / p
}

// ToStr 将任意类型转换为字符串
func ToStr(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case bool:
		return strconv.FormatBool(val)
	case []byte:
		return string(val)
	case float64:
		return strconv.FormatFloat(val, 'f', 2, 64)
	case time.Time:
		return val.Format("2006-01-02 15:04:05")
	case int8:
		return strconv.Itoa(int(val))
	case int32:
		return strconv.Itoa(int(val))
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case fmt.Stringer:
		return val.String()
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return ""
		}
		return ToStr(rv.Elem().Interface())
	}

	rv = reflect.Indirect(rv)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', 2, 64)
	case reflect.String:
		return rv.String()
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}

	return ""
}

var nums = [...]string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
var secs = [...]string{"", "万", "亿", "万亿", "亿亿"}
var chns = [...]string{"", "十", "百", "千"}

// GetZHNum 将数字转换为中文数字字符串
func GetZHNum(num int) (str string) {
	if num < 1 {
		return nums[0]
	}
	pos, needzero := 0, false
	for num > 0 {
		sec := num % 10000
		if needzero {
			str = nums[0] + str
		}
		secstr := secString(sec)
		if sec != 0 {
			str = secstr + secs[pos] + str
		} else {
			str = secstr + str
		}
		needzero = (sec < 1000) && (sec > 0)
		num /= 10000
		pos++
	}
	if strings.HasPrefix(str, "一十") {
		return str[3:]
	}
	return str
}

func secString(sec int) string {
	str, pos, zero := "", 0, true
	for sec > 0 {
		v := sec % 10
		if v == 0 {
			if sec == 0 || !zero {
				zero, str = true, nums[v]+str
			}
		} else {
			ins := nums[v] + chns[pos]
			zero, str = false, ins+str
		}
		pos++
		sec /= 10
	}
	return str
}
