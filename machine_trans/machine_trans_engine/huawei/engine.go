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
		return "", err
	}

	to, err := transToLanguage(toLanguage)

	if err != nil {
		return "", err
	}

	request := &model.RunTextTranslationRequest{}
	request.Body = &model.TextTranslationReq{
		To:   to,
		From: from,
		Text: text,
	}
	response, err := c.client.RunTextTranslation(request)

	if err != nil {
		return "", err
	}

	if response.ErrorCode != nil {
		return "", fmt.Errorf("error code: %v, error msg: %v", *response.ErrorCode, *response.ErrorMsg)
	}

	return *response.TranslatedText, nil
}

func transFromLanguage(fromLanguage trans.LanguageType) (model.TextTranslationReqFrom, error) {
	switch fromLanguage {
	case trans.LanguageTypeZh:
		return model.GetTextTranslationReqFromEnum().ZH, nil
	case trans.LanguageTypeZhTW:
		return model.GetTextTranslationReqFromEnum().ZH_TW, nil
	case trans.LanguageTypeEn:
		return model.GetTextTranslationReqFromEnum().EN, nil
	case trans.LanguageTypeDe:
		return model.GetTextTranslationReqFromEnum().DE, nil
	case trans.LanguageTypeFr:
		return model.GetTextTranslationReqFromEnum().FR, nil
	case trans.LanguageTypePt:
		return model.GetTextTranslationReqFromEnum().PT, nil
	case trans.LanguageTypeEs:
		return model.GetTextTranslationReqFromEnum().ES, nil
	case trans.LanguageTypeTh:
		return model.GetTextTranslationReqFromEnum().TH, nil
	default:
		return model.TextTranslationReqFrom{}, fmt.Errorf("fromLanguage: %s is not supported", fromLanguage)
	}
}

func transToLanguage(toLanguage trans.LanguageType) (model.TextTranslationReqTo, error) {
	switch toLanguage {
	case trans.LanguageTypeZh:
		return model.GetTextTranslationReqToEnum().ZH, nil
	//case trans.LanguageTypeZhTW: // 华为不支持繁体中文
	//	return model.GetTextTranslationReqToEnum().ZH_TW, nil
	case trans.LanguageTypeEn:
		return model.GetTextTranslationReqToEnum().EN, nil
	case trans.LanguageTypeDe:
		return model.GetTextTranslationReqToEnum().DE, nil
	case trans.LanguageTypeFr:
		return model.GetTextTranslationReqToEnum().FR, nil
	case trans.LanguageTypePt:
		return model.GetTextTranslationReqToEnum().PT, nil
	case trans.LanguageTypeEs:
		return model.GetTextTranslationReqToEnum().ES, nil
	case trans.LanguageTypeTh:
		return model.GetTextTranslationReqToEnum().TH, nil
	default:
		return model.TextTranslationReqTo{}, trans.ErrLanguageTypeNotSupported
	}
}
