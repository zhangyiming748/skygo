package orm

import (
	"reflect"
	"strconv"
)

/**
Query can be viewed as a container, it store a request's params in it's standard format.
The format is very simple, it is a struct that has series of filed below:
		Conditions []Condition
		offset int
		limit int
		fields []string
		Sorters []Sorter
		includes []string

Parsers like gorm_parser or xorm_parser can parse a Query Object to it's db object.  So, the Query object is
a standard format, every different ORM parser can get their own DB object base on the same Query.
*/

const OPERATOR_IS_EQUAL = 0

const OPERATOR_IS_GREATER_THAN = 1

const OPERATOR_IS_GREATER_THAN_OR_EQUAL = 2

const OPERATOR_IS_IN = 3

const OPERATOR_IS_NOT_IN = 4

const OPERATOR_IS_LESS_THAN = 5

const OPERATOR_IS_LESS_THAN_OR_EQUAL = 6

const OPERATOR_IS_LIKE = 7

const OPERATOR_IS_NOT_LIKE = 8

const OPERATOR_IS_JSON_CONTAINS = 9

const OPERATOR_IS_NOT_EQUAL = 10

// const OPERATOR_CONTAINS                 = 11;
// const OPERATOR_NOT_CONTAINS             = 12;
const OPERATOR_IS_IS_NULL = 13

const OPERATOR_IS_IS_NOT_NULL = 14

const OPERATOR_EXISTS = 15

// --------------Find Conditions-------------//
const CTYPE_AND = 0
const CTYPE_OR = 1
const CTYPE_IN = 3

type Condition struct {
	ctype    int
	field    string
	operator int
	value    interface{}
}

func (c *Condition) GetType() int {
	return c.ctype
}

func (c *Condition) GetOperator() int {
	return c.operator
}

func (c *Condition) GetField() string {
	return c.field
}

func (c *Condition) GetValue() interface{} {
	return c.value
}

func (c *Condition) GetValueString() string {
	switch reflect.TypeOf(c.value).Kind() {
	case reflect.String:
		return c.value.(string)
	case reflect.Int:
		return strconv.Itoa(c.value.(int))
	case reflect.Int8:
		return strconv.Itoa(int(c.value.(int8)))
	case reflect.Int16:
		return strconv.Itoa(int(c.value.(int16)))
	case reflect.Int32:
		return strconv.Itoa(int(c.value.(int32)))
	case reflect.Int64:
		return strconv.Itoa(int(c.value.(int64)))
	case reflect.Float32:
		_c := float64(c.value.(float32))
		return strconv.FormatFloat(_c, 'f', -1, 64)
		//url上传的数据，如果经过json处理，数字类型的会被自动转为float64，我们要获取小数位数
	case reflect.Float64:
		_c := c.value.(float64)
		return strconv.FormatFloat(_c, 'f', -1, 64)
	default:
		panic(reflect.TypeOf(c.value).Kind().String() + " condition value can not turn to string")
	}
	return ""
}

// --------------Sorter-------------//
const ASCENDING = 0
const DESCENDING = 1

type Sorter struct {
	field     string
	direction int
}

func (s *Sorter) GetDirection() int {
	return s.direction
}

func (s *Sorter) GetField() string {
	return s.field
}

//当我们获取一些特殊字段的时候,会把它拆为其他字段的条件，然后这个字段的条件要废除掉，采用置为nil的方式，因此注释掉了
//--------------Join-------------//
//type Join struct {
//	joinType JoinType
//	tableName string
//	relation  string
//}
//type JoinType int8
//const (
//	innerJoin JoinType = iota
//	leftJoin
//	rightJoin
//)

// --------------Query-------------//
type Query struct {
	Conditions []*Condition
	offset     int
	limit      int
	fields     []string
	Sorters    []Sorter
	includes   []string
	page       int
	//joins      []Join
}

func (q *Query) Merge(newQuery *Query) {
	q.AddCondition(newQuery.Conditions...)
	q.SetOffset(newQuery.GetOffset())
	q.SetLimit(newQuery.GetLimit())
	q.AddField(newQuery.GetFields())
	q.AddSorter(newQuery.GetSorter()...)
	q.SetPage(newQuery.GetPage())

}

func (q *Query) AddField(fields []string) {
	q.fields = append(q.fields, fields...)
}

func (q *Query) HasFields() bool {
	if len(q.fields) > 0 {
		return true
	}
	return false
}

func (q *Query) GetFields() []string {
	return q.fields
}

