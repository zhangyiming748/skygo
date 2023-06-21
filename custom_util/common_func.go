package custom_util

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"skygo_detection/common"
	"skygo_detection/custom_error"
	"skygo_detection/guardian/src/net/qmap"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/thinkeridea/go-extend/exnet"
)

// GetTopicLogId 获取Topic log_id
func GetTopicLogId(topicName string) string {
	u1 := uuid.NewV4()
	return topicName + "_" + u1.String()
}

// GetHostName 获取主机名
func GetHostName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	return hostname, nil
}

// GetPageOffsetLimit ...
func GetPageOffsetLimit(req *qmap.QM) (offset int, limit int, err error) {
	err = nil
	page := req.DefaultInt("page", 1)
	limit = req.DefaultInt("limit", common.DefaultLimitNum)

	if page < 1 {
		err = custom_error.ErrPageParam
		return
	}

	if limit < 1 {
		err = custom_error.ErrPageLimitParam
		return
	}

	if limit > common.MaxLimitNum {
		err = custom_error.ErrPageLimitGtMax
		return
	}
	offset = (page - 1) * limit
	return
}

// GetLimit ...
func GetLimit(req *qmap.QM, params ...int) (limit int, err error) {
	if len(params) > 0 {
		limit = req.DefaultInt("limit", params[0])
	} else {
		limit = req.DefaultInt("limit", common.DefaultLimitNum)
	}
	if limit < 1 {
		err = custom_error.ErrPageLimitParam
		return
	}

	if limit > common.MaxLimitNum {
		err = custom_error.ErrPageLimitGtMax
		return
	}
	return
}

// CheckPortNo is 核对端口号格式
func CheckPortNo(port string) (int, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return -1, err
	}
	if portInt < 0 || portInt > 65535 {
		return -1, errors.New("端口超出返回")
	}

	return portInt, nil
}

// RemoveRepByLoop 通过两重循环过滤重复元素
func RemoveRepByLoop(slc []int) []int {
	result := []int{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// RemoveRepByMap 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// RemoveRep 元素去重
func RemoveRep(slc []int) []int {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return RemoveRepByLoop(slc)
	}

	// 大于的时候，通过map来过滤
	return RemoveRepByMap(slc)
}

// Create 创建文件句柄
func Create(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)

	_, err := os.Stat(dir)

	// 目录不存在时, 新建目录
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
	}

	return os.Create(filePath)
}

// Mkdir 创建目录
func Mkdir(dirName string) error {
	_, err := os.Stat(dirName)

	// 目录不存在时, 新建目录
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, 0777); err != nil {
			return err
		}
	}

	return nil
}

// 私有ip网址范围
var privateIPRanges = [][2]uint{
	{167772160, 184549375},   // A类：10.0.0.0-10.255.255.255
	{2886729728, 2887778303}, // B类：172.16.0.0-172.31.255.255
	{3232235520, 3232301055}, // C类：192.168.0.0-192.168.255.255
}

// IsPrivateIP 核对ip地址是否是私有ip
func IsPrivateIP(ipStr string) bool {
	ipNum, _ := exnet.IPString2Long(ipStr)
	for _, tmpRange := range privateIPRanges {
		if tmpRange[0] <= ipNum && ipNum <= tmpRange[1] {
			return true
		}
	}

	return false
}

// Sorter 排序结构
type Sorter struct {
	Field     string
	Direction int
}

// GetSorter 获取排序参数，支持单个参数与多个参数  {"id":1} 或 [{"id":1},{"category_id":-1}]
func GetSorter(req *qmap.QM, abledFields map[string]bool, arg ...string) (sliceSorter []Sorter, err error) {
	key := "sort"
	if len(arg) > 0 {
		key = arg[0]
	}
	var sort interface{}
	tmpSort := req.Interface(key)
	if sortStr, ok := tmpSort.(string); ok {
		json.Unmarshal([]byte(sortStr), &sort)
	} else {
		sort = tmpSort
	}

	// 单字段排序
	m, mFlag := sort.(map[string]interface{})
	if mFlag {
		if sorter, has := GetOneSorter(m, abledFields); has {
			sliceSorter = append(sliceSorter, sorter)
		}
		return
	}

	// 多字段排序
	s, sFlag := sort.([]interface{})
	if sFlag {
		for _, s := range s {
			if sorter, has := GetOneSorter(s, abledFields); has {
				sliceSorter = append(sliceSorter, sorter)
			}
		}
		return
	}

	return
}

// GetOneSorter ...
func GetOneSorter(sort interface{}, abledFields map[string]bool) (sorter Sorter, has bool) {
	has = false

	m, ok := sort.(map[string]interface{})
	if ok {
		for k, v := range m {
			if _, abled := abledFields[k]; !abled {
				continue
			}

			if _v, ok := v.(float64); ok {
				sorter = Sorter{
					Field:     k,
					Direction: int(_v),
				}
				has = true
				break
			}
		}
	}

	return
}

// GetSortVal ...
func GetSortVal(sortInt int) string {
	if sortInt == 1 {
		return "asc"
	}
	return "desc"
}

// Sha256 ...
func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	b := h.Sum(nil)

	return fmt.Sprintf("%x", b)
}

// StringSliceToIntSlice 字符串数组转化为Int数组
func StringSliceToIntSlice(strs []string) ([]int, error) {
	result := make([]int, len(strs))
	fmt.Println("strs:", strs)
	for i, str := range strs {
		intValue, err := strconv.Atoi(str)
		if err == nil {
			result[i] = intValue
		} else {
			fmt.Println(err)
			return nil, err
		}
	}
	return result, nil
}

