package u64

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestRandomN(t *testing.T) {
	RegisterTestingT(t)

	Ω(RandomN(0)).Should(BeEquivalentTo(0))
	Ω(RandomN(1)).Should(BeEquivalentTo(0))

	for i := 0; i < 100; i++ {
		Ω(RandomN(50) < 50).Should(BeTrue())
	}

	Ω(RandomRange(0, 0)).Should(BeEquivalentTo(0))
	Ω(RandomRange(1, 0)).Should(BeEquivalentTo(0))
	Ω(RandomRange(0, 1)).Should(BeEquivalentTo(0))
	Ω(RandomRange(1, 1)).Should(BeEquivalentTo(1))

	for i := 0; i < 100; i++ {
		n := RandomRange(30, 70)
		Ω(n >= 30 && n < 70).Should(BeTrue())

		n = RandomRange(70, 30)
		Ω(n >= 30 && n < 70).Should(BeTrue())
	}
}
