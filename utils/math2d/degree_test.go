package math2d

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestDegree(t *testing.T) {
	RegisterTestingT(t)

	Ω(Atan2Degree(0, 0)).Should(BeEquivalentTo(0))
	Ω(Atan2Degree(1, 0)).Should(BeEquivalentTo(0))
	Ω(Atan2Degree(-1, 0)).Should(BeEquivalentTo(180))
	Ω(Atan2Degree(0, 1)).Should(BeEquivalentTo(90))
	Ω(Atan2Degree(0, -1)).Should(BeEquivalentTo(270))
}

func TestToDegree360(t *testing.T) {
	RegisterTestingT(t)

	for i := int64(0); i < 10000; i++ {
		if i < 360 {
			Ω(Degree360(i)).Should(BeEquivalentTo(i))

			if i > 0 {
				Ω(Degree360(-i)).Should(BeEquivalentTo(360 - i))
			} else {
				Ω(Degree360(-i)).Should(BeEquivalentTo(0))
			}
		}

		mod := i % 360
		Ω(Degree360(i)).Should(BeEquivalentTo(mod))

		if mod > 0 {
			Ω(Degree360(-i)).Should(BeEquivalentTo(360 - mod))
		} else {
			Ω(Degree360(-i)).Should(BeEquivalentTo(0))
		}
	}
}

func TestDegreeBetween(t *testing.T) {
	RegisterTestingT(t)

	for i := 0; i < 1000; i++ {
		d1 := int64(i)

		for i := 0; i <= 180; i++ {
			diff := int64(i)
			d2 := d1 + diff
			d3 := d1 - diff

			Ω(GetAngleBetween(d1, d2)).Should(BeEquivalentTo(diff))
			Ω(GetAngleBetween(d1, d3)).Should(BeEquivalentTo(diff))
		}

		for i := 181; i < 360; i++ {
			diff := 360 - int64(i)
			d2 := d1 + diff
			d3 := d1 - diff

			Ω(GetAngleBetween(d1, d2)).Should(BeEquivalentTo(diff))
			Ω(GetAngleBetween(d1, d3)).Should(BeEquivalentTo(diff))
		}
	}

}
