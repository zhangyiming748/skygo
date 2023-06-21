package custom_util

import (
	"strconv"
	"strings"
)

// 数字组成的切片，转为string, 使用逗号分割
func IntsToString(ids []int) string {
	sList := make([]string, 0)
	for _, id := range ids {
		s := strconv.Itoa(id)
		if s != "" {
			sList = append(sList, s)
		}
	}
	return strings.Join(sList, ",")
}
