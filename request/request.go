package request

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/livexy/pkg/logx"

	"go.uber.org/zap"
)

// GET请求数据
func HttpGet(api string, timeout int, headers ...Header) (string, error) {
	logx.Logger.Debug("GET请求接口", zap.String("api", api))
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
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
		logx.Error.Error("GET请求数据错误：", zap.String("api", api))
		return "", err
	}
	logx.Logger.Debug("GET请求数据", zap.String("data", data))
	return data, nil
}

// POST请求数据
func HttpPost(api string, body string, timeout int, headers ...Header) (string, error) {
	logx.Logger.Debug("POST请求接口", zap.String("api", api))
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
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
		logx.Error.Error("POST请求数据错误：", zap.String("api", api))
		return "", err
	}
	logx.Logger.Debug("POST请求数据", zap.String("data", data))
	return data, nil
}
func HttpPostJson(api string, body string, timeout int, headers ...Header) (string, error) {
	headers = append(headers, Header{ Name: "Content-Type", Value: "application/json" })
	return HttpPost(api, body, timeout, headers...)
}
