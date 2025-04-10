package utils

import (
	"io"
	"net/http"
	"net/url"
)

// DoGET 发送GET请求
func DoGET(targetUrl string, values map[string]string) (string, error) {
	// 解析的URL，可以拆分或添加各部分内容，Scheme（协议）、Host（主机）、Path（路径）、RawQuery（查询字符串）等
	parsedUrl, err := url.ParseRequestURI(targetUrl)
	if err != nil {
		return "", err
	}
	urlValues := url.Values{}
	for k, v := range values {
		urlValues.Set(k, v)
	}
	parsedUrl.RawQuery = urlValues.Encode()
	// http请求获取返回信息
	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		return "", err
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
