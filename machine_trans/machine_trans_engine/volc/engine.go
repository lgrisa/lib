package volc

import (
	"context"
	"encoding/json"
	"fmt"
	trans "github.com/lgrisa/lib/machine_trans/machine_trans_engine"
	"github.com/lgrisa/lib/utils/logutil"
	"github.com/pkg/errors"
	"github.com/volcengine/volc-sdk-golang/base"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewClient(appId, appSecret string) *Engine {
	serviceInfo := &base.ServiceInfo{
		Timeout: 5 * time.Second,
		Host:    host,
		Header: http.Header{
			"Accept": []string{"application/json"},
		},
		Credentials: base.Credentials{Region: base.RegionCnNorth1, Service: "translate"},
	}
	apiInfoList := map[string]*base.ApiInfo{
		"TranslateText": {
			Method: http.MethodPost,
			Path:   "/",
			Query: url.Values{
				"Action":  []string{"TranslateText"},
				"Version": []string{kServiceVersion},
			},
		},
	}

	client := base.NewClient(serviceInfo, apiInfoList)
	client.SetAccessKey(appId)
	client.SetSecretKey(appSecret)

	return &Engine{
		client:  client,
		limiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(5)), 5),
	}
}

func (e *Engine) TranslateFor(text string, fromLanguage, toLanguage trans.LanguageType) (string, error) {
	// 如果出现换行符，替换为\n
	text = strings.ReplaceAll(text, "\n", "\\n")

	for {
		logutil.LogTraceF("火山翻译:%s", text)

		if res, err := e.Translate(text, fromLanguage, toLanguage); err == nil {

			isInTag := false

			length := len(res)

			trimSpaceStr := []string{` \ n`, `\n `, `\ n`}

			for i := 0; i < length; i++ {

				//删除<>内的空格
				if res[i] == '<' {
					isInTag = true
				} else if res[i] == '>' {
					isInTag = false
				} else if res[i] == ' ' && isInTag {
					res = res[:i] + res[i+1:]
					i--
					length--
				}

				for _, str := range trimSpaceStr {
					lenStr := len(str)

					trimStr := strings.ReplaceAll(str, " ", "")

					if i+lenStr < length && res[i:i+lenStr] == str {
						res = res[:i] + trimStr + res[i+lenStr:]
						i--
						length -= lenStr - len(trimStr)
					}
				}
			}

			return res, nil
		} else {

			logutil.LogErrorF("火山翻译(%s)错误:(%v)", text, err)

			// 如果错误包含 超过了每秒频率上限，那么等待1s
			if strings.Contains(err.Error(), "超过了每秒频率上限") ||
				strings.Contains(err.Error(), "limit") {

			} else {
				return "", err
			}

			time.Sleep(time.Second)
		}
	}
}

func (e *Engine) Translate(text string, fromLanguage, toLanguage trans.LanguageType) (string, error) {
	if e.client == nil {
		return "", errors.New("火山引擎未初始化")
	}

	if err := e.limiter.Wait(context.Background()); err != nil {
		return "", errors.Errorf("限流器等待错误: %v", err)
	}

	sourceLanguage := e.languageCode(fromLanguage)
	targetLanguage := e.languageCode(toLanguage)
	request := Req{
		SourceLanguage: sourceLanguage,
		TargetLanguage: targetLanguage,
		TextList:       []string{text},
	}

	requestJson, _ := json.Marshal(request)

	resp, code, err := e.client.Json("TranslateText", nil, string(requestJson))
	if err != nil {
		return "", errors.WithStack(err)
	}
	if code != 200 {
		return "", fmt.Errorf("火山翻译错误，返回错误码:%d", code)
	}
	data := translationData{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return "", err
	}

	if data.ResponseMetadata.Error.Code != "" {
		return "", fmt.Errorf("火山翻译错误，返回错误码:%s，错误原因:%s", data.ResponseMetadata.Error.Code, data.ResponseMetadata.Error.Message)
	}

	if len(data.TranslationList) == 0 {
		return "", fmt.Errorf("火山翻译错误，返回内容为空")
	}

	return data.TranslationList[0].Translation, nil
}
