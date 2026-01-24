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

var (
	reUserName     = regexp.MustCompile(`^[0-9@A-Z_a-z]*$`)
	reEmail        = regexp.MustCompile(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z].){1,4}[a-z]{2,4}$`)
	reBase64       = regexp.MustCompile(`^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`)
	reWinPath      = regexp.MustCompile(`^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`)
	reUnixPath     = regexp.MustCompile(`^((?:\/[a-zA-Z0-9\.\:]+(?:_[a-zA-Z0-9\:\.]+)*(?:\-[\:a-zA-Z0-9\.]+)*)+\/?)$`)
	reSemver       = regexp.MustCompile(`^v?(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$`)
	reUrl          = regexp.MustCompile(`^(\/\/|https?:\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[01]\d|22[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`)
	reUUID         = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	reDataURI      = regexp.MustCompile(`^data:.+\/(.+);base64$`)
	reMobile       = regexp.MustCompile(`^1[3-9]\d{9}$`)
	reIP           = regexp.MustCompile(`(([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).){3}([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])`)
	reChineseName  = regexp.MustCompile(`^[\x{4e00}-\x{9fa5}]{2,8}$`)
	reStrongPasses = []*regexp.Regexp{
		regexp.MustCompile(`[0-9]+`),
		regexp.MustCompile(`[a-z]+`),
		regexp.MustCompile(`[A-Z]+`),
		regexp.MustCompile(`[~!@#$%^&*?_-]+`),
	}
)

// IsUserName 校验用户名格式
func IsUserName(userName string) bool {
	if len(userName) < 3 {
		return false
	}
	return reUserName.MatchString(userName)
}

// IsEmail 校验邮箱格式
func IsEmail(email string) bool {
	return reEmail.MatchString(email)
}

// IsBase64 校验是否为 Base64 编码
func IsBase64(data string) bool {
	if len(data) == 0 {
		return false
	}
	return reBase64.MatchString(data)
}

// IsAlpha 校验是否全由字母组成
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

// IsAlphanumeric 校验是否由字母和数字组成
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

// IsNumeric 校验是否全由数字组成
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

// IsBase64Char 校验 Base64 解码后的内容是否为字母数字
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

// IsWinPath 校验是否为 Windows 路径格式
func IsWinPath(data string) bool {
	if len(data) == 0 {
		return false
	}
	if !reWinPath.MatchString(data) {
		return false
	}
	if len(data) > 3 && len(data[3:]) > 32767 {
		return false
	}
	return true
}

// IsUnixPath 校验是否为 Unix 路径格式
func IsUnixPath(data string) bool {
	return reUnixPath.MatchString(data)
}

// IsSemver 校验是否为语义化版本号格式
func IsSemver(data string) bool {
	return reSemver.MatchString(data)
}

// IsUrl 校验是否为 URL 格式
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
	return reUrl.MatchString(data)
}

// IsUUID 校验是否为 UUID 格式
func IsUUID(data string) bool {
	return reUUID.MatchString(data)
}

// IsJSON 校验是否为有效的 JSON 字符串
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal(strx.ToBytes(str), &js) == nil
}

// IsDataURI 校验是否为有效的 Data URI 格式
func IsDataURI(str string) bool {
	dataURI := strings.Split(str, ",")
	if len(dataURI) < 2 {
		return false
	}
	if !reDataURI.MatchString(dataURI[0]) {
		return false
	}
	return IsBase64(dataURI[1])
}

// IsMobile 校验是否为中国大陆手机号格式
func IsMobile(mobile string) bool {
	return reMobile.MatchString(mobile)
}

// IsIP 校验是否为 IP 地址格式
func IsIP(ip string) bool {
	return reIP.MatchString(ip)
}

// IsInternalIPv4 校验是否为内部或本地回环 IPv4 地址
func IsInternalIPv4(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}

// IsPrivateIPv4 校验是否为私有网络 IPv4 地址
func IsPrivateIPv4(ip net.IP) bool {
	return ip != nil && ip.IsPrivate()
}

// IsCName 校验是否为 2-8 个汉字的中文姓名格式
func IsCName(nickName string) bool {
	return reChineseName.MatchString(nickName)
}

// IsStrongPass 校验密码强度并返回得分（数字、小写、大写、特殊字符各占1分）
func IsStrongPass(pass string) int {
	if len(pass) < 6 {
		return 0
	}
	level := 0
	for _, reg := range reStrongPasses {
		if reg.MatchString(pass) {
			level++
		}
	}
	return level
}

// IsNum 校验字符串是否可以解析为浮点数
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsInt 校验字符串是否可以解析为整数
func IsInt(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// IsDate 校验字符串是否符合常见的日期格式
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
