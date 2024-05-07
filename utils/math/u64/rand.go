package u64

import "math/rand"

func RandomN(n uint64) uint64 {
	switch n {
	case 0:
		return 0
	case 1:
		return 0
	default:
		return rand.Uint64() % n
	}
}

func RandomRange(min, max uint64) uint64 {
	a, b := min, max
	if min > max {
		a, b = max, min
	}
	return a + RandomN(b-a)
}
