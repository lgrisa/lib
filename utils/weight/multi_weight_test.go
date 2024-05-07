package weight

import (
	"fmt"
	. "github.com/onsi/gomega"
	"math/rand"
	"testing"
	"time"
)

func TestRandomN(t *testing.T) {
	RegisterTestingT(t)

	rand.Seed(time.Now().UnixNano())

	for _, weight := range [][]uint64{
		{0, 0}, {1, 0}, {0, 1}, {1, 0, 1},
	} {
		index, err := RandomN(weight, 1)
		Ω(err).Should(Equal(ErrZeroWeight))
		Ω(index).Should(BeNil())
	}

	for _, weight := range [][]uint64{
		{0}, {1, 1}, {2, 2, 2}, {3, 0, 3, 4},
	} {
		count := len(weight)

		index, err := RandomN(weight, count)
		Ω(err).Should(Succeed())

		a := make([]int, count)
		for i := 0; i < count; i++ {
			a[i] = i
		}
		Ω(index).Should(BeEquivalentTo(a))
	}

	for _, weight := range [][]uint64{
		{}, {0}, {1, 1}, {2, 2, 2}, {3, 0, 3, 4},
	} {
		count := len(weight)

		index, err := RandomN(weight, count+1)
		Ω(err).Should(Equal(ErrChoiceNotEnough))
		Ω(index).Should(BeNil())
	}

	weight := []uint64{0}
	index, err := RandomN(weight, 1)
	Ω(err).Should(Succeed())
	Ω(index).Should(BeEquivalentTo([]int{0}))

	weight = []uint64{}
	index, err = RandomN(weight, 1)
	Ω(err).Should(Equal(ErrChoiceNotEnough))
	Ω(index).Should(BeNil())

	for i := 0; i < 1000; i++ {
		weight := []uint64{10, 4, 1, 5, 6}
		index, err := RandomN(weight, rand.Intn(len(weight)-1)+1)
		Ω(err).Should(Succeed())

		// 没有重复元素
		for i, v0 := range index {
			for j, v1 := range index {
				if i != j {
					Ω(v0).ShouldNot(BeEquivalentTo(v1))
				}
			}

			Ω(v0 >= 0 && v0 < len(weight)).Should(BeTrue())
		}
	}
}

func TestRandomN2(t *testing.T) {
	RegisterTestingT(t)

	rand.Seed(time.Now().UnixNano())

	weight := []uint64{1, 2, 3, 4, 5}

	for n := 0; n < 5; n++ {
		timesMap := make(map[int]int)
		for i := 0; i < 10000; i++ {
			index, err := RandomN(weight, n+1)
			Ω(err).Should(Succeed())

			for _, v := range index {
				timesMap[v]++
			}
		}
		fmt.Println(timesMap)
	}

}

func TestRandomFilterN(t *testing.T) {
	RegisterTestingT(t)

	rand.Seed(time.Now().UnixNano())

	type WeightObject struct {
		index int
	}

	count := 100

	var array []*WeightObject
	for i := 0; i < count; i++ {
		array = append(array, &WeightObject{
			index: i,
		})
	}

	newArray, err := RandomFilterN(array, 10, func(w *WeightObject) uint64 {
		return 0
	})
	Ω(err).Should(Succeed())
	Ω(newArray).Should(BeEmpty())

	for i := 0; i < 10; i++ {
		newArray, err = RandomFilterN(array, 10, func(w *WeightObject) uint64 {
			if w.index < i {
				return 1
			}
			return 0
		})
		Ω(err).Should(Succeed())
		Ω(newArray).Should(HaveLen(i))

		indexMap := make(map[int]bool)
		for _, v := range newArray {
			Ω(v.index < i).Should(BeTrue())

			Ω(indexMap[v.index]).Should(BeFalse())
			indexMap[v.index] = true
		}
	}

	for i := 0; i < 1000; i++ {
		newArray, err = RandomFilterN(array, 10, func(w *WeightObject) uint64 {
			if w.index%2 == 0 {
				return 1
			}
			return 0
		})
		Ω(err).Should(Succeed())

		indexMap := make(map[int]bool)
		for _, v := range newArray {
			Ω(v.index % 2).Should(BeEquivalentTo(0))

			Ω(indexMap[v.index]).Should(BeFalse())
			indexMap[v.index] = true
		}

		newArray, err = RandomFilterN(array, 10, func(w *WeightObject) uint64 {
			if w.index%2 != 0 {
				return 1
			}
			return 0
		})
		Ω(err).Should(Succeed())

		indexMap = make(map[int]bool)
		for _, v := range newArray {
			Ω(v.index % 2).Should(BeEquivalentTo(1))

			Ω(indexMap[v.index]).Should(BeFalse())
			indexMap[v.index] = true
		}
	}

	//indexMap := make(map[int]int)
	//for i := 0; i < 10000; i++ {
	//	newArray, err = RandomFilterN(array, 10, func(w *WeightObject) uint64 {
	//		return uint64(w.index + 1)
	//	})
	//	Ω(err).Should(Succeed())
	//
	//	for _, v := range newArray {
	//		indexMap[v.index]++
	//	}
	//}
	//
	//for i := 0; i < count; i++ {
	//	fmt.Println(i, indexMap[i])
	//}

}
