package sortkeys

import (
	"sort"
)

func NewU64TopN(n uint64) *U64TopN {
	return &U64TopN{
		n:     n,
		array: make([]uint64, 0, n),
	}
}

type U64TopN struct {
	n uint64

	array []uint64

	minKey      uint64
	minKeyIndex uint64
}

func (t *U64TopN) Sum() uint64 {
	var sum uint64
	for _, v := range t.array {
		sum += v
	}
	return sum
}

func (t *U64TopN) Array() []uint64 {
	return t.array
}

func (t *U64TopN) CopyArray() []uint64 {
	array := make([]uint64, len(t.array))
	copy(array, t.array)
	return array
}

func (t *U64TopN) SortAsc() []uint64 {
	array := t.CopyArray()
	Uint64s(array)
	return array
}

func (t *U64TopN) SortDesc() []uint64 {
	array := t.CopyArray()
	sort.Sort(sort.Reverse(Uint64Slice(array)))
	return array
}

func (t *U64TopN) Size() int {
	return len(t.array)
}

func (t *U64TopN) Add(k uint64) {

	oldLen := uint64(len(t.array))
	if oldLen < t.n {

		t.array = append(t.array, k)

		if oldLen == 0 || k < t.minKey {
			t.minKey = k
			t.minKeyIndex = oldLen
		}
		return
	}

	// topN已经满了，如果超过最小值，那么跟最小值进行交换，然后重新获取最小值
	if t.minKey < k {
		t.array[t.minKeyIndex] = k

		for i, v := range t.array {

			if i == 0 || v < t.minKey {
				t.minKey = v
				t.minKeyIndex = uint64(i)
			}
		}
	}

}
