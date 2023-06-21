package custom_util

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"skygo_detection/guardian/src/net/qmap"
	"strconv"
)

// 存在且必须是字符型
func FetchJsonString(m gin.H, str string) (string, bool) {
	if v, has := m[str]; has {
		if _v, ok := v.(string); ok {
			return _v, true
		}
	}
	return "", false
}

// 存在且必须是数字类型，转int
func FetchJsonInt(m gin.H, str string) (int, bool) {
	if v, has := m[str]; has {
		switch v.(type) {
		case float64:
			return int(v.(float64)), true
		case int:
			return v.(int), true
		}
	}
	return 0, false
}

func FetchJsonMap(m gin.H, str string) (gin.H, bool) {
	if v, has := m[str]; has {
		if _v, ok := v.(gin.H); ok {
			return _v, true
		}
	}
	return nil, false
}

func FetchJsonSlice(m gin.H, str string) ([]interface{}, bool) {
	if v, has := m[str]; has {
		if _v, ok := v.([]interface{}); ok {
			return _v, true
		}
	}
	return nil, false
}

// 从map中复制指定列的内容
//
//columns:map[string]string 形如{"column_name":"string"}, key:列名, val:内容类型(可选:string/int)
func CopyMapColumns(rawMap qmap.QM, columns map[string]string) qmap.QM {
	result := qmap.QM{}
	for column, columnType := range columns {
		switch columnType {
		case "int":
			if val, exist := rawMap.TryInt(column); exist {
				result[column] = val
			}
		case "string":
			if val, exist := rawMap.TryString(column); exist {
				result[column] = val
			}
		default:
			if _, exist := rawMap[column]; exist {
				result[column] = rawMap.Interface(column)
			}
		}
	}
	return result
}

func MapToSlice(m map[int]qmap.QM) []interface{} {
	s := make([]interface{}, len(m), len(m))
	for index, v := range m {
		s[index-1] = v
	}
	return s
}

func MapHasString(m map[string]interface{}, key string) (string, bool) {
	v, ok := m[key]
	if ok == true {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			return v.(string), true
		case reflect.Float64:
			return strconv.Itoa(int(v.(float64))), true
		}
	}
	return "", false
}

func MapHasInt(m map[string]interface{}, key string) (int, bool) {
	v, ok := m[key]
	if ok == true {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Int:
			return v.(int), true
		case reflect.Float64:
			return int(v.(float64)), true
		case reflect.String:
			i, err := strconv.Atoi(v.(string))
			if err != nil {
				panic(err)
			}
			return i, true
		}
	}
	return 0, false
}
