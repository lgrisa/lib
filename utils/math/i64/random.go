package i64

import "math/rand"

func Random(min, max int64) int64 {
	if min < max {
		return min + rand.Int63n(1+max-min)
	}
	return min
}
