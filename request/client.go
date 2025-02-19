package request

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Client struct {
	httpclient *http.Client
	timeout    time.Duration
}

func CreateClient(timeout time.Duration) *Client {
	client := &Client{timeout: timeout, httpclient: &http.Client{Timeout: timeout}}
	return client
}

func (a *Client) GetRepetition(url string, params []QueryParameter, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParamsRepetition(params)
	return a.Request("GET", fullurl, nil, headers...)
}
func (a *Client) Get(url string, params map[string]any, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("GET", fullurl, nil, headers...)
}
func (a *Client) Delete(url string, params map[string]any, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("DELETE", fullurl, nil, headers...)
}

func (a *Client) Post(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("POST", fullurl, strings.NewReader(body), headers...)
}
func (a *Client) PostJson(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	headers = append(headers, Header{"Content-Type", "application/json"})
	return a.Request("POST", fullurl, strings.NewReader(body), headers...)
}

func (a *Client) Put(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("PUT", fullurl, strings.NewReader(body), headers...)
}
func (a *Client) PutJson(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	headers = append(headers, Header{"Content-Type", "application/json"})
	return a.Request("PUT", fullurl, strings.NewReader(body), headers...)
}

func (a *Client) Del(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	return a.Request("DELETE", fullurl, strings.NewReader(body), headers...)
}
func (a *Client) DelJson(url string, params map[string]any, body string, headers ...Header) (string, error) {
	fullurl := url + ConvertToQueryParams(params)
	headers = append(headers, Header{"Content-Type", "application/json"})
	return a.Request("DELETE", fullurl, strings.NewReader(body), headers...)
}

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

func (a *Client) UploadRequest(uri, fieldname string, params map[string]string, fullpath, name string, headers ...Header) (string, error) {
	if len(name) == 0 {
		name = fullpath
	}
	file, err := os.Open(filepath.Clean(fullpath))
	if err != nil {
		return "", err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldname, name)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	if params == nil {
		params = make(map[string]string)
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", uri, body)
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
	result, err := responseHandle(resp, err)
	return result, err
}

func (a *Client) UploadFiles(uri, fieldname string, params map[string]string, files *multipart.FileHeader, headers ...Header) (string, error) {
	file, err := files.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldname, files.Filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	if params == nil {
		params = make(map[string]string)
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", uri, body)
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
	result, err := responseHandle(resp, err)
	return result, err
}

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
