package custom_util

import (
	"fmt"
	"sort"
)

// 求并集
func UnionInt(slice1, slice2 []int) []int {
	m := make(map[int]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// 求交集
func IntersectInt(slice1, slice2 []int) []int {
	m := make(map[int]int)
	nn := make([]int, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times >= 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求差集 slice1-并集
func DifferenceInt(slice1, slice2 []int) []int {
	m := make(map[int]int)
	nn := make([]int, 0)
	inter := IntersectInt(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

func InStrSlice(arr []string, find string) bool {
	for _, v := range arr {
		if v == find {
			return true
		}
	}
	return false
}

func UnionStr(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// 求交集
func IntersectStr(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求差集 slice1-并集
func DifferenceStr(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := IntersectStr(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

// 将分片map转为key-value格式
func SliceMapToOneColumnMapInt(mapData []map[string]int, key string, valKey string) map[int]int {
	returnData := make(map[int]int)
	for _, tmpV := range mapData {
		if v, ok := tmpV[key]; ok {
			returnData[v] = tmpV[valKey]
		}
	}
	return returnData
}

// 将分片map id name 转为key-value格式
func SliceMapToIdNameMap(mapData []map[string]interface{}) map[int]string {
	returnData := make(map[int]string)
	for _, tmpV := range mapData {
		if id, ok := tmpV["id"].(int64); ok {
			returnData[int(id)], _ = tmpV["name"].(string)
		}
	}
	return returnData
}

func SliceRemoveDuplicates(slice []string) []string {
	sort.Strings(slice)
	i := 0
	var j int
	for {
		if i >= len(slice)-1 {
			break
		}
		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {
		}
		slice = append(slice[:i+1], slice[j:]...)
		i++
	}
	return slice
}

func SliceIntMin(slice []int) int {
	if len(slice) == 0 {
		return 0
	}
	checkNum := slice[0]

	for _, v := range slice {
		if checkNum > v {
			checkNum = v
		}
	}
	return checkNum
}

func SliceIntMax(slice []int) int {
	if len(slice) == 0 {
		return 0
	}
	checkNum := slice[0]

	for _, v := range slice {
		if checkNum < v {
			checkNum = v
		}
	}
	return checkNum
}

// SliceUnique 数组去重
func SliceUnique(input []string) []string {
	result := []string{}
	// 存放不重复主键
	tempMap := map[string]byte{}
	for _, e := range input {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			// 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}

	return result
}

// SliceInt2Str 整型slice车string slice
func SliceInt2Str(input []int) []string {
	sliceStr := []string{}
	if input == nil {
		return sliceStr
	}

	for _, v := range input {
		sliceStr = append(sliceStr, fmt.Sprintf("%d", v))
	}

	return sliceStr
}

func IndexOfSlice(needle string, haystack []string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}
	return false
}

// InStringSlice is 在分片中查找字符串是否存在
func InStringSlice(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}
