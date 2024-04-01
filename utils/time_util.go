package utils

import (
	"time"
)

const offsetHour = 8 // 东八区

var (
	SecondsLayout = "2006-01-02 15:04:05"
	East8         = time.FixedZone("East-8", offsetHour*60*60)
	GameZone      = East8
)

func IsSameMonth(t1, t2 time.Time) bool {
	return MonthDay(t1) == MonthDay(t2) && YearDay(t1) == YearDay(t2)
}

func MonthDay(t time.Time) time.Month {
	return t.In(GameZone).Month()
}

func YearDay(t time.Time) int {
	return t.In(GameZone).Year()
}

func ParseSecondsLayout(value string) (time.Time, error) {
	return time.ParseInLocation(SecondsLayout, value, GameZone)
}
