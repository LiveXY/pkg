package num

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/livexy/linq"

	"github.com/shopspring/decimal"
)

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}
func UIntToStr(num uint) string {
	return strconv.FormatUint(uint64(num), 10)
}
func IntToStr(num int) string {
	return strconv.Itoa(num)
}
func FloatToStr(num float32, precs ...int) string {
	prec := 2
	if len(precs) > 0 {
		prec = precs[0]
	}
	return strconv.FormatFloat(float64(num), 'f', prec, 64)
}
func Float64ToStr(num float64, precs ...int) string {
	prec := 2
	if len(precs) > 0 {
		prec = precs[0]
	}
	return strconv.FormatFloat(num, 'f', prec, 64)
}
func Int8ToStr(num int8) string {
	return strconv.Itoa(int(num))
}
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
func Int8ToStrArray(data []int8, gtzero bool) (list []string) {
	for _, v := range data {
		if gtzero && v <= 0 {
			continue
		}
		list = append(list, Int8ToStr(v))
	}
	list = linq.Uniq(list)
	return
}
func ArrayAddInt8(source []int8, val int8) []int8 {
	for i := range source {
		source[i] += val
	}
	return source
}

func FloatRound(val float64, n int32) float64 {
	val, _ = decimal.NewFromFloat(val).Round(n).Float64()
	return val
}


// 转字符串
func ToStr(v any) string {
	var tmp = reflect.Indirect(reflect.ValueOf(v)).Interface()
	switch v := tmp.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.Itoa(int(v))
	case int16:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 2, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	case fmt.Stringer:
		return v.String()
	default:
		return ""
	}
}


var nums = [...]string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
var secs = [...]string{"", "万", "亿", "万亿", "亿亿"}
var chns = [...]string{"", "十", "百", "千"}

// 数字转中文数字
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
	if strings.Index(str, "一十") == 0 {
		str = str[3:]
	}
	return
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
