package orm

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"reflect"
	"strings"
)

type ExtraParams = map[string]interface{}

const DefaultPageSize = 20

type XormSession struct {
	Session *xorm.Session
	Query   *Query
	apply   bool //标识是否执行了 ApplyQuery操作，即基于Query去操作Session

	defaultPageSize int
	defaultSortId   bool                                     //按照id排序是否开启
	transform       *Transformer                             //用于执行转换map[string]interface{}的逻辑
	structToMapFunc func(interface{}) map[string]interface{} //用于把xorm查询结果（某struct）转为map的方法，orm包提供默认方法
	excludeFields   []string
}

func (x *XormSession) SetDefaultPageSize(i int) {
	x.defaultPageSize = i
}

func (x *XormSession) ApplyQuery() {
	//if x.apply == true {
	//	return
	//}
	//x.apply = true

	//------- 对offset、 limit、page三个参数进行处理，得到最终的limit和offset --------
	//limit必须有，否则取默认值
	if x.Query.HasLimit() == false {
		if x.defaultPageSize > 0 {
			x.Query.SetLimit(x.defaultPageSize)
		} else {
			x.Query.SetLimit(DefaultPageSize)
		}
	}
	//page可以有，有的话，通过它计算额外offset
	if x.Query.HasPage() {
		offset := x.Query.GetLimit() * (x.Query.GetPage() - 1)
		if x.Query.HasOffset() {
			x.Query.SetOffset(x.Query.GetOffset() + offset)
		} else {
			x.Query.SetOffset(offset)
		}
	}
	//limit 和 offset并不在这里应用，而是交给all、one等方法
	//x.Session = x.Session.Limit(x.Query.GetLimit(), x.Query.GetOffset())

	if x.Query.HasFields() {
		x.Session = x.Session.Cols(x.Query.GetFields()...)
	}

	if x.Query.HasCondition() {
		//fill out the conditions object, we support "and"、"or"
		conditionPtrs := x.Query.GetCondition()
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
		conditions = append(andConditions, inConditions...)
		conditions = append(conditions, orConditions...)

		//change a condition object to a real ORM’s function's params
		//for example, gorm has functions like 'where'、'or', they both need params
		for _, condition := range conditions {
			operator := x.getOperator(condition.GetOperator())
			if operator == "" {
				continue
			}

			var param interface{}
			var args interface{}

			switch operator {
			case XORM_OPERATOR_IS_IS_NOT_NULL, XORM_OPERATOR_IS_IS_NULL:
				//db.Model(vehicleMount).Where("vin is not null").Limit(3).Find(&vehicleMount)
				//db.Model(vehicleMount).Where("vin is null").Limit(3).Find(&vehicleMount)
				param = fmt.Sprintf("%s %s", condition.GetField(), operator)
			case XORM_OPERATOR_IS_IN, XORM_OPERATOR_IS_NOT_IN:
				//db.Where("name in (?)", []string{"jinzhu", "jinzhu 2"}).Find(&users)
				param = fmt.Sprintf("%s", condition.GetField())
				args = condition.GetValue()
			case XORM_OPERATOR_IS_LIKE:
				param = fmt.Sprintf("%s LIKE ?", condition.GetField())
				args = "%" + condition.GetValue().(string) + "%"
			default:
				//db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
				//db.Where("email LIKE ?", "%jinzhu%").Delete(Email{})
				param = fmt.Sprintf("%s %s ?", condition.GetField(), operator)
				args = condition.GetValueString()
			}

			switch condition.GetType() {
			case CTYPE_AND:
				if args == nil {
					x.Session = x.Session.Where(param)
				} else {
					x.Session = x.Session.Where(param, args)
				}
			case CTYPE_OR:
				if args == nil {
					x.Session = x.Session.Or(param)
				} else {
					x.Session = x.Session.Or(param, args)
				}
			case CTYPE_IN:
				if args == nil {
					x.Session = x.Session.In(param.(string))
				} else {
					x.Session = x.Session.In(param.(string), args)
				}
			}
		}
	}

	//db.Order("age desc").Order("name").Find(&users)
	if x.Query.HasSorter() {
		sorters := x.Query.GetSorter()
		for _, sorter := range sorters {
			switch sorter.GetDirection() {
			case DESCENDING:
				x.Session = x.Session.Desc(sorter.GetField())
			case ASCENDING:
				x.Session = x.Session.Asc(sorter.GetField())
			default:
				x.Session = x.Session.Asc(sorter.GetField())
			}
		}
	}

	//当我们获取一些特殊字段的时候,会把它拆为其他字段的条件，然后这个字段的条件要废除掉，采用置为nil的方式，因此注释掉了
	//if x.Query.HasJoins() {
	//	for _, join := range x.Query.GetJoins() {
	//		_joinType := "LEFT"
	//		switch join.joinType {
	//		case innerJoin:
	//			_joinType = "INNER"
	//		case leftJoin:
	//			_joinType = "LEFT"
	//		case rightJoin:
	//			_joinType = "RIGHT"
	//		}
	//		x.Session = x.Session.Join(_joinType, join.tableName, join.relation)
	//	}
	//}
}

