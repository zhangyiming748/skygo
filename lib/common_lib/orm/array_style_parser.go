package orm

import (
	"strings"
)

type ArrayStyleParser struct{}

func (this *ArrayStyleParser) Parse(params map[string]interface{}, query Query) Query {
	for key, val := range params {
		if operate, key := this.getOperateType(key); operate != "" {
			switch operate {
			case MG_EQUAL:
				query.AndWhereEqual(key, val)
			case MG_GREATER_THAN:
				query.AndWhereGt(key, val)
			case MG_GREATER_THAN_OR_EQUAL:
				query.AndWhereGte(key, val)
			case MG_LESS_THAN:
				query.AndWhereLt(key, val)
			case MG_LESS_THAN_OR_EQUAL:
				query.AndWhereLte(key, val)
			case MG_NOT_EQUAL:
				query.AndWhereNotEqual(key, val)
			case MG_LIKE:
				query.AndWhereLike(key, val)
			case MG_IN:
				query.AndWhereIn(key, val)
			case MG_EXISTS:
				query.AndExists(key, val)
			case MG_ORLIKE:
				query.OrLike(key, val)
			}
		}
	}

	return query

}

const MG_EQUAL = "e"
const MG_GREATER_THAN = "gt"
const MG_GREATER_THAN_OR_EQUAL = "gte"
const MG_LESS_THAN = "lt"
const MG_LESS_THAN_OR_EQUAL = "lte"
const MG_NOT_EQUAL = "ne"
const MG_IN = "in"
const MG_LIKE = "l"
const MG_EXISTS = "exists"
const MG_ORLIKE = "ol"

func (this *ArrayStyleParser) getOperateType(op string) (operator, key string) {
	if splited := strings.Split(op, "_"); len(splited) > 1 {
		return splited[0], strings.Join(splited[1:], "_")
	} else {
		return "", op
	}
}
