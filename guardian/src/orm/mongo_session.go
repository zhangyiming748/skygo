package orm

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MongoSession struct {
	Session         *mgo.Collection
	Query           *Query
	Match           bson.M
	sort            bson.D
	defaultPageSize int
	defaultSortId   bool                                     //按照id排序是否开启
	transform       *Transformer                             //用于执行转换map[string]interface{}的逻辑
	structToMapFunc func(interface{}) map[string]interface{} //用于把mgo查询结果（某struct）转为map的方法，orm包提供默认方法
	LastMaxId       string                                   //上次查询的最大id(用于优化分页查询效率)
	LastMinId       string                                   //上次查询的最小id(用于优化分页查询效率)
}

func (x *MongoSession) SetDefaultPageSize(i int) {
	x.defaultPageSize = i
}

func (x *MongoSession) ApplyQuery() *MongoSession {
	x.Match = bson.M{}
	x.sort = bson.D{}
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

	//condition
	if x.Query.HasCondition() {
		//fill out the conditions object, we support "and"、"or"

		conditionPtrs := x.Query.GetCondition()

		conditions := []Condition{}
		andConditions := []Condition{}
		orConditions := []Condition{}
		for _, conditionPtr := range conditionPtrs {
			if conditionPtr != nil {
				condition := *conditionPtr
				switch condition.GetType() {
				case CTYPE_AND:
					andConditions = append(andConditions, condition)
				case CTYPE_OR:
					orConditions = append(orConditions, condition)
				case CTYPE_IN:
					andConditions = append(andConditions, condition)
				}
			}
		}
		conditions = append(andConditions, orConditions...)
		//change a condition object to a real ORM’s function's params
		//for example, gorm has functions like 'where'、'or', they both need params
		for _, condition := range conditions {
			operator := x.getOperator(condition.GetOperator())
			if operator == "" {
				continue
			}
			switch operator {
			case MG_OPERATOR_IS_EQUAL:
				x.Match[condition.GetField()] = bson.M{"$eq": condition.GetValue()}
			case MG_OPERATOR_IS_NOT_EQUAL:
				x.Match[condition.GetField()] = bson.M{"$ne": condition.GetValue()}
			case ES_OPERATOR_IS_GREATER_THAN:
				if val, has := x.Match[condition.GetField()]; has {
					val.(bson.M)["$gt"] = condition.GetValue()
				} else {
					x.Match[condition.GetField()] = bson.M{"$gt": condition.GetValue()}
				}
			case MG_OPERATOR_IS_GREATER_THAN_OR_EQUAL:
				if val, has := x.Match[condition.GetField()]; has {
					val.(bson.M)["$gte"] = condition.GetValue()
				} else {
					x.Match[condition.GetField()] = bson.M{"$gte": condition.GetValue()}
				}
			case MG_OPERATOR_IS_LESS_THAN:
				if val, has := x.Match[condition.GetField()]; has {
					val.(bson.M)["$lt"] = condition.GetValue()
				} else {
					x.Match[condition.GetField()] = bson.M{"$lt": condition.GetValue()}
				}
			case MG_OPERATOR_IS_LESS_THAN_OR_EQUAL:
				if val, has := x.Match[condition.GetField()]; has {
					val.(bson.M)["$lte"] = condition.GetValue()
				} else {
					x.Match[condition.GetField()] = bson.M{"$lte": condition.GetValue()}
				}
			case MG_OPERATOR_IS_IN:
				x.Match[condition.GetField()] = bson.M{"$in": condition.GetValue()}
			case MG_OPERATOR_IS_LIKE:
				x.Match[condition.GetField()] = bson.M{"$regex": bson.RegEx{condition.GetValueString(), "."}}
			case MG_OPERATOR_EXISTS:
				x.Match[condition.GetField()] = bson.M{"$exists": condition.GetValue()}
			default:
			}
		}
	}

	//sort
	if x.Query.HasSorter() {
		sorters := x.Query.GetSorter()
		for _, sorter := range sorters {
			switch sorter.GetDirection() {
			case ASCENDING:
				tempSort := bson.DocElem{
					sorter.GetField(),
					1,
				}
				x.sort = append(x.sort, tempSort)
			default:
				tempSort := bson.DocElem{
					sorter.GetField(),
					-1,
				}
				x.sort = append(x.sort, tempSort)
			}
		}
	}
	return x
}

