package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/livexy/pkg/num"
)

// 格式输出文件大小
func FormatFileSize(fileSize uint64) (size string) {
	if fileSize < 1024 {
		size = fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		size = fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		size = fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		size = fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		size = fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		size = fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
	return
}

// 格式化用时
func FormatUseTime(t int64) (size string) {
	if t < 1 {
		return ""
	}
	if t < 60 {
		return strconv.FormatInt(t, 10) + "s"
	} else if t < (60 * 60) {
		if t%60 == 0 {
			return strconv.FormatInt(t/60, 10) + "m"
		}
		return strconv.FormatInt(t/60, 10) + "m" + strconv.FormatInt(t%60, 10) + "s"
	} else if t < (60 * 60 * 60) {
		if (t/60-60)%60 == 0 {
			return strconv.FormatInt(t/(60*60), 10) + "h"
		}
		return strconv.FormatInt(t/(60*60), 10) + "h" + strconv.FormatInt((t/60-60)%60, 10) + "m"
	} else {
		if (t/60/60-60)%60 == 0 {
			return strconv.FormatInt(t/(60*60*60), 10) + "d"
		}
		return strconv.FormatInt(t/(60*60*60), 10) + "d" + strconv.FormatInt((t/60/60-60)%60, 10) + "h"
	}
}

// 格式化用时
func FormatCNUseTime(t int64) (size string) {
	if t < 1 {
		return ""
	}
	if t < 60 {
		return strconv.FormatInt(t, 10) + "秒"
	} else if t < (60 * 60) {
		if t%60 == 0 {
			return strconv.FormatInt(t/60, 10) + "分钟"
		}
		return strconv.FormatInt(t/60, 10) + "分" + strconv.FormatInt(t%60, 10) + "秒"
	} else if t < (60 * 60 * 60) {
		if (t/60-60)%60 == 0 {
			return strconv.FormatInt(t/(60*60), 10) + "小时"
		}
		return strconv.FormatInt(t/(60*60), 10) + "时" + strconv.FormatInt((t/60-60)%60, 10) + "分"
	} else {
		if (t/60/60-60)%60 == 0 {
			return strconv.FormatInt(t/(60*60*60), 10) + "天"
		}
		return strconv.FormatInt(t/(60*60*60), 10) + "天" + strconv.FormatInt((t/60/60-60)%60, 10) + "时"
	}
}

// 三目运算int
func IF[T any](cond bool, suc, fail T) T {
	if cond {
		return suc
	} else {
		return fail
	}
}

// 用时
func UseTime(title string, fn func(int), n ...int) string {
	start := time.Now()
	num := 1
	if len(n) > 0 {
		num = n[0]
	}
	for i := 0; i < num; i++ {
		fn(i)
	}
	cost := time.Since(start)
	return title + "执行" + strconv.Itoa(num) + "次用时：" + cost.String()
}

func ToStr(params ...any) string {
	len := len(params)
	var key string
	for i := 0; i < len; i++ {
		param := params[i]
		switch t := param.(type) {
		case string:
			key += t + "_"
		case int8:
			key += num.Int8ToStr(t) + "_"
		case int:
			key += num.IntToStr(t) + "_"
		case int32:
			key += num.IntToStr(int(t)) + "_"
		case int64:
			key += num.Int64ToStr(t) + "_"
		case uint8:
			key += num.UIntToStr(uint(t)) + "_"
		case uint:
			key += num.UIntToStr(t) + "_"
		case uint32:
			key += num.UIntToStr(uint(t)) + "_"
		case uint64:
			key += num.Int64ToStr(int64(t)) + "_"
		case float32:
			key += strings.ReplaceAll(num.FloatToStr(t), ".00", "") + "_"
		case float64:
			key += strings.ReplaceAll(num.Float64ToStr(t), ".00", "") + "_"
		case bool:
			if bool(t) {
				key += "True_"
			} else {
				key += "False_"
			}
		}
	}
	return strings.TrimSuffix(key, "_")
}