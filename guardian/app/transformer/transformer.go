package transformer

import (
	"reflect"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/guardian/src/orm"
	"skygo_detection/guardian/util"
)

//type TransformerFunc func(interface{}) interface{}
//type TransformerFunc func(map[string]interface{}) map[string]interface{}

type Transformer interface {
	GetOrmTransformer(*orm.Query, Transformer) *orm.Transformer
	AdditionalFields(qm qmap.QM) qmap.QM //添加额外的字段
	ExcludeFields() []string             //剔除指定的字段
}

type BaseTransformer struct{}

func (this *BaseTransformer) GetOrmTransformer(q *orm.Query, t Transformer) *orm.Transformer {
	//if user's request URL has params "include", we can find it in Query object
	includeFuncValue := map[string]reflect.Value{} //key is function name, value is a reflect.Value
	if q.HasInclude() {
		include := q.GetInclude()
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

func (this *BaseTransformer) AdditionalFields(qm qmap.QM) qmap.QM {
	return qm
}

func (this *BaseTransformer) ExcludeFields() []string {
	return nil
}
