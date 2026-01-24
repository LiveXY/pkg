package request

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/livexy/pkg/logx"

	"go.uber.org/zap"
)

var (
	defaultClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          1000,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   100,
			MaxConnsPerHost:       0,
		},
		Timeout: 30 * time.Second,
	}
)

// Get 发起 GET 请求并返回响应字符串
func Get(api string, timeout int, headers ...Header) (string, error) {
	logx.Logger.Debug("GET请求接口", zap.String("api", api))
	client := defaultClient
	if timeout > 0 {
		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		logx.Error.Error("GET请求数据错误：", zap.String("api", api), zap.Error(err))
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	for _, header := range headers {
		req.Header.Set(header.Name, header.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		logx.Error.Error("GET请求数据错误：", zap.String("api", api), zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Error.Error("GET请求数据错误：", zap.String("api", api), zap.Error(err))
		return "", err
	}
	data := string(body)
	if strings.HasPrefix(data, "<") {
		logx.Error.Error("GET请求数据错误：返回了HTML", zap.String("api", api), zap.String("data", data))
		return "", fmt.Errorf("api returns html content: %w", http.ErrNoLocation)
	}
	logx.Logger.Debug("GET请求数据", zap.String("data", data))
	return data, nil
}

// Post 发起 POST 请求并返回响应字符串
func Post(api string, body string, timeout int, headers ...Header) (string, error) {
	logx.Logger.Debug("POST请求接口", zap.String("api", api))
	client := defaultClient
	if timeout > 0 {
		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}
	req, err := http.NewRequest("POST", api, strings.NewReader(body))
	if err != nil {
		logx.Error.Error("POST请求数据错误：", zap.String("api", api), zap.Error(err))
		return "", err
	}
	for _, header := range headers {
		req.Header.Set(header.Name, header.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		logx.Error.Error("POST请求数据错误：", zap.String("api", api), zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Error.Error("POST请求数据错误：", zap.String("api", api), zap.Error(err))
		return "", err
	}
	data := string(b)
	if strings.HasPrefix(data, "<") {
		logx.Error.Error("POST请求数据错误：返回了HTML", zap.String("api", api), zap.String("data", data))
		return "", fmt.Errorf("api returns html content")
	}
	logx.Logger.Debug("POST请求数据", zap.String("data", data))
	return data, nil
}

// PostJSON 发起带有 Content-Type: application/json 的 POST 请求
func PostJSON(api string, body string, timeout int, headers ...Header) (string, error) {
	headers = append(headers, Header{Name: "Content-Type", Value: "application/json"})
	return Post(api, body, timeout, headers...)
}
