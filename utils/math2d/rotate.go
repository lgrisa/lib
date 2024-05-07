package math2d

import "math"

func Rotate(ox, oy, tx, ty, angrad float64) (float64, float64) {
	ax, ay := tx-ox, ty-oy
	bx, by := RotateOrigin(ax, ay, angrad)
	return bx + ox, by + oy
}

// Angrad 为弧度
func RotateOrigin(ax, ay, angrad float64) (bx, by float64) {
	//假设o点为圆心(原点)，则有计算公式：
	//b.x = a.x*cos(angle)  - a.y*sin(angle)
	//b.y = a.x*sin(angle) + a.y*cos(angle)

	sin, cos := math.Sincos(angrad)
	bx = ax*cos - ay*sin
	by = ax*sin + ay*cos
	return
}
