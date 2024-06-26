package i32

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestUint64(t *testing.T) {
	RegisterTestingT(t)

	Ω(Uint64(0)).Should(Equal(uint64(0)))
	Ω(Uint64(10)).Should(Equal(uint64(10)))
	Ω(Uint64(-1)).Should(Equal(uint64(0)))
	Ω(Uint64(-100)).Should(Equal(uint64(0)))
}

func TestMultiF64(t *testing.T) {
	RegisterTestingT(t)

	f20 := 2.0
	f18 := 1.8
	f02 := f20 - f18

	// 实际上这货 != 0.2
	fmt.Println(f02)
	Ω(f02).ShouldNot(Equal(0.2))

	for i := 0; i < 9; i++ {

		d := 10
		e := 2

		n := i
		for i := 0; i < n; i++ {
			d *= 10
			e *= 10
		}
		Ω(MultiF64(uint64(d), f02)).Should(Equal(int32(e)))
	}

}

func TestIndex(t *testing.T) {
	RegisterTestingT(t)

	Ω(Must1(-1, 10)).Should(BeEquivalentTo(0))
	Ω(Must1(0, 10)).Should(BeEquivalentTo(0))
	Ω(Must1(1, 10)).Should(BeEquivalentTo(0))
	Ω(Must1(2, 10)).Should(BeEquivalentTo(1))
	Ω(Must1(10, 10)).Should(BeEquivalentTo(9))
	Ω(Must1(11, 10)).Should(BeEquivalentTo(9))

	Ω(Must0(-1, 10)).Should(BeEquivalentTo(0))
	Ω(Must0(0, 10)).Should(BeEquivalentTo(0))
	Ω(Must0(1, 10)).Should(BeEquivalentTo(1))
	Ω(Must0(5, 10)).Should(BeEquivalentTo(5))
	Ω(Must0(9, 10)).Should(BeEquivalentTo(9))
	Ω(Must0(10, 10)).Should(BeEquivalentTo(9))

}
