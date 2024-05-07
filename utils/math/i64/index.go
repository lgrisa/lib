package i64

func Must1(i int64, n int) int {
	ii := int(i)
	if ii > 0 && ii <= n {
		return ii - 1
	}

	if ii <= 0 {
		return 0
	}

	return n - 1
}

func Must0(i int64, n int) int {
	ii := int(i)
	if ii >= 0 && ii < n {
		return ii
	}

	if ii < 0 {
		return 0
	}

	return n - 1
}
