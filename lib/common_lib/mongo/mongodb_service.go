package mongo

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/service"

	"skygo_detection/guardian/src/net/qmap"
)

var mgoSession = map[string]*mgo.Database{}
var mgoMu sync.Mutex

const MGO_PAGE_LIMIT = 10000 //MongoDB分页count上限

func GetDefaultMongodbDatabase() *mgo.Database {
	mongoConfig := service.LoadConfig().MongoDB
	return GetMongodbDatabase(mongoConfig.DBName)
}

func GetMongodbDatabase(db string) *mgo.Database {
	mgoMu.Lock()
	defer mgoMu.Unlock()
	if mgoSession[db] == nil {
		mgoSession[db] = newMongoSession().DB(db)
	} else {
		if err := mgoSession[db].Session.Ping(); err != nil {
			mgoSession[db] = newMongoSession().DB(db)
		}
	}
	return mgoSession[db]
}

func newMongoSession() *mgo.Session {
	mongoConfig := service.LoadConfig().MongoDB
	dialUrl := fmt.Sprintf("mongodb://%s:%s@%s:%d%s?maxIdleTimeMS=10000", mongoConfig.Username, mongoConfig.Password, mongoConfig.Host, mongoConfig.Port, mongoConfig.ExtraUrl)
	if mongoConfig.AuthSource != "" {
		dialUrl += "&authSource=" + mongoConfig.AuthSource
	}
	if mongoConfig.ReplicaSet != "" {
		dialUrl += "&replicaSet=" + mongoConfig.ReplicaSet
	}
	if session, err := mgo.Dial(dialUrl); err == nil {
		return session
	} else {
		panic(err)
	}
}

type MgoOrmSession struct {
	MongoSession
}

func (t *MgoOrmSession) GetQuery() *Query {
	return t.Query
}

func NewMgoSession(table string) *MgoOrmSession {
	return &MgoOrmSession{*NewMgoSessionInner(GetDefaultMongodbDatabase().C(table), UrlParams{})}
}

