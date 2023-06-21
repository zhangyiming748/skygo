package mongo

import (
	"math"
	"reflect"
	"strings"
)

// 分页方法，根据传递过来的页数，每页数，总数，返回分页的内容 7个页数 前 1，2，3，4，5 后 的格式返回,小于5页返回具体页数
// page 当前页码
// pergage 每页数量
func Paginator(page, pergage int, nums int64) map[string]interface{} {

	var lastpage int //后一页地址
	//根据nums总数，和prepage每页数量 生成分页总数
	totalpages := int(math.Ceil(float64(nums) / float64(pergage))) //page总数
	//有个问题，前端访问不存在的页，返回页数出错，所以注销掉
	//if page > totalpages {
	//	page = totalpages
	//}
	if page <= 0 {
		page = 1
	}
	paginatorMap := make(map[string]interface{})
	paginatorMap["total_pages"] = lastpage
	paginatorMap["current_page"] = page      //当前页数
	paginatorMap["per_page"] = pergage       //每页条数
	paginatorMap["total"] = nums             //总记录数
	paginatorMap["total_pages"] = totalpages //总页数

	return paginatorMap
}

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
