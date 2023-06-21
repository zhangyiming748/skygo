package orm_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/orm"
)

func NewWidgetWithParams(collection string, params qmap.QM) *Widget {
	w := new(Widget)
	coll := GetDefaultMongoDatabase().Collection(collection)
	w.collection = coll
	w.SetParams(params)
	return w
}

func NewWidgetWithCollectionName(collection string) *Widget {
	w := new(Widget)
	coll := GetDefaultMongoDatabase().Collection(collection)
	w.collection = coll
	return w
}

func NewWidget() *Widget {
	return new(Widget)
}

type AllResult = []map[string]interface{}
type PaginationResult = map[string]interface{}
type OneResult = map[string]interface{}

type Widget struct {
	queryStr        string
	transformerFunc *orm.TransformerFunc
	collection      *mongo.Collection      // mongodb集合
	p1              orm.PhalconStyleParser // 一种解析器，把查询字符串解析到orm.Query
	p2              orm.ArrayStyleParser   // 一种解析器，把查询array解析到orm.Query
	orm.Query                              // 所有查询条件的最终数据结构

	Match bson.M
	sort  bson.D
}

func (w *Widget) SetCollection(collection *mongo.Collection) *Widget {
	w.collection = collection
	return w
}

// 首次执行生效
// 把用户请求的queryString按照phalcon风格解析后, 填充到Query
func (w *Widget) SetQueryStr(queryStr string) *Widget {
	w.Query = w.p1.ParseQueryStrAgain(queryStr, w.Query)
	return w
}

func (w *Widget) SetParams(params map[string]interface{}) *Widget {
	w.Query = w.p2.Parse(params, w.Query)
	return w
}

// 设置transformerFunc， 参数即可TransformerFunc函数类型的对象
func (w *Widget) SetTransformerFunc(f orm.TransformerFunc) *Widget {
	w.transformerFunc = &f
	return w
}

// 设置transformerFunc， 参数为实现了TransformerIf接口的对象
func (w *Widget) SetTransformer(i orm.TransformerIf) *Widget {
	w.transformerFunc = orm.GetTransformerFunc(i)
	return w
}

// 使用p.transformerFunc
func (w *Widget) transformerRun(result OneResult) OneResult {
	if w.transformerFunc != nil {
		result = (*w.transformerFunc)(result)
	}
	return result
}

