package check

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/livexy/pkg/bytex"
	"github.com/livexy/pkg/strx"
)

// 用户名
func IsUserName(userName string) bool {
	len := len(userName)
	if len < 3 {
		return false
	}
	match, err := regexp.MatchString(`^[0-9@A-Z_a-z]*$`, userName)
	if err != nil {
		return false
	}
	return match
}

// 邮箱检测
func IsEmail(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z].){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// BASE64
func IsBase64(data string) bool {
	if len(data) == 0 {
		return false
	}
	pattern := `^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(data)
}
func IsAlpha(data string) bool {
	if data == "" {
		return false
	}
	for _, v := range data {
		if !unicode.IsLetter(v) {
			return false
		}
	}
	return true
}
func IsAlphanumeric(data string) bool {
	if data == "" {
		return false
	}
	for _, v := range data {
		if !(unicode.IsDigit(v) || unicode.IsLetter(v)) {
			return false
		}
	}
	return true
}
func IsNumeric(data string) bool {
	if data == "" {
		return false
	}
	for _, v := range data {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
func IsBase64Char(data string) bool {
	if len(data) == 0 {
		return false
	}
	bs, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return false
	}
	return IsAlphanumeric(bytex.ToStr(bs))
}

func IsWinPath(data string) bool {
	if len(data) == 0 {
		return false
	}
	pattern := `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(data) {
		return false
	}
	if len(data) > 3 && len(data[3:]) > 32767 {
		return false
	}
	return true
}
func IsUnixPath(data string) bool {
	pattern := `^((?:\/[a-zA-Z0-9\.\:]+(?:_[a-zA-Z0-9\:\.]+)*(?:\-[\:a-zA-Z0-9\.]+)*)+\/?)$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(data)
}
func IsSemver(data string) bool {
	pattern := "^v?(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(data)
}
func IsUrl(data string) bool {
	if data == "" || len(data) >= 2083 || len(data) <= 3 || strings.HasPrefix(data, ".") {
		return false
	}
	u, err := url.Parse(data)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	pattern := `^(\/\/|https?:\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(data)
}
func IsUUID(data string) bool {
	pattern := "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(data)
}
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal(strx.ToBytes(str), &js) == nil
}
func IsDataURI(str string) bool {
	dataURI := strings.Split(str, ",")
	pattern := "^data:.+\\/(.+);base64$"
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(dataURI[0]) {
		return false
	}
	return IsBase64(dataURI[1])
}

// 手机号检测
func IsMobile(mobile string) bool {
	regular := `^1\d{10}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}

// IP地址
func IsIP(ip string) bool {
	match, err := regexp.MatchString(`(([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).){3}([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])`, ip)
	if err != nil {
		return false
	}
	return match
}

// 检测是否本机IP
func IsInternalIPv4(ip string) bool {
	if len(ip) == 0 {
		return false
	}
	if strings.HasPrefix(ip, "127.") {
		return true
	}
	if strings.HasPrefix(ip, "10.") {
		return true
	}
	if strings.HasPrefix(ip, "192.168.") {
		return true
	}
	if strings.HasPrefix(ip, "169.254.") {
		return true
	}
	if strings.HasPrefix(ip, "172.") {
		ips := strings.Split(ip, ".")
		ip2 := strx.ToInt(ips[1])
		if ip2 >= 16 && ip2 <= 31 {
			return true
		}
	}
	return false
}

func IsPrivateIPv4(ip net.IP) bool {
	return ip != nil && (ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

// 姓名
func IsCName(nickName string) bool {
	reg := regexp.MustCompile("^[\u4e00-\u9fa5]{2,8}$")
	return reg.MatchString(nickName)
}

// 密码强度
func IsStrongPass(pass string) int {
	if len(pass) < 6 {
		return 0
	}
	level := 0
	list := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range list {
		match, _ := regexp.MatchString(pattern, pass)
		if match {
			level++
		}
	}
	return level
}

// 是否数字
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// 是否整数
func IsInt(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// 是否日期
func IsDate(str string) bool {
	var format string
	switch len(str) {
	case 8:
		format = "20060102"
	case 10:
		format = "2006-01-02"
	case 12:
		format = "200601021504"
	case 14:
		format = "20060102150405"
	case 17:
		format = "2006-01-02 15:04"
	case 19:
		format = "2006-01-02 15:04:05"
	}
	if _, err := time.Parse(format, str); err != nil {
		return false
	}
	return true
}
