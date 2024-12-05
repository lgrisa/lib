package main

import "github.com/lgrisa/lib/utils/holiday"

func main() {
	holidayPath := "config/holiday/" //导出节假日数据路径
	holiday.Export(map[int]string{2024: "https://www.gov.cn/yaowen/liebiao/202310/content_6911540.htm"}, holidayPath)
}
