package math2d

import (
	. "github.com/onsi/gomega"
	"math"
	"testing"
)

func TestMove(t *testing.T) {
	RegisterTestingT(t)

	checkPos := func(targetX, targetY, exceptX, exceptY float64) {
		Ω(multi1000(targetX)).Should(BeEquivalentTo(multi1000(exceptX)))
		Ω(multi1000(targetY)).Should(BeEquivalentTo(multi1000(exceptY)))
	}

	type data struct {
		posX, posY       float64
		degree           int64
		distance         float64
		targetX, targetY float64
	}

	newData := func(posX, posY float64, degree int64, distance, targetX, targetY float64) *data {
		return &data{
			posX:     posX,
			posY:     posY,
			degree:   degree,
			distance: distance,
			targetX:  targetX,
			targetY:  targetY,
		}
	}

	datas := []*data{
		newData(0, 0, 0, 1, 1, 0),
		newData(0, 0, 90, 1, 0, 1),
		newData(0, 0, 180, 1, -1, 0),
		newData(0, 0, 270, 1, 0, -1),
		newData(0, 0, 360, 1, 1, 0),

		newData(0, 0, 0, 3, 3, 0),
		newData(0, 0, 90, 3, 0, 3),
		newData(0, 0, 180, 3, -3, 0),
		newData(0, 0, 270, 2, 0, -2),
		newData(0, 0, 360, 2, 2, 0),

		newData(0, 0, 45, 1, 1/math.Sqrt(2), 1/math.Sqrt(2)),
	}

	for _, d := range datas {
		targetX, targetY := Move(d.posX, d.posY, DegreeToAngrad(d.degree), d.distance)
		checkPos(targetX, targetY, d.targetX, d.targetY)
	}

}

func TestCurrentLinePoint(t *testing.T) {
	RegisterTestingT(t)

	type data struct {
		startX, startY, targetX, targetY   float64
		startFrame, endFrame, currentFrame int64

		exceptX, exceptY float64
	}

	var array = []data{
		{0, 0, 1, 1, 0, 10, -1, 0, 0},
		{0, 0, 1, 1, 0, 10, 0, 0, 0},

		{0, 0, 1, 2, 0, 10, 1, 0.1, 0.2},
		{0, 0, 1, 2, 0, 10, 2, 0.2, 0.4},
		{0, 0, 1, 2, 0, 10, 5, 0.5, 1},
		{0, 0, 1, 2, 0, 10, 8, 0.8, 1.6},
		{0, 0, 1, 2, 0, 10, 9, 0.9, 1.8},

		{0, 0, 1, 1, 0, 10, 10, 1, 1},
		{0, 0, 1, 1, 0, 10, 100, 1, 1},

		// 异常数据
		{0, 0, 1, 1, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 0, 0, 1, 1, 1},

		{0, 0, 1, 1, 0, -1, -1, 0, 0},
		{0, 0, 1, 1, 0, -1, 1, 1, 1},
	}

	for _, v := range array {
		x, y := CurrentLinePoint(v.startX, v.startY, v.targetX, v.targetY, v.startFrame, v.endFrame, v.currentFrame)
		Ω(multi1000(x)).Should(BeEquivalentTo(multi1000(v.exceptX)))
		Ω(multi1000(y)).Should(BeEquivalentTo(multi1000(v.exceptY)))
	}
}
