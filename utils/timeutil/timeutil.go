package timeutil

import (
	"github.com/lgrisa/lib/utils/math/i64"
	"github.com/lgrisa/lib/utils/math/imath"
	"math"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

func MarshalArray64(array []time.Time) []int64 {
	out := make([]int64, len(array))
	for i, v := range array {
		out[i] = Marshal64(v)
	}

	return out
}

func MarshalArray32(array []time.Time) []int32 {
	out := make([]int32, len(array))
	for i, v := range array {
		out[i] = Marshal32(v)
	}

	return out
}

func CopyUnix32Array(array []time.Time, intArray []int32) {
	n := imath.Min(len(array), len(intArray))
	if n > 0 {
		for i := 0; i < n; i++ {
			array[i] = Unix32(intArray[i])
		}
	}
}

func CopyUnix64Array(dest []time.Time, src []int64) {
	n := imath.Min(len(dest), len(src))
	if n > 0 {
		for i := 0; i < n; i++ {
			dest[i] = Unix64(src[i])
		}
	}
}

func Marshal64(t time.Time) int64 {
	if t.IsZero() {
		// 为了防止当前时间就是 unixZeroTime, 这里不用 IsZero(t)
		return 0
	}

	if unix := t.Unix(); unix < math.MinInt32 {
		logrus.WithField("stack", string(debug.Stack())).WithField("unix", unix).Error("timeutil.Marshal t.unix < math.MinInt32")
		return 0
	} else {
		return unix
	}
}

func Marshal32(t time.Time) int32 {
	return int32(Marshal64(t))
}

func Unix32(second int32) time.Time {
	return Unix64(int64(second))
}

func Unix64(second int64) time.Time {
	return time.Unix(second, 0)
}

func IsSameDay(t1, t2 time.Time) bool {
	ns := i64.Abs(t1.Sub(t2).Nanoseconds())
	if ns > Day.Nanoseconds() {
		return false
	}

	return DailyTime.PrevTime(t1).Equal(DailyTime.PrevTime(t2))
}

func IsSameMonth(t1, t2 time.Time) bool {
	return MonthDay(t1) == MonthDay(t2) && YearDay(t1) == YearDay(t2)
}

func DivideTimes(x, y time.Duration) uint64 {
	if x <= 0 || y <= 0 {
		return 0
	}

	return uint64((x + y - 1) / y)
}

// 将duration数组转成seconds数组
func DurationArrayToSecondArray(array []time.Duration) (result []int32) {
	result = make([]int32, 0, len(array))

	for _, duration := range array {
		result = append(result, int32(duration/time.Second))
	}

	return result
}

func Duration32(seconds int32) time.Duration {
	return time.Duration(seconds) * time.Second
}

func Duration64(seconds int64) time.Duration {
	return time.Duration(seconds) * time.Second
}

func Duration2Second(d time.Duration) int64 {
	return int64(d / time.Second)
}

func Duration2SecondI32(d time.Duration) int32 {
	return int32(d / time.Second)
}

var unixZeroTime = Unix64(0)

func IsZero(time time.Time) bool {
	return time.IsZero() || time.Equal(unixZeroTime)
}

func Min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func Max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func MinDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func MaxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func MultiDuration(multi float64, d time.Duration) time.Duration {
	return time.Duration(multi * float64(d))
}

func NextTickTime(prevTime, ctime time.Time, d time.Duration) time.Time {
	if ctime.Before(prevTime) {
		return prevTime
	}
	return prevTime.Add((ctime.Sub(prevTime)/d)*d + d)
}

func Between(t, start, end time.Time) bool {
	return !t.Before(start) && t.Before(end)
}

func BetweenClosed(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}

func Rate(startTime, endTime, ctime time.Time) float64 {
	return i64.Rate(startTime.Unix(), endTime.Unix(), ctime.Unix())
}

func Midnight(t time.Time) time.Time {
	return DailyTime.PrevTime(t)
}

// 获取某年某个月有多少天
func GetYearMonthHaveDay(inYear, inMonth int) int {
	// 2000-02-00 = 2000-01-31
	return time.Date(inYear, time.Month(inMonth+1), 0, 0, 0, 0, 0, time.Local).Day()
}

// 计算日期相差多少天
func SubDays(t1, t2 time.Time) (day int) {
	if t1.Before(t2) {
		t_ := t1
		t1 = t2
		t2 = t_
	}

	day = int(Midnight(t1).Sub(Midnight(t2)) / Day)

	return
}
