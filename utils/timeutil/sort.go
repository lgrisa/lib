package timeutil

import (
	"sort"
	"time"
)

func Sort(l []time.Duration) {
	sort.Sort(DurationSlice(l))
}

type DurationSlice []time.Duration

func (p DurationSlice) Len() int           { return len(p) }
func (p DurationSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p DurationSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func Contain(poses []int64, pos int64) bool {
	for _, p := range poses {
		if p == pos {
			return true
		}
	}

	return false
}
