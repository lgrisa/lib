package slice

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestSlice(t *testing.T) {
	RegisterTestingT(t)

	arr := []int{1, 1, 2, 3, 4}
	n := len(arr)

	Ω(Contains(arr, 4)).Should(BeTrue())
	arr = RemoveIndex(arr, 4)
	n -= 1
	Ω(len(arr)).Should(BeEquivalentTo(n))
	Ω(Contains(arr, 4)).Should(BeFalse())

	Ω(Contains(arr, 1)).Should(BeTrue())
	arr = Remove(arr, 1)
	n -= 2
	Ω(len(arr)).Should(BeEquivalentTo(n))
	Ω(Contains(arr, 1)).Should(BeFalse())

	Ω(Contains(arr, 2)).Should(BeTrue())
	arr = Remove(arr, 2)
	n -= 1
	Ω(len(arr)).Should(BeEquivalentTo(n))
	Ω(Contains(arr, 2)).Should(BeFalse())

	Ω(Contains(arr, 3)).Should(BeTrue())
	arr = RemoveIndex(arr, 0)
	n -= 1
	Ω(len(arr)).Should(BeEquivalentTo(n))
	Ω(Contains(arr, 3)).Should(BeFalse())

	Ω(arr).Should(BeEmpty())

	for i := 0; i < 10; i++ {
		arr = RemoveIndex(arr, i)
		Ω(arr).Should(BeEmpty())

		arr = Remove(arr, i)
		Ω(arr).Should(BeEmpty())

		Ω(Contains(arr, i)).Should(BeFalse())
	}
}
