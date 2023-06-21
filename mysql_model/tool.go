package mysql_model

import "skygo_detection/lib/common_lib/mysql"

type Tool struct {
	Id            int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	ToolNumber    int    `xorm:"not null comment('主键id') INT(11)" json:"tool_number"`
	Name          string `xorm:"not null default '' comment('工具名称') VARCHAR(255)" json:"name"`
	Brand         string `xorm:"not null default '' comment('品牌') VARCHAR(255)" json:"brand"`
	CategoryName  string `xorm:"not null default '' comment('类型名称') VARCHAR(255)" json:"category_name"`
	CategoryID    int    `xorm:"not null comment('工具分类id，暂时不用') INT(11)" json:"category_id"`
	TestPic       string `xorm:"not null default '' comment('测试工具图片,xxx|xxx|xxx') VARCHAR(255)" json:"test_pic"`
	ToolDetail    string `xorm:"not null default '' comment('工具介绍') VARCHAR(255)" json:"tool_detail"`
	UseDetail     string `xorm:"not null default '' comment('使用方法') VARCHAR(255)" json:"use_detail"`
	UseManualLink string `xorm:"not null default '' comment('使用手册连接') VARCHAR(255)" json:"use_manual_link"`
	UseManual     string `xorm:"not null default '' comment('使用手册文件id,xxx|xxx|xxx') VARCHAR(255)" json:"use_manual"`
	LinkPic       string `xorm:"not null default '' comment('工具连接示意图id，xxx|xxx') VARCHAR(255)" json:"link_pic"`
	Script        string `xorm:"not null default '' comment('工具脚本id，xxx|xxx') VARCHAR(255)" json:"script"`
	ParamsJson    string `xorm:"not null default '' comment('工具配置') VARCHAR(255)" json:"params_json"`
	CreateTime    int    `xorm:"not null default 0 comment('创建时间（秒）') INT(11)" json:"create_time"`
	UpdateTime    int    `xorm:"not null default 0 comment('更新时间') INT(11)" json:"update_time"`
	//工具版本
	ToolVersion     string `xorm:"not null default '' comment('版本型号') VARCHAR(255)" json:"tool_version"`
	SoftwareVersion string `xorm:"not null default '' comment('软件版本') VARCHAR(255)" json:"software_version"`
	HardwareVersion string `xorm:"not null default '' comment('硬件版本') VARCHAR(255)" json:"hardware_version"`
	SystemVersion   string `xorm:"not null default '' comment('系统版本') VARCHAR(255)" json:"system_version"`
}

var CategoryNameList = map[int]string{
	1: "硬件安全",
	2: "系统安全",
	3: "应用安全",
	4: "无线电安全",
	5: "车载网络安全",
	6: "代码安全",
	7: "固件安全",
	8: "云端安全",
}

// 增删改查
func (this *Tool) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *Tool) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (this *Tool) Remove() (int64, error) {
	return mysql.GetSession().ID(this.Id).Delete(this)
}
