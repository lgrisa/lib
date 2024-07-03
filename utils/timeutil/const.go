package timeutil

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	SecondsPerHour = 60 * 60
	SecondsPerDay  = 24 * SecondsPerHour
	offsetHour     = 8 // 东八区

)

var (
	East8 = time.FixedZone("East-8", offsetHour*60*60)
	//East8Offset = 8 * time.Hour // (duration east of UTC).

	GameZone = East8
	//GameZoneOffset = East8Offset

	StartTime = time.Date(2000, 1, 1, 0, 0, 0, 0, GameZone)

	DailyTime = NewDailyCycleTime(StartTime.Unix())

	Sunday    = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Sunday).Unix())
	Monday    = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Monday).Unix())
	Tuesday   = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Tuesday).Unix())
	Wednesday = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Wednesday).Unix())
	Thursday  = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Thursday).Unix())
	Friday    = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Friday).Unix())
	Saturday  = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Saturday).Unix())

	weeks = [...]*CycleTime{
		Sunday,
		Monday,
		Tuesday,
		Wednesday,
		Thursday,
		Friday,
		Saturday,
	}

	DayLayout         = "2006-01-02"
	HourLayout        = "2006-01-02-15"
	MinuteLayout      = "2006-01-02-15-04"
	SecondsLayout     = "2006-01-02_15:04:05"
	DaySlashLayout    = "2006/01/02"
	DaySlashLayoutLen = len(DaySlashLayout)
	TimeNumLayout     = "150405"

	Day  = 24 * time.Hour
	Week = 7 * Day
)

func GetDayDuration(day int64) time.Duration {
	return time.Duration(day) * Day
}

// 每日0点重置，返回下一天的0点（传入今日0点，也返回下一天0点）

func GetNextResetDailyTime(ctime time.Time, offset time.Duration) time.Time {
	return GetNextResetTime(ctime, DailyTime, Day, offset)
}

// 每周重置，周日算新的一周

func GetNextResetWeeklyTime(ctime time.Time, offset time.Duration) time.Time {
	return GetNextResetTime(ctime, Sunday, Week, offset)
}

// 每周重置，周一算新的一周

func GetNextResetWeeklyTimeIsMonday(ctime time.Time, offset time.Duration) time.Time {
	return GetNextResetTime(ctime, Monday, Week, offset)
}

// 每月重置，返回下一月1日0点

func GetNextResetMonthlyTime(ctime time.Time, offset time.Duration) time.Time {
	return GetNextResetMonthlyDayTime(ctime, 1, offset)
}

// 每月重置，返回下一月指定日期的0点

func GetNextResetMonthlyDayTime(ctime time.Time, monthDay int, offset time.Duration) time.Time {
	// 每个月的天数是不是不固定的，所以循环时间需要每次计算
	year, month, _ := ctime.Date()

	// 本月判断
	currMonthDay := GetYearMonthHaveDay(year, int(month))
	if currMonthDay > monthDay {
		currMonthDay = monthDay
	}
	currMonthData := time.Date(year, month, currMonthDay, 0, 0, 0, 0, GameZone)
	currMonthData = currMonthData.Add(offset)
	if ctime.Before(currMonthData) {
		return currMonthData
	}

	// 找下个月
	nextMonthDay := GetYearMonthHaveDay(year, int(month+1))
	if nextMonthDay > monthDay {
		nextMonthDay = monthDay
	}
	nextMonthUnix := time.Date(year, month+1, nextMonthDay, 0, 0, 0, 0, GameZone)
	return nextMonthUnix.Add(offset)
}

func GetNextResetTime(ctime time.Time, cycleTime *CycleTime, cycleDuration, offset time.Duration) time.Time {

	prevTime := cycleTime.PrevTime(ctime)
	resetTime := prevTime.Add(offset)

	if ctime.Before(resetTime) {
		return resetTime
	}

	return resetTime.Add(cycleDuration)
}

func WeekCycleTime(d time.Weekday) *CycleTime {
	return weeks[d]
}

func NextWeekTime(t time.Time, weekday time.Weekday) time.Time {

	wd := t.Weekday()
	diff := weekday - wd
	if diff < 0 {
		diff += 7
	}

	return t.Add(time.Duration(int64(diff)*SecondsPerDay) * time.Second)
}

func ParseMidnightTime(value string) (time.Time, error) {
	if t, err := ParseDayLayout(value); err != nil {
		if t, err := ParseHourLayout(value); err != nil {
			return time.Time{}, err
		} else {
			return DailyTime.PrevTime(t), nil
		}
	} else {
		return t, nil
	}
}

func ParseDayLayout(value string) (time.Time, error) {
	return time.ParseInLocation(DayLayout, value, GameZone)
}

func ParseDaySlashLayout(value string) (time.Time, error) {
	return time.ParseInLocation(DaySlashLayout, value, GameZone)
}

func ParseHourLayout(value string) (time.Time, error) {
	return time.ParseInLocation(HourLayout, value, GameZone)
}

func ParseMinuteLayout(value string) (time.Time, error) {
	return time.ParseInLocation(MinuteLayout, value, GameZone)
}

func ParseSecondsLayout(value string) (time.Time, error) {
	return time.ParseInLocation(SecondsLayout, value, GameZone)
}

// 自动补全日期

func CompletionMMDD(value, sep string) string {
	// 2018-1-2 -> 2018-01-02

	if len(value) < DaySlashLayoutLen {
		array := strings.SplitN(value, sep, 3)
		if len(array) == 3 {
			if len(array[1]) == 1 {
				array[1] = "0" + array[1]
			}
			if len(array[2]) == 1 {
				array[2] = "0" + array[2]
			}

			return array[0] + sep + array[1] + sep + array[2]
		}
	}
	return value
}

func ParseHMS(value string) (hour, minute, second int, err error) {
	if len(value) <= 0 {
		return
	}

	hhmmss := strings.Split(value, ":")
	n := len(hhmmss)
	if n <= 0 {
		return 0, 0, 0, errors.Errorf("parse hms fail, [%v]", value)
	}

	// 小时分钟秒
	if hour, err = strconv.Atoi(hhmmss[0]); err != nil {
		return 0, 0, 0, errors.Wrapf(err, "parse hms fail, %v", value)
	}

	if hour < 0 || hour >= 24 {
		return 0, 0, 0, errors.Errorf("parse hms fail, invalid hour, %v", value)
	}

	if n <= 1 {
		return
	}

	// 分钟
	if minute, err = strconv.Atoi(hhmmss[1]); err != nil {
		return 0, 0, 0, errors.Wrapf(err, "parse hms fail, %v", value)
	}

	if minute < 0 || minute >= 60 {
		return 0, 0, 0, errors.Errorf("parse hms fail, invalid minute, %v", value)
	}

	if n <= 2 {
		return
	}

	// 秒
	if second, err = strconv.Atoi(hhmmss[2]); err != nil {
		return 0, 0, 0, errors.Wrapf(err, "parse hms fail, %v", value)
	}

	if second < 0 || second >= 60 {
		return 0, 0, 0, errors.Errorf("parse hms fail, invalid second, %v", value)
	}

	return
}
