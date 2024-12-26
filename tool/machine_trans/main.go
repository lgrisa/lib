package main

import (
	"flag"
	"github.com/lgrisa/lib/machine_trans/machine_trans_engine/huawei"
	"github.com/lgrisa/lib/machine_trans/machine_trans_engine/volc"
	"github.com/lgrisa/lib/machine_trans/machine_trans_manager"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
	"time"
)

var transPath = flag.String("transPath", "D:/WarGameProject/Excel/Localization", "翻译配置表路径") //"../../Excel/Localization"

var configPath = flag.String("configPath", "D:/WarGameProject/tools/excel-tool/LocalizationTools工具翻译配置表.xlsx", "翻译配置表路径") //"../excel-tool/LocalizationTools工具翻译配置表.xlsx"

var logLevel = flag.Int("logLevel", 0, "日志等级")

func main() {

	startTime := time.Now()

	flag.Parse()

	utils.LogInfoF("翻译配置表路径:（%v） 翻译配置表路径:（%v）", *transPath, *configPath)

	utils.InitLog(*logLevel)

	if *transPath == "" || *configPath == "" {
		utils.LogErrorF("transPath or configPath is empty")
		return
	}

	m, err := machine_trans_manager.NewClient(*configPath, *transPath)

	if err != nil {
		utils.LogErrorF(errors.WithStack(err).Error())
		return
	}

	huaWeiClient, err := huawei.NewClient(
		"YCEXVLOW4BEXJQ1Z3QQU",
		"gMrClgUnkFKIDFeaXA7sqghawCfRt6ILJTKTQLBe")

	if err != nil {
		utils.LogErrorF(errors.WithStack(err).Error())
		return
	}

	if errRegister := m.RegisterClient(machine_trans_manager.HuaWei, huaWeiClient); errRegister != nil {
		utils.LogErrorF(errors.WithStack(errRegister).Error())
		return
	}

	volcEngineClient := volc.NewClient(
		"AKLTMzJiZGE1MzAwODY0NDg5ODhmZjAzODQ4YWY5ZmEzZTI",
		"WWpCbE1EZGlPRGt4WVRFNU5EQXhOMkpqTVRsak4yVmxNR1kzWkdZNFlUaw==")

	if errRegister := m.RegisterClient(machine_trans_manager.VolcEngine, volcEngineClient); errRegister != nil {
		utils.LogErrorF(errors.WithStack(errRegister).Error())
		return
	}

	err = m.Start()

	if err != nil {
		utils.LogErrorF(errors.WithStack(err).Error())
		return
	}

	utils.LogInfoF("翻译配置表结束, 耗时:%v", time.Since(startTime))
}
