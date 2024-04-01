package main

import (
	"fmt"
	"github.com/packer/config"
	"time"
)

func main() {
	//utils.Export(map[int]string{2024: "https://www.gov.cn/yaowen/liebiao/202310/content_6911540.htm"}, config.StartConfig.ConfigPath.HolidayPath)

	if err := config.ReadStartUpConfig(); err != nil {
		fmt.Printf("读取配置失败: %v", err)
		return
	}

	holidayData, err := config.LoadHolidaysJson()

	if err != nil {
		fmt.Printf("LoadHolidaysJson error: %s\n", err)
		return
	}

	if holidayData != nil {
		ctime := time.Now()
		fmt.Printf("Today %v is holiday:%v ", ctime.String(), holidayData.IsHoliday(ctime))
	}

	for i, v := range *holidayData {
		fmt.Printf("%v:%v\n", i, v)
	}
}
