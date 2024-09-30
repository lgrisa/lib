package main

import (
	"github.com/lgrisa/lib/machine_translation/machine_translation_engine/volc"
	"github.com/lgrisa/lib/machine_translation/machine_translation_manager"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
)

func main() {

	excelPath := "D:/WarGameProject/Excel/Localization"
	//excelPath := "../../Excel/Localization"

	configPath := "../excel-tool/LocalizationTools工具翻译配置表.xlsx"

	m, err := machine_translation_manager.NewClient(configPath, excelPath)

	if err != nil {
		utils.LogErrorF(errors.WithStack(err).Error())
		return
	}

	volcEngineClient := volc.NewVolcEngineTransClient(
		"AKLTMzJiZGE1MzAwODY0NDg5ODhmZjAzODQ4YWY5ZmEzZTI",
		"WWpCbE1EZGlPRGt4WVRFNU5EQXhOMkpqTVRsak4yVmxNR1kzWkdZNFlUaw==")

	if errRegister := m.RegisterClient(machine_translation_manager.VolcEngine, volcEngineClient); errRegister != nil {
		utils.LogErrorF(errors.WithStack(errRegister).Error())
		return
	}

	err = m.Start()

	if err != nil {
		utils.LogErrorF(errors.WithStack(err).Error())
		return
	}
}
