package i64

import (
	"math"
)

func Int32(x int64) int32 {
	if x >= 0 {
		return positiveInt32(x)
	} else {
		return negativeInt32(x)
	}
}

func Int32Array(x []int64) []int32 {
	ia := make([]int32, len(x))
	for i, v := range x {
		ia[i] = Int32(v)
	}

	return ia
}

func FromInt32Array(x []int32) []int64 {
	ia := make([]int64, len(x))
	for i, v := range x {
		ia[i] = int64(v)
	}

	return ia
}

func positiveInt32(x int64) int32 {
	return int32(Min(Max(x, 0), math.MaxInt32))
}

func negativeInt32(x int64) int32 {
	return int32(Max(Min(x, 0), math.MinInt32))
}

type GetU64 func(k int64) uint64

func NewGetU64(m map[int64]uint64) GetU64 {
	return func(k int64) uint64 {
		return m[k]
	}
}

func EmptyGetU64() GetU64 {
	return func(k int64) uint64 {
		return 0
	}
}

func Rate(start, end, cur int64) float64 {

	// 特殊情况处理
	if end <= cur {
		return 1
	}

	if cur <= start {
		return 0
	}

	diff := cur - start
	total := end - start

	return float64(diff) / float64(total)
}

func MultiF64(d int64, coef float64) int64 {
	if d == 0 || coef <= 0 {
		return 0
	}

	fd := float64(d)
	return int64((coef + (1 / (fd * 10))) * fd)
}

func ToI32KeyMap[K int64 | uint64, V any](from map[K]V) (to map[int32]V) {
	if n := len(from); n > 0 {
		to = make(map[int32]V, n)
		for k, v := range from {
			to[int32(k)] = v
		}
	}
	return
}

func ToI32ValMap[K comparable, V int64 | uint64](from map[K]V) (to map[K]int32) {
	if n := len(from); n > 0 {
		to = make(map[K]int32, n)
		for k, v := range from {
			to[k] = int32(v)
		}
	}
	return
}

func ToI32KeyValueMap[K int64 | uint64, V int64 | uint64](from map[K]V) (to map[int32]int32) {
	if n := len(from); n > 0 {
		to = make(map[int32]int32, n)
		for k, v := range from {
			to[int32(k)] = int32(v)
		}
	}
	return
}
