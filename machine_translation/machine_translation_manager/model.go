package machine_translation_manager

import (
	"github.com/lgrisa/lib/machine_translation/machine_translation_engine"
)

type Manager struct {
	OriginalExcelMap map[string]map[string]*OriginalExcelInfo

	TransConfigMap map[machine_translation_engine.LanguageType]*TranslatedLanguageInfo

	TransClient map[EngineType]TranslatedClient

	excelPath string

	configPath string
}

type EngineType string

const (
	VolcEngine EngineType = "volcEngine" // 火山引擎
	Tencent    EngineType = "tencent"    // 腾讯
)

type OriginalExcelInfo struct {
	keyID      string
	valueCn    string
	valueLocal string
}

type TranslatedLanguageInfo struct {
	languageType machine_translation_engine.LanguageType
	path         string
}

type TranslatedClient interface {
	Translate(text string, fromLanguage, toLanguage machine_translation_engine.LanguageType) (string, error)
	TranslateFor(text string, fromLanguage, toLanguage machine_translation_engine.LanguageType) (string, error)
}
