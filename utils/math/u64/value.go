package u64

func NewValue(v uint64) *Value {
	return &Value{value: v}
}

type Value struct {
	value uint64
}

func (v *Value) Get() uint64 {
	return v.value
}

func (v *Value) Set(toSet uint64) {
	v.value = toSet
}

func (v *Value) Add(toAdd uint64) uint64 {
	v.value += toAdd
	return v.value
}

func (v *Value) Sub(toSub uint64) uint64 {
	v.value = Sub(v.value, toSub)
	return v.value
}
