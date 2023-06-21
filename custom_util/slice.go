package custom_util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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

func Intersect(slice1, slice2 []string) []string { // 取两个切片的交集
	m := make(map[string]int)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}
	return n
}

func Difference(slice1, slice2 []string) []string { // 取要校验的和已经校验过的差集
	m := make(map[string]int)
	n := make([]string, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range slice1 {
		if m[value] == 0 {
			n = append(n, value)
		}
	}

	for _, v := range slice2 {
		if m[v] == 0 {
			n = append(n, v)
		}
	}
	return n
}

func DifferenceDel(slice1, slice2 []string) []string { // 取出slice1中有，slice2中没有的
	m := make(map[string]int)
	n := make([]string, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range slice1 {
		if m[value] == 0 {
			n = append(n, value)
		}
	}
	return n
}

func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

func InterfaceToString(slice []interface{}) []string {
	result := []string{}
	for _, item := range slice {
		result = append(result, item.(string))
	}
	return result
}

func InterfaceToInt(slice []interface{}) []int {
	result := []int{}
	for _, item := range slice {
		result = append(result, item.(int))
	}
	return result
}

func InterfaceToInt64(slice []interface{}) []int64 {
	result := []int64{}
	for _, item := range slice {
		result = append(result, item.(int64))
	}
	return result
}

func Implode(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

func SliceToString(obj []interface{}) string {
	jsonObj, _ := json.Marshal(obj)
	return string(jsonObj)
}

func StringToSlice(str string) ([]interface{}, error) {
	s := []interface{}{}
	err := json.Unmarshal([]byte(str), &s)
	return s, err
}

func SliceToSqlInString(slice []int, column string) string {
	str := strings.Replace(Join(slice, ","), ",", "\",\"", -1)
	return fmt.Sprint(column, " in (", "\"", str, "\")")
}
