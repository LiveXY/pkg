package ipx

import (
	"net"
	"strings"
	"sync"

	"github.com/livexy/pkg/logx"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
	"go.uber.org/zap"
)

var searcher *service.Ip2Region

// Init 初始化 IP2Region 地址库
func Init(dbfile string) {
	v4Config, err := service.NewV4Config(service.VIndexCache, dbfile, 20)
	if err != nil {
		logx.Error.Error("未配置IPv4地址库："+dbfile, zap.Error(err))
		panic("未配置IP地址库->" + dbfile)
	}
	v6Config, err := service.NewV6Config(service.VIndexCache, dbfile+".v6", 20)
	if err != nil {
		logx.Error.Error("未配置IPv6地址库："+dbfile+".v6", zap.Error(err))
		panic("未配置IP地址库->" + dbfile + ".v6")
	}
	searcher, err = service.NewIp2Region(v4Config, v6Config)
	if err != nil {
		logx.Error.Error("未配置IP地址库：", zap.Error(err))
		panic("未配置IP地址库->" + dbfile)
	}
}

// Close 关闭地址库资源
func Close() {
	if searcher != nil {
		searcher.Close()
	}
}

// GetIPAddress 根据 IP 查询具体的物理地址描述
func GetIPAddress(ip string) string {
	address, _ := searcher.SearchByStr(ip)
	add := strings.Split(address, "|")
	if len(add) < 4 || add[2] == "0" && add[3] == "0" {
		return ""
	}
	if add[2] == "0" && add[3] != "0" {
		return add[3]
	}
	if add[3] == "0" && add[2] != "0" {
		return add[2]
	}
	return add[2] + " " + add[3]
}

var localIP string
var once2 sync.Once

// GetPrivateIPv4 获取当前主机的私有 IPv4 地址
func GetPrivateIPv4() string {
	once2.Do(func() {
		localIP = getPrivateIPv4()
	})
	return localIP
}

func getPrivateIPv4() string {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		ip := ipnet.IP.To4()
		if ip != nil && (ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168) {
			return ip.String()
		}
	}
	return "127.0.0.1"
}
