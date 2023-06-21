package mysql_model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	// "skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
)

type KnowledgeTestCaseHistory struct {
	Id         int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TestCaseId int    `xorm:"not null comment('测试用例真实id') INT(11)" json:"test_case_id"`
	Version    string `xorm:"not null comment('版本') VARCHAR(255)" json:"version"`
	Content    string `xorm:"not null comment('内容') VARCHAR(255)" json:"content"`
	OPId       int    `xorm:"not null comment('用户id') INT(11)" json:"op_id"`
	TimeStamp  int    `xorm:"not null comment('时间戳') INT(11)" json:"time_stamp"`
}

func (this *KnowledgeTestCaseHistory) GetAllByCaseId(testCaseId int) []KnowledgeTestCaseHistory {
	testCases := []KnowledgeTestCaseHistory{}
	mysql.GetSession().Where("test_case_id = ?", testCaseId).Find(&testCases)
	return testCases
}

func calculateVersion(version string) string {
	if version == "" {
		return "v1.0"
	}
	if strings.Contains(version, "v") {
		ver := strings.TrimLeft(version, "v")
		if vver, err := strconv.ParseFloat(ver, 64); err == nil {
			vver += 0.1
			return fmt.Sprintf("v%0.1f", vver)
		}
	}
	return version
}

func (this *KnowledgeTestCaseHistory) FindLastVersion(testCaseId int) (string, bool) {
	session := mysql.GetSession()
	session.Where("test_case_id = ?", testCaseId)
	session.OrderBy("time_stamp desc")
	has, _ := session.Get(this)
	return this.Version, has
}

func (this *KnowledgeTestCaseHistory) Create() error {
	now := int(time.Now().Unix())
	this.Version = calculateVersion(this.Version)
	this.TimeStamp = now
	_, err := mysql.GetSession().InsertOne(this)
	return err
}
