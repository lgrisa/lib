package volc

import (
	"github.com/volcengine/volc-sdk-golang/base"
	"golang.org/x/time/rate"
)

const (
	host            = "open.volcengineapi.com"
	kServiceVersion = "2020-06-01"
)

type Engine struct {
	client  *base.Client
	limiter *rate.Limiter
}

type translationData struct {
	TranslationList []struct {
		Translation            string `json:"Translation"`
		DetectedSourceLanguage string `json:"DetectedSourceLanguage"`
	} `json:"TranslationList"`
	ResponseMetadata struct {
		RequestID string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Error     struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
	} `json:"ResponseMetadata"`
}

type Req struct {
	SourceLanguage string   `json:"SourceLanguage"`
	TargetLanguage string   `json:"TargetLanguage"`
	TextList       []string `json:"TextList"`
}
