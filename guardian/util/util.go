package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

/*
*
结构体对象转map(key由驼峰转蛇形)
*/
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		snakeName, _ := SnakeString(obj1.Field(i).Name)
		data[snakeName] = obj2.Field(i).Interface()
	}
	return data
}

func ResolvePointValue(value interface{}) interface{} {
	valueType := reflect.TypeOf(value).Kind()

	if valueType == reflect.Ptr {
		direct := reflect.Indirect(reflect.ValueOf(value))
		if direct.CanAddr() {
			return direct.Interface()
		} else {
			return new(interface{})
		}
	} else {
		return value
	}
}

/*
*
结构体对象转map(key由驼峰转蛇形)(只转换指定列)
*/
func StructToMapWithColumns(obj interface{}, columns map[string]bool) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		fieldName := obj1.Field(i).Name
		if _, exist := columns[fieldName]; exist {
			snakeName, _ := SnakeString(fieldName)
			data[snakeName] = obj2.Field(i).Interface()
		}
	}
	return data
}

func SnakeString(s string) (string, bool) {
	data := make([]byte, 0, len(s)*2)
	change := false
	j := false
	pre := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if d >= 'A' && d <= 'Z' {
			if i > 0 && j && pre {
				change = true
				data = append(data, '_')
			}
		} else {
			pre = true
		}

		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:])), change
}

func MapToString(obj map[string]interface{}) string {
	jsonObj, _ := json.Marshal(obj)
	return string(jsonObj)
}

func SliceToString(obj []interface{}) string {
	jsonObj, _ := json.Marshal(obj)
	return string(jsonObj)
}

// 将数组里面的对象key从驼峰LoadAve转化为蛇形load_ave
func SliceToSnakeSlice(obj []interface{}) []interface{} {
	for _, item := range obj {
		item := item.(map[string]interface{})
		for key, value := range item {
			snakeKey, _ := SnakeString(key)
			if snakeKey != key {
				delete(item, key)
				item[snakeKey] = value
			}
		}
	}
	return obj
}

func StringToMap(obj string) interface{} {
	rawByte := []byte(obj)
	var result map[string]interface{}
	json.Unmarshal(rawByte, &result)
	return result
}

func StringToStringArray(obj string) []string {
	var lists []string
	dec := json.NewDecoder(strings.NewReader(obj))
	if err := dec.Decode(&lists); err != nil {
		panic(err)
	}
	return lists
}

func StringToMapArray(obj string) []interface{} {
	var lists []interface{}
	dec := json.NewDecoder(strings.NewReader(obj))
	if err := dec.Decode(&lists); err != nil {
		panic(err)
	}
	return lists
}

func FetchMapUnixTime(obj map[string]interface{}, key string) (i int) {
	if t, exist := obj[key]; exist {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			i = int(t.(float64) / 1000)
		case reflect.Int64:
			i = int(t.(int64) / 1000)
		}

	} else {
		i = 0
	}
	return
}

func FetchMapInterface(obj map[string]interface{}, key string) (s interface{}) {
	if t, exist := obj[key]; exist {
		s = t.(interface{})
	}
	return
}

func FetchMapString(obj map[string]interface{}, key string) (s string) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = strconv.FormatFloat(t.(float64), 'f', 6, 64)
		case reflect.Int64:
			s = strconv.FormatInt(t.(int64), 10)
		case reflect.Int32:
			s = strconv.FormatInt(int64(t.(int32)), 10)
		case reflect.Int:
			s = strconv.Itoa(t.(int))
		default:
			s = t.(string)
		}
	}
	return
}

func MustFetchMapString(obj map[string]interface{}, key string) (s string) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = strconv.FormatFloat(t.(float64), 'f', 6, 64)
		case reflect.Int64:
			s = strconv.FormatInt(t.(int64), 10)
		case reflect.Int32:
			s = strconv.FormatInt(int64(t.(int32)), 10)
		case reflect.Int:
			s = strconv.Itoa(t.(int))
		default:
			s = t.(string)
		}
		return
	} else {
		panic("Key: " + key + " does not exist")
	}
}

func MustFetchMapInt(obj map[string]interface{}, key string) (s int) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = int(t.(float64))
		case reflect.Int64:
			s = int(t.(int64))
		case reflect.Int32:
			s = int(t.(int32))
		case reflect.String:
			if temp, err := strconv.Atoi(t.(string)); err != nil {
				panic(err)
			} else {
				s = temp
			}
		default:
			s = t.(int)
		}
	} else {
		panic("Key: " + key + " does not exist")
	}
	return
}

