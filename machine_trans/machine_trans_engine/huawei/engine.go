package huawei

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	nlp "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/region"
	trans "github.com/lgrisa/lib/machine_trans/machine_trans_engine"
	"github.com/pkg/errors"
)

func NewClient(ak, sk string) (*Client, error) {

	auth, err := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		SafeBuild()

	if err != nil {
		return nil, errors.Errorf("NewCredentialsBuilder fail, ak: %s, sk: %s, error: %v", ak, sk, err)
	}

	curRegion, err := region.SafeValueOf("cn-north-4")

	if err != nil {
		return nil, errors.Errorf("SafeValueOf fail, region: %s, error: %v", "cn-north-4", err)
	}

	hcHttpClient, err := nlp.NlpClientBuilder().
		WithRegion(curRegion).
		WithCredential(auth).
		SafeBuild()

	if err != nil {
		return nil, errors.Errorf("NlpClientBuilder fail, ak: %s, sk: %s, error: %v", ak, sk, err)
	}

	client := nlp.NewNlpClient(hcHttpClient)

	return &Client{
		ak: ak,
		sk: sk,

		client: client,
	}, nil
}

func (c *Client) TranslateFor(text string, fromLanguage, toLanguage trans.LanguageType) (string, error) {
	for {
		transResult, err := c.Translate(text, fromLanguage, toLanguage)
		if err != nil {
			return "", err
		} else {
			return transResult, nil
		}
	}
}

func (c *Client) Translate(text string, fromLanguage, toLanguage trans.LanguageType) (string, error) {

	from, err := transFromLanguage(fromLanguage)

	if err != nil {
		return "", fmt.Errorf("transFromLanguage Err: (%v)", err)
	}

	to, err := transToLanguage(toLanguage)

	if err != nil {
		return "", fmt.Errorf("transToLanguage Err: (%v)", err)
	}

	request := &model.RunTextTranslationRequest{}
	request.Body = &model.TextTranslationReq{
		To:   to,
		From: from,
		Text: text,
	}
	response, err := c.client.RunTextTranslation(request)

	if err != nil {
		return "", fmt.Errorf("RunTextTranslation Err: (%v)", err)
	}

	if response.ErrorCode != nil {
		return "", fmt.Errorf("ErrorCode Err: (%v)", *response.ErrorCode)
	}

	return *response.TranslatedText, nil
}