func NewMgoSessionWithCond(table string, params qmap.QM) *MgoOrmSession {
	newMgoSession := MgoOrmSession{*NewMgoSessionInner(GetDefaultMongodbDatabase().C(table), UrlParams{})}
	return newMgoSession.AddCondition(params)
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

func (this *MgoOrmSession) AddCondition(params qmap.QM) *MgoOrmSession {
	for key, val := range params {
		if operate, key := this.getOperateType(key); operate != "" && !this.isZeroVal(key, val) {
			switch operate {
			case MG_EQUAL:
				this.Query.AndWhereEqual(key, val)
			case MG_GREATER_THAN:
				this.Query.AndWhereGt(key, val)
			case MG_GREATER_THAN_OR_EQUAL:
				this.Query.AndWhereGte(key, val)
			case MG_LESS_THAN:
				this.Query.AndWhereLt(key, val)
			case MG_LESS_THAN_OR_EQUAL:
				this.Query.AndWhereLte(key, val)
			case MG_NOT_EQUAL:
				this.Query.AndWhereNotEqual(key, val)
			case MG_LIKE:
				this.Query.AndWhereLike(key, val)
			case MG_IN:
				this.Query.AndWhereIn(key, val)
			case MG_EXISTS:
				this.Query.AndExists(key, val)
			case MG_ORLIKE:
				this.Query.OrLike(key, val)
			}
		}
	}

	return this
}

func (this *MgoOrmSession) AddUrlQueryCondition(queryParams string) *MgoOrmSession {
	urlParams := UrlParams{}

	u := url.URL{RawQuery: queryParams}
	for k, v := range u.Query() {
		if len(v) != 1 {
			continue
		}
		urlParams[k] = v[0]
	}
	urlQuery := NewQueryPhalconStyle(urlParams)
	this.Query.Merge(urlQuery)

	return this
}

func (this *MgoOrmSession) SetTransformFunc(fun func(qmap.QM) qmap.QM) *MgoOrmSession {
	transformFunc := TransformerFunc(func(result orm.OneResult) orm.OneResult {
		return map[string]interface{}(fun(qmap.QM(result)))
	})
	this.MongoSession.SetTransformer(&transformFunc)
	return this
}

func (this *MgoOrmSession) SetTransformer(trans Transformer) *MgoOrmSession {
	this.MongoSession.SetTransformer(this.getOrmTransformer(trans))
	return this
}

func (this *MgoOrmSession) getOrmTransformer(t Transformer) *TransformerFunc {
	//if user's request URL has params "include", we can find it in Query object
	includeFuncValue := map[string]reflect.Value{} //key is function name, value is a reflect.Value
	if this.Query.HasInclude() {
		include := this.Query.GetInclude()
		for _, v := range include {
			//use reflect, get the Value corresponding to the function name
			value := reflect.ValueOf(t).MethodByName(CamelString("include_" + v))
			if value.IsValid() == true {
				includeFuncValue[v] = value
			}
		}
	}

	f := TransformerFunc(func(result OneResult) OneResult {
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

func (this *MgoOrmSession) getOperateType(op string) (operator, key string) {
	if splited := strings.Split(op, "_"); len(splited) > 1 {
		return splited[0], strings.Join(splited[1:], "_")
	} else {
		return "", op
	}
}

// 判断是否是零值
// 目前主要过滤渠道号
func (this *MgoOrmSession) isZeroVal(key string, val interface{}) bool {
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

func (this *MgoOrmSession) SetLimit(limit int) *MgoOrmSession {
	this.MongoSession.SetLimit(limit)
	return this
}

func (this *MgoOrmSession) SetPage(page int) *MgoOrmSession {
	this.MongoSession.SetPage(page)
	return this
}

func (this *MgoOrmSession) Pagination(count int) map[string]interface{} {
	return this.MongoSession.Pagination(count)
}

func (this *MgoOrmSession) Get() (*[]map[string]interface{}, error) {
	all := this.All()
	return &all, nil
}

func (this *MgoOrmSession) QueryGet(operations []bson.M) (*[]map[string]interface{}, error) {
	all := this.MTCHAll(operations)
	return &all, nil
}

func (this *MgoOrmSession) MATCHGetPage(operations []bson.M) (*qmap.QM, error) {
	all := this.MTCHAll(operations)
	allNum := this.MATCHCOUNT(operations)
	pages := this.MongoSession.MatchPagination(len(all), allNum)

	qm := qmap.QM{
		"list":       all,
		"pagination": pages,
	}
	return &qm, nil
}

func (this *MgoOrmSession) MATCHALL(operations []bson.M) (*qmap.QM, error) {
	all := this.MTCHAll(operations)
	qm := qmap.QM{
		"match_all": all,
	}
	return &qm, nil
}

func (this *MgoOrmSession) GetPage() (*qmap.QM, error) {
	all := this.All()
	pages := this.MongoSession.Pagination(len(all))

	qm := qmap.QM{
		"list":       all,
		"pagination": pages,
	}
	return &qm, nil
}

// 此方法用于提高MongoDB海量数据集的分页查询效率（MongoDB在对海量数据集进行count时很耗时）
// 如果数据集总数total能够从外部传入，则本方法不再进行额外的count操作
// 如果数据集总数total传参-1，则本方法使用有上限的count统计操作
func (this *MgoOrmSession) GetPageWithLimit(total int) (*qmap.QM, error) {
	all := this.All()
	pages := this.MongoSession.PaginationWithLimit(total, len(all), MGO_PAGE_LIMIT)
	qm := qmap.QM{
		"list":       all,
		"pagination": pages,
	}
	return &qm, nil
}

func (this *MgoOrmSession) GetOne() (*qmap.QM, error) {
	tmpRes, err := this.MongoSession.GetOne()
	var res = qmap.QM(tmpRes)
	return &res, err
}

func (this *MgoOrmSession) One(result interface{}) error {
	this.ApplyQuery()
	return this.Session.Find(this.Match).One(result)
}

func (this *MgoOrmSession) Count() (n int, err error) {
	this.ApplyQuery()
	return this.Session.Find(this.Match).Count()
}
func (this *MgoOrmSession) Update(selector interface{}, update interface{}) error {
	return this.Session.Update(selector, update)
}

func (this *MgoOrmSession) UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return this.Session.UpdateAll(selector, update)
}

func (this *MgoOrmSession) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return this.Session.Upsert(selector, update)
}

func (this *MgoOrmSession) Insert(docs ...interface{}) error {
	return this.Session.Insert(docs...)
}

func (this *MgoOrmSession) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	return this.Session.RemoveAll(selector)
}

func (this *MgoOrmSession) Iter() *mgo.Iter {
	this.ApplyQuery()
	return this.Session.Find(this.Match).Iter()
}

func GridFSUpload(prefix, filename string, fileContent []byte) (fileId string, err error) {
	if fi, createErr := GetDefaultMongodbDatabase().GridFS(prefix).Create(filename); err == nil {
		if _, writeErr := fi.Write(fileContent); writeErr == nil {
			if closeErr := fi.Close(); closeErr == nil {
				return fi.Id().(bson.ObjectId).Hex(), nil
			} else {
				err = closeErr
			}
		} else {
			err = writeErr
		}
	} else {
		err = createErr
	}
	return
}

func GridFSOpenId(prefix string, fileId bson.ObjectId) (file *mgo.GridFile, err error) {
	return GetDefaultMongodbDatabase().GridFS(prefix).OpenId(fileId)
}

func GridFSOpen(prefix, filePath string) (file *mgo.GridFile, err error) {
	return GetDefaultMongodbDatabase().GridFS(prefix).Open(filePath)
}

func GridFSRename(prefix, newFilename string, fileId bson.ObjectId) (string, error) {
	if oldFile, err := GetDefaultMongodbDatabase().GridFS(prefix).OpenId(fileId); err == nil {
		defer oldFile.Close()
		if fileContent, readErr := ioutil.ReadAll(oldFile); readErr == nil {
			if newFileId, renameErr := GridFSUpload(prefix, newFilename, fileContent); renameErr == nil {
				if removeErr := GetDefaultMongodbDatabase().GridFS(prefix).RemoveId(fileId); removeErr != nil {
					return "", renameErr
				}
				return newFileId, nil
			} else {
				return "", renameErr
			}
		} else {
			return "", readErr
		}

	} else {
		return "", err
	}
}
func GridFSRemoveFile(prefix string, fileId bson.ObjectId) error {
	err := GetDefaultMongodbDatabase().GridFS(prefix).RemoveId(fileId)
	if err != nil {
		return err
	}
	return nil
}
