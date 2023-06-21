package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type AllResult = []map[string]interface{}
type PaginationResult = map[string]interface{}

type OneResult = map[string]interface{}

func NewQueryPhalconStyle(urlParams UrlParams) *Query {
	return new(PhalconStyleParser).Parse(urlParams)
}

func NewMgoSessionInner(session *mgo.Collection, urlParams UrlParams) *MongoSession {
	return &MongoSession{
		Session:         session,
		Query:           NewQueryPhalconStyle(urlParams),
		defaultPageSize: DefaultPageSize,
		structToMapFunc: StructToMap,
		Match:           bson.M{},
		sort:            bson.D{},
	}
}
