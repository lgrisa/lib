package volc

import trans "github.com/lgrisa/lib/machine_trans/machine_trans_engine"

func (e *Engine) languageCode(code trans.LanguageType) string {

	switch code {
	case "auto":
		return ""
	case "zh-tw":
		return "zh-Hant"
	case trans.LanguageTypeZhTW:
		return "zh-Hant"
	case trans.LanguageTypeKr:
		return "ko"
	}

	return string(code)
}
