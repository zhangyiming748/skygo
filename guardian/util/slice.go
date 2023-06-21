package util

import (
	"fmt"
	"strconv"
	"strings"
)

func InStringSlice(val string, slice []string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func Join(a interface{}, sep string) (j string) {
	switch a.(type) {
	case []string:
		j = strings.Join(a.([]string), sep)
	case []int:
		t := a.([]int)
		if len(t) > 0 {
			j = fmt.Sprintf("%v", t[0])
			for i := 1; i < len(t); i++ {
				j += sep + fmt.Sprintf("%v", t[i])
			}
		}
	case []int64:
		t := a.([]int64)
		if len(t) > 0 {
			j = fmt.Sprintf("%v", t[0])
			for i := 1; i < len(t); i++ {
				j += sep + fmt.Sprintf("%v", t[i])
			}
		}
	case []int32:
		t := a.([]int32)
		if len(t) > 0 {
			j = fmt.Sprintf("%v", t[0])
			for i := 1; i < len(t); i++ {
				j += sep + fmt.Sprintf("%v", t[i])
			}
		}
	default:
		panic("unsupported join type")
	}

	return
}

func SplitInt(s, sep string) []int {
	split := []int{}
	for _, s := range strings.Split(s, sep) {
		if t, err := strconv.Atoi(s); err == nil {
			split = append(split, t)
		}
	}
	return split
}
