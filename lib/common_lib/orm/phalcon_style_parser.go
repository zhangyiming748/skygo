package orm

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type PhalconStyleParser struct{}

func (p *PhalconStyleParser) ParseQueryStr(queryStr string) Query {
	m := QueryStrToMap(queryStr)
	return p.ParseQueryMap(m)
}

func (p *PhalconStyleParser) ParseQueryStrAgain(queryStr string, q Query) Query {
	m := QueryStrToMap(queryStr)
	return p.ParseQueryMapAgain(m, q)
}

func (p *PhalconStyleParser) ParseQueryMap(params map[string]interface{}) Query {
	// 我们接受的数据，本质都是 map[string]string, 也就是我们定义的 UrlParmas类型
	data := UrlParams{}
	for k, v := range params {
		if _v, ok := v.(string); ok {
			data[k] = _v
		}
	}

	var err error
	fields := []string{}
	if p.isEnabled(FIELDS) {
		fields = p.extractCommaSeparatedValues(data, FIELDS)
	}

	offset := 0
	if p.isEnabled(OFFSET) {
		offset, _ = p.extractInt(data, OFFSET)
	}

	limit := 0
	if p.isEnabled(LIMIT) {
		if limit, err = p.extractInt(data, LIMIT); err != nil {
			panic(err)
		}
	}

	page := 1
	if p.isEnabled(PAGE) {
		page, _ = p.extractInt(data, PAGE)
	}

	where := map[string]interface{}{}
	if p.isEnabled(WHERE) {
		where = p.extractArray(data, WHERE)
	}

	sort := []Sorter{}
	if p.isEnabled(SORT) {
		sort = p.extractSort(data, SORT)
	}

	include := []string{}
	if p.isEnabled(INCLUDE) {
		include = p.extractCommaSeparatedValues(data, INCLUDE)
	}

	q := Query{}

	if len(fields) > 0 {
		q.AddField(fields)
	}

	if offset > 0 {
		q.SetOffset(offset)
	}

	if page > 0 {
		q.SetPage(page)
	}

	if limit > 0 {
		q.SetLimit(limit)
	}

	if len(include) > 0 {
		q.AddInclude(include)
	}

	if len(where) > 0 {
		for field, item := range where {
			for oper, value := range item.(map[string]interface{}) {
				operatorNum, ok := p.extractOperator(oper)
				if ok {
					q.AddCondition(&Condition{
						ctype:    CTYPE_AND,
						field:    field,
						value:    value,
						operator: operatorNum,
					})
				}
			}
		}
	}

	if len(sort) > 0 {
		for _, sortTtem := range sort {
			var direction int
			switch sortTtem.direction {
			case SORT_DESCENDING:
				direction = DESCENDING
			case SORT_ASCENDING:
				direction = ASCENDING
			default:
				direction = ASCENDING
			}

			q.AddSorter(Sorter{field: sortTtem.field, direction: direction})
		}
	}

	return q
}

func (p *PhalconStyleParser) ParseQueryMapAgain(params map[string]interface{}, q Query) Query {
	// 我们接受的数据，本质都是 map[string]string, 也就是我们定义的 UrlParmas类型
	data := UrlParams{}
	for k, v := range params {
		if _v, ok := v.(string); ok {
			data[k] = _v
		}
	}

	var err error
	fields := []string{}
	if p.isEnabled(FIELDS) {
		fields = p.extractCommaSeparatedValues(data, FIELDS)
	}

	offset := 0
	if p.isEnabled(OFFSET) {
		offset, _ = p.extractInt(data, OFFSET)
	}

	limit := 0
	if p.isEnabled(LIMIT) {
		if limit, err = p.extractInt(data, LIMIT); err != nil {
			panic(err)
		}
	}

	page := 1
	if p.isEnabled(PAGE) {
		page, _ = p.extractInt(data, PAGE)
	}

	where := map[string]interface{}{}
	if p.isEnabled(WHERE) {
		where = p.extractArray(data, WHERE)
	}

	sort := []Sorter{}
	if p.isEnabled(SORT) {
		sort = p.extractSort(data, SORT)
	}

	include := []string{}
	if p.isEnabled(INCLUDE) {
		include = p.extractCommaSeparatedValues(data, INCLUDE)
	}

	if len(fields) > 0 {
		q.AddField(fields)
	}

	if offset > 0 {
		q.SetOffset(offset)
	}

	if page > 0 {
		q.SetPage(page)
	}

	if limit > 0 {
		q.SetLimit(limit)
	}

	if len(include) > 0 {
		q.AddInclude(include)
	}

	if len(where) > 0 {
		for field, item := range where {
			for oper, value := range item.(map[string]interface{}) {
				operatorNum, ok := p.extractOperator(oper)
				if ok {
					q.AddCondition(&Condition{
						ctype:    CTYPE_AND,
						field:    field,
						value:    value,
						operator: operatorNum,
					})
				}
			}
		}
	}

	if len(sort) > 0 {
		for _, sortTtem := range sort {
			var direction int
			switch sortTtem.direction {
			case SORT_DESCENDING:
				direction = DESCENDING
			case SORT_ASCENDING:
				direction = ASCENDING
			default:
				direction = ASCENDING
			}

			q.AddSorter(Sorter{field: sortTtem.field, direction: direction})
		}
	}

	return q
}

