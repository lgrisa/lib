package math2d

import (
	. "github.com/onsi/gomega"
	"math"
	"testing"
)

func TestRotate(t *testing.T) {
	RegisterTestingT(t)

	bx, by := RotateOrigin(1, 0, DegreeToAngrad(90))
	Ω(multi1000(bx)).Should(BeEquivalentTo(0))
	Ω(multi1000(by)).Should(BeEquivalentTo(1000))

	bx, by = RotateOrigin(1, 0, DegreeToAngrad(180))
	Ω(multi1000(bx)).Should(BeEquivalentTo(-1000))
	Ω(multi1000(by)).Should(BeEquivalentTo(0))

	bx, by = RotateOrigin(1, 0, DegreeToAngrad(270))
	Ω(multi1000(bx)).Should(BeEquivalentTo(0))
	Ω(multi1000(by)).Should(BeEquivalentTo(-1000))

	for i := 0; i < 1000; i++ {
		degree := int64(i)
		angard := DegreeToAngrad(degree)
		bx, by := RotateOrigin(1, 0, angard)

		sin, cos := math.Sincos(angard)

		Ω(multi1000(bx)).Should(BeEquivalentTo(multi1000(cos)))
		Ω(multi1000(by)).Should(BeEquivalentTo(multi1000(sin)))

	}
}

func TestRotate2(t *testing.T) {
	RegisterTestingT(t)

	bx, by := Rotate(1, 1, 2, 1, DegreeToAngrad(90))
	Ω(multi1000(bx)).Should(BeEquivalentTo(1000))
	Ω(multi1000(by)).Should(BeEquivalentTo(2000))

	bx, by = Rotate(1, 1, 2, 1, DegreeToAngrad(180))
	Ω(multi1000(bx)).Should(BeEquivalentTo(0))
	Ω(multi1000(by)).Should(BeEquivalentTo(1000))

	bx, by = Rotate(1, 1, 2, 1, DegreeToAngrad(270))
	Ω(multi1000(bx)).Should(BeEquivalentTo(1000))
	Ω(multi1000(by)).Should(BeEquivalentTo(0))
}

func multi1000(coef float64) int64 {
	fd := float64(1000)
	tiny := 1 / (fd * 10)

	if coef < 0 {
		tiny = -tiny
	}

	return int64((coef + tiny) * fd)
}