func (x *MongoSession) getOperator(operator int) string {
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
const MG_OPERATOR_IS_EQUAL = "="
const MG_OPERATOR_IS_GREATER_THAN = ">"
const MG_OPERATOR_IS_GREATER_THAN_OR_EQUAL = ">="
const MG_OPERATOR_IS_LESS_THAN = "<"
const MG_OPERATOR_IS_LESS_THAN_OR_EQUAL = "<="
const MG_OPERATOR_IS_IN = "IN"
const MG_OPERATOR_IS_NOT_IN = "NOT IN"
const MG_OPERATOR_IS_LIKE = "LIKE"
const MG_OPERATOR_IS_NOT_LIKE = "NOT LIKE"
const MG_OPERATOR_IS_JSON_CONTAINS = "JSON_CONTAINS"
const MG_OPERATOR_IS_NOT_EQUAL = "<>"
const MG_OPERATOR_IS_IS_NULL = "IS NULL"
const MG_OPERATOR_IS_IS_NOT_NULL = "IS NOT NULL"
const MG_OPERATOR_EXISTS = "EXIST"
const MG_DEFAULT_KEY = "value"

func (x *MongoSession) operatorMap() map[int]string {
	return map[int]string{
		OPERATOR_IS_EQUAL:                 MG_OPERATOR_IS_EQUAL,
		OPERATOR_IS_GREATER_THAN:          MG_OPERATOR_IS_GREATER_THAN,
		OPERATOR_IS_GREATER_THAN_OR_EQUAL: MG_OPERATOR_IS_GREATER_THAN_OR_EQUAL,
		OPERATOR_IS_IN:                    MG_OPERATOR_IS_IN,
		OPERATOR_IS_NOT_IN:                MG_OPERATOR_IS_NOT_IN,
		OPERATOR_IS_LESS_THAN:             MG_OPERATOR_IS_LESS_THAN,
		OPERATOR_IS_LESS_THAN_OR_EQUAL:    MG_OPERATOR_IS_LESS_THAN_OR_EQUAL,
		OPERATOR_IS_LIKE:                  MG_OPERATOR_IS_LIKE,
		OPERATOR_IS_NOT_LIKE:              MG_OPERATOR_IS_NOT_LIKE,
		OPERATOR_IS_JSON_CONTAINS:         MG_OPERATOR_IS_JSON_CONTAINS,
		OPERATOR_IS_NOT_EQUAL:             MG_OPERATOR_IS_NOT_EQUAL,
		OPERATOR_IS_IS_NULL:               MG_OPERATOR_IS_IS_NULL,
		OPERATOR_IS_IS_NOT_NULL:           MG_OPERATOR_IS_IS_NOT_NULL,
		OPERATOR_EXISTS:                   MG_OPERATOR_EXISTS,
	}
}

// -----------------------API for users to use-------------
func (x *MongoSession) SetTransformer(f *Transformer) *MongoSession {
	x.transform = f
	return x
}

func (x *MongoSession) SetPage(page int) *MongoSession {
	x.Query.SetPage(page)
	return x
}

func (x *MongoSession) SetLimit(limit int) *MongoSession {
	x.Query.SetLimit(limit)
	return x
}

// 0升序， 1降序
func (x *MongoSession) AddSorter(field string, direction int) {
	s := Sorter{
		field:     field,
		direction: direction,
	}
	x.Query.AddSorter(s)
}

//------------------------ Func----------------------------

func (x *MongoSession) Pagination(count int) PaginationResult {
	if total, err := x.Session.Find(x.Match).Count(); err == nil {
		pages := Paginator(x.Query.GetPage(), x.Query.GetLimit(), int64(total))
		pages["count"] = count
		return pages
	} else {
		panic(err)
	}
}

// 带count统计上限的分页
// 由于MongoDB在count时，性能较低，部分分页场景下可以给count添加统计上限来提高分页查询效率
func (x *MongoSession) PaginationWithLimit(totalCount, count, countLimit int) PaginationResult {
	if totalCount < 0 {
		if total, err := x.Session.Find(x.Match).Limit(countLimit).Count(); err == nil {
			pages := Paginator(x.Query.GetPage(), x.Query.GetLimit(), int64(total))
			pages["count"] = count
			return pages
		} else {
			panic(err)
		}
	} else {
		pages := Paginator(x.Query.GetPage(), x.Query.GetLimit(), int64(totalCount))
		pages["count"] = count
		return pages
	}
}

func (x *MongoSession) MatchPagination(count, total int) PaginationResult {
	pages := Paginator(x.Query.GetPage(), x.Query.GetLimit(), int64(total))
	pages["count"] = count
	return pages
}

func (x *MongoSession) MTCHAll(operations []bson.M) AllResult {
	x.ApplyQuery()
	if len(x.Match) > 0 {
		operations = append(operations, bson.M{"$match": x.Match})
	}
	if len(x.sort) > 0 {
		operations = append(operations, bson.M{"$sort": x.sort})
	} else {
		operations = append(operations, bson.M{"$sort": bson.D{bson.DocElem{"_id", -1}}})
	}
	operations = append(operations, bson.M{"$skip": x.Query.GetOffset()})
	operations = append(operations, bson.M{"$limit": x.Query.GetLimit()})
	result := AllResult{}

	if err := x.Session.Pipe(operations).All(&result); err == nil {
		//transformer
		if x.transform != nil {
			transformer := *x.transform
			for _, one := range result {
				transformer(one)
			}
		}
		return result
	} else {
		panic(err)
	}
}

func (x *MongoSession) MATCHCOUNT(operations []bson.M) int {
	x.ApplyQuery()
	result := AllResult{}

	if err := x.Session.Pipe(operations).All(&result); err == nil {
		//transformer
		if x.transform != nil {
			transformer := *x.transform
			for _, one := range result {
				transformer(one)
			}
		}
		return len(result)
	} else {
		return 0
	}
}

func (x *MongoSession) All() AllResult {
	x.ApplyQuery()
	operations := []bson.M{}
	if len(x.Match) > 0 {
		operations = append(operations, bson.M{"$match": x.Match})
	}
	if len(x.sort) > 0 {
		operations = append(operations, bson.M{"$sort": x.sort})
	} else {
		operations = append(operations, bson.M{"$sort": bson.D{bson.DocElem{"_id", -1}}})
	}
	operations = append(operations, bson.M{"$skip": x.Query.GetOffset()})
	operations = append(operations, bson.M{"$limit": x.Query.GetLimit()})
	result := AllResult{}

	if err := x.Session.Pipe(operations).All(&result); err == nil {
		//transformer
		if x.transform != nil {
			transformer := *x.transform
			for _, one := range result {
				transformer(one)
			}
		}
		return result
	} else {
		panic(err)
	}
}

func (x *MongoSession) GetOne() (OneResult, error) {
	x.ApplyQuery()
	result := OneResult{}
	if err := x.Session.Find(x.Match).One(&result); err == nil {
		//transformer
		if x.transform != nil {
			transformer := *x.transform
			transformer(result)
		}
		return result, nil
	} else {
		return nil, err
	}
}