func FetchMapInt(obj map[string]interface{}, key string) (s int) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = int(t.(float64))
		case reflect.Int64:
			s = int(t.(int64))
		case reflect.Int32:
			s = int(t.(int32))
		case reflect.String:
			if temp, err := strconv.Atoi(t.(string)); err != nil {
				panic(err)
			} else {
				s = temp
			}
		default:
			s = t.(int)
		}
	}
	return
}

func FetchMapInt32(obj map[string]interface{}, key string) (s int32) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = int32(t.(float64))
		case reflect.Int64:
			s = int32(t.(int64))
		case reflect.Int32:
			s = int32(t.(int32))
		case reflect.String:
			if temp, err := strconv.Atoi(t.(string)); err != nil {
				panic(err)
			} else {
				s = int32(temp)
			}
		default:
			s = t.(int32)
		}
	}
	return
}

func FetchMapInt64(obj map[string]interface{}, key string) (s int64) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Int:
			s = int64(t.(int))
		case reflect.Float64:
			s = int64(t.(float64))
		case reflect.Int32:
			s = int64(t.(int32))
		case reflect.String:
			if temp, err := strconv.Atoi(t.(string)); err != nil {
				panic(err)
			} else {
				s = int64(temp)
			}
		default:
			s = t.(int64)
		}
	}
	return
}

func FetchMapFloat64(obj map[string]interface{}, key string) (s float64) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float32:
			s = float64(t.(float32))
		case reflect.Int64:
			s = float64(t.(int64))
		default:
			s = t.(float64)
		}
	}
	return
}

func FetchMapFloat32(obj map[string]interface{}, key string) (s float32) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = float32(t.(float64))
		default:
			s = t.(float32)
		}
	}
	return
}

func FetchMapBool(obj map[string]interface{}, key string) (b bool) {
	if t, exist := obj[key]; exist && t != nil {
		b = t.(bool)
	}
	return
}

func FetchMapSlice(obj map[string]interface{}, key string) (s []interface{}) {
	if t, exist := obj[key]; exist {
		s = t.([]interface{})
	}
	return
}

func MustFetchMapSlice(obj map[string]interface{}, key string) (s []interface{}) {
	if t, exist := obj[key]; exist {
		s = t.([]interface{})
		return
	} else {
		panic("Key: " + key + " does not exist")
	}
}

func FetchMapMap(obj map[string]interface{}, key string) map[string]interface{} {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Ptr:
			temp := ResolvePointValue(t)
			if temp != nil {
				return temp.(map[string]interface{})
			}
		default:
			return t.(map[string]interface{})
		}
	}
	return nil
}

func BoolToInt(obj bool) int {
	if obj == true {
		return 1
	} else {
		return 0
	}
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func CalcHMACMd5(key, content []byte) string {
	hmacs := hmac.New(md5.New, []byte(key))
	hmacs.Write(content)
	return base64.StdEncoding.EncodeToString(hmacs.Sum([]byte(nil)))
}

/*
*
从map结构体中提取所有key，作为slice返回
keyMap 键值映射表，将从map中提取的key进行简单的映射替换
*/
func GetMapKeysToSlice(obj map[string]interface{}, keyMap map[string]string) (slice []string) {
	if obj != nil {
		for k, _ := range obj {
			if keyMap != nil {
				if val, ok := keyMap[k]; ok {
					slice = append(slice, val)
					continue
				}
			}
			slice = append(slice, k)
		}
	}
	return
}

var ip = ""

func GetIPAddr() (ip string) {
	if ip == "" {
		if addrs, err := net.InterfaceAddrs(); err == nil {
			for _, address := range addrs {
				// 检查ip地址判断是否回环地址
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()
						return
					}
				}
			}
		}
		ip = "Cannot fetch ip addr"
	}
	return
}

func FormatBytes(size float64) string {
	units := [6]string{" B", " KB", " MB", " GB", " TB", " PB"}
	index := 0
	for index = 0; size >= 1024; index++ {
		size = size / 1024
	}
	return fmt.Sprintf("%.2f%s", size, units[index])
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func snakeString(s string) (string, bool) {
	data := make([]byte, 0, len(s)*2)
	change := false
	j := false
	pre := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if d >= 'A' && d <= 'Z' {
			if i > 0 && j && pre {
				change = true
				data = append(data, '_')
			}
		} else {
			pre = true
		}

		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:])), change
}

func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func CreateAnyTypeSlice(slice interface{}) ([]interface{}, bool) {
	val, ok := IsSlice(slice)

	if !ok {
		return nil, false
	}

	sliceLen := val.Len()

	out := make([]interface{}, sliceLen)

	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	return out, true
}

func IsSlice(arg interface{}) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)

	if val.Kind() == reflect.Slice {
		ok = true
	}

	return
}
