package machine_trans_manager

import (
	"github.com/lgrisa/lib/machine_trans/machine_trans_engine"
)

type Manager struct {
	OriginalExcelMap map[string]map[string]*OriginalExcelInfo

	TransConfigMap map[machine_trans_engine.LanguageType]*TranslatedLanguageInfo

	TransClient []TranslatedClient

	excelPath string

	configPath string
}

type EngineType string

const (
	VolcEngine EngineType = "volcEngine" // 火山引擎
	Tencent    EngineType = "tencent"    // 腾讯
	HuaWei     EngineType = "huawei"     // 华为
)

type OriginalExcelInfo struct {
	keyID      string
	valueCn    string
	valueLocal string
}

type TranslatedLanguageInfo struct {
	languageType machine_trans_engine.LanguageType
	path         string
}

type TranslatedClient interface {
	Translate(text string, fromLanguage, toLanguage machine_trans_engine.LanguageType) (string, error)
	TranslateFor(text string, fromLanguage, toLanguage machine_trans_engine.LanguageType) (string, error)
}
