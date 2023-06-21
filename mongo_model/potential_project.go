package mongo_model

import (
	"errors"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
)

type PotentialProject struct {
	Id                    bson.ObjectId `bson:"_id,omitempty"`
	Name                  string        `bson:"name"`                   //潜在项目名称
	Company               string        `bson:"company"`                //潜在项目所属车厂
	Description           string        `bson:"description"`            //潜在项目描述
	Budget                int           `bson:"budget"`                 //客户预算(单位:元)
	Type                  int           `bson:"type"`                   //潜在项目类型(战略项目、正常项目)
	AcceptanceProbability int           `bson:"acceptance_probability"` //中标概率(低、中、高)
	Status                int           `bson:"status"`                 //潜在项目状态（已立项、未立项）
	UpdateTime            int64         `bson:"update_time"`
	CreateTime            int64         `bson:"create_time"`
	CreateUid             int           `bson:"create_uid"`        //创建者ID
	CreateName            string        `bson:"create_name"`       //创建者
	ProManagerUid         int           `bson:"pro_manager_uid"`   //项目经理UID
	ProManagerUname       string        `bson:"pro_manager_uname"` //项目经理UID
}

func (this *PotentialProject) Create(rawInfo *qmap.QM) (*PotentialProject, error) {
	this.Id = bson.NewObjectId()
	this.Name = rawInfo.String("name")
	this.Company = rawInfo.String("company")
	this.Description = rawInfo.String("description")
	this.Budget = rawInfo.Int("budget")
	this.CreateUid = rawInfo.Int("user_id")
	this.CreateName = rawInfo.String("user_name")
	this.Type = rawInfo.Int("type")
	this.AcceptanceProbability = rawInfo.Int("acceptance_probability")
	this.Status = rawInfo.Int("status")
	this.ProManagerUid = rawInfo.Int("pro_manager_uid")
	this.ProManagerUname = rawInfo.String("pro_manager_uname")
	this.UpdateTime = custom_util.GetCurrentMilliSecond()
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	mongoClient := mongo.NewMgoSession(common.MC_POTENTIAL_PROJECT)
	if err := mongoClient.Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *PotentialProject) Update(potentialProjectId string, rawInfo qmap.QM) (*PotentialProject, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(potentialProjectId),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_POTENTIAL_PROJECT, params)
	if err := mongoClient.One(&this); err == nil {
		if val, has := rawInfo.TryString("name"); has {
			this.Name = val
		}
		if val, has := rawInfo.TryString("company"); has {
			this.Company = val
		}
		if val, has := rawInfo.TryString("description"); has {
			this.Description = val
		}
		if val, has := rawInfo.TryInt("budget"); has {
			this.Budget = val
		}
		if val, has := rawInfo.TryInt("type"); has {
			this.Type = val
		}
		if val, has := rawInfo.TryInt("acceptance_probability"); has {
			this.AcceptanceProbability = val
		}
		if val, has := rawInfo.TryInt("status"); has {
			this.Status = val
		}
		if val, has := rawInfo.TryInt("user_id"); has {
			this.CreateUid = val
		}
		if val, has := rawInfo.TryString("user_name"); has {
			this.CreateName = val
		}
		if val, has := rawInfo.TryInt("pro_manager_uid"); has {
			this.ProManagerUid = val
		}
		if val, has := rawInfo.TryString("pro_manager_uname"); has {
			this.ProManagerUname = val
		}
		this.UpdateTime = custom_util.GetCurrentMilliSecond()
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Potential Project not found")
	}
	return this, nil
}
