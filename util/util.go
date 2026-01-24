package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/livexy/pkg/num"
)

// FormatFileSize 格式化输出文件大小（B, KB, MB, GB, TB, EB）
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

// FormatUseTime 格式化输出用时（s, m, h, d）
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

// FormatCNUseTime 格式化中文输出用时（秒, 分, 时, 天）
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

// UseTime 测试带索引参数的函数执行耗时
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

// ToStr 将多个参数转换为以下划线连接的字符串
func ToStr(params ...any) string {
	var key strings.Builder
	for _, param := range params {
		switch t := param.(type) {
		case string:
			key.WriteString(t + "_")
		case int8:
			key.WriteString(num.Int8ToStr(t) + "_")
		case int:
			key.WriteString(num.IntToStr(t) + "_")
		case int32:
			key.WriteString(num.IntToStr(int(t)) + "_")
		case int64:
			key.WriteString(num.Int64ToStr(t) + "_")
		case uint8:
			key.WriteString(num.UIntToStr(uint(t)) + "_")
		case uint:
			key.WriteString(num.UIntToStr(t) + "_")
		case uint32:
			key.WriteString(num.UIntToStr(uint(t)) + "_")
		case uint64:
			key.WriteString(num.UInt64ToStr(t) + "_")
		case float32:
			key.WriteString(strings.ReplaceAll(num.FloatToStr(t), ".00", "") + "_")
		case float64:
			key.WriteString(strings.ReplaceAll(num.Float64ToStr(t), ".00", "") + "_")
		case bool:
			if bool(t) {
				key.WriteString("True_")
			} else {
				key.WriteString("False_")
			}
		}
	}
	return strings.TrimSuffix(key.String(), "_")
}
