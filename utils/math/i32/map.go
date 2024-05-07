package i32

import "github.com/lgrisa/lib/utils/math/imath"

func Map2Array(m map[int32]int32) (keys, values []int32) {
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return
}

func ArrayToMap(keys, values []int32) (m map[int32]int32) {
	m = make(map[int32]int32)
	if n := imath.Min(len(keys), len(values)); n > 0 {
		for i := 0; i < n; i++ {
			m[keys[i]] = values[i]
		}
	}
	return
}

func MapKey2Arrary(m map[int32]struct{}) (keys []int32) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func ArrayToMapKey(keys []int32) (m map[int32]struct{}) {
	m = make(map[int32]struct{})
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return
}

func CopyMap(a map[int32]int32) map[int32]int32 {
	out := make(map[int32]int32, len(a))
	CopyMapTo(out, a)
	return out
}

func CopyMapTo(dest, src map[int32]int32) {
	for k, v := range src {
		dest[k] = v
	}
}
