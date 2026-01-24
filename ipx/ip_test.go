package ipx

import (
	"testing"
)

// TestGetPrivateIPv4 测试获取本机私有 IPv4 地址的功能
func TestGetPrivateIPv4(t *testing.T) {
	ip := GetPrivateIPv4()
	if ip == "" {
		t.Errorf("GetPrivateIPv4 返回了空字符串")
	}
	// 简单核查是否为 127.0.0.1 或私有 IP 段
	if ip != "127.0.0.1" && !isPrivate(ip) {
		t.Errorf("GetPrivateIPv4 返回了非私有 IP：%v", ip)
	}
}

func isPrivate(ip string) bool {
	// Simple check for private IP ranges
	return (len(ip) >= 3 && ip[0:3] == "10.") ||
		(len(ip) >= 7 && ip[0:7] == "192.168") ||
		(len(ip) >= 4 && ip[0:4] == "172.")
}
