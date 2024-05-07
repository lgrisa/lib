package i64

import (
	"strconv"
	"strings"
)

var Split = strings.Split

var FromString = strconv.ParseInt

func FromStringArray(sa []string) ([]int64, error) {
	var as []int64
	for _, s := range sa {
		a, err := FromString(s, 10, 64)
		if err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	return as, nil
}

func ParsingIntParam(toAdd, split string) ([]int64, error) {
	var list []int64
	strList := strings.Split(toAdd, split)
	if len(strList) == 0 {
		return list, nil
	}

	return FromStringArray(strList)
}
