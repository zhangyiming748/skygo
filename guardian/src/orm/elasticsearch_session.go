package orm

import (
	"context"

	"gopkg.in/olivere/elastic.v5"

	"skygo_detection/guardian/src/net/qmap"
)

type ElasticSearchSession struct {
	*elastic.Client
	OrmQuery    *Query
	EsBoolQuery *elastic.BoolQuery
	EsFieldSort []*elastic.FieldSort
	apply       bool // 标识是否执行了 ApplyQuery操作，即基于Query去操作Session

	defaultPageSize int
	defaultSortId   bool                                     // 按照id排序是否开启
	transform       *Transformer                             // 用于执行转换map[string]interface{}的逻辑
	structToMapFunc func(interface{}) map[string]interface{} // 用于把xorm查询结果（某struct）转为map的方法，orm包提供默认方法
}

func (x *ElasticSearchSession) SetDefaultPageSize(i int) {
	x.defaultPageSize = i
}

func (x *ElasticSearchSession) ApplyQuery() *ElasticSearchSession {
	if x.apply == true {
		return x
	}
	x.apply = true

	// ------- 对offset、 limit、page三个参数进行处理，得到最终的limit和offset --------
	// limit必须有，否则取默认值
	if x.OrmQuery.HasLimit() == false {
		if x.defaultPageSize > 0 {
			x.OrmQuery.SetLimit(x.defaultPageSize)
		} else {
			x.OrmQuery.SetLimit(DefaultPageSize)
		}
	}
	// page可以有，有的话，通过它计算额外offset
	if x.OrmQuery.HasPage() {
		offset := x.OrmQuery.GetLimit() * (x.OrmQuery.GetPage() - 1)
		if x.OrmQuery.HasOffset() {
			x.OrmQuery.SetOffset(x.OrmQuery.GetOffset() + offset)
		} else {
			x.OrmQuery.SetOffset(offset)
		}
	}
	// condition
	if x.OrmQuery.HasCondition() {
		// fill out the conditions object, we support "and"、"or"
		conditionPtrs := x.OrmQuery.GetCondition()
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
					orConditions = append(orConditions, condition)
				}
			}
		}
		conditions = append(andConditions, orConditions...)
		// change a condition object to a real ORM’s function's params
		// for example, gorm has functions like 'where'、'or', they both need params
		for _, condition := range conditions {
			operator := x.getOperator(condition.GetOperator())
			if operator == "" {
				continue
			}

			switch operator {
			case ES_OPERATOR_IS_EQUAL:
				x.EsBoolQuery.Filter(elastic.NewTermsQuery(condition.GetField(), condition.GetValue()))
			case ES_OPERATOR_IS_IN:
				x.EsBoolQuery.Filter(elastic.NewTermsQuery(condition.GetField(), condition.GetValue().([]interface{})...))
			case ES_OPERATOR_IS_NOT_IN:
				x.EsBoolQuery.MustNot(elastic.NewTermsQuery(condition.GetField(), condition.GetValue().([]interface{})...))
			case ES_OPERATOR_IS_GREATER_THAN:
				x.EsBoolQuery.Filter(elastic.NewRangeQuery(condition.GetField()).Gt(condition.GetValue()))
			case ES_OPERATOR_IS_GREATER_THAN_OR_EQUAL:
				x.EsBoolQuery.Filter(elastic.NewRangeQuery(condition.GetField()).Gte(condition.GetValue()))
			case ES_OPERATOR_IS_LESS_THAN:
				x.EsBoolQuery.Filter(elastic.NewRangeQuery(condition.GetField()).Lt(condition.GetValue()))
			case ES_OPERATOR_IS_LESS_THAN_OR_EQUAL:
				x.EsBoolQuery.Filter(elastic.NewRangeQuery(condition.GetField()).Lte(condition.GetValue()))
			case ES_OPERATOR_IS_LIKE:
				x.EsBoolQuery.Filter(elastic.NewFuzzyQuery(condition.GetField(), condition.GetValueString()))
			case ES_OPERATOR_IS_NOT_EQUAL:
				x.EsBoolQuery.MustNot(elastic.NewTermQuery(condition.GetField(), condition.GetValue()))
			}
		}
	}
	// sort
	if x.OrmQuery.HasSorter() {
		sorters := x.OrmQuery.GetSorter()
		for _, sorter := range sorters {
			switch sorter.GetDirection() {
			case DESCENDING:
				x.EsFieldSort = append(x.EsFieldSort, elastic.NewFieldSort(sorter.GetField()).Desc())
			case ASCENDING:
				x.EsFieldSort = append(x.EsFieldSort, elastic.NewFieldSort(sorter.GetField()).Asc())
			}
		}
	}
	return x
}

func (x *ElasticSearchSession) getOperator(operator int) string {
	m := x.operatorMap()
	return m[operator]
}

