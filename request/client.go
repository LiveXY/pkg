package request

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Client HTTP 客户端结构体
type Client struct {
	httpclient *http.Client
	timeout    time.Duration
}

// CreateClient 创建并初始化 HTTP 客户端
func CreateClient(timeout time.Duration) *Client {
	client := &Client{timeout: timeout, httpclient: &http.Client{Timeout: timeout}}
	return client
}

// GetWithParams 发起带参数请求的 GET 请求
func (a *Client) GetWithParams(url string, params []QueryParameter, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParamsRepetition(params)
	return a.Request("GET", fullurl, nil, headers...)
}

// Get 发起带参数的 GET 请求
func (a *Client) Get(url string, params map[string]any, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("GET", fullurl, nil, headers...)
}

// Delete 发起 DELETE 请求
func (a *Client) Delete(url string, params map[string]any, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("DELETE", fullurl, nil, headers...)
}

// Post 发起 POST 请求
func (a *Client) Post(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("POST", fullurl, strings.NewReader(body), headers...)
}

// PostJSON 发起发送 JSON 数据的 POST 请求
func (a *Client) PostJSON(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	headers = append(headers, Header{"Content-Type", "application/json"})
	return a.Request("POST", fullurl, strings.NewReader(body), headers...)
}

// Put 发起 PUT 请求
func (a *Client) Put(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("PUT", fullurl, strings.NewReader(body), headers...)
}

// PutJSON 发起发送 JSON 数据的 PUT 请求
func (a *Client) PutJSON(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	headers = append(headers, Header{"Content-Type", "application/json"})
	return a.Request("PUT", fullurl, strings.NewReader(body), headers...)
}

// DeleteWithBody 发起带有正文的 DELETE 请求
func (a *Client) DeleteWithBody(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("DELETE", fullurl, strings.NewReader(body), headers...)
}

// DeleteJSON 发起发送 JSON 数据的 DELETE 请求
func (a *Client) DeleteJSON(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	headers = append(headers, Header{"Content-Type", "application/json"})
	return a.Request("DELETE", fullurl, strings.NewReader(body), headers...)
}

// Request 执行通用的 HTTP 请求
func (a *Client) Request(method, url string, body io.Reader, headers ...Header) (string, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	for _, header := range headers {
		req.Header.Set(header.Name, header.Value)
	}
	resp, err := a.httpclient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := responseHandle(resp, err)
	return result, err
}

// UploadRequest 上传文件请求
func (a *Client) UploadRequest(uri, fieldname string, params map[string]string, fullpath, name string, headers ...Header) (string, error) {
	if len(name) == 0 {
		name = fullpath
	}
	file, err := os.Open(filepath.Clean(fullpath))
	if err != nil {
		return "", err
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer file.Close()
		defer pw.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile(fieldname, name)
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		if _, err = io.Copy(part, file); err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		for key, val := range params {
			if err = writer.WriteField(key, val); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
	}()

	req, err := http.NewRequest("POST", uri, pr)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	for _, header := range headers {
		req.Header.Set(header.Name, header.Value)
	}
	resp, err := a.httpclient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return responseHandle(resp, err)
}

// UploadFiles 上传多文件请求
func (a *Client) UploadFiles(uri, fieldname string, params map[string]string, files *multipart.FileHeader, headers ...Header) (string, error) {
	file, err := files.Open()
	if err != nil {
		return "", err
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer file.Close()
		defer pw.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile(fieldname, files.Filename)
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		if _, err = io.Copy(part, file); err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		for key, val := range params {
			if err = writer.WriteField(key, val); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
	}()

	req, err := http.NewRequest("POST", uri, pr)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	for _, header := range headers {
		req.Header.Set(header.Name, header.Value)
	}
	resp, err := a.httpclient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return responseHandle(resp, err)
}

// ResetParams 重置并过滤参数
func (a *Client) ResetParams(params any) map[string]any {
	out, jout := make(map[string]any), make(map[string]any)
	bs, err := json.Marshal(params)
	if err != nil {
		return out
	}
	err = json.Unmarshal(bs, &jout)
	if err != nil {
		return out
	}
	keys := []string{}
	for k := range jout {
		keys = append(keys, k)
	}
	for _, v := range keys {
		p, ok := jout[v].(string)
		if !ok {
			out[v] = jout[v]
		}
		if ok && len(p) != 0 {
			out[v] = jout[v]
		}
	}
	return out
}

func responseHandle(resp *http.Response, err error) (string, error) {
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBody := string(b)
	return respBody, nil
}