func (q *Query) AddCondition(conditions ...*Condition) {
	q.Conditions = append(q.Conditions, conditions...)
}

func (q *Query) HasCondition() bool {
	if len(q.Conditions) > 0 {
		return true
	}
	return false
}

func (q *Query) GetCondition() []*Condition {
	return q.Conditions
}

// 当我们获取一些特殊字段的时候,会把它拆为其他字段的条件，然后这个字段的条件要废除掉，采用置为nil的方式，因此注释掉了
func (q *Query) UnsetCondition(key int) {
	if key < len(q.Conditions) {
		q.Conditions[key] = nil
	}
}

func (q *Query) AddSorter(sorter ...Sorter) {
	q.Sorters = append(q.Sorters, sorter...)
}

func (q *Query) HasSorter() bool {
	if len(q.Sorters) > 0 {
		return true
	}
	return false
}

func (q *Query) GetSorter() []Sorter {
	return q.Sorters
}

func (q *Query) SetOffset(offset int) {
	q.offset = offset
}

func (q *Query) HasOffset() bool {
	if q.offset > 0 {
		return true
	}
	return false
}

func (q *Query) GetOffset() int {
	return q.offset
}

func (q *Query) SetLimit(limit int) {
	q.limit = limit
}

func (q *Query) HasLimit() bool {
	if q.limit > 0 {
		return true
	}
	return false
}

func (q *Query) GetLimit() int {
	return q.limit
}

func (q *Query) AddInclude(include []string) {
	q.includes = append(q.includes, include...)
}

func (q *Query) HasInclude() bool {
	if len(q.includes) > 0 {
		return true
	}
	return false
}

func (q *Query) GetInclude() []string {
	return q.includes
}

func (q *Query) SetPage(page int) {
	q.page = page
}

func (q *Query) HasPage() bool {
	if q.page > 0 {
		return true
	}
	return false
}

func (q *Query) GetPage() int {
	return q.page
}

//如下代码，为了扩展功能，让可以支持连表查询，但是为了减少复杂度，建议直接使用xorm的session.Join()去连表，因此注释掉了
//func (q *Query) HasJoins() bool {
//	if len(q.joins) > 0 {
//		return true
//	}
//	return false
//}
//func (q *Query) GetJoins() []Join {
//	return q.joins
//}
//func (q *Query) AddJoins(joins ...Join) {
//	q.joins = append(q.joins, joins...)
//}

// -----------------------API for users to use-------------
func (q *Query) Desc(name string) {
	q.AddSorter(Sorter{
		field:     name,
		direction: DESCENDING,
	})
}

func (q *Query) Asc(name string) {
	q.AddSorter(Sorter{
		field:     name,
		direction: ASCENDING,
	})
}

func (q *Query) AndWhereEqual(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_AND,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_EQUAL,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) AndWhereLike(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_AND,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_LIKE,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) AndWhereGt(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_AND,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_GREATER_THAN,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) AndWhereGte(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_AND,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_GREATER_THAN_OR_EQUAL,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) AndWhereLt(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_AND,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_LESS_THAN,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) AndWhereLte(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_AND,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_LESS_THAN_OR_EQUAL,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) AndWhereNotEqual(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_AND,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_NOT_EQUAL,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) AndWhereIn(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_IN,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_IN,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) OrLike(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_OR,
		field:    field,
		value:    value,
		operator: OPERATOR_IS_LIKE,
	}
	q.AddCondition(c)
	return q
}

func (q *Query) AndExists(field string, value interface{}) *Query {
	c := &Condition{
		ctype:    CTYPE_AND,
		field:    field,
		value:    value,
		operator: OPERATOR_EXISTS,
	}
	q.AddCondition(c)
	return q
}

//如下代码，为了扩展功能，让可以支持连表查询，但是为了减少复杂度，建议直接使用xorm的session.Join()去连表，因此注释掉了
//example： ormSession.GetQuery().SetLeftJoin("vehicle_mount_auth_list", "vehicle_mount_auth_list.sn = vehicle_mount.sn")
//func (q *Query) SetLeftJoin(tableName, relation string) {
//	q.AddJoins(Join{
//		leftJoin,
//		tableName,
//		relation,
//	})
//}
//func (q *Query) SetInnerJoin(tableName, relation string) {
//	q.AddJoins(Join{
//		innerJoin,
//		tableName,
//		relation,
//	})
//}
//func (q *Query) SetRightJoin(tableName, relation string) {
//	q.AddJoins(Join{
//		rightJoin,
//		tableName,
//		relation,
//	})
//}
