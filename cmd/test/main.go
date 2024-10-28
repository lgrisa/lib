package main

import (
	"fmt"
	"github.com/lgrisa/lib/machine_translation/machine_translation_engine"
	"github.com/lgrisa/lib/machine_translation/machine_translation_engine/volc"
	"github.com/lgrisa/lib/utils"
)

func main() {
	fmt.Printf("%s@%s", utils.GetUsername(), utils.GetHostname())

	// Output: lgrisa@localhost

	volcEngineClient := volc.NewVolcEngineTransClient(
		"AKLTMzJiZGE1MzAwODY0NDg5ODhmZjAzODQ4YWY5ZmEzZTI",
		"WWpCbE1EZGlPRGt4WVRFNU5EQXhOMkpqTVRsak4yVmxNR1kzWkdZNFlUaw==")

	res, err := volcEngineClient.TranslateFor("前一天没有挑战记录时，6点重置后增加<color='#FBB800'>1次</color>重置次数，最多可累计<color=\"FBB800\">{0}次</color>。", machine_translation_engine.LanguageTypeZh, machine_translation_engine.LanguageTypeTh)

	if err != nil {
		utils.LogErrorF(err.Error())
		return
	}

	fmt.Println("返回结果:", res)
}
