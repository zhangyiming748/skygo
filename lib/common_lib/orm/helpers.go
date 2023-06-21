package orm

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strings"

	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/xorm"
)

// queryStr为url请求的参数部分(?后的部分)， 本函数返回解析后的键值对组成的map
func QueryStrToMap(queryStr string) map[string]interface{} {
	u := url.URL{RawQuery: queryStr}
	params := map[string]interface{}{}
	for k, v := range u.Query() {
		if len(v) != 1 {
			continue
		}
		params[k] = v[0]
	}
	return params
}

func ApplyQuery(q Query, session *xorm.Session) *xorm.Session {
	if q.HasFields() {
		session = session.Cols(q.GetFields()...)
	}

	if q.HasCondition() {
		// fill out the conditions object, we support "and"、"or"
		conditionPtrs := q.GetCondition()
		conditions := []Condition{}
		andConditions := []Condition{}
		orConditions := []Condition{}
		inConditions := []Condition{}
		for _, conditionPtr := range conditionPtrs {
			if conditionPtr != nil {
				condition := *conditionPtr
				switch condition.GetType() {
				case CTYPE_AND:
					andConditions = append(andConditions, condition)
				case CTYPE_OR:
					orConditions = append(orConditions, condition)
				case CTYPE_IN:
					inConditions = append(inConditions, condition)
				}
			}
		}
		conditions = append(andConditions, orConditions...)
		conditions = append(conditions, inConditions...)

		//change a condition object to a real ORM’s function's params
		//for example, gorm has functions like 'where'、'or', they both need params
		for _, condition := range conditions {
			operator := getXormOperator(condition.GetOperator())
			if operator == "" {
				continue
			}

			var param interface{}
			var args interface{}

			switch operator {
			case XormOperatorIsIsNotNull, XormOperatorIsIsNull:
				//db.Model(vehicleMount).Where("vin is not null").Limit(3).Find(&vehicleMount)
				//db.Model(vehicleMount).Where("vin is null").Limit(3).Find(&vehicleMount)
				param = fmt.Sprintf("%s %s", condition.GetField(), operator)
			case XormOperatorIsIn, XormOperatorIsNotIn:
				//db.Where("name in (?)", []string{"jinzhu", "jinzhu 2"}).Find(&users)
				param = fmt.Sprintf("%s", condition.GetField())
				args = condition.GetValue()
			case XormOperatorIsLike, XormOperatorIsNotLike:
				param = fmt.Sprintf("%s %s ?", condition.GetField(), operator)
				args = fmt.Sprintf("%%%s%%", condition.GetValueString())
			default:
				//db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
				//db.Where("email LIKE ?", "%jinzhu%").Delete(Email{})
				param = fmt.Sprintf("%s %s ?", condition.GetField(), operator)
				args = condition.GetValueString()
			}

			switch condition.GetType() {
			case CTYPE_AND:
				if args == nil {
					session = session.Where(param)
				} else {
					session = session.Where(param, args)
				}
			case CTYPE_OR:
				if args == nil {
					session = session.Or(param)
				} else {
					session = session.Or(param, args)
				}
			case CTYPE_IN:
				if args == nil {
					session = session.In(param.(string))
				} else {
					session = session.In(param.(string), args)
				}
			}
		}
	}

	//db.Order("age desc").Order("name").Find(&users)
	if q.HasSorter() {
		sorters := q.GetSorter()
		for _, sorter := range sorters {
			switch sorter.GetDirection() {
			case DESCENDING:
				session = session.Desc(sorter.GetField())
			case ASCENDING:
				session = session.Asc(sorter.GetField())
			default:
				session = session.Asc(sorter.GetField())
			}
		}
	}

	//当我们获取一些特殊字段的时候,会把它拆为其他字段的条件，然后这个字段的条件要废除掉，采用置为nil的方式，因此注释掉了
	//if q.HasJoins() {
	//	for _, join := range q.GetJoins() {
	//		_joinType := "LEFT"
	//		switch join.joinType {
	//		case innerJoin:
	//			_joinType = "INNER"
	//		case leftJoin:
	//			_joinType = "LEFT"
	//		case rightJoin:
	//			_joinType = "RIGHT"
	//		}
	//		session = session.Join(_joinType, join.tableName, join.relation)
	//	}
	//}

	return session
}

func getXormOperator(operator int) string {
	return QueryOperatorToXorm[operator]
}

// 基于TransformerIf接口类型的对象，把它变为TransformerFunc函数类型的值
func GetTransformerFunc(l TransformerIf) *TransformerFunc {
	trans := func(qm qmap.QM) qmap.QM {
		// 集成 ModifyItem()
		qm = l.ModifyItem(qm)

		// 集成 ExcludeFields()
		if fields := l.ExcludeFields(); len(fields) > 0 {
			for _, field := range fields {
				if _, exist := qm[field]; exist {
					delete(qm, field)
				}
			}
		}
		return qm
	}
	t := TransformerFunc(trans)
	return &t
}

// 分页方法
// 根据传递过来的页数，每页数，总数，返回分页的内容 7个页数 前 1，2，3，4，5 后 的格式返回,小于5页返回具体页数
// page 当前页码
// perPage 每页数量
// nums 总数
func Paginator(page, perPage int, nums int64) map[string]interface{} {
	var lastPage int // 后一页地址
	// 根据nums总数，和perPage每页数量 生成分页总数
	totalPage := int(math.Ceil(float64(nums) / float64(perPage))) // page总数
	if page > totalPage {
		page = totalPage
	}
	if page <= 0 {
		page = 1
	}
	paginatorMap := make(map[string]interface{})
	paginatorMap["total_pages"] = lastPage
	paginatorMap["current_page"] = page     // 当前页数
	paginatorMap["per_page"] = perPage      // 每页条数
	paginatorMap["total"] = nums            // 总记录数
	paginatorMap["total_pages"] = totalPage // 总页数

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
