package timex

import (
	"strconv"
	"time"

	"github.com/livexy/pkg/strx"
)

// 用时
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

// 格式日期
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

// 日期转换为YMD int
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

// Unix日期转字符串
func UnixToStr(t int64, format ...Format) string {
	return ToStr(UinxToTime(t), format...)
}

func IntHMToStr(t int) string {
	ts := strconv.Itoa(t)
	switch len(ts) {
	case 1:
		ts = "00:0" + ts
	case 2:
		ts = "00:" + ts
	case 3:
		ts = "0" + ts[:1] + ":" + ts[1:]
	case 4:
		ts = ts[:2] + ":" + ts[2:]
	}
	return ts
}

// Uinx转时间
func UinxToTime(t int64) time.Time { return time.Unix(t, 0) }

// 获取周一日期
func GetMondayInt() (weekmonday int) {
	now := time.Now()
	weekmonday, err := strconv.Atoi(GetMondayTime(now).Format("20060102"))
	if err != nil {
		weekmonday = 0
	}
	return
}
func GetMondayTime(now time.Time) (monday time.Time) {
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}