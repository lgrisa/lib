package main

import (
	"github.com/lgrisa/library/config"
	"github.com/lgrisa/library/utils"
	"time"
)

func main() {
	//utils.Export(map[int]string{2024: "https://www.gov.cn/yaowen/liebiao/202310/content_6911540.htm"}, config.StartConfig.ConfigPath.HolidayPath)

	if err := config.ReadStartUpConfig(); err != nil {
		utils.LogPrintf("读取配置失败: %v", err)
		return
	}

	holidayData, err := config.LoadHolidaysJson()

	if err != nil {
		utils.LogPrintf("LoadHolidaysJson error: %s\n", err)
		return
	}

	if holidayData != nil {
		ctime := time.Now()
		utils.LogPrintf("Today %v is holiday:%v ", ctime.String(), holidayData.IsHoliday(ctime))
	}

	for i, v := range *holidayData {
		utils.LogPrintf("%v:%v\n", i, v)
	}
}