func (x *XormSession) getOperator(operator int) string {
	m := x.operatorMap()
	return m[operator]
}

/*
This is a xorm parser.
When a user's request is analysed and turn to be a Query object.Any parser program can parse it to a DB object.

Note:
	Query 's include filed is different, it will not parse by current parser, it is usd by transformer program.
*/
//
const XORM_OPERATOR_IS_EQUAL = "="
const XORM_OPERATOR_IS_GREATER_THAN = ">"
const XORM_OPERATOR_IS_GREATER_THAN_OR_EQUAL = ">="
const XORM_OPERATOR_IS_LESS_THAN = "<"
const XORM_OPERATOR_IS_LESS_THAN_OR_EQUAL = "<="
const XORM_OPERATOR_IS_IN = "IN"
const XORM_OPERATOR_IS_NOT_IN = "NOT IN"
const XORM_OPERATOR_IS_LIKE = "LIKE"
const XORM_OPERATOR_IS_NOT_LIKE = "NOT LIKE"
const XORM_OPERATOR_IS_JSON_CONTAINS = "JSON_CONTAINS"
const XORM_OPERATOR_IS_NOT_EQUAL = "<>"
const XORM_OPERATOR_IS_IS_NULL = "IS NULL"
const XORM_OPERATOR_IS_IS_NOT_NULL = "IS NOT NULL"

//const XORM_DEFAULT_KEY = "value"

func (x *XormSession) operatorMap() map[int]string {
	return map[int]string{
		OPERATOR_IS_EQUAL:                 XORM_OPERATOR_IS_EQUAL,
		OPERATOR_IS_GREATER_THAN:          XORM_OPERATOR_IS_GREATER_THAN,
		OPERATOR_IS_GREATER_THAN_OR_EQUAL: XORM_OPERATOR_IS_GREATER_THAN_OR_EQUAL,
		OPERATOR_IS_IN:                    XORM_OPERATOR_IS_IN,
		OPERATOR_IS_NOT_IN:                XORM_OPERATOR_IS_NOT_IN,
		OPERATOR_IS_LESS_THAN:             XORM_OPERATOR_IS_LESS_THAN,
		OPERATOR_IS_LESS_THAN_OR_EQUAL:    XORM_OPERATOR_IS_LESS_THAN_OR_EQUAL,
		OPERATOR_IS_LIKE:                  XORM_OPERATOR_IS_LIKE,
		OPERATOR_IS_NOT_LIKE:              XORM_OPERATOR_IS_NOT_LIKE,
		OPERATOR_IS_JSON_CONTAINS:         XORM_OPERATOR_IS_JSON_CONTAINS,
		OPERATOR_IS_NOT_EQUAL:             XORM_OPERATOR_IS_NOT_EQUAL,
		OPERATOR_IS_IS_NULL:               XORM_OPERATOR_IS_IS_NULL,
		OPERATOR_IS_IS_NOT_NULL:           XORM_OPERATOR_IS_IS_NOT_NULL,
	}
}

// -----------------------API for users to use-------------
func (x *XormSession) SetTransformer(f *Transformer) *XormSession {
	x.transform = f
	return x
}

func (x *XormSession) SetExcludeFields(f []string) *XormSession {
	x.excludeFields = f
	return x
}

func (x *XormSession) AndWhereEqual(params map[string]interface{}) *XormSession {
	for k, v := range params {
		x.Query.AndWhereEqual(k, v)
	}
	return x
}

func (x *XormSession) SetPage(page int) *XormSession {
	x.Query.SetPage(page)
	return x
}

func (x *XormSession) SetLimit(limit int) *XormSession {
	x.Query.SetLimit(limit)
	return x
}

