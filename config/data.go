package config

import (
	"encoding/json"
	"github.com/lgrisa/lib/utils/holiday"
	"strings"
	"time"
)

type HolidayData map[time.Time]bool

func LoadHolidaysJson() (*HolidayData, error) {

	holidayMap := &HolidayData{}

	if fss, err := HolidayFs.ReadDir("holiday"); err != nil {
		return nil, err
	} else {
		for _, fs := range fss {
			if fs.IsDir() {
				continue
			}

			if strings.Contains(fs.Name(), "holiday") {
				if data, err := HolidayFs.ReadFile("holiday/" + fs.Name()); err != nil {
					return nil, err
				} else {
					var dates []holiday.HolidayDate

					if err := json.Unmarshal(data, &dates); err != nil {
						return nil, err
					}

					for _, date := range dates {
						(*holidayMap)[date.Date] = date.IsOffDay
					}
				}
			}
		}
	}

	return holidayMap, nil
}

func (h *HolidayData) IsHoliday(date time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	if v, ok := (*h)[date]; ok {
		return v
	}

	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday || date.Weekday() == time.Friday {
		return true
	}

	return false
}
