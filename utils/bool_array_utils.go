package utils

func NewBoolArray(maxN int) *BoolArray {
	if maxN < 0 {
		panic("NewBoolArray maxN must >= 0")
	}

	n := (maxN + 32 - 1) / 32

	ba := &BoolArray{}
	ba.array = make([]uint32, n)
	ba.maxN = maxN

	return ba
}

type BoolArray struct {
	array []uint32
	maxN  int
}

func (d *BoolArray) GetBool(i int) bool {
	if i < 0 || i >= d.maxN {
		return false
	}

	arrayIndex := i / 32
	offsetIndex := uint(i % 32)

	return d.array[arrayIndex]&(1<<offsetIndex) != 0
}

func (d *BoolArray) SetTrue(i int) {
	if i < 0 || i >= d.maxN {
		return
	}

	arrayIndex := i / 32
	offsetIndex := uint(i % 32)
	d.array[arrayIndex] = d.array[arrayIndex] | (1 << offsetIndex)
}

func (d *BoolArray) SetFalse(i int) {
	if i < 0 || i >= d.maxN {
		return
	}

	arrayIndex := i / 32
	offsetIndex := uint(i % 32)
	d.array[arrayIndex] = d.array[arrayIndex] & (^(1 << offsetIndex))
}

func (d *BoolArray) Encode() []uint32 {
	v := make([]uint32, len(d.array))
	copy(v, d.array)
	return v
}

func (d *BoolArray) Decode(v []uint32) {
	copy(d.array, v)
}
