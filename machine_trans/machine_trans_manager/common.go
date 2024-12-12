package machine_trans_manager

import (
	"github.com/lgrisa/lib/machine_trans/machine_trans_engine"
	"github.com/pkg/errors"
)

func NewClient(configPath string, excelPath string) (*Manager, error) {

	m := &Manager{
		TransConfigMap:   make(map[machine_trans_engine.LanguageType]*TranslatedLanguageInfo),
		OriginalExcelMap: make(map[string]map[string]*OriginalExcelInfo),
		TransClient:      make(map[EngineType]TranslatedClient),
		excelPath:        excelPath,
		configPath:       configPath,
	}

	return m, nil
}

func (m *Manager) RegisterClient(engineType EngineType, client TranslatedClient) error {

	if client == nil {
		return errors.Errorf("RegisterClient:%v client is nil", engineType)
	}

	m.TransClient[engineType] = client

	return nil
}

func (m *Manager) Start() error {

	if errInitConfig := m.initConfig(); errInitConfig != nil {
		return errors.Errorf("m.initConfig => %v", errInitConfig)
	}

	if errLoadOriginExcel := m.LoadOriginExcel(); errLoadOriginExcel != nil {
		return errors.Errorf("m.LoadOriginExcel => %v", errLoadOriginExcel)
	}

	if errTranslate := m.translate(); errTranslate != nil {
		return errors.Errorf("m.translate => %v", errTranslate)
	}

	return nil
}
