package ipx

import (
	"net"
	"strings"
	"sync"

	"github.com/livexy/pkg/logx"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/zap"
)

var searcher *xdb.Searcher

// 初始化IP库
func Init(dbfile string) {
	cbuff, err := xdb.LoadContentFromFile(dbfile)
	if err != nil {
		logx.Error.Error("未配置IP地址库：", zap.Error(err))
		panic("未配置IP地址库->" + dbfile)
	}
	searcher, err = xdb.NewWithBuffer(cbuff)
	if err != nil {
		logx.Error.Error("未配置IP地址库：", zap.Error(err))
		panic("未配置IP地址库->" + dbfile)
	}
}

// 获取IP地址
func GetIPAddress(ip string) string {
	address, _ := searcher.SearchByStr(ip)
	add := strings.Split(address, "|")
	if len(add) < 4 || add[2] == "0" && add[3] == "0" {
		return ""
	}
	if add[2] == "0" {
		return add[3]
	}
	return add[2] + " " + add[3]
}

var localIP string
var once2 sync.Once

// 获取本机IP
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
