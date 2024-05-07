package i64

func Max(x, y int64) int64 {
	if x < y {
		return y
	}

	return x
}

func MaxSlice(slice ...int64) int64 {
	var max int64
	for i, v := range slice {
		if i > 0 {
			max = Max(max, v)
		} else {
			max = v
		}
	}
	return max
}
