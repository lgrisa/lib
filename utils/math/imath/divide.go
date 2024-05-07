package imath

func DivideTimes(x, y int) int {
	if x <= 0 || y <= 0 {
		return 0
	}

	return (x + y - 1) / y
}
