package imath

func Must1(i, n int) int {
	ii := i
	if ii > 0 && ii <= n {
		return ii - 1
	}

	if ii <= 0 {
		return 0
	}

	return n - 1
}

func Must0(i, n int) int {
	ii := i
	if ii >= 0 && ii < n {
		return ii
	}

	if ii < 0 {
		return 0
	}

	return n - 1
}
