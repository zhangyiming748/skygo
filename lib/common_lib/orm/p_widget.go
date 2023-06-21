package orm

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type AllResult = []map[string]interface{}
type PaginationResult = map[string]interface{}
type OneResult = map[string]interface{}

// 解析请求参数，从中得到Query
func DecodeQuery(ctx gin.Context) {

}

type PWidget struct {
	queryStr        string
	transformerFunc *TransformerFunc
	session         xorm.Engine
	Query

	isSet bool
	PhalconStyleParser
}

// 首次执行生效
// 把用户请求的queryString按照phalcon风格解析后, 填充到Query
func (p *PWidget) SetQueryStr(queryStr string) *PWidget {
	if p.isSet == false {
		p.Query = p.ParseQueryStr(queryStr)
	}
	return p
}

// 设置transformerFunc， 参数即可TransformerFunc函数类型的对象
func (p *PWidget) SetTransformerFunc(f TransformerFunc) *PWidget {
	p.transformerFunc = &f
	return p
}

// 设置transformerFunc， 参数为实现了TransformerIf接口的对象
func (p *PWidget) SetTransformer(i TransformerIf) *PWidget {
	p.transformerFunc = GetTransformerFunc(i)
	return p
}

// 使用p.transformerFunc
func (p *PWidget) transformerRun(result OneResult) OneResult {
	if p.transformerFunc != nil {
		result = (*p.transformerFunc)(result)
	}
	return result
}

// 把widget中的Query应用到session中
func (p *PWidget) ApplyQuery(session *xorm.Session) *xorm.Session {
	session = ApplyQuery(p.Query, session)
	return session
}

func (p *PWidget) PaginatorFind(session *xorm.Session, modelsPtr interface{}) map[string]interface{} {
	session = p.ApplyQuery(session)

	// 分页设置
	// 对offset、 limit、page三个参数进行处理，得到最终的limit和offset
	// page可以有，有的话，通过它计算额外offset
	limit := p.Query.GetLimit()
	page := p.Query.GetPage()
	offset := p.Query.GetOffset()
	if limit <= 0 {
		limit = DefaultPageSize
	}
	if page > 0 {
		offset = limit * (page - 1)
	}
	session = session.Limit(limit, offset)

	// 查询
	total, err := session.FindAndCount(modelsPtr)
	if err != nil {
		panic(err)
	}

	// struct to map
	all := AllResult{}
	for i := 0; i < reflect.ValueOf(modelsPtr).Elem().Len(); i++ {
		// reflect.ValueOf(modelsPtr).Elem().Index(i).Type().Kind() is reflect.Struct
		one := reflect.ValueOf(modelsPtr).Elem().Index(i).Interface()
		all = append(all, StructToMap(one))
	}

	for key, one := range all {
		all[key] = p.transformerRun(one)
	}

	// 获取分页内容
	paginator := Paginator(page, limit, total)

	// 返回最终数据
	result := gin.H{
		"list":       all,
		"pagination": paginator,
	}

	return result
}

func (p *PWidget) Find(session *xorm.Session) map[string]interface{} {
	session = p.ApplyQuery(session)
	models := make([]map[string]interface{}, 0)

	// 查询
	err := session.Find(&models)
	if err != nil {
		panic(err)
	}
	for key, value := range models {
		models[key] = p.transformerRun(value)
	}

	// 返回最终数据
	result := gin.H{
		"list": models,
	}

	return result
}

func (p *PWidget) All(session *xorm.Session, modelsPtr interface{}) ([]map[string]interface{}, error) {
	session = p.ApplyQuery(session)

	// 查询
	err := session.Find(modelsPtr)
	if err != nil {
		return nil, err
	}

	// struct to map
	all := AllResult{}
	for i := 0; i < reflect.ValueOf(modelsPtr).Elem().Len(); i++ {
		// reflect.ValueOf(modelsPtr).Elem().Index(i).Type().Kind() is reflect.Struct
		one := reflect.ValueOf(modelsPtr).Elem().Index(i).Interface()
		all = append(all, StructToMap(one))
	}

	for key, one := range all {
		all[key] = p.transformerRun(one)
	}

	return all, nil
}

func (p *PWidget) One(session *xorm.Session, modelsPtr interface{}) (OneResult, error) {
	session = p.ApplyQuery(session)

	// 查询
	has, err := session.Get(modelsPtr)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("记录不存在")
	}

	// struct to map
	one := OneResult{}
	one = StructToMap(reflect.ValueOf(modelsPtr).Elem().Interface())
	one = p.transformerRun(one)

	return one, nil
}

// todo mysql貌似不行
func (p *PWidget) Get(session *xorm.Session) (OneResult, error) {
	session = p.ApplyQuery(session)
	result := OneResult{}
	has, err := session.Get(&result)
	if err != nil {
		panic(err)
	}

	if has {
		result = p.transformerRun(result)
		return result, nil
	} else {
		return result, errors.New("记录不存在")
	}
}
