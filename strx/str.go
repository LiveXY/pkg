package strx

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"
)

// ToInt 将字符串转换为 int
func ToInt(valstr string) int {
	val, err := strconv.Atoi(valstr)
	if err != nil {
		val = 0
	}
	return val
}

// ToUInt 将字符串转换为 uint
func ToUInt(valstr string) uint {
	val, err := strconv.ParseUint(valstr, 10, 64)
	if err != nil {
		return 0
	}
	return uint(val)
}

// ToBool 将字符串转换为 bool
func ToBool(valstr string) bool {
	valstr = strings.ToLower(valstr)
	switch valstr {
	case "true":
		return true
	case "1":
		return true
	}
	return false
}

// ToInt64 将字符串转换为 int64
func ToInt64(valstr string) int64 {
	val, err := strconv.ParseInt(valstr, 10, 64)
	if err != nil {
		val = 0
	}
	return val
}

// ToFloat 将字符串转换为 float64
func ToFloat(v string) float64 {
	v2, _ := strconv.ParseFloat(v, 64)
	return v2
}

// ToInt8 将字符串转换为 int8
func ToInt8(valstr string) int8 {
	val, err := strconv.Atoi(valstr)
	if err != nil {
		val = 0
	}
	if val < math.MinInt8 || val > math.MaxInt8 {
		return 0
	}
	return int8(val)
}

// Format 格式化字符串，支持 {0}, {1} 占位符
func Format(tpl string, args ...any) string {
	if len(args) == 0 {
		return tpl
	}
	for i := 0; i < len(args); i++ {
		placeholder := "{" + strconv.Itoa(i) + "}"
		if !strings.Contains(tpl, placeholder) {
			continue
		}

		var val string
		switch v := args[i].(type) {
		case int:
			val = strconv.Itoa(v)
		case int64:
			val = strconv.FormatInt(v, 10)
		case float32:
			val = strconv.FormatFloat(float64(v), 'f', 2, 64)
		case float64:
			val = strconv.FormatFloat(v, 'f', 2, 64)
		case time.Time:
			val = v.Format("2006-01-02 15:04:05")
		case error:
			val = v.Error()
		case string:
			val = v
		default:
			val = fmt.Sprint(v)
		}
		tpl = strings.ReplaceAll(tpl, placeholder, val)
	}
	return tpl
}

// ToBytes 将字符串转换为字节切片（零拷贝）
func ToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	// #nosec G103
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Copy 复制字符串
func Copy(s string) string {
	return string(ToBytes(s))
}

// Int8Contains 检查以指定分隔符分隔的字符串中是否包含指定的 int8 数值
func Int8Contains(str string, i int8, s string) bool {
	target := strconv.Itoa(int(i))
	list := strings.Split(str, s)
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

// ToTime 将字符串转换为 time.Time
func ToTime(str string) (time.Time, error) {
	if len(str) == 0 {
		return time.Now(), errors.New("日期不能为空！")
	}
	checks := [][]int{{1000, 9999}, {1, 12}, {1, 31}}
	var d, t string
	var ss []string
	if strings.Contains(str, "/") {
		if strings.Contains(str, " ") {
			s2 := strings.Split(str, " ")
			d, t = s2[0], s2[1]
			ss = strings.Split(d, "/")
		} else {
			ss = strings.Split(str, "/")
		}
		if len(ss) != 3 {
			return time.Now(), errors.New("日期格式错误！")
		}
		for i, v := range ss {
			v1 := ToInt(v)
			if v1 < 1 || !(checks[i][0] <= v1 && v1 <= checks[i][1]) {
				return time.Now(), errors.New("日期格式错误！")
			}
			if len(v) == 1 {
				ss[i] = "0" + v
			}
		}
		str = strings.Join(ss, "-")
		if len(t) > 0 {
			str += " " + t
		}
	} else if strings.Contains(str, "-") {
		if strings.Contains(str, " ") {
			s2 := strings.Split(str, " ")
			d, t = s2[0], s2[1]
			ss = strings.Split(d, "-")
		} else {
			ss = strings.Split(str, "-")
		}
		if len(ss) > 0 {
			if len(ss) != 3 {
				return time.Now(), errors.New("日期格式错误！")
			}
			for i, v := range ss {
				v1 := ToInt(v)
				if v1 < 1 || !(checks[i][0] <= v1 && v1 <= checks[i][1]) {
					return time.Now(), errors.New("日期格式错误！")
				}
				if len(v) == 1 {
					ss[i] = "0" + v
				}
			}
			str = strings.Join(ss, "-")
			if len(t) > 0 {
				str += " " + t
			}
		}
	}
	var format string
	switch len(str) {
	case 8:
		format = "20060102"
	case 10:
		format = "2006-01-02"
	case 12:
		format = "200601021504"
	case 13:
		format = "2006-01-02 15"
	case 14:
		format = "20060102150405"
	case 16:
		format = "2006-01-02 15:04"
	case 19:
		format = "2006-01-02 15:04:05"
	}
	return time.ParseInLocation(format, str, time.Local)
}

// ToUnix 将日期字符串转换为 Unix 时间戳
func ToUnix(str string) int64 {
	t, err := ToTime(str)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// ToMap 将字符串切片转换为 map
func ToMap(data []string) map[string]struct{} {
	list := make(map[string]struct{})
	for _, v := range data {
		list[strings.TrimSpace(v)] = struct{}{}
	}
	return list
}

// MapToStr 将 map 转换为字符串切片
func MapToStr(list map[string]struct{}) (data []string) {
	for id := range list {
		data = append(data, id)
	}
	return
}

// PadLeft 在字符串左侧填充字符
func PadLeft(s string, size int, ch rune) string {
	return pad(s, size, ch, true)
}

// PadRight 在字符串右侧填充字符
func PadRight(s string, size int, ch rune) string {
	return pad(s, size, ch, false)
}

func pad(s string, size int, ch rune, isLeft bool) string {
	if size <= 0 {
		return s
	}
	count := utf8.RuneCountInString(s)
	pads := size - count
	if pads <= 0 {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s) + pads*utf8.RuneLen(ch))

	if isLeft {
		for i := 0; i < pads; i++ {
			sb.WriteRune(ch)
		}
		sb.WriteString(s)
	} else {
		sb.WriteString(s)
		for i := 0; i < pads; i++ {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

// Repeat 重复输出指定字符
func Repeat(ch rune, repeat int) string {
	if repeat <= 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(repeat)
	for i := 0; i < repeat; i++ {
		sb.WriteRune(ch)
	}
	return sb.String()
}

// Len 获取字符串的字符数
func Len(s string) int {
	return utf8.RuneCountInString(s)
}

// Sub 截取字符串
func Sub(s string, start, end int) string {
	if s == "" {
		return ""
	}
	ulen := utf8.RuneCountInString(s)
	if end < 1000000 {
		if end < 0 {
			end += ulen
		}
		if end > ulen {
			end = ulen
		}
		if start < 0 {
			start += ulen
		}
	}

	if start > end || start >= ulen || end < 0 {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start == 0 && end == ulen {
		return s
	}

	sb := strings.Builder{}
	index := 0
	for _, v := range s {
		if index >= end {
			break
		}
		if index >= start {
			sb.WriteRune(v)
		}
		index++
	}
	return sb.String()
}
