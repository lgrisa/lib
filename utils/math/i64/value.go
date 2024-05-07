package i64

func NewValue(v int64) *Value {
	return &Value{V: v}
}

type Value struct {
	V int64
}