// Days 计算两个时间戳之间的天数
func Days(timestampFrom, timestampTo int64) int {
	var midnightUnix = func(t time.Time) int64 {
		y, m, d := t.Date()
		return time.Date(y, m, d+1, 0, 0, 0, 0, time.Local).Unix()
	}
	var days = 0
	for {
		if midnightUnix(time.Unix(timestampFrom, 0).AddDate(0, 0, days)) >= timestampTo {
			days++
			break
		}
		days++
	}
	return days
}

// GetTimeStr ...
func GetTimeStr(timeInt int64) string {
	if timeInt == 0 {
		return ""
	}
	return time.Unix(timeInt, 0).Format("2006-01-02 15:04")
}

// GetTimeLoc ...
func GetTimeLoc() *time.Location {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return loc
}

// DateFormatYmdHis 年月日时分秒
func DateFormatYmdHis(timeUnix int64) string {
	return DateFormat(timeUnix, "2006-01-02 15:04:05")
}

// DateFormatYmdHi 年月日时分
func DateFormatYmdHi(timeUnix int64) string {
	return DateFormat(timeUnix, "2006-01-02 15:04")
}

// DateFormatYmd 年月日
func DateFormatYmd(timeUnix int64) string {
	return DateFormat(timeUnix, "2006-01-02")
}

func DateYmdTimestamp(timeUnix int64) int64 {
	t, _ := time.ParseInLocation("2006-01-02", DateFormatYmd(timeUnix), GetTimeLoc())
	return t.Unix()
}

// DateFormat ...
func DateFormat(timeUnix int64, formatStr string) string {
	tm := time.Unix(timeUnix, 0).In(GetTimeLoc())
	// tm.Format("2006-01-02 15:04:05")
	return tm.Format(formatStr)
}

// GetNowTimeUnix 获取当前时间戳
func GetNowTimeUnix() int64 {
	return time.Now().In(GetTimeLoc()).Unix()
}

// GetNowTimeUnix13 获取毫秒
func GetNowTimeUnix13() int64 {
	return (time.Now().In(GetTimeLoc()).UnixNano() / 1e6)
}

// GetTimeUnix13ByForwardDay 获取几天前的对应毫秒
func GetTimeUnix13ByForwardDay(day int) int64 {
	return (time.Now().In(GetTimeLoc()).Unix() - int64(day*3600*24)) * 1000
}

// TimeStrToUnix 时间字符串转时间戳
func TimeStrToUnix(timeStr string) (int, error) {
	tm := 0
	tmp, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, GetTimeLoc())
	if err != nil {
		return tm, err
	}
	tm = int(tmp.Unix())

	return tm, nil
}

// TimeYmdStrToUnix 时间字符串转时间戳
func TimeYmdStrToUnix(timeStr string) (int, error) {
	tm := 0
	tmp, err := time.ParseInLocation("2006-01-02", timeStr, GetTimeLoc())
	if err != nil {
		return tm, err
	}
	tm = int(tmp.Unix())

	return tm, nil
}

// TimeStrToUnix13 时间字符串转时间戳毫秒
func TimeStrToUnix13(timeStr string) (int64, error) {
	i, err := TimeStrToUnix(timeStr)
	if err != nil {
		return 0, err
	}

	return int64(i * 1000), nil
}

// InIntSlice ...
func InIntSlice(val int, slice []int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

// ErrRecover ...
func ErrRecover() {
	// 发生宕机时，获取panic传递的上下文并打印
	err := recover()
	if err == nil {
		return
	}
	switch err.(type) {
	case runtime.Error: // 运行时错误
		fmt.Println("runtime error:", err)
	default: // 非运行时错误
		fmt.Println("error:", err)
	}
}

// StructToMap 结构体转值为map
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return
	}

	if err = json.Unmarshal(jsonBytes, &newMap); err != nil {
		return
	}

	return
}

// StructToMap 结构体转值为list
func StructToList(obj interface{}) (newMap []map[string]interface{}, err error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return
	}

	if err = json.Unmarshal(jsonBytes, &newMap); err != nil {
		return
	}

	return
}

func JsonStrToStruct(jsonStr string, res interface{}) (err error) {
	jsonBytes := []byte(jsonStr)

	if err = json.Unmarshal(jsonBytes, res); err != nil {
		return
	}

	return
}

func InterfaceToJsonStr(data interface{}) (string, error) {
	tmpBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(tmpBytes), nil
}

func CheckKeyExistMap(val int, mapInfo map[int]string) bool {
	_, ok := mapInfo[val]
	return ok
}

func JsonInterfaceToStruct(data interface{}, newData interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, newData)
	if err != nil {
		return err
	}

	return nil
}

/*
*
结构体对象转map(key由驼峰转蛇形)
*/
func StructToMap2(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		snakeName, _ := SnakeString(obj1.Field(i).Name)
		data[snakeName] = obj2.Field(i).Interface()
	}
	return data
}

// Sep 23, 2021 06:42:47.778473844 UTC
// Sep 18, 2021 17:45:16.460053843 CST
// 转换为北京时间
func TimeTransLocal(strTime string) string {
	t, _ := time.Parse("Jan 2, 2006 15:04:05.000000000 MST", strTime)
	if strings.Contains(strTime, "UTC") {
		t = t.Add(time.Duration(8) * time.Hour)
	} else if strings.Contains(strTime, "CST") {
		t = t.Add(time.Duration(14) * time.Hour)
	}
	return t.Format("2006-01-02 15:04:05")
}
