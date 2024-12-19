package huawei

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/model"
	trans "github.com/lgrisa/lib/machine_trans/machine_trans_engine"
)

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
	case trans.LanguageTypeKr:
		return model.GetTextTranslationReqFromEnum().KO, nil
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
	case trans.LanguageTypeKr:
		return model.GetTextTranslationReqToEnum().KO, nil
	default:
		return model.TextTranslationReqTo{}, fmt.Errorf("toLanguage: %s is not supported", toLanguage)
	}
}
