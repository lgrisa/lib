package utils

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// NewRequest 请求包装
func NewRequest(method, url string, data []byte) (body []byte, err error) {

	if method == "GET" {
		url = fmt.Sprint(url, "?", string(data))
		data = nil
	}

	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return body, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return body, err
	}

	return body, err
}

func SendRequest(ctx context.Context, method string, uri string, header map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// set Header
	if len(header) > 0 {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	// with context
	req = req.WithContext(ctx)

	// 发送请求
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func HttpRequest(ctx context.Context, method string, url string, headers map[string]string, requestBody io.Reader) (responseHeader http.Header, body []byte, err error, returnCode int) {
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, nil, err, 0
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	resp, err := ctxhttp.Do(ctx, nil, req)
	if err != nil {
		return nil, nil, err, 0
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err, 0
	}

	return resp.Header, body, nil, resp.StatusCode
}
