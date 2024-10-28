package main

import (
	"flag"
	"github.com/lgrisa/lib/machine_translation/machine_translation_engine/volc"
	"github.com/lgrisa/lib/machine_translation/machine_translation_manager"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
)

var transPath = flag.String("transPath", "D:/WarGameProject/Excel/Localization", "翻译配置表路径") //"../../Excel/Localization"

var configPath = flag.String("configPath", "D:/WarGameProject/tools/excel-tool/LocalizationTools工具翻译配置表.xlsx", "翻译配置表路径") //"../excel-tool/LocalizationTools工具翻译配置表.xlsx"

func main() {

	flag.Parse()

	utils.LogInfoF("翻译配置表路径:（%v） 翻译配置表路径:（%v）", *transPath, *configPath)

	utils.InitLog()

	if *transPath == "" || *configPath == "" {
		utils.LogErrorF("transPath or configPath is empty")
		return
	}

	m, err := machine_translation_manager.NewClient(*configPath, *transPath)

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
