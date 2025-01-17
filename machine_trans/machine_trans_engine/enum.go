package machine_trans_engine

import "github.com/pkg/errors"

type LanguageType string

const (
	LanguageTypeAuto LanguageType = "auto"  //	自动检测
	LanguageTypeZh   LanguageType = "zh"    //	中文
	LanguageTypeZhTW LanguageType = "zh-TW" //	繁体中文
	LanguageTypeEn   LanguageType = "en"    //	英语
	LanguageTypeDe   LanguageType = "de"    //	德语
	LanguageTypeFr   LanguageType = "fr"    //	法语
	LanguageTypePt   LanguageType = "pt"    //	葡萄牙语
	LanguageTypeEs   LanguageType = "es"    // 西班牙语
	LanguageTypeTh   LanguageType = "th"    // 泰语
	LanguageTypeKr   LanguageType = "kr"    // 韩语
	LanguageTypeRu   LanguageType = "ru"    // 俄语
)

func GetLanguageType(languageName string) (LanguageType, error) {
	switch languageName {
	case "auto":
		return LanguageTypeAuto, nil
	case "zh":
		return LanguageTypeZh, nil
	case "zh-tw":
		return LanguageTypeZhTW, nil
	case "zh-TW":
		return LanguageTypeZhTW, nil
	case "en":
		return LanguageTypeEn, nil
	case "de":
		return LanguageTypeDe, nil
	case "fr":
		return LanguageTypeFr, nil
	case "pt":
		return LanguageTypePt, nil
	case "es":
		return LanguageTypeEs, nil
	case "th":
		return LanguageTypeTh, nil
	case "kr":
		return LanguageTypeKr, nil
	case "ru":
		return LanguageTypeRu, nil
	default:
		return "", errors.Errorf("languageName: %s is not supported", languageName)
	}
}