func (p *PhalconStyleParser) isEnabled(s string) bool {
	for _, v := range enableFeatures {
		if v == s {
			return true
		}
	}
	return false
}

func (p *PhalconStyleParser) extractCommaSeparatedValues(data UrlParams, key string) []string {
	if _, ok := data[key]; ok == true {
		str := data[key]
		return strings.Split(str, ",")
	}
	return nil
}

func (p *PhalconStyleParser) extractInt(data UrlParams, key string) (int, error) {
	if _, ok := data[key]; ok == true {
		return strconv.Atoi(data[key])
	}
	return 0, nil
}

func (p *PhalconStyleParser) extractSort(data UrlParams, key string) []Sorter {
	sort := []Sorter{}
	if _, ok := data[key]; ok == true {
		var val interface{}
		if err := json.Unmarshal([]byte(data[key]), &val); err == nil {
			if reflect.TypeOf(val).Kind() == reflect.Map {
				// 如果是map，只取第一个
				if m, ok := val.(map[string]interface{}); ok {
					for k, v := range m {
						if _v, ok := v.(float64); ok {
							sort = append(sort, Sorter{
								k,
								int(_v),
							})
							break
						}
					}
				}
			} else if reflect.TypeOf(val).Kind() == reflect.Slice {
				// 如果是slice，说明有多个参数顺序值
				if i, ok := val.([]interface{}); ok {
					for _, s := range i {
						if m, ok := s.(map[string]interface{}); ok {
							for k, v := range m {
								if _v, ok := v.(float64); ok {
									sort = append(sort, Sorter{
										k,
										int(_v),
									})
									break
								}
							}
						}
					}
				}
			}
		}
	}
	return sort
}

func (p *PhalconStyleParser) extractArray(data UrlParams, key string) map[string]interface{} {
	if _, ok := data[key]; ok == true {
		arr := map[string]interface{}{}
		json.Unmarshal([]byte(data[key]), &arr)
		return arr
	}
	return nil
}

func (p *PhalconStyleParser) extractOperator(oper string) (int, bool) {
	operators := p.operatorMap()
	i, ok := operators[oper]
	return i, ok
}

func (p *PhalconStyleParser) operatorMap() map[string]int {
	return map[string]int{
		URL_OPERATOR_IS_EQUAL:                 Operator_is_equal,
		URL_OPERATOR_IS_GREATER_THAN:          OperatorIsGreaterThan,
		URL_OPERATOR_IS_GREATER_THAN_OR_EQUAL: OperatorIsGreaterThanOrEqual,
		URL_OPERATOR_IS_LESS_THAN:             OperatorIsLessThan,
		URL_OPERATOR_IS_LESS_THAN_OR_EQUAL:    OperatorIsLessThanOrEqual,
		URL_OPERATOR_IS_LIKE:                  OperatorIsLike,
		URL_OPERATOR_IS_JSON_CONTAINS:         OperatorIsJsonContains,
		URL_OPERATOR_IS_NOT_EQUAL:             OperatorIsNotEqual,
	}
}

const FIELDS = "fields"
const OFFSET = "offset"
const LIMIT = "limit"
const HAVING = "having"
const WHERE = "where"
const SORT = "sort"
const PAGE = "page"

// const EXCLUDES = "excludes"
const URL_OPERATOR_IS_EQUAL = "e"
const URL_OPERATOR_IS_GREATER_THAN = "gt"
const URL_OPERATOR_IS_GREATER_THAN_OR_EQUAL = "gte"
const URL_OPERATOR_IS_LESS_THAN = "lt"
const URL_OPERATOR_IS_LESS_THAN_OR_EQUAL = "lte"
const URL_OPERATOR_IS_LIKE = "l"
const URL_OPERATOR_IS_JSON_CONTAINS = "jc"
const URL_OPERATOR_IS_NOT_EQUAL = "ne"

// const URL_OPERATOR_CONTAINS = "c"
// const URL_OPERATOR_NOT_CONTAINS = "nc"
const SORT_ASCENDING = 1
const SORT_DESCENDING = -1
const INCLUDE = "include"

// url参数
type UrlParams = map[string]string

// queryParamStr为请求的参数部分， 本函数返回解析后的键值对组成的map
func GetParamsFromRawQuery(queryParamStr string) map[string]interface{} {
	u := url.URL{RawQuery: queryParamStr}
	params := map[string]interface{}{}
	for k, v := range u.Query() {
		if len(v) != 1 {
			continue
		}
		params[k] = v[0]
	}
	return params
}

var enableFeatures = []string{
	FIELDS,
	OFFSET,
	LIMIT,
	PAGE,
	WHERE,
	SORT,
	// EXCLUDES,
	INCLUDE,
}
