package decimal

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestYuan2Fen(t *testing.T) {
	RegisterTestingT(t)

	Ω(Yuan2Fen("0.15")).Should(BeEquivalentTo(15))
	Ω(Yuan2Fen("0.14")).Should(BeEquivalentTo(14))
	Ω(Yuan2Fen("0.10")).Should(BeEquivalentTo(10))
	Ω(Yuan2Fen("0.1")).Should(BeEquivalentTo(10))
	Ω(Yuan2Fen("1.15")).Should(BeEquivalentTo(115))
}
