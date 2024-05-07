package utils

import (
	"github.com/sirupsen/logrus"
	"math"
	"runtime/debug"
	"time"
)

func NewFrameRate(rate int64) *FrameRate {
	durationPerFrame := time.Second / time.Duration(rate)

	return &FrameRate{
		rate:             rate,
		rateF64:          float64(rate),
		durationPerFrame: durationPerFrame,
	}
}

type FrameRate struct {
	// 帧率，每秒多少帧
	rate             int64
	rateF64          float64
	durationPerFrame time.Duration
}

func (fr *FrameRate) GetRate() int64 {
	return fr.rate
}

func (fr *FrameRate) GetRateF64() float64 {
	return fr.rateF64
}

func (fr *FrameRate) GetDurationPerFrame() time.Duration {
	return fr.durationPerFrame
}

func (fr *FrameRate) TimeToFrameAndAccumTime(time float64) (int64, float64) {
	frame := fr.TimeToFrame(time)

	// 计算实际使用的时间 = frame * rate
	t := float64(frame) / fr.rateF64
	return frame, t - time
}

func (fr *FrameRate) TimeToFrame(time float64) int64 {
	if time > 0 {
		return int64(math.Max(1, math.Ceil(time*fr.rateF64))) // 至少需要1帧
	}
	return 0
}
func (fr *FrameRate) DurationToFrame(duration time.Duration) int64 {
	if duration > 0 {
		return fr.TimeToFrame(duration.Seconds())
	}
	return 0
}

//func (fr *FrameRate) FrameToSecond(frame int64) float64 {
//	return float64(frame) / fr.rateF64
//}
//
//func (fr *FrameRate) MillisToFrame(millis int64) int64 {
//	time := float64(millis) / 1000
//	return fr.TimeToFrame(time)
//}
//
//func (fr *FrameRate) FrameToMillis(frame int64) int64 {
//	return frame * 1000 / fr.rate
//}

func (fr *FrameRate) MoveFrame(startX, startY, endX, endY, speedPerSecond float64) int64 {
	if speedPerSecond <= 0 {
		logrus.WithField("stack", string(debug.Stack())).Errorf("计算需要的frame时, 速度<=0: %f", speedPerSecond)
		return 1
	}
	diffX := startX - endX
	diffY := startY - endY
	distance := math.Sqrt(diffX*diffX + diffY*diffY)

	needTime := distance / speedPerSecond

	return fr.TimeToFrame(needTime)
}
