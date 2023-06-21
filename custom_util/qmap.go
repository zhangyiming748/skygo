package custom_util

import (
	"fmt"
	"reflect"
	"skygo_detection/guardian/src/net/qmap"
	"sort"
	"time"
)

/*
*
时间参数我们有一个统一的请求参数标准
基于请求，返回起止时间
*/
func GetStartEndTimeFromRequestQmap(req qmap.QM) (sTime, eTime time.Time, sStr, eStr string) {
	//根据时间参数（两种可选其一），得到具体时间范围
	endTime := req.DefaultInt("end_time", 0)
	timeSpan := req.DefaultInt("time_span", 1)
	timeInterval := req.DefaultString("time_interval", "day")
	startDate := req.DefaultInt("start_date", 0) //查询起始时间(时间戳)
	endDate := req.DefaultInt("end_date", 0)     //查询结束时间(时间戳)
	if startDate == 0 && endDate == 0 {
		sTime, eTime = CalculateTimeRange(endTime, timeSpan, timeInterval)
		sStr = sTime.Format("2006-01-02")
		eStr = eTime.Format("2006-01-02")
	} else {
		sTime = time.Unix(int64(startDate), 0)
		eTime = time.Unix(int64(endDate), 0)
		//时间戳转为上海时间的字符串Y-m-d形式
		loc, _ := time.LoadLocation("Asia/Shanghai") //上海时区
		sStr = sTime.In(loc).Format("2006-01-02")
		eStr = eTime.In(loc).Format("2006-01-02")
	}
	return
}

func KStringSort(req qmap.QM) []interface{} {
	var keys []string
	var result []interface{}
	for k, _ := range req {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		result = append(result, req[k])
	}
	return result
}

func KSort(req map[interface{}]interface{}) []interface{} {
	var keys []interface{}
	var result []interface{}
	var reqMap = map[interface{}]interface{}{}
	for k, val := range req {
		if reflect.TypeOf(k).Kind() == reflect.Int64 {
			k = int(k.(int64))
		}
		keys = append(keys, k)
		reqMap[k] = val
	}

	if len(keys) > 0 {
		switch reflect.TypeOf(keys[0]).Kind() {
		case reflect.String:
			stringKeys := InterfaceToString(keys)
			sort.Strings(stringKeys)
			for _, k := range stringKeys {
				result = append(result, reqMap[k])
			}
			break
		case reflect.Int:
			stringKeys := InterfaceToInt(keys)
			sort.Ints(stringKeys)
			for _, k := range stringKeys {
				result = append(result, reqMap[k])
			}
			break
		default:
			fmt.Println("unknow")
		}
	}

	return result
}

func ArrayValues(elements map[interface{}]interface{}) []interface{} {
	i, vals := 0, make([]interface{}, len(elements))
	for _, val := range elements {
		vals[i] = val
		i++
	}
	return vals
}

func ArrayReverse(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
