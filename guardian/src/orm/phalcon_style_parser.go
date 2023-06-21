package orm

import (
	"encoding/json"
	"strconv"
	"strings"
)

/*
PhalconStyle

The request URL's params parts describe the conditions we need focus on.

Params [where] : it must be a standard json string describe the find conditions, it should use with operator like "e"、"gt" egs.
				for example:  ?where={"sn":{"e":"89860617070021807390"}}
Params [fields] : it must be a  string separated  by comma, describe the fields that must be shown.
				for example: ?where={"sn":{"e":"89860617070021807390"}}&fields=sn,brand
Params [sort]:  it must be a standard json string describe the order of the result. -1 represents descend, 1 represents Ascend
				for example: ?where={"brand":%20"{"e":"BYD"}}&limit=10&sort={"id":-1}
Params [offset] :  it is a number that tell us how many records should be skipped before search.
				for example: ?where={"brand":"{"e":"BYD"}}&offset=2&limit=10
Params [limit] :  it is a number that tell us how many records should be return
				for example: ?where={"brand":"{"e":"BYD"}}&offset=2&limit=10
Params [include] : it must be a string separated by comma, describe the corresponding models that must be search together
				for example: ?where={"brand":"{"e":"BYD"}}&offset=2&limit=10&include=vehicle_location
*/

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
const URL_OPERATOR_IS_IN = "in"

// const URL_OPERATOR_CONTAINS = "c"
// const URL_OPERATOR_NOT_CONTAINS = "nc"
const SORT_ASCENDING = 1
const SORT_DESCENDING = -1
const INCLUDE = "include"

var enableFeatures = []string{
	FIELDS,
	OFFSET,
	LIMIT,
	PAGE,
	WHERE,
	SORT,
	//EXCLUDES,
	INCLUDE,
}

type PhalconStyleParser struct {
}

func (p *PhalconStyleParser) Parse(data UrlParams) *Query {
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
		limit, _ = p.extractInt(data, LIMIT)
	}

	page := 1
	if p.isEnabled(PAGE) {
		page, _ = p.extractInt(data, PAGE)
	}

	where := map[string]interface{}{}
	if p.isEnabled(WHERE) {
		where = p.extractArray(data, WHERE)
	}

	sort := map[string]interface{}{}
	if p.isEnabled(SORT) {
		sort = p.extractArray(data, SORT)
	}

	include := []string{}
	if p.isEnabled(INCLUDE) {
		include = p.extractCommaSeparatedValues(data, INCLUDE)
	}

	q := new(Query)

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
					if operatorNum == OPERATOR_IS_IN {
						q.AddCondition(&Condition{
							ctype:    CTYPE_IN,
							field:    field,
							value:    value,
							operator: OPERATOR_IS_IN,
						})
					} else {
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
	}

	if len(sort) > 0 {
		for field, RawDrection := range sort {
			var direction int
			_RowDirection := RawDrection.(float64)
			switch _RowDirection {
			case SORT_DESCENDING:
				direction = DESCENDING
			case SORT_ASCENDING:
				direction = ASCENDING
			default:
				direction = ASCENDING
			}

			q.AddSorter(Sorter{field: field, direction: direction})
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
		URL_OPERATOR_IS_EQUAL:                 OPERATOR_IS_EQUAL,
		URL_OPERATOR_IS_GREATER_THAN:          OPERATOR_IS_GREATER_THAN,
		URL_OPERATOR_IS_GREATER_THAN_OR_EQUAL: OPERATOR_IS_GREATER_THAN_OR_EQUAL,
		URL_OPERATOR_IS_LESS_THAN:             OPERATOR_IS_LESS_THAN,
		URL_OPERATOR_IS_LESS_THAN_OR_EQUAL:    OPERATOR_IS_LESS_THAN_OR_EQUAL,
		URL_OPERATOR_IS_LIKE:                  OPERATOR_IS_LIKE,
		URL_OPERATOR_IS_JSON_CONTAINS:         OPERATOR_IS_JSON_CONTAINS,
		URL_OPERATOR_IS_NOT_EQUAL:             OPERATOR_IS_NOT_EQUAL,
		URL_OPERATOR_IS_IN:                    OPERATOR_IS_IN,
	}
}