// 0升序， 1降序
func (x *XormSession) AddSorter(field string, direction int) {
	s := Sorter{
		field:     field,
		direction: direction,
	}
	x.Query.AddSorter(s)
}

//------------------------ Func----------------------------

func (x *XormSession) Pagination(modelPtr interface{}, count int) PaginationResult {
	//totalNumber
	var totalNum int64
	totalNum, err := x.Session.Count(modelPtr)
	if err != nil {
		panic(err)
	}
	pages := Paginator(x.Query.GetPage(), x.Query.GetLimit(), int64(totalNum))
	pages["count"] = count
	return pages
}

func (x *XormSession) All(modelsPtr interface{}) (AllResult, *xorm.Session) {
	if !(reflect.TypeOf(modelsPtr).Kind() == reflect.Ptr &&
		reflect.ValueOf(modelsPtr).Elem().Type().Kind() == reflect.Slice &&
		reflect.ValueOf(modelsPtr).Elem().Type().Elem().Kind() == reflect.Struct) {
		panic("Xorm的SessionFind参数格式不对")
	}

	//默认排序
	//except Query's sorter, we use the primary key to sort
	//1、get the primary key name, default to ""
	fieldName := ""
	t := reflect.ValueOf(modelsPtr).Elem().Type() //get the fields number of a struct
	//如果model有方法TableName，据此获取表名
	modelName := ""
	if method, has := reflect.ValueOf(modelsPtr).Elem().Type().Elem().MethodByName("TableName"); has {
		param := make([]reflect.Value, 1)
		param[0] = reflect.New(reflect.ValueOf(modelsPtr).Elem().Type().Elem()).Elem()
		modelName = method.Func.Call(param)[0].String()
	}
	filedNum := t.Elem().NumField()
	for i := 0; i < filedNum; i++ {
		fieldTag := t.Elem().Field(i).Tag.Get("xorm")
		if strings.Contains(fieldTag, "pk") {
			fieldName, _ = SnakeString(t.Elem().Field(i).Name)
		}
	}
	if fieldName != "" {
		sorters := x.Query.GetSorter()
		sortersNum := len(sorters)
		for _, sorter := range sorters {
			if sorter.GetField() != fieldName {
				sortersNum--
			}
		}
		if sortersNum == 0 {
			if modelName != "" {
				x.AddSorter(modelName+"."+fieldName, 1)
			} else {
				x.AddSorter(fieldName, 1)
			}
		}
	}

	//apply
	x.ApplyQuery()

	//分页
	x.Session = x.Session.Limit(x.Query.GetLimit(), x.Query.GetOffset())

	//在查询之前备份一份session对象，作为函数返回值，后续可以用于查询分页总数
	copySession := *x.Session

	err := x.Session.Find(modelsPtr)
	if err != nil {
		panic(err)
	}

	//struct to map
	all := AllResult{}
	for i := 0; i < reflect.ValueOf(modelsPtr).Elem().Len(); i++ {
		//reflect.ValueOf(modelsPtr).Elem().Index(i).Type().Kind() is reflect.Struct
		one := reflect.ValueOf(modelsPtr).Elem().Index(i).Interface()
		all = append(all, x.structToMapFunc(one))
	}

	for key, one := range all {
		all[key] = x.transformItem(one)
	}

	return all, &copySession
}

func (x *XormSession) transformItem(item map[string]interface{}) map[string]interface{} {
	//transformer
	if x.transform != nil {
		item = (*x.transform)(item)
	}
	//excludeFields
	if len(x.excludeFields) > 0 {
		for _, field := range x.excludeFields {
			if _, exist := item[field]; exist {
				delete(item, field)
			}
		}
	}

	return item
}

func (x *XormSession) One(modelPtr interface{}) (has bool, one OneResult) {
	if !(reflect.TypeOf(modelPtr).Kind() == reflect.Ptr &&
		reflect.ValueOf(modelPtr).Elem().Type().Kind() == reflect.Struct) {
		panic("Xorm的SessionFind参数格式不对")
	}

	x.ApplyQuery()
	has, err := x.Session.Get(modelPtr)
	if err != nil {
		panic(err)
	}
	if has == false {
		return
	}

	//struct to map
	one = OneResult{}
	one = x.structToMapFunc(reflect.ValueOf(modelPtr).Elem().Interface()) //reflect.ValueOf(modelsPtr).Elem().Type().Kind() is reflect.Struct
	one = x.transformItem(one)

	return
}
