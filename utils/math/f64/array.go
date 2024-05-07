package f64

func IsAscSorted(a []float64) bool {
	var prev float64
	for i, v := range a {
		if i > 0 && prev > v {
			return false
		}
		prev = v
	}
	return true
}

func IsDescSorted(a []float64) bool {
	var prev float64
	for i, v := range a {
		if i > 0 && prev < v {
			return false
		}
		prev = v
	}
	return true
}
