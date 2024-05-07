package math2d

import "math"

func Move(posX, posY, angrad, distance float64) (targetX, targetY float64) {
	if distance == 0 {
		return posX, posY
	}

	sin, cos := math.Sincos(angrad)
	targetX = posX + distance*cos
	targetY = posY + distance*sin
	return
}

func MoveToNearby(startX, startY, targetX, targetY, distance float64) (float64, float64) {
	radians := Point2Radians(targetX, targetY, startX, startY)
	return Move(targetX, targetY, radians, distance)
}

func CurrentLinePoint(startX, startY, targetX, targetY float64, startFrame, endFrame, currentFrame int64) (float64, float64) {
	if currentFrame <= startFrame {
		return startX, startY
	}

	if currentFrame >= endFrame {
		return targetX, targetY
	}

	// 计算位移的百分比
	percent := float64(currentFrame-startFrame) / float64(endFrame-startFrame)

	// 计算位移的坐标
	x := startX + (targetX-startX)*percent
	y := startY + (targetY-startY)*percent
	return x, y
}