// 把widget中的Query应用到widget中
func (this *Widget) ApplyQuery() {
	this.Match = bson.M{}
	this.sort = bson.D{}
	//------- 对offset、 limit、page三个参数进行处理，得到最终的limit和offset --------
	//limit必须有，否则取默认值
	if this.Query.HasLimit() == false {
		this.Query.SetLimit(orm.DefaultPageSize)
	}
	//page可以有，有的话，通过它计算额外offset
	if this.Query.HasPage() {
		offset := this.Query.GetLimit() * (this.Query.GetPage() - 1)
		if this.Query.HasOffset() {
			this.Query.SetOffset(this.Query.GetOffset() + offset)
		} else {
			this.Query.SetOffset(offset)
		}
	}

	//condition
	if this.Query.HasCondition() {
		//fill out the conditions object, we support "and"、"or"

		conditionPtrs := this.Query.GetCondition()

		conditions := []orm.Condition{}
		andConditions := []orm.Condition{}
		orConditions := []orm.Condition{}
		for _, conditionPtr := range conditionPtrs {
			if conditionPtr != nil {
				condition := *conditionPtr
				switch condition.GetType() {
				case orm.CTYPE_AND:
					andConditions = append(andConditions, condition)
				case orm.CTYPE_OR:
					orConditions = append(orConditions, condition)
				case orm.CTYPE_IN:
					andConditions = append(andConditions, condition)
				}
			}
		}
		conditions = append(andConditions, orConditions...)
		//change a condition object to a real ORM’s function's params
		//for example, gorm has functions like 'where'、'or', they both need params
		for _, condition := range conditions {
			operator := this.getOperator(condition.GetOperator())
			if operator == "" {
				continue
			}
			switch operator {
			case MG_OPERATOR_IS_EQUAL:
				this.Match[condition.GetField()] = bson.M{"$eq": condition.GetValue()}
			case MG_OPERATOR_IS_NOT_EQUAL:
				this.Match[condition.GetField()] = bson.M{"$ne": condition.GetValue()}
			case MG_OPERATOR_IS_GREATER_THAN:
				if val, has := this.Match[condition.GetField()]; has {
					val.(bson.M)["$gt"] = condition.GetValue()
				} else {
					this.Match[condition.GetField()] = bson.M{"$gt": condition.GetValue()}
				}
			case MG_OPERATOR_IS_GREATER_THAN_OR_EQUAL:
				if val, has := this.Match[condition.GetField()]; has {
					val.(bson.M)["$gte"] = condition.GetValue()
				} else {
					this.Match[condition.GetField()] = bson.M{"$gte": condition.GetValue()}
				}
			case MG_OPERATOR_IS_LESS_THAN:
				if val, has := this.Match[condition.GetField()]; has {
					val.(bson.M)["$lt"] = condition.GetValue()
				} else {
					this.Match[condition.GetField()] = bson.M{"$lt": condition.GetValue()}
				}
			case MG_OPERATOR_IS_LESS_THAN_OR_EQUAL:
				if val, has := this.Match[condition.GetField()]; has {
					val.(bson.M)["$lte"] = condition.GetValue()
				} else {
					this.Match[condition.GetField()] = bson.M{"$lte": condition.GetValue()}
				}
			case MG_OPERATOR_IS_IN:
				this.Match[condition.GetField()] = bson.M{"$in": condition.GetValue()}
			case MG_OPERATOR_IS_LIKE:
				this.Match[condition.GetField()] = bson.M{"$regex": primitive.Regex{condition.GetValueString(), "."}}
			case MG_OPERATOR_EXISTS:
				this.Match[condition.GetField()] = bson.M{"$exists": condition.GetValue()}
			default:
			}
		}
	}

	//sort
	if this.Query.HasSorter() {
		sorters := this.Query.GetSorter()
		for _, sorter := range sorters {
			switch sorter.GetDirection() {
			case orm.ASCENDING:
				tempSort := bson.E{
					sorter.GetField(),
					1,
				}
				this.sort = append(this.sort, tempSort)
			default:
				tempSort := bson.E{
					sorter.GetField(),
					-1,
				}
				this.sort = append(this.sort, tempSort)
			}
		}
	}
}

const MG_OPERATOR_IS_EQUAL = "="
const MG_OPERATOR_IS_GREATER_THAN = ">"
const MG_OPERATOR_IS_GREATER_THAN_OR_EQUAL = ">="
const MG_OPERATOR_IS_LESS_THAN = "<"
const MG_OPERATOR_IS_LESS_THAN_OR_EQUAL = "<="
const MG_OPERATOR_IS_IN = "IN"
const MG_OPERATOR_IS_NOT_IN = "NOT IN"
const MG_OPERATOR_IS_LIKE = "LIKE"
const MG_OPERATOR_IS_NOT_LIKE = "NOT LIKE"
const MG_OPERATOR_IS_JSON_CONTAINS = "JSON_CONTAINS"
const MG_OPERATOR_IS_NOT_EQUAL = "<>"
const MG_OPERATOR_IS_IS_NULL = "IS NULL"
const MG_OPERATOR_IS_IS_NOT_NULL = "IS NOT NULL"
const MG_OPERATOR_EXISTS = "EXIST"
const MG_DEFAULT_KEY = "value"

func (this *Widget) getOperator(operator int) string {
	m := this.operatorMap()
	return m[operator]
}

