package custom_util

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

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

func StringToIntArray(obj string) []int {
	var lists []int
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

func IntToBool(i int) bool {
	if i == 0 {
		return false
	} else {
		return true
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

// 获取指定精度的浮点数
func Float64Round(val float64, p int) (s float64) {
	format := "%." + fmt.Sprintf("%d", p) + "f"
	s, _ = strconv.ParseFloat(fmt.Sprintf(format, val), p)
	return
}

// 流量转化(字节转为最大显示单位)
func NetFlowUnitConverter(total int) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB"}
	unitIndex := 0
	totalFloat := float64(total)
	for totalFloat > 1024 {
		unitIndex++
		totalFloat = totalFloat / 1024
	}
	return fmt.Sprintf("%.2f%s", totalFloat, units[unitIndex])
}

func ToInt(i interface{}) (int, error) {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Int:
		return i.(int), nil
	case reflect.Float64:
		return int(i.(float64)), nil
	case reflect.String:
		i, err := strconv.Atoi(i.(string))
		if err == nil {
			return i, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("当前值%s无法转换为int类型", fmt.Sprint(i)))
}

// 判断字符串是否包含特殊字符,如果包含在特殊字符前方增加反斜杠
func SpeCharsAddBackslash(chars string) string {
	htmlSpe := `-+*.?/\$%^()<>`
	speStr := ``
	for _, v := range chars {
		//判断当前run 是否是特殊字符
		if strings.ContainsRune(htmlSpe, v) {
			//存在特殊字符
			speStr += `\` + string(v)
		} else {
			speStr += string(v)
		}
	}
	return speStr
}