/*
This is a xorm parser.
When a user's request is analysed and turn to be a OrmQuery object.Any parser program can parse it to a DB object.

Note:
	OrmQuery 's include filed is different, it will not parse by current parser, it is usd by transformer program.
*/
//
const ES_OPERATOR_IS_EQUAL = "="
const ES_OPERATOR_IS_GREATER_THAN = ">"
const ES_OPERATOR_IS_GREATER_THAN_OR_EQUAL = ">="
const ES_OPERATOR_IS_LESS_THAN = "<"
const ES_OPERATOR_IS_LESS_THAN_OR_EQUAL = "<="
const ES_OPERATOR_IS_IN = "IN"
const ES_OPERATOR_IS_NOT_IN = "NOT IN"
const ES_OPERATOR_IS_LIKE = "LIKE"
const ES_OPERATOR_IS_NOT_LIKE = "NOT LIKE"
const ES_OPERATOR_IS_JSON_CONTAINS = "JSON_CONTAINS"
const ES_OPERATOR_IS_NOT_EQUAL = "<>"
const ES_OPERATOR_IS_IS_NULL = "IS NULL"
const ES_OPERATOR_IS_IS_NOT_NULL = "IS NOT NULL"
const ES_DEFAULT_KEY = "value"

func (x *ElasticSearchSession) operatorMap() map[int]string {
	return map[int]string{
		OPERATOR_IS_EQUAL:                 ES_OPERATOR_IS_EQUAL,
		OPERATOR_IS_GREATER_THAN:          ES_OPERATOR_IS_GREATER_THAN,
		OPERATOR_IS_GREATER_THAN_OR_EQUAL: ES_OPERATOR_IS_GREATER_THAN_OR_EQUAL,
		OPERATOR_IS_IN:                    ES_OPERATOR_IS_IN,
		OPERATOR_IS_NOT_IN:                ES_OPERATOR_IS_NOT_IN,
		OPERATOR_IS_LESS_THAN:             ES_OPERATOR_IS_LESS_THAN,
		OPERATOR_IS_LESS_THAN_OR_EQUAL:    ES_OPERATOR_IS_LESS_THAN_OR_EQUAL,
		OPERATOR_IS_LIKE:                  ES_OPERATOR_IS_LIKE,
		OPERATOR_IS_NOT_LIKE:              ES_OPERATOR_IS_NOT_LIKE,
		OPERATOR_IS_JSON_CONTAINS:         ES_OPERATOR_IS_JSON_CONTAINS,
		OPERATOR_IS_NOT_EQUAL:             ES_OPERATOR_IS_NOT_EQUAL,
		OPERATOR_IS_IS_NULL:               ES_OPERATOR_IS_IS_NULL,
		OPERATOR_IS_IS_NOT_NULL:           ES_OPERATOR_IS_IS_NOT_NULL,
	}
}

// -----------------------API for users to use-------------
func (x *ElasticSearchSession) SetTransformer(f *Transformer) *ElasticSearchSession {
	x.transform = f
	return x
}

func (x *ElasticSearchSession) SetPage(page int) *ElasticSearchSession {
	x.OrmQuery.SetPage(page)
	return x
}

func (x *ElasticSearchSession) SetLimit(limit int) *ElasticSearchSession {
	x.OrmQuery.SetLimit(limit)
	return x
}

// 0升序， 1降序
func (x *ElasticSearchSession) AddSorter(field string, direction int) {
	s := Sorter{
		field:     field,
		direction: direction,
	}
	x.OrmQuery.AddSorter(s)
}

// ------------------------ Func----------------------------
func (x *ElasticSearchSession) Pagination(index, indexType string, count int) PaginationResult {
	total, err := x.Count(index, indexType)
	if err != nil {
		panic(err)
	}
	pages := Paginator(x.OrmQuery.GetPage(), x.OrmQuery.GetLimit(), total)
	pages["count"] = count
	return pages
}

func (x *ElasticSearchSession) All(index, indexType string) (AllResult, error) {
	x.ApplyQuery()
	searchService := elastic.NewSearchService(x.Client).Query(x.EsBoolQuery)
	for _, sort := range x.EsFieldSort {
		searchService.SortBy(sort)
	}
	searchService.Size(x.OrmQuery.GetLimit())
	searchService.From(x.OrmQuery.GetOffset())
	searchResult, err := searchService.Index(index).Type(indexType).Do(context.Background())
	if err != nil {
		return nil, err
	}
	resultSlice := []map[string]interface{}{}
	for _, v := range searchResult.Hits.Hits {
		b, err := v.Source.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if m, err := qmap.NewWithString(string(b)); err == nil {
			m["_id"] = v.Id
			resultSlice = append(resultSlice, m)
		} else {
			return nil, err
		}
	}
	// transformer
	if x.transform != nil {
		transformer := *x.transform
		for key, one := range resultSlice {
			resultSlice[key] = transformer(one)
		}
	}

	return resultSlice, nil
}

func (x *ElasticSearchSession) Count(index, indexType string) (int64, error) {
	x.ApplyQuery()
	return x.Client.Count().Query(x.EsBoolQuery).Index(index).Type(indexType).Do(context.Background())
}

func (x *ElasticSearchSession) One(index, indexType string) (OneResult, error) {
	x.ApplyQuery()
	searchResult, err := x.Client.Search().Index(index).Type(indexType).Query(x.EsBoolQuery).Size(1).Do(context.Background())
	if err != nil {
		panic(err)
	}
	for _, v := range searchResult.Hits.Hits {
		b, err := v.Source.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if one, err := qmap.NewWithString(string(b)); err == nil {
			one["_id"] = v.Id
			return one, nil
		} else {
			return nil, err
		}
	}
	return qmap.QM{}, nil
}
