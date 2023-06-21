package sys_service

import (
	"net/url"
	"strings"

	"gopkg.in/olivere/elastic.v5"

	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/guardian/src/orm"
)

func NewElasticSession(session *elastic.Client) *ElasticSearchSession {
	r := new(ElasticSearchSession)
	r.Xs = orm.NewElasticSearchSession(session, orm.UrlParams{})
	return r
}

type ElasticSearchSession struct {
	Xs *orm.ElasticSearchSession
}

const ES_EQUAL = "e"
const ES_GREATER_THAN = "gt"
const ES_GREATER_THAN_OR_EQUAL = "gte"
const ES_LESS_THAN = "lt"
const ES_LESS_THAN_OR_EQUAL = "lte"
const ES_NOT_EQUAL = "ne"
const ES_IN = "in"
const ES_LIKE = "l"
const ES_EXISTS = "exists"
const ES_ORLIKE = "ol"

func (this *ElasticSearchSession) AddCondition(params qmap.QM) *ElasticSearchSession {
	for key, val := range params {
		if operate, key := this.getOperateType(key); operate != "" && !this.isZeroVal(key, val) {
			switch operate {
			case ES_EQUAL:
				this.Query().AndWhereEqual(key, val)
			case ES_GREATER_THAN:
				this.Query().AndWhereGt(key, val)
			case ES_GREATER_THAN_OR_EQUAL:
				this.Query().AndWhereGte(key, val)
			case ES_LESS_THAN:
				this.Query().AndWhereLt(key, val)
			case ES_LESS_THAN_OR_EQUAL:
				this.Query().AndWhereLte(key, val)
			case ES_NOT_EQUAL:
				this.Query().AndWhereNotEqual(key, val)
			case ES_LIKE:
				this.Query().AndWhereLike(key, val)
			case ES_IN:
				this.Query().AndWhereIn(key, val)
			case ES_EXISTS:
				this.Query().AndExists(key, val)
			case ES_ORLIKE:
				this.Query().OrLike(key, val)
			}
		}
	}

	return this
}

func (this *ElasticSearchSession) getOperateType(op string) (operator, key string) {
	if splited := strings.Split(op, "_"); len(splited) > 1 {
		return splited[0], strings.Join(splited[1:], "_")
	} else {
		return "", op
	}
}

// 判断是否是零值
// 目前主要过滤渠道号
func (this *ElasticSearchSession) isZeroVal(key string, val interface{}) bool {
	if key != "channel_id" {
		return false
	}
	switch val.(type) {
	case string:
		if val.(string) == "" {
			return true
		}
	}
	return false
}

func (t *ElasticSearchSession) Query() *orm.Query {
	return t.Xs.OrmQuery
}

func (t *ElasticSearchSession) GetPage(index, indexType string) (*qmap.QM, error) {
	all, err := t.Xs.All(index, indexType)
	if err != nil {
		return nil, err
	}
	qm := qmap.QM{
		"list":       all,
		"pagination": t.Xs.Pagination(index, indexType, len(all)),
	}
	return &qm, nil
}

func (t *ElasticSearchSession) GetAll(index, indexType string) ([]map[string]interface{}, error) {
	return t.Xs.All(index, indexType)
}

func (t *ElasticSearchSession) GetOne(index, indexType string) (*qmap.QM, error) {
	one, err := t.Xs.One(index, indexType)
	var oneQM qmap.QM = one
	return &oneQM, err
}

func (t *ElasticSearchSession) Count(index, indexType string) (int64, error) {
	return t.Xs.Count(index, indexType)
}

func (t *ElasticSearchSession) AddSorter(field string, isDescending bool) *ElasticSearchSession {
	if isDescending {
		t.Xs.AddSorter(field, 1)
	} else {
		t.Xs.AddSorter(field, 0)
	}
	return t
}

func (this *ElasticSearchSession) AddUrlQueryCondition(queryParams string) *ElasticSearchSession {
	urlParams := orm.UrlParams{}

	u := url.URL{RawQuery: queryParams}
	for k, v := range u.Query() {
		if len(v) != 1 {
			continue
		}
		urlParams[k] = v[0]
	}
	urlQuery := orm.NewQueryPhalconStyle(urlParams)
	this.Xs.OrmQuery.Merge(urlQuery)

	return this
}

func (t *ElasticSearchSession) SetTransformFunc(fun func(qmap.QM) qmap.QM) *ElasticSearchSession {
	transformFunc := orm.Transformer(func(result orm.OneResult) orm.OneResult {
		return fun(result)
	})
	t.Xs.SetTransformer(&transformFunc)
	return t
}
