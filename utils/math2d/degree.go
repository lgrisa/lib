package math2d

import "math"

func Point2Degree(fromX, fromY, toX, toY float64) int64 {
	dx := toX - fromX
	dy := toY - fromY
	return Atan2Degree(dx, dy)
}

func Point2Radians(fromX, fromY, toX, toY float64) float64 {
	dx := toX - fromX
	dy := toY - fromY
	return math.Atan2(dy, dx)
}


// 平面直角坐标系，计算角度
func Atan2Degree(x, y float64) int64 {
	angrad := math.Atan2(y, x)
	return Degree360(AngradToDegree(angrad))
}

// 弧度转角度
func AngradToDegree(angrad float64) int64 {
	return int64(angrad * 180.0 / math.Pi)
}

// 角度转弧度
func DegreeToAngrad(degrees int64) float64 {
	return float64(degrees) * math.Pi / 180
}

// 变成0-359
func Degree360(degrees int64) int64 {
	return (degrees%360 + 360) % 360
}

// 获取2个角度的差值（返回锐角角度），0-180
func GetAngleBetween(d1, d2 int64) int64 {
	d1 = Degree360(d1)
	d2 = Degree360(d2)
	if d1 == d2 {
		return 0
	}

	var diff int64
	if d1 > d2 {
		diff = d1 - d2
	} else {
		diff = d2 - d1
	}

	if diff > 180 {
		// 转成锐角
		return 360 - diff
	}

	return diff
}