func (this *Widget) operatorMap() map[int]string {
	return map[int]string{
		orm.Operator_is_equal:            MG_OPERATOR_IS_EQUAL,
		orm.OperatorIsGreaterThan:        MG_OPERATOR_IS_GREATER_THAN,
		orm.OperatorIsGreaterThanOrEqual: MG_OPERATOR_IS_GREATER_THAN_OR_EQUAL,
		orm.OperatorIsIn:                 MG_OPERATOR_IS_IN,
		orm.OperatorIsNotIn:              MG_OPERATOR_IS_NOT_IN,
		orm.OperatorIsLessThan:           MG_OPERATOR_IS_LESS_THAN,
		orm.OperatorIsLessThanOrEqual:    MG_OPERATOR_IS_LESS_THAN_OR_EQUAL,
		orm.OperatorIsLike:               MG_OPERATOR_IS_LIKE,
		orm.OperatorIsNotLike:            MG_OPERATOR_IS_NOT_LIKE,
		orm.OperatorIsJsonContains:       MG_OPERATOR_IS_JSON_CONTAINS,
		orm.OperatorIsNotEqual:           MG_OPERATOR_IS_NOT_EQUAL,
		orm.OperatorIsIsNull:             MG_OPERATOR_IS_IS_NULL,
		orm.OperatorIsIsNotNull:          MG_OPERATOR_IS_IS_NOT_NULL,
		orm.OperatorExists:               MG_OPERATOR_EXISTS,
	}
}

// ---------------------------------------------

func (this *Widget) SetLimit(limit int) *Widget {
	this.Query.SetLimit(limit)
	return this
}

func (this *Widget) AddSorter(fieldName string, i int) {
	direction := orm.ASCENDING
	if i < 0 {
		direction = orm.DESCENDING
	}

	sort := orm.NewSorter(fieldName, direction)
	this.Query.AddSorter(*sort)
}

// ---------------------------------------------

// 分页查询，基于Find，得到一组数据
func (this *Widget) PaginatorFind() (qmap.QM, error) {
	all, err := this.Find()
	if err != nil {
		return nil, nil
	}

	total, err := this.collection.CountDocuments(context.Background(), this.Match)
	if err == nil {
		pages := orm.Paginator(this.Query.GetPage(), this.Query.GetLimit(), int64(total))
		pages["count"] = len(all)

		qm := qmap.QM{
			"list":       all,
			"pagination": pages,
		}
		return qm, nil
	} else {
		return nil, err
	}
}

// 查询结果为[]map[string]interface{}
// 会使用transformer更新数据
func (this *Widget) Find() (AllResult, error) {
	this.ApplyQuery()
	operations := []bson.M{}
	if len(this.Match) > 0 {
		operations = append(operations, bson.M{"$match": this.Match})
	}
	if len(this.sort) > 0 {
		operations = append(operations, bson.M{"$sort": this.sort})
	} else {
		operations = append(operations, bson.M{"$sort": bson.D{bson.E{"_id", -1}}})
	}
	operations = append(operations, bson.M{"$skip": this.Query.GetOffset()})
	operations = append(operations, bson.M{"$limit": this.Query.GetLimit()})
	result := AllResult{}

	ct := context.Background()
	c, err := this.collection.Aggregate(ct, operations)
	if err != nil {
		return result, err
	}

	if err := c.All(ct, &result); err != nil {
		return result, err
	}

	for k, one := range result {
		result[k] = this.transformerRun(one)
	}
	return result, nil
}

// 查询结果为map[string]interface{}
// 会使用transformer更新数据
func (this *Widget) Get() (qmap.QM, error) {
	this.ApplyQuery()

	result := qmap.QM{}

	err := this.collection.FindOne(context.Background(), this.Match).Decode(&result)
	if err == nil {
		result = this.transformerRun(result)
		return result, nil
	} else {
		return nil, err
	}
}

func (this *Widget) Count() (int64, error) {
	this.ApplyQuery()
	return this.collection.CountDocuments(context.Background(), this.Match)
}

// 查询单条记录，不会使用transformer
func (this *Widget) One(modelPtr interface{}) error {
	this.ApplyQuery()

	err := this.collection.FindOne(context.Background(), this.Match).Decode(modelPtr)
	return err
}

// 查询多条记录，不会使用transformer
func (this *Widget) All(modelPtr interface{}) error {
	this.ApplyQuery()

	cur, err := this.collection.Find(context.Background(), this.Match)
	if err != nil {
		return err
	}

	return cur.All(context.Background(), modelPtr)
}
