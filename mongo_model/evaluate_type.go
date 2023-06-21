package mongo_model

import (
	"errors"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type EvaluateType struct {
	Id    bson.ObjectId       `bson:"_id,omitempty"`
	Name  string              `bson:"name"`  //名称
	Attrs []*EvaluateTypeAttr `bson:"attrs"` //测试项属性
}

type EvaluateTypeAttr struct {
	AttrName   string `bson:"attr_name"`   //属性名称
	AttrKey    string `bson:"attr_key"`    //属性关键字
	AttrType   string `bson:"attr_type"`   //属性类型(字符串:string,整形:int, 浮点型:float, 日期:date, 布尔型:bool)
	IsRequired int    `bson:"is_required"` //是否必填(0:否，1:是)
}

func (this *EvaluateType) Create(rawInfo *qmap.QM) (*EvaluateType, error) {
	this.Id = bson.NewObjectId()
	this.Name = rawInfo.MustString("name")
	if attrs, has := rawInfo.TrySlice("attrs"); has {
		for _, item := range attrs {
			var attrJSON qmap.QM = item.(map[string]interface{})
			attr := EvaluateTypeAttr{
				AttrName:   attrJSON.String("attr_name"),
				AttrKey:    attrJSON.String("attr_key"),
				AttrType:   attrJSON.String("attr_type"),
				IsRequired: attrJSON.Int("is_required"),
			}
			this.Attrs = append(this.Attrs, &attr)
		}
	}
	mongoClient := mongo.NewMgoSession(common.MC_EVALUATE_TYPE)
	if err := mongoClient.Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateType) Update(id string, rawInfo qmap.QM) (*EvaluateType, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TYPE, params)
	if err := mongoClient.One(&this); err == nil {
		if val, has := rawInfo.TryString("name"); has {
			this.Name = val
		}
		if attrs, has := rawInfo.TrySlice("attrs"); has {
			var tempAttrs []*EvaluateTypeAttr
			for _, item := range attrs {
				var attrJSON qmap.QM = item.(map[string]interface{})
				attr := EvaluateTypeAttr{
					AttrName:   attrJSON.String("attr_name"),
					AttrKey:    attrJSON.String("attr_key"),
					AttrType:   attrJSON.String("attr_type"),
					IsRequired: attrJSON.Int("is_required"),
				}
				tempAttrs = append(tempAttrs, &attr)
			}
			this.Attrs = tempAttrs
		}
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Item not found")
	}
	return this, nil
}

func (this *EvaluateType) GetOne(id, name string) (*EvaluateType, error) {
	params := qmap.QM{}
	if id != "" {
		params["e__id"] = bson.ObjectIdHex(id)
	}
	if name != "" {
		params["e_name"] = name
	}
	if len(params) == 0 {
		return nil, errors.New("Provide at least one id and name!")
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TYPE, params).One(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

// 提取测试类型的属性值
func (this *EvaluateType) ExtraAttributeMap(rawInfo qmap.QM) (qmap.QM, error) {
	attrMaps := qmap.QM{}
	for _, attr := range this.Attrs {
		var exist bool
		var val interface{}
		switch attr.AttrType {
		case "string":
			val, exist = rawInfo.TryString(attr.AttrKey)
		case "int", "date":
			val, exist = rawInfo.TryInt(attr.AttrKey)
		case "float":
			val, exist = rawInfo.TryFloat64(attr.AttrKey)
		case "bool":
			val, exist = rawInfo.TryBool(attr.AttrKey)
		default:
			return nil, errors.New("Unknown attr type!")
		}
		if attr.IsRequired == 1 && exist == false {
			return nil, errors.New(fmt.Sprintf("attribute of %s is required to the evaluate type of %s", attr.AttrKey, this.Name))
		}
		attrMaps[attr.AttrKey] = val
	}
	return attrMaps, nil
}
