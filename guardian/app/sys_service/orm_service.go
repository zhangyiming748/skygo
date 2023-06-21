package sys_service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/url"
	"reflect"
	"skygo_detection/guardian/app/transformer"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/guardian/src/orm"
	"skygo_detection/guardian/util"
	"strings"
	"sync"
	"time"
)

const DEFAULT_DB = "default_db"

var engine = map[string]*xorm.Engine{}
var mu sync.Mutex

func NewOrm() *xorm.Engine {
	return NewMysqlOrm(DEFAULT_DB)
}

func NewMysqlOrm(db string) *xorm.Engine {
	mu.Lock()
	defer mu.Unlock()
	if engine[db] == nil {
		engine[db] = getNewConnection(db)
	} else {
		if err := engine[db].Ping(); err != nil {
			engine[db] = getNewConnection(db)
		}
	}
	return engine[db]
}

func getNewConnection(db string) *xorm.Engine {
	var dbConfig *DBConfig
	if db == DEFAULT_DB {
		dbConfig = GetDefaultDBConfig()
		dbConfigStr := GetDBConnectionStr(dbConfig)
		if engine, err := xorm.NewEngine("mysql", dbConfigStr); err == nil {
			engine.SetMaxIdleConns(dbConfig.MaxIdleConnection)
			engine.SetMaxOpenConns(dbConfig.MaxOpenConnection)
			engine.SetConnMaxLifetime(time.Second * dbConfig.MaxLifeTime)

			InitDefaultXORMLogger(engine)
			return engine
		} else {
			panic(err)
		}
	} else {
		panic("unknown db name")
	}
}

func GetDBConnectionStr(dbConfig *DBConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s", dbConfig.UserName, dbConfig.Password, dbConfig.HostName, dbConfig.Port, dbConfig.DBName, dbConfig.Charset)
}

func InitDefaultXORMLogger(engine *xorm.Engine) {
	logConfig := GetDefaultDBConfig().Log
	engine.SetLogger(xorm.NewSimpleLogger(GetLogFile(logConfig.FilePath)))
	engine.SetLogLevel(logConfig.Level)
	engine.ShowSQL(GetDefaultDBConfig().ShowSql)
}

type OrmSession struct {
	orm.XormSession
}

func (t *OrmSession) GetQuery() *orm.Query {
	return t.Query
}

func NewSession() *OrmSession {
	return &OrmSession{*orm.NewXormSession(NewOrm().NewSession(), orm.UrlParams{})}
}

func NewSessionWithCond(params qmap.QM) *OrmSession {
	newSession := OrmSession{*orm.NewXormSession(NewOrm().NewSession(), orm.UrlParams{})}
	return newSession.AddCondition(params)
}

func (this *OrmSession) Table(table interface{}) *OrmSession {
	this.Session.Table(table)
	return this
}

const EQUAL = "e"
const GREATER_THAN = "gt"
const GREATER_THAN_OR_EQUAL = "gte"
const LESS_THAN = "lt"
const LESS_THAN_OR_EQUAL = "lte"
const NOT_EQUAL = "ne"
const IN = "in"
const LIKE = "lk"
const OR_LIKE = "orlk"

func (this *OrmSession) AddCondition(params qmap.QM) *OrmSession {
	for key, val := range params {
		if operate, key := this.getOperateType(key); operate != "" && !this.isZeroVal(key, val) {
			switch operate {
			case EQUAL:
				this.Query.AndWhereEqual(key, val)
			case GREATER_THAN:
				this.Query.AndWhereGt(key, val)
			case GREATER_THAN_OR_EQUAL:
				this.Query.AndWhereGte(key, val)
			case LESS_THAN:
				this.Query.AndWhereLt(key, val)
			case LESS_THAN_OR_EQUAL:
				this.Query.AndWhereLte(key, val)
			case NOT_EQUAL:
				this.Query.AndWhereNotEqual(key, val)
			case IN:
				this.Query.AndWhereIn(key, val)
			case LIKE:
				this.Query.AndWhereLike(key, val)
			case OR_LIKE:
				this.Query.OrLike(key, val)
			}
		}
	}

	return this
}

func (this *OrmSession) AddUrlQueryCondition(queryParams string) *OrmSession {
	urlParams := orm.UrlParams{}

	u := url.URL{RawQuery: queryParams}
	for k, v := range u.Query() {
		if len(v) != 1 {
			continue
		}
		urlParams[k] = v[0]
	}
	urlQuery := orm.NewQueryPhalconStyle(urlParams)
	this.Query.Merge(urlQuery)

	return this
}

func (this *OrmSession) SetTransformFunc(fun func(qmap.QM) qmap.QM) *OrmSession {
	transformFunc := orm.Transformer(func(result orm.OneResult) orm.OneResult {
		return map[string]interface{}(fun(qmap.QM(result)))
	})
	this.XormSession.SetTransformer(&transformFunc)
	return this
}

func (this *OrmSession) SetTransformer(trans transformer.Transformer) *OrmSession {
	this.XormSession.SetTransformer(this.getOrmTransformer(trans))
	return this
}

