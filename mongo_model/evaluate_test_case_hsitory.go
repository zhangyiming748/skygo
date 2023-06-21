package mongo_model

import (
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type TestCaseHistoryContent struct {
	Id         bson.ObjectId `bson:"_id,omitempty"  json:"_id"`         //测试用例修改记录id
	TestCaseId string        `bson:"test_case_id"  json:"test_case_id"` //测试用例id
	Version    string        `bson:"version"  json:"version"`           //版本
	Content    string        `bson:"content"  json:"content"`           //内容
	OPId       int           `bson:"op_id"  json:"op_id"`               //用户id
	UserName   string        `bson:"user_name"  json:"user_name"`       //用户名称
	TimeStamp  int           `bson:"timestamp"  json:"timestamp"`       //时间戳
}

func (this *TestCaseHistoryContent) Create(rawInfo qmap.QM) (*TestCaseHistoryContent, error) {
	this.Id = bson.NewObjectId()
	this.TestCaseId = rawInfo.String("test_case_id")
	this.Version = rawInfo.String("content_version")
	this.Content = rawInfo.String("content")
	this.OPId = rawInfo.Int("op_id")
	this.UserName = rawInfo.String("user_name")
	this.TimeStamp = rawInfo.Int("timestamp")
	if err := mongo.NewMgoSession(common.MC_TEST_CASE_CONTENT).Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *TestCaseHistoryContent) Get(rawInfo qmap.QM) (*[]map[string]interface{}, error) {
	testCaseId := rawInfo.MustString("test_case_id")
	params := qmap.QM{
		"e_test_case_id": testCaseId,
	}
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_TEST_CASE_CONTENT, params)
	return mgoSession.SetLimit(10000).Get()
}

func (this *TestCaseHistoryContent) BulkDelete(rawIds []string) error {
	if len(rawIds) > 0 {
		match := bson.M{
			"test_case_id": bson.M{"$in": rawIds},
		}
		_, err := mongo.NewMgoSession(common.MC_TEST_CASE_CONTENT).RemoveAll(match)
		return err
	}
	return nil
}
