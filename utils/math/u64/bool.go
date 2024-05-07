package u64

func NewBool64() *Bool64 {
	return &Bool64{}
}

type Bool64 struct {
	value uint64
}

func (b *Bool64) GetValue() uint64 {
	return b.value
}

func (b *Bool64) SetValue(toSet uint64) {
	b.value = toSet
}

func (b *Bool64) Get(index int) bool {
	if index >= 0 && index < 64 {
		return b.value&(1<<uint64(index)) != 0
	}

	return false
}

func (b *Bool64) TrySet(index int) bool {
	if !b.Get(index) {
		b.SetTrue(index)
		return true
	}
	return false
}

func (b *Bool64) SetTrue(index int) {
	if index >= 0 && index < 64 {
		b.value = b.value | 1<<uint64(index)
	}
}

func (b *Bool64) SetFalse(index int) {
	if index >= 0 && index < 64 {
		b.value = b.value & ^(1 << uint64(index))
	}
}


func BoolToU64(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}