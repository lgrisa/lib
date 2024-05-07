package f64

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestSorted(t *testing.T) {
	RegisterTestingT(t)

	Ω(IsAscSorted([]float64{})).Should(BeTrue())
	Ω(IsAscSorted([]float64{1})).Should(BeTrue())
	Ω(IsAscSorted([]float64{1, 3, 9})).Should(BeTrue())
	Ω(IsAscSorted([]float64{1, 1, 2, 9})).Should(BeTrue())
	Ω(IsAscSorted([]float64{1, 3, 2, 9})).Should(BeFalse())
	Ω(IsAscSorted([]float64{4, 3, 5, 9})).Should(BeFalse())

	Ω(IsDescSorted([]float64{})).Should(BeTrue())
	Ω(IsDescSorted([]float64{1})).Should(BeTrue())
	Ω(IsDescSorted([]float64{5, 2, 1})).Should(BeTrue())
	Ω(IsDescSorted([]float64{5, 2, 2, 1})).Should(BeTrue())

	Ω(IsDescSorted([]float64{9, 5, 6, 1})).Should(BeFalse())
	Ω(IsDescSorted([]float64{1, 5, 2, 2, 1})).Should(BeFalse())
}
