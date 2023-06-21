package mongo_model

import (
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/mongo_model_tmp"
)

type EvaluateVulScanner struct {
	ID            bson.ObjectId   `bson:"_id,omitempty" json:"_id"`
	TaskID        string          `bson:"task_id" json:"task_id"`               //任务id
	DeviceInfos   DeviceInfo      `bson:"device_infos" json:"device_infos"`     //设备配置信息
	ScanResult    []ScannerResult `bson:"scan_result" json:"scan_result"`       //漏洞扫描结果
	SearchContent string          `bson:"search_content" json:"search_content"` //
}

type DeviceInfo struct {
	Company    string `bson:"company" json:"company"`         //车机厂商
	SysVersion string `bson:"sys_version" json:"sys_version"` //系统版本
	CPUMode    string `bson:"cpu_mode" json:"cpu_mode"`       //芯片型号
	CPUVersion string `bson:"cpu_version" json:"cpu_version"` //芯片版本
	Platform   string `bson:"platform" json:"platform"`       //平台
	SysSdkVer  string `bson:"sys_sdk_ver" json:"sys_sdk_ver"` //sdk版本
	CarMode    string `bson:"car_mode" json:"car_mode"`       //车机型号
	Brand      string `bson:"brand" json:"brand"`             //车机品牌
}

type ScannerResult struct {
	CveId               string `bson:"cve_id" json:"cve_id"`                               //cve ID编号
	GoogleSeverityLevel string `bson:"google_severity_level" json:"google_severity_level"` //漏洞级别
	DateExposure        string `bson:"date_exposure" json:"date_exposure"`                 //披露时间
	DateBulletin        string `bson:"date_bulletin" json:"date_bulletin"`                 //发布时间
	Sketch              string `bson:"sketch" json:"sketch"`                               //漏洞简述
	Description         string `bson:"description" json:"description"`                     //漏洞详细描述
}

func (this *EvaluateVulScanner) Create(rawInfo qmap.QM) (*EvaluateVulScanner, error) {
	//查询 任务库里有没有 此任务id
	taskId := rawInfo.MustString("task_id")
	vulTask := new(mongo_model_tmp.EvaluateVulTask)
	params := qmap.QM{
		"e_task_id": taskId,
	}
	var isTaskLinkData bool
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TASK, params).One(&vulTask); err != nil {
		//会去判断工具库有没有这个测试任务
		if !new(ToolTaskResultBindTest).ScriptIsSetToolTask(taskId) {
			return nil, err
		} else {
			isTaskLinkData = true
		}
	}
	//存在这个任务ID，把漏洞扫描结果入库
	err := this.create(rawInfo)
	if err != nil {
		return nil, err
	}
	//存在这个任务ID，把漏洞结果 更新到 任务库里
	//不存在这个任务ID，但是项目中的工具存在，就忽略更新这个任务状态
	if !isTaskLinkData {
		vulTask.VulScannerID = this.ID.Hex()
		err = this.uploadVulTask(vulTask)
		if err != nil {
			return nil, err
		}
	}

	//把数据分拆 写入多个库里去
	//删除多余的库
	//删除存在的信息
	deviceInfo1 := new(EvaluateVulDeviceInfo)
	deviceInfo := rawInfo["device_infos"]
	tmpDeviceInfo := deviceInfo.(map[string]interface{})
	tmpDeviceInfo["task_id"] = taskId
	if _, err := deviceInfo1.GetOne(taskId); err == nil {
		deviceInfo1.BulkDelete([]string{taskId})
	}
	deviceInfo1.Create(tmpDeviceInfo)

	//删除存在的内容
	vulInfo := new(EvaluateVulInfo)
	if _, err := vulInfo.GetOne(taskId); err == nil {
		vulInfo.BulkDelete([]string{taskId})
	}
	scanResult := rawInfo["scan_result"]
	for _, result := range scanResult.([]interface{}) {
		tmpRsult := qmap.QM(result.(map[string]interface{}))
		tmpRsult["task_id"] = taskId
		if google := tmpRsult.Int("google_severity_level"); google > 3 || google < 1 {
			continue
		}
		new(EvaluateVulInfo).Create(tmpRsult)
	}
	//通知脚本任务，该漏洞已完成
	new(ToolTaskResultBindTest).ScriptUpdateToolTaskStatus(taskId)
	return this, nil
}

func (this *EvaluateVulScanner) create(rawInfo qmap.QM) error {
	taskId := rawInfo.MustString("task_id")
	params := qmap.QM{
		"e_task_id": taskId,
	}
	ormSession := mongo.NewMgoSession(common.MC_EVALUATE_VUL_SCANNER)
	if err := ormSession.AddCondition(params).One(this); err != nil {
		if err := ormSession.Insert(rawInfo); err == nil {
			return nil
		} else {
			return err
		}
	} else {
		selector := qmap.QM{
			"_id": this.ID,
		}
		if err := ormSession.Update(selector, rawInfo); err == nil {
			ormSession.AddCondition(params).One(this)
			return nil
		} else {
			return err
		}
	}

}
func (this *EvaluateVulScanner) uploadVulTask(vulTask *mongo_model_tmp.EvaluateVulTask) error {
	selector := qmap.QM{
		"_id": vulTask.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"vul_scanner_id": vulTask.VulScannerID,
			"test_time":      int(custom_util.GetCurrentMilliSecond()),
			"status":         common.VUL_PRELIMINARY_END,
		},
	}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_TASK).Update(selector, update); err == nil {
		return nil
	} else {
		return err
	}
}
