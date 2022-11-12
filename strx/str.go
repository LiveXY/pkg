package strx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"
)

// 字符串转换int
func ToInt(valstr string) int {
	val, err := strconv.Atoi(valstr)
	if err != nil {
		val = 0
	}
	return val
}

// 字符串转uint
func ToUInt(valstr string) uint {
	val, err := strconv.ParseInt(valstr, 10, 64)
	if err != nil {
		val = 0
	}
	return uint(val)
}
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

// 字符串转换int64
func ToInt64(valstr string) int64 {
	val, err := strconv.ParseInt(valstr, 10, 64)
	if err != nil {
		val = 0
	}
	return val
}

// 字符串转浮点数
func ToFloat(v string) float64 {
	v2, _ := strconv.ParseFloat(v, 64)
	return v2
}

// 字符串转换int8
func ToInt8(valstr string) int8 {
	val, err := strconv.Atoi(valstr)
	if err != nil {
		val = 0
	}
	return int8(val)
}

// 字符串格式输出
func Format(tpl string, args ...any) string {
	if len(args) == 0 {
		return tpl
	}
	for i := 0; i < len(args); i++ {
		switch args[i].(type) {
		case int:
			args[i] = strconv.Itoa(args[i].(int))
		case int64:
			args[i] = strconv.FormatInt(args[i].(int64), 10)
		case float32:
			args[i] = strconv.FormatFloat(float64(args[i].(float32)), 'f', 2, 64)
		case time.Time:
			args[i] = args[i].(time.Time).Format("2006-01-02 15:04:05")
		case error:
			args[i] = args[i].(error).Error()
		}
		tpl = strings.ReplaceAll(tpl, fmt.Sprintf("{%d}", i), "%v")
	}
	return fmt.Sprintf(tpl, args...)
}

func ToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func Copy(s string) string {
	return string(ToBytes(s))
}
func Int8Contains(str string, i int8, s string) bool {
	return strings.Contains(s+str+s, s+strconv.Itoa(int(i))+s)
}

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
func ToUnix(str string) int64 {
	t, err := ToTime(str)
	if err != nil {
		return 0
	}
	return t.Unix()
}

func ToMap(data []string) map[string]struct{} {
	list := make(map[string]struct{})
	for _, v := range data {
		list[strings.TrimSpace(v)] = struct{}{}
	}
	return list
}
func MapToStr(list map[string]struct{}) (data []string) {
	for id := range list {
		data = append(data, id)
	}
	return
}

func PadLeft(s string, size int, ch rune) string {
	return pad(s, size, ch, true)
}
func PadRight(s string, size int, ch rune) string {
	return pad(s, size, ch, false)
}
func pad(s string, size int, ch rune, isLeft bool) string {
	if size <= 0 {
		return s
	}
	pads := size - utf8.RuneCountInString(s)
	if pads <= 0 {
		return s
	}
	if isLeft {
		return Repeat(ch, pads) + s
	}
	return s + Repeat(ch, pads)
}
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
func Len(s string) int {
	return utf8.RuneCountInString(s)
}
func Sub(s string, start, end int) string {
	if s == "" {
		return ""
	}
	ulen := utf8.RuneCountInString(s)
	if end < 0 {
		end += ulen
	}
	if end > ulen {
		end = ulen
	}
	if start < 0 {
		start += ulen
	}
	if start > end {
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
	sb.Grow(end - start)
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
