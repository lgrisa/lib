package holiday

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{}

type HolidayDate struct {
	Year     int
	Name     string
	Date     time.Time
	IsOffDay bool
}

type Holiday struct {
	Year         int
	HolidayDates []HolidayDate
	Container    []soup.Root
}

func (s *Holiday) ParseRules() {
	for _, p := range s.Container {
		if p.Text() == "" {
			continue
		}

		//判断是否为大写数字开头的序号，大写数字开头的序号为具体放假安排
		mRegex := regexp.MustCompile(`[一二三四五六七八九十]、(.+?)：(.+)`)
		match := mRegex.FindStringSubmatch(p.FullText())
		if len(match) <= 2 {
			continue
		}

		//分段处理，降低匹配复杂度
		for _, str := range regexp.MustCompile("[，。；]").Split(match[2], -1) {
			if str == "" {
				continue
			}

			//获取休息日
			rest := regexp.MustCompile(`(.+)(放假|补休|调休|公休)+(?:\d+天)?$`).FindStringSubmatch(str)
			if len(rest) > 2 {
				//解析具体日期
				s.ExtractDates(match[1], rest[1], true)
				continue
			}

			//获取工作日
			work := regexp.MustCompile(`(.+)上班$`).FindStringSubmatch(str)
			if len(work) > 1 {
				// 解析具体日期
				s.ExtractDates(match[1], work[1], false)
				continue
			}
		}

	}
}

func (s *Holiday) ExtractDates(name, txt string, offDay bool) {
	txt = strings.ReplaceAll(txt, "(", "（")
	txt = strings.ReplaceAll(txt, ")", "）")

	//[xxxx年][x月]x日
	matches := regexp.MustCompile(`(?:(\d+)年)?(?:(\d+)月)?(\d+)日`).FindAllStringSubmatch(txt, -1)
	for _, match := range matches {
		if match[2] == "" {
			continue
		}

		if match[1] == "" {
			match[1] = strconv.Itoa(s.Year)
		}

		if s.IsExist(fmt.Sprintf("%s-%s-%s", match[1], match[2], match[3])) {
			continue
		}
		s.HolidayDates = append(s.HolidayDates, HolidayDate{
			s.Year,
			name,
			s.GetDate(match[1], match[2], match[3]),
			offDay,
		})
	}

	//[xxxx年]x月x日至[xxxx年][x月]x日
	ext2txt := regexp.MustCompile(`（.+?）`).ReplaceAllString(txt, "")
	matches = regexp.MustCompile(`(?:(\d+)年)?(?:(\d+)月)?(\d+)日(?:至|-|—)(?:(\d+)年)?(?:(\d+)月)?(\d+)日`).
		FindAllStringSubmatch(ext2txt, -1)
	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		if match[1] == "" {
			match[1] = strconv.Itoa(s.Year)
		}
		if match[4] == "" {
			match[4] = strconv.Itoa(s.Year)
		}

		if match[5] == "" {
			match[5] = match[2]
		}

		start := s.GetDate(match[1], match[2], match[3])
		end := s.GetDate(match[4], match[5], match[6])
		//解析日期范围
		for i := 0; i <= int(end.Sub(start).Hours()/24); i++ {
			d := s.GetDate(match[1], match[2], match[3]).AddDate(0, 0, i)

			if s.IsExist(d.Format("2006-1-2")) {
				continue
			}
			s.HolidayDates = append(s.HolidayDates, HolidayDate{
				s.Year,
				name,
				d,
				offDay,
			})
		}
	}

	//x月x日(星期x)、x月x日(星期x)
	ext3txt := regexp.MustCompile(`（.+?）`).ReplaceAllString(txt, "")
	matches = regexp.MustCompile(
		`(?:(\d+)年)?(?:(\d+)月)?(\d+)日(?:（[^）]+）)?(?:、(?:(\d+)年)?(?:(\d+)月)?(\d+)日(?:（[^）]+）)?)+`,
	).FindAllStringSubmatch(ext3txt, -1)
	for _, match := range matches {

		if len(match) < 6 {
			continue
		}

		if match[1] == "" {
			match[1] = strconv.Itoa(s.Year)
		}
		if match[4] == "" {
			match[4] = strconv.Itoa(s.Year)
		}

		if match[5] == "" {
			match[5] = match[2]
		}
		d := s.GetDate(match[1], match[2], match[3])
		if !s.IsExist(d.Format("2006-1-2")) {
			s.HolidayDates = append(s.HolidayDates, HolidayDate{
				s.Year,
				name,
				d,
				offDay,
			})
		}

		d = s.GetDate(match[4], match[5], match[6])
		if !s.IsExist(d.Format("2006-1-2")) {
			s.HolidayDates = append(s.HolidayDates, HolidayDate{
				s.Year,
				name,
				d,
				offDay,
			})
		}
	}

}

// IsExist 判断是否已存在
func (s *Holiday) IsExist(date string) bool {
	for _, i := range s.HolidayDates {
		if i.Date.Format("2006-1-2") == date {
			return true
		}
	}
	return false
}

// GetDate 日期字符串转成日期对象
func (s *Holiday) GetDate(y, m, d string) time.Time {
	t, _ := time.ParseInLocation("2006-1-2", fmt.Sprintf("%s-%s-%s", y, m, d), time.Local)
	return t
}

func (s *Holiday) FetchPage(url string) error {
	r, err := HTTPGet(url, map[string]interface{}{})
	if err != nil {
		return err
	}
	//定位到 id = UCAP-CONTENT 的 div 容器，读取所有的 p 标签条目
	s.Container = soup.HTMLParse(r).Find("div", "id", "UCAP-CONTENT").FindAll("p")

	if len(s.Container) == 0 {
		return fmt.Errorf("Page parse error ")
	}
	return nil
}

func InitHolidayParse(year int) *Holiday {
	return &Holiday{
		Year: year,
	}
}

func Export(export map[int]string, exportPath string) {

	//urls, err := holiday.SearchPageUrls()
	//if err != nil {
	//	log.Fatalln(fmt.Sprintf("查询 %d 放假通知异常", year))
	//}

	for year, url := range export {

		holiday := InitHolidayParse(year)
		// 请求具体通知页面，并分析放假安排
		log.Printf("[ %d ] ====> %s", year, url)

		if err := holiday.FetchPage(url); err != nil {
			fmt.Printf("获取并分析 %d 放假通知页面，异常\n%s\n", year, url)
			continue
		}
		holiday.ParseRules()

		marshal, err := json.Marshal(holiday.HolidayDates)
		if err != nil {
			fmt.Printf("序列化 %d 放假数据，异常\n%s\n", year, err)
			continue
		}

		os.Mkdir(exportPath, 0755)

		err = os.WriteFile(fmt.Sprintf(exportPath+"/holiday-%d.json", year), marshal, 0644)
		if err != nil {
			fmt.Printf("写入文件异常\n%s\n", err)
			continue
		}
	}
}
