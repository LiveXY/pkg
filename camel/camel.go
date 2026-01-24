package camel

import (
	"bytes"
	"strings"
	"unicode"
)

var commonInitialisms = []string{"ACL", "API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "TCP", "TLS", "TTL", "UDP", "UI", "UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XMPP", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer
var uncommonInitialismsReplacer *strings.Replacer

func init() {
	var commonInitialismsForReplacer []string
	var uncommonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		lower := strings.ToLower(initialism)
		title := string(unicode.ToUpper(rune(lower[0]))) + lower[1:]
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, title)
		uncommonInitialismsForReplacer = append(uncommonInitialismsForReplacer, title, initialism)
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
	uncommonInitialismsReplacer = strings.NewReplacer(uncommonInitialismsForReplacer...)
}

// BigCamel 将下划线命名转换为大驼峰命名（UpperCamelCase）
func BigCamel(name string) string {
	if name == "" {
		return ""
	}
	temp := strings.Split(name, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if vv[0] >= 'a' && vv[0] <= 'z' {
				vv[0] -= 32
			}
			s += string(vv)
		}
	}
	s = uncommonInitialismsReplacer.Replace(s)
	return s
}

// SmallCamel 将命名转换为小驼峰命名（lowerCamelCase）
func SmallCamel(name string) string {
	if name == "" {
		return ""
	}
	value := commonInitialismsReplacer.Replace(name)
	strs := []rune(value)
	if strs[0] >= 'A' && strs[0] <= 'Z' {
		strs[0] = strs[0] + 32
	}
	return string(strs)
}

// UnBigCamel 将驼峰命名转换回下划线命名（snake_case）
func UnBigCamel(name string) string {
	const (
		lower = false
		upper = true
	)
	if name == "" {
		return ""
	}
	var (
		value                                    = commonInitialismsReplacer.Replace(name)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)
	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] >= 'A' && value[i+1] <= 'Z'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}
	buf.WriteByte(value[len(value)-1])
	s := strings.ToLower(buf.String())
	return s
}
