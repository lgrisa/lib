package i64

func Between(v, min, max int64) int64 {

	if max < min {
		max, min = min, max
	}

	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func Between32(v int32, min, max int64) int64 {
	return Between(int64(v), min, max)
}
