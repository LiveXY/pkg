package timex

import (
	"fmt"
	"strconv"
	"time"

	"github.com/livexy/pkg/strx"
)

// UseTime 测试函数执行耗时
func UseTime(title string, fn func(), n ...int) string {
	start := time.Now()
	num := 1
	if len(n) > 0 {
		num = n[0]
	}
	for i := 0; i < num; i++ {
		fn()
	}
	cost := time.Since(start)
	return title + "执行" + strconv.Itoa(num) + "次用时：" + cost.String()
}

// Format 时间格式化枚举
type Format uint

const (
	IntYMD Format = iota
	IntYM
	IntY
	IntYMDHMS
	IntHM
	YMD
	YMDHM
	YMDHMS
	CNYMD
	CNYMDHM
	CNYMDHMS
	ENYMD
	ENYMDHM
	ENYMDHMS
)

// ToStr 将 time.Time 转换为字符串
func ToStr(t time.Time, format ...Format) string {
	var f Format
	if len(format) > 0 {
		f = format[0]
	}
	switch f {
	case IntYM:
		return t.Format("200601")
	case IntY:
		return t.Format("2006")
	case IntYMDHMS:
		return t.Format("20060102150405")
	case IntHM:
		return t.Format("1504")
	case YMD:
		return t.Format("2006-01-02")
	case YMDHM:
		return t.Format("2006-01-02 15:04")
	case YMDHMS:
		return t.Format("2006-01-02 15:04:05")
	case CNYMD:
		return t.Format("2006年1月2日")
	case CNYMDHM:
		return t.Format("2006年1月2日 15:04")
	case CNYMDHMS:
		return t.Format("2006年1月2日 15:04:05")
	case ENYMD:
		return t.Format("2006/01/02")
	case ENYMDHM:
		return t.Format("2006/01/02 15:04")
	case ENYMDHMS:
		return t.Format("2006/01/02 15:04:05")
	default:
		return t.Format("20060102")
	}
}

// ToInt 将时间转换为整数表示（如 YMD 整数）
func ToInt(t time.Time, format ...Format) int {
	var f Format
	if len(format) > 0 {
		f = format[0]
	}
	if f >= YMD {
		return 0
	}
	return strx.ToInt(ToStr(t, format...))
}

// UnixToStr 将 Unix 时间戳转换为字符串
func UnixToStr(t int64, format ...Format) string {
	return ToStr(UnixToTime(t), format...)
}

// IntHMToStr 将小时分钟整数转换为 HH:mm 格式字符串
func IntHMToStr(t int) string {
	h := t / 100
	m := t % 100
	return fmt.Sprintf("%02d:%02d", h, m)
}

// UnixToTime 将 Unix 时间戳转换为 time.Time
func UnixToTime(t int64) time.Time { return time.Unix(t, 0) }

// GetMondayInt 获取本周一的日期整数表示
func GetMondayInt() (weekmonday int) {
	now := time.Now()
	weekmonday, err := strconv.Atoi(GetMondayTime(now).Format("20060102"))
	if err != nil {
		weekmonday = 0
	}
	return
}

// GetMondayTime 获取指定时间所在周的周一时间
func GetMondayTime(now time.Time) (monday time.Time) {
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}
