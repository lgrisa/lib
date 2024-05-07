package i64

type Map map[int64]int64

func Map2Int32Array(m map[int64]int64) (keys, values []int32) {

	for k, v := range m {
		keys = append(keys, Int32(k))
		values = append(values, Int32(v))
	}

	return
}

func MapKey2Int64Array(m map[int64]struct{}) (keys []int64) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func Int64ArrayToMapKey(keys []int64) (m map[int64]struct{}) {
	m = make(map[int64]struct{})
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return
}

func CopyMap(a map[int64]int64) map[int64]int64 {
	out := make(map[int64]int64, len(a))
	CopyMapTo(out, a)
	return out
}

func CopyMapTo[K int64 | uint64 | int32 | uint32, V any](dest, src map[K]V) {
	for k, v := range src {
		dest[k] = v
	}
}
