package huawei

import nlp "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2"

type Client struct {
	ak     string
	sk     string
	client *nlp.NlpClient
}
