package i64

func Min(x, y int64) int64 {
	if x > y {
		return y
	}

	return x
}

func MinSlice(slice ...int64) int64 {
	var min int64
	for i, v := range slice {
		if i > 0 {
			min = Min(min, v)
		} else {
			min = v
		}
	}
	return min
}
