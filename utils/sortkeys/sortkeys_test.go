package sortkeys

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestSort(t *testing.T) {
	RegisterTestingT(t)

	arr := []uint64{1, 5, 6}

	Ω(SearchU64s(arr, 0)).Should(BeEquivalentTo(0))
	Ω(SearchU64s(arr, 1)).Should(BeEquivalentTo(0))
	Ω(SearchU64s(arr, 2)).Should(BeEquivalentTo(1))
	Ω(SearchU64s(arr, 5)).Should(BeEquivalentTo(1))
	Ω(SearchU64s(arr, 6)).Should(BeEquivalentTo(2))
	Ω(SearchU64s(arr, 7)).Should(BeEquivalentTo(3))

}
