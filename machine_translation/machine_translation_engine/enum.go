package machine_translation_engine

import "github.com/pkg/errors"

type LanguageType string

const (
	LanguageTypeAuto LanguageType = "auto"
	LanguageTypeZh   LanguageType = "zh"
	LanguageTypeZhTW LanguageType = "zh-TW"
	LanguageTypeEn   LanguageType = "en"
	LanguageTypeDe   LanguageType = "de"
	LanguageTypeFr   LanguageType = "fr"
	LanguageTypePt   LanguageType = "pt"
	LanguageTypeEs   LanguageType = "es"
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
	default:
		return "", errors.Errorf("languageName: %s is not supported", languageName)
	}
}
