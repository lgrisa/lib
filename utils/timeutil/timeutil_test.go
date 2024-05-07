package timeutil

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestUnixTime(t *testing.T) {
	RegisterTestingT(t)

	var second int32 = 100
	tt := Unix32(second)
	Ω(tt.Unix()).Should(Equal(int64(second)))

	second = int32(time.Now().Unix())
	tt = Unix32(second)
	Ω(tt.Unix()).Should(Equal(int64(second)))
}

func TestZero(t *testing.T) {
	RegisterTestingT(t)

	tt := time.Unix(0, 0)
	Ω(tt.IsZero()).Should(BeFalse())

	tt2 := time.Time{}
	Ω(tt2.IsZero()).Should(BeTrue())

	tt3 := time.Unix(tt2.Unix(), int64(tt2.Nanosecond()))
	Ω(tt3.IsZero()).Should(BeTrue())

	Ω(IsZero(time.Time{})).Should(BeTrue())
	Ω(IsZero(Unix64(0))).Should(BeTrue())
	Ω(IsZero(Unix64(1))).Should(BeFalse())
	Ω(IsZero(Unix64(-1))).Should(BeFalse())
}

func TestMarshal(t *testing.T) {
	RegisterTestingT(t)

	zero := time.Time{}
	unixZero := time.Unix(0, 0)

	Ω(Marshal64(zero)).Should(Equal(int64(0)))
	Ω(Marshal32(zero)).Should(Equal(int32(0)))

	Ω(Marshal64(unixZero)).Should(Equal(int64(0)))
	Ω(Marshal32(unixZero)).Should(Equal(int32(0)))

	Ω(Marshal64(zero.Add(time.Second))).Should(Equal(int64(0)))
	Ω(Marshal32(zero.Add(time.Second))).Should(Equal(int32(0)))

}

func TestDuration(t *testing.T) {
	RegisterTestingT(t)

	d32 := Duration32(10)
	Ω(d32).Should(Equal(10 * time.Second))

	d64 := Duration64(10)
	Ω(d32).Should(Equal(d64))

	i32 := Duration2SecondI32(d32)
	Ω(i32).Should(Equal(int32(10)))

	i64 := Duration2Second(d64)
	Ω(i64).Should(Equal(int64(10)))
}

func TestNextTickTime(t *testing.T) {
	RegisterTestingT(t)

	now := time.Now()
	prev := now.Add(-30 * time.Second)

	d := time.Minute

	Ω(NextTickTime(prev, now, d)).Should(Equal(prev.Add(time.Minute)))
	Ω(NextTickTime(now, now, d)).Should(Equal(now.Add(time.Minute)))
	Ω(NextTickTime(now.Add(time.Second), now, d)).Should(Equal(now.Add(time.Second)))
	Ω(NextTickTime(prev, now.Add(time.Hour), d)).Should(Equal(prev.Add(time.Hour).Add(time.Minute)))
}

func TestCompletionMMDD(t *testing.T) {
	RegisterTestingT(t)

	Ω(CompletionMMDD("2015-1-2", "-")).Should(Equal("2015-01-02"))
	Ω(CompletionMMDD("2015/1/2", "/")).Should(Equal("2015/01/02"))

	Ω(CompletionMMDD("2015-11-2", "-")).Should(Equal("2015-11-02"))
	Ω(CompletionMMDD("2015/1/22", "/")).Should(Equal("2015/01/22"))

	Ω(CompletionMMDD("2015-11-22", "-")).Should(Equal("2015-11-22"))
	Ω(CompletionMMDD("2015/11/22", "/")).Should(Equal("2015/11/22"))

	Ω(CompletionMMDD("2015-1-2", "/")).Should(Equal("2015-1-2"))
	Ω(CompletionMMDD("2015/1/2", "-")).Should(Equal("2015/1/2"))
}

func TestGetYearMonthHaveDay(t *testing.T) {
	RegisterTestingT(t)

	Ω(GetYearMonthHaveDay(2022, 1)).Should(Equal(int(31)))
	Ω(GetYearMonthHaveDay(2022, 2)).Should(Equal(int(28)))
	Ω(GetYearMonthHaveDay(2022, 4)).Should(Equal(int(30)))
	Ω(GetYearMonthHaveDay(2022, 12)).Should(Equal(int(31)))

	Ω(GetYearMonthHaveDay(2020, 2)).Should(Equal(int(29)))
}

func TestSubDays(t *testing.T) {
	RegisterTestingT(t)

	t1, _ := ParseSecondsLayout("2022-01-02_15:04:05")
	t2, _ := ParseSecondsLayout("2022-01-01_15:04:05")
	t3, _ := ParseSecondsLayout("2022-01-03_00:04:05")

	Ω(SubDays(t1, t2)).Should(Equal(1))
	Ω(SubDays(t1, t3)).Should(Equal(1))
	Ω(SubDays(t2, t3)).Should(Equal(2))
}
