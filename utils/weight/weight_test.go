package weight

import (
	. "github.com/onsi/gomega"
	"math/rand"
	"testing"
	"time"
)

func TestWeightRandomer_Index(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	weight := []uint64{}
	r, err := NewWeightRandomer(weight)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{0}
	r, err = NewWeightRandomer(weight)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{0, 1}
	r, err = NewWeightRandomer(weight)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{1, 0, 1}
	r, err = NewWeightRandomer(weight)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{1, 2, 3}
	r, err = NewWeightRandomer(weight)
	Ω(err).Should(Succeed())

	Ω(r.Index(0)).Should(Equal(0))
	Ω(r.Index(1)).Should(Equal(1))
	Ω(r.Index(2)).Should(Equal(1))
	Ω(r.Index(3)).Should(Equal(2))
	Ω(r.Index(4)).Should(Equal(2))
	Ω(r.Index(5)).Should(Equal(2))

	Ω(r.Index(6)).Should(Equal(2))
	Ω(r.Index(7)).Should(Equal(2))
}

func TestRandomFilter(t *testing.T) {
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

	for i := 0; i < count; i++ {
		v := RandomFilter(array, func(w *WeightObject) uint64 {
			if w.index == i {
				return uint64(i + 1)
			}
			return 0
		})
		Ω(v.index).Should(BeEquivalentTo(i))
	}

	//indexMap := make(map[int]int)
	//for i := 0; i < 100000; i++ {
	//	v := RandomFilter(array, func(w *WeightObject) uint64 {
	//		//return 0
	//		return uint64(w.index + 1)
	//	})
	//	indexMap[v.index]++
	//}
	//for i := 0; i < count; i++ {
	//	fmt.Println(i, indexMap[i])
	//}

}