func (this *OrmSession) getOrmTransformer(t transformer.Transformer) *orm.Transformer {
	//if user's request URL has params "include", we can find it in Query object
	includeFuncValue := map[string]reflect.Value{} //key is function name, value is a reflect.Value
	if this.Query.HasInclude() {
		include := this.Query.GetInclude()
		for _, v := range include {
			//use reflect, get the Value corresponding to the function name
			value := reflect.ValueOf(t).MethodByName(util.CamelString("include_" + v))
			if value.IsValid() == true {
				includeFuncValue[v] = value
			}
		}
	}

	f := orm.Transformer(func(result orm.OneResult) orm.OneResult {
		//call include function which the transformer has
		if len(includeFuncValue) > 0 {
			for key, value := range includeFuncValue {
				param := make([]reflect.Value, 1)
				param[0] = reflect.ValueOf(result)
				result[key] = value.Call(param)[0].Interface()
			}
		}

		//AdditionalFields
		result = map[string]interface{}(t.AdditionalFields(qmap.QM(result)))

		//ExcludeFields
		if fields := t.ExcludeFields(); len(fields) > 0 {
			for _, field := range fields {
				if _, exist := result[field]; exist {
					delete(result, field)
				}
			}
		}
		return result
	})
	return &f
}

func (this *OrmSession) SetExcludeFields(excludeFields []string) *OrmSession {
	this.XormSession.SetExcludeFields(excludeFields)
	return this
}

func (this *OrmSession) getOperateType(op string) (operator, key string) {
	if splited := strings.Split(op, "_"); len(splited) > 1 {
		return splited[0], strings.Join(splited[1:], "_")
	} else {
		return "", op
	}
}

// 判断是否是零值
// 目前主要过滤渠道号
func (this *OrmSession) isZeroVal(key string, val interface{}) bool {
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

func (this *OrmSession) SetLimit(limit int) *OrmSession {
	this.XormSession.SetLimit(limit)
	return this
}

func (this *OrmSession) SetPage(page int) *OrmSession {
	this.XormSession.SetPage(page)
	return this
}

func (this *OrmSession) Pagination(table interface{}, count int) qmap.QM {
	//totalNumber
	var totalNum int64
	totalNum, err := this.Count(table)
	if err != nil {
		panic(err)
	}
	pages := orm.Paginator(this.Query.GetPage(), this.Query.GetLimit(), int64(totalNum))
	pages["count"] = count
	return pages
}

func (this *OrmSession) Get(result interface{}) (*[]map[string]interface{}, error) {
	all, _ := this.All(result)

	return &all, nil
}

func (this *OrmSession) Find(result interface{}) error {
	this.ApplyQuery()
	return this.XormSession.Session.Find(result)
}

func (this *OrmSession) GetPage(table interface{}, result interface{}) (*qmap.QM, error) {
	all, session := this.All(result)
	//totalNumber
	var totalNum int64
	totalNum, err := session.Count(table)
	if err != nil {
		panic(err)
	}
	pages := orm.Paginator(this.Query.GetPage(), this.Query.GetLimit(), int64(totalNum))
	pages["count"] = len(all)

	qm := qmap.QM{
		"list":       all,
		"pagination": pages,
	}

	return &qm, nil
}

func (this *OrmSession) GetOne(result interface{}) (bool, *qmap.QM) {
	has, item := this.One(result)
	temp := qmap.QM(item)
	return has, &temp
}

func (this *OrmSession) Rows(result interface{}) (*xorm.Rows, error) {
	this.ApplyQuery()
	return this.XormSession.Session.Rows(result)
}

func (this *OrmSession) Count(beans ...interface{}) (int64, error) {
	this.ApplyQuery()
	return this.XormSession.Session.Count(beans...)
}

func (this *OrmSession) Cols(columns ...string) *OrmSession {
	this.XormSession.Session.Cols(columns...)
	return this
}

func (this *OrmSession) Select(str string) *OrmSession {
	this.XormSession.Session.Select(str)
	return this
}

func (this *OrmSession) OrderBy(order string) *OrmSession {
	this.XormSession.Session.OrderBy(order)
	return this
}

func (this *OrmSession) GroupBy(keys string) *OrmSession {
	this.XormSession.Session.GroupBy(keys)
	return this
}

func (this *OrmSession) In(column string, args ...interface{}) *OrmSession {
	this.Session.In(column, args...)
	return this
}

func (this *OrmSession) Create(table interface{}) (int64, error) {
	return this.Session.InsertOne(table)
}

func (this *OrmSession) Insert(table interface{}) (int64, error) {
	return this.Session.Insert(table)
}

func (this *OrmSession) Update(table interface{}, id int) (int64, error) {
	return this.Session.ID(id).Update(table)
}

func (this *OrmSession) UpdateById(table interface{}, id int, data qmap.QM) (int64, error) {
	params := qmap.QM{
		"e_id": id,
	}
	this.AddCondition(params)
	this.ApplyQuery()
	return this.Session.Table(table).Update(data)
}

func (this *OrmSession) DeleteByIds(table interface{}, ids []int) (int64, error) {
	this.AddCondition(qmap.QM{"in_id": ids})
	this.ApplyQuery()
	return this.Session.Delete(table)
}

func (this *OrmSession) Delete(table interface{}) (int64, error) {
	this.ApplyQuery()
	return this.Session.Delete(table)
}
