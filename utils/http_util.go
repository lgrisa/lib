package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ContentTypeXWwwFormUrlencodedText = "application/x-www-form-urlencoded"
	ContentTypeJsonText               = "application/json;charset=UTF-8"
)

func BuildHttpBody(contentType string, dataMap map[string]interface{}) (io.Reader, error) {
	var requestBody io.Reader

	if contentType == ContentTypeXWwwFormUrlencodedText {
		var slice1 []string
		for k, v := range dataMap {
			value, ok := v.(string)
			if ok {
				slice1 = append(slice1, k+"="+url.QueryEscape(value))
			}
		}

		requestBody = strings.NewReader(strings.Join(slice1, "&"))
	} else if contentType == ContentTypeJsonText {
		b, err := json.Marshal(dataMap)
		if err != nil {
			return nil, err
		}

		requestBody = bytes.NewBuffer(b)
	}

	return requestBody, nil
}

func Request(ctx context.Context, url string, method string, headers map[string]string, requestBody io.Reader) (responseHeader http.Header, returnCode int, body []byte, err error) {
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, 0, nil, errors.Errorf("http.NewRequest fail: %v", err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	resp, err := ctxhttp.Do(ctx, nil, req)
	if err != nil {
		return nil, 0, nil, errors.Errorf("cxhttp.Do fail: %v", err)
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			LogErrorF("Body.Close fail: %v", err)
		}
	}(resp.Body)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, nil, errors.Errorf("io.ReadAll fail: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil, errors.Errorf("http status code: %d", resp.StatusCode)
	}

	return resp.Header, resp.StatusCode, body, nil
}
