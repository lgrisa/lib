package math2d

import "math"

func DistanceSquared(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}

func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(DistanceSquared(x1, y1, x2, y2))
}

func IsInRange(x1, y1, x2, y2, dist float64) bool {
	return DistanceSquared(x1, y1, x2, y2) <= dist*dist
}

// 点到线段的距离 https://stackoverflow.com/a/6853926
func PointToLineDistance(pointX, pointY, lineX1, lineY1, lineX2, lineY2 float64) float64 {
	A := pointX - lineX1
	B := pointY - lineY1
	C := lineX2 - lineX1
	D := lineY2 - lineY1

	dot := A*C + B*D
	len_sq := C*C + D*D
	var param float64 = -1
	if len_sq != 0 { //in case of 0 length line
		param = dot / len_sq
	}

	var xx, yy float64

	if param < 0 {
		xx = lineX1
		yy = lineY1
	} else if param > 1 {
		xx = lineX2
		yy = lineY2
	} else {
		xx = lineX1 + param*C
		yy = lineY1 + param*D
	}

	dx := pointX - xx
	dy := pointY - yy
	return math.Sqrt(dx*dx + dy*dy)
}

// 点p是否在矩形中
func IsInRectRange(px, py, x1, y1, x2, y2 float64) bool {
	minX := math.Min(x1, x2)
	maxX := math.Max(x1, x2)
	minY := math.Min(y1, y2)
	maxY := math.Max(y1, y2)

	return IsInRectMinMaxRange(px, py, minX, minY, maxX, maxY)
}

func IsInRectMinMaxRange(px, py, minX, minY, maxX, maxY float64) bool {
	if px < minX || px > maxX || py < minY || py > maxY {
		return false
	}
	return true
}

func IsInRectIntersect(rect1MinX, rect1MinY, rect1MaxX, rect1MaxY, rect2MinX, rect2MinY, rect2MaxX, rect2MaxY float64) bool {

	if IsInRectMinMaxRange(rect1MinX, rect1MinY, rect2MinX, rect2MinY, rect2MaxX, rect2MaxY) ||
		IsInRectMinMaxRange(rect1MinX, rect1MaxY, rect2MinX, rect2MinY, rect2MaxX, rect2MaxY) ||
		IsInRectMinMaxRange(rect1MaxX, rect1MinY, rect2MinX, rect2MinY, rect2MaxX, rect2MaxY) ||
		IsInRectMinMaxRange(rect1MaxX, rect1MaxY, rect2MinX, rect2MinY, rect2MaxX, rect2MaxY) ||
		IsInRectMinMaxRange(rect2MinX, rect2MinY, rect1MinX, rect1MinY, rect1MaxX, rect1MaxY) ||
		IsInRectMinMaxRange(rect2MinX, rect2MaxY, rect1MinX, rect1MinY, rect1MaxX, rect1MaxY) ||
		IsInRectMinMaxRange(rect2MaxX, rect2MinY, rect1MinX, rect1MinY, rect1MaxX, rect1MaxY) ||
		IsInRectMinMaxRange(rect2MaxX, rect2MaxY, rect1MinX, rect1MinY, rect1MaxX, rect1MaxY) {
		return true
	}
	return false
}

func IsInRectCenterRange(px, py, rectCenterX, rectCenterY, rectLength, rectWidth, negativeRectAngrad float64) bool {

	// 先判断是否在矩形外接圆内
	dx := px - rectCenterX
	dy := py - rectCenterY

	halfLength := rectLength / 2
	halfWidth := rectWidth / 2
	if dx*dx+dy*dy > halfLength*halfLength+halfWidth*halfWidth {
		return false
	}

	// 点p绕矩形中心点旋转到0度
	// 此时矩形坐标变成 (cx - len/2, cy - width/2), (cx + len/2, cy + width/2)
	nx, ny := Rotate(rectCenterX, rectCenterY, px, py, negativeRectAngrad)

	return IsInRectRange(nx, ny, rectCenterX-halfLength, rectCenterY-halfWidth, rectCenterX+halfLength, rectCenterY+halfWidth)
}
