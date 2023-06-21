package orm

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-xorm/xorm"
	"gopkg.in/olivere/elastic.v5"
)

type AllResult = []map[string]interface{}
type PaginationResult = map[string]interface{}
type OneResult = map[string]interface{}

func NewQueryPhalconStyle(urlParams UrlParams) *Query {
	return new(PhalconStyleParser).Parse(urlParams)
}

func NewXormSession(session *xorm.Session, urlParams UrlParams) *XormSession {
	return &XormSession{
		Session:         session,
		Query:           NewQueryPhalconStyle(urlParams),
		defaultPageSize: DefaultPageSize,
		structToMapFunc: StructToMap,
	}
}

func NewElasticSearchSession(session *elastic.Client, urlParams UrlParams) *ElasticSearchSession {
	return &ElasticSearchSession{
		Client:          session,
		OrmQuery:        NewQueryPhalconStyle(urlParams),
		defaultPageSize: DefaultPageSize,
		structToMapFunc: StructToMap,
		EsBoolQuery:     elastic.NewBoolQuery(),
		EsFieldSort:     []*elastic.FieldSort{},
	}
}

func NewMgoSession(session *mgo.Collection, urlParams UrlParams) *MongoSession {
	return &MongoSession{
		Session:         session,
		Query:           NewQueryPhalconStyle(urlParams),
		defaultPageSize: DefaultPageSize,
		structToMapFunc: StructToMap,
		Match:           bson.M{},
		sort:            bson.D{},
	}
}
