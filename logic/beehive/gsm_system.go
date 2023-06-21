package beehive

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/blog"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/mysql_model"
	"strconv"
	"strings"
)

type GsmSystemStandResp struct {
	Status    bool   `json:"status"`
	MessageId int    `json:"message_id"`
	Message   string `json:"message"`
}

// 设置配置前或开启设备前需要配置防火墙
func gsmSystemBeforeStart() (GsmSystemStandResp, error) {
	args := make(map[string]string)
	args["iface"] = common.IFACE
	host, _ := GetHost(GsmSystem)
	resp, err := custom_util.HttpPostJson(nil, args, strings.Join([]string{host, common.BEFORE_START}, "/"))
	blog.Info(common.BEFORE_START, zap.Any("request:", args), zap.Any("responds:", string(resp)))
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	var gsmSystemStandResp GsmSystemStandResp
	err = json.Unmarshal(resp, &gsmSystemStandResp)
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	return gsmSystemStandResp, nil
}

// 设置启动前设备参数-设备
func gsmSystemSetConfig2Device(cid int) (GsmSystemStandResp, error) {
	args := make(map[string]int)
	args["id"] = cid
	host, _ := GetHost(GsmSystem)
	resp, err := custom_util.HttpPostJson(nil, args, strings.Join([]string{host, common.CONFIG_SYSTEM}, "/"))
	blog.Info(common.CONFIG_SYSTEM, zap.Any("request:", args), zap.Any("responds:", string(resp)))
	if err != nil {
		log.Println(err)
		return GsmSystemStandResp{}, err
	}
	var gsmSystemStandResp GsmSystemStandResp
	log.Println(string(resp))
	err = json.Unmarshal(resp, &gsmSystemStandResp)
	if err != nil {
		log.Println(err)
		return GsmSystemStandResp{}, err
	}
	blog.Info(common.CONFIG_SYSTEM, zap.Any("gsmSystemStandResp:", gsmSystemStandResp))
	return gsmSystemStandResp, err
}

// 设置启动前设备参数-日志
func gsmSystemSetConfig2Log(tid, cid, mid int) error {
	task := new(mysql_model.BeehiveLog)
	task.Title = strings.Join([]string{"尝试配置设备启动前参数", strconv.Itoa(cid)}, ":")
	switch mid {
	case 0:
		task.Content = "未知失败,需要查阅日志"
	case 1:
		task.Content = "修改配置成功，请手动启动系统"
	case 2:
		task.Content = "关闭系统失败，需要手动请求stop api关闭"
	case 3:
		task.Content = "无法找到数据库文件，请检查/etc/OpenBTS/OpenBTS.db文件"
	}
	task.TaskId = tid
	err := task.SetLog()
	if err != nil {
		return err
	}
	return nil
}

// 设置启动前设备参数-任务表
func GsmSystemSetConfig2Task(tid, cid int) (int64, error) {
	if mysql_model.GsmTaskNotExist(tid) {
		mysql_model.CreateGsmTask(tid)
	}
	task := new(mysql_model.BeehiveGsmSystem)
	task.TaskId = tid
	task.ConfigId = cid
	return task.ForceUpdateConfigId()
}

// 设置启动前参数
func GsmSystemSetConfig(tid, cid int) (GsmSystemStandResp, error) {
	// 设置硬件启动前配置
	prefix, err := gsmSystemBeforeStart()
	if err != nil && prefix.Status {
		log.Printf("设置硬件启动前配置错误%v\n", err)
		return GsmSystemStandResp{}, err
	}
	// 更新对应的任务表
	_, err = GsmSystemSetConfig2Task(tid, cid)
	if err != nil {
		log.Printf("更新对应的任务表错误%v\n", err)
		return GsmSystemStandResp{}, err
	}
	// 配置硬件参数
	resp, err := gsmSystemSetConfig2Device(cid)
	if err != nil && !resp.Status {
		log.Printf("配置硬件参数错误%v\n", err)
		return GsmSystemStandResp{}, err
	}
	// 事件写入日志
	err = gsmSystemSetConfig2Log(tid, cid, resp.MessageId)
	if err != nil {
		log.Printf("事件写入日志错误%v\n", err)
		return GsmSystemStandResp{}, err
	}
	return resp, err
}

// 启动系统-设备
func gsmSystemStart2Device() (GsmSystemStandResp, error) {
	args := make(map[string]string)
	host, _ := GetHost(GsmSystem)
	resp, err := custom_util.HttpPostJson(nil, args, strings.Join([]string{host, common.START_2G_SYSTEM}, "/"))
	blog.Info(common.START_2G_SYSTEM, zap.Any("request:", args), zap.Any("responds:", string(resp)))
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	var gsmSystemStandResp GsmSystemStandResp
	err = json.Unmarshal(resp, &gsmSystemStandResp)
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	return gsmSystemStandResp, err
}

// 启动系统-任务表
func gsmSystemStart2Task(tid, mid int) (int64, error) {
	task := new(mysql_model.BeehiveGsmSystem)
	task.TaskId = tid
	if mid == 1 {
		task.SystemStatus = 1
	}
	return task.ForceUpdateSystemStatus()
}

// 启动系统-日志
func gsmSystemStart2Log(tid, mid int) error {
	task := new(mysql_model.BeehiveLog)
	task.TaskId = tid
	task.Title = "尝试启动系统"
	switch mid {
	case common.START_FAILED:
		task.Content = "未知启动失败,需要查阅日志"
	case common.START_SUCCESS:
		task.Content = "启动成功"
	case common.ISRUNNING:
		task.Content = "设备正在运行中"
	case common.DEVICE_NOT_CONNECT:
		task.Content = "未连接USRP,请接入设备后在尝试"
	}
	return task.SetLog()
}

// 启动系统
func GsmSystemStart(tid int) (GsmSystemStandResp, error) {
	prefix, err := gsmSystemBeforeStart()
	if err != nil && !prefix.Status {
		return GsmSystemStandResp{}, err
	}
	start, err := gsmSystemStart2Device()
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	_, err = gsmSystemStart2Task(tid, start.MessageId)
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	err = gsmSystemStart2Log(tid, start.MessageId)
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	return start, err
}

// 关闭系统-设备
func gsmSystemStop2Device() (GsmSystemStandResp, error) {
	args := make(map[string]string)
	host, _ := GetHost(GsmSystem)
	resp, err := custom_util.HttpPostJson(nil, args, strings.Join([]string{host, common.STOP_2G_SYSTEM}, "/"))
	blog.Info(common.STOP_2G_SYSTEM, zap.Any("request:", args), zap.Any("responds:", string(resp)))
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	var gsmSystemStandResp GsmSystemStandResp

	err = json.Unmarshal(resp, &gsmSystemStandResp)
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	return gsmSystemStandResp, err
}

// 关闭系统-任务
func gsmSystemStop2Task(tid, mid int) (int64, error) {
	task := new(mysql_model.BeehiveGsmSystem)
	task.TaskId = tid
	if mid == 1 || mid == 0 {
		task.SystemStatus = 2
	}
	return task.Update()
}

// 关闭系统-日志
func gsmSystemStop2Log(tid, mid int) error {
	task := new(mysql_model.BeehiveLog)
	task.TaskId = tid
	task.Title = "尝试关闭系统"
	switch mid {
	case common.STOP_FAILED:
		task.Content = "关闭设备失败,需要查阅日志"
	case common.STOP_SUCCESS:
		task.Content = "关闭设备成功"
	case common.NOT_RUNNING:
		task.Content = "程序未在运行在,无需关闭"
	}
	return task.SetLog()
}

// 关闭系统
func GsmSystemStop(tid int) (GsmSystemStandResp, error) {
	stop, err := gsmSystemStop2Device()
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	_, err = gsmSystemStop2Task(tid, stop.MessageId)
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	err = gsmSystemStop2Log(tid, stop.MessageId)
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	return stop, err
}

// 获取短信结构体
type GsmSystemSmsInfo struct {
	Status    bool       `json:"status"`
	MessageId int        `json:"message_id"`
	Message   string     `json:"message"`
	Infos     [][]string `json:"infos"`
}

// 获取短信-设备
func gsmSystemGetSms2Device() (GsmSystemSmsInfo, error) {
	args := make(map[string]string)
	host, _ := GetHost(GsmSystem)
	resp, err := custom_util.HttpPostJson(nil, args, strings.Join([]string{host, common.GET_SMS}, "/"))
	blog.Info(common.GET_SMS, zap.Any("request:", args), zap.Any("responds:", string(resp)))
	if err != nil {
		return GsmSystemSmsInfo{}, err
	}
	var gsmSystemSmsInfo GsmSystemSmsInfo
	err = json.Unmarshal(resp, &gsmSystemSmsInfo)
	if err != nil {
		return GsmSystemSmsInfo{}, err
	}
	return gsmSystemSmsInfo, err
}

// 获取短信-任务表
func gsmSystemGetSms2Task(tid int) (int64, error) {
	this := new(mysql_model.BeehiveGsmSystem)
	this.TaskId = tid
	return this.Update()
}

// 获取短信-日志
func gsmSystemGetSms2Log(tid int, s GsmSystemSmsInfo) {
	this := new(mysql_model.BeehiveLog)
	this.TaskId = tid
	this.Title = "尝试获取短信"
	this.Content = fmt.Sprintf("获取到%d条短信", len(s.Infos))
	//switch s.MessageId {
	//case 0:
	//	this.Content = strings.Join([]string{content, "未知失败,需要查阅日志"}, ":")
	//case 1:
	//	this.Content = strings.Join([]string{content, "读取成功"}, ":")
	//case 2:
	//	this.Content = strings.Join([]string{content, "无法找到日志文件，请检查/var/log/syslog文件"}, ":")
	//case 3:
	//	this.Content = strings.Join([]string{content, "未在日志文件中找到短信信息"}, ":")
	//}
	err := this.SetLog()
	if err != nil {
		return
	}
}

// 获取短信-保存到短信表
func gsmSystemBatchesSaveSms(tid int, sms GsmSystemSmsInfo) (int, error) {
	smss := make([]mysql_model.BeehiveGsmSystemSms, 0)
	for _, sms := range sms.Infos {
		this := new(mysql_model.BeehiveGsmSystemSms)
		this.TaskId = tid
		this.Time = sms[0]
		this.SmsContent = sms[1]
		this.SendMobile = sms[2]
		this.SendImsi = sms[3]
		this.RecvMobile = sms[4]
		this.RecvImsi = sms[5]
		smss = append(smss, *this)
	}
	success, err := mysql.GetSession().Insert(smss)
	if err != nil {
		return int(success), err
	}
	return int(success), nil
}

// 获取短信
func GsmSystemGetSms(tid int) (int, error) {
	device, err := gsmSystemGetSms2Device()
	if err != nil {
		return 0, err
	}
	_, err = gsmSystemGetSms2Task(tid)
	if err != nil {
		return 0, err
	}
	gsmSystemGetSms2Log(tid, device)
	sms, err := gsmSystemBatchesSaveSms(tid, device)
	if err != nil {
		return 0, err
	}
	return sms, nil
}

// 获取已经存入数据库的短信数量填充角标
func GsmSystemGetSmsNum(tid int) (int, error) {
	this := new(mysql_model.BeehiveGsmSystemSms)
	this.TaskId = tid
	return this.GsmSystemGetSMSNum()
}

// 获取当前系统中所有终端设备信息
type Devices struct {
	Status    bool       `json:"status"`
	MessageId int        `json:"message_id"`
	Message   string     `json:"message"`
	Infos     [][]string `json:"infos"`
}

// 获取当前系统中所有终端设备信息-设备
func gsmSystemGetDevices2Device() (Devices, error) {
	args := make(map[string]string)
	host, _ := GetHost(GsmSystem)
	resp, err := custom_util.HttpPostJson(nil, args, strings.Join([]string{host, common.GET_DEVICES}, "/"))
	blog.Info(common.GET_DEVICES, zap.Any("request:", args), zap.Any("responds:", string(resp)))
	if err != nil {
		return Devices{}, err
	}
	var devices Devices
	err = json.Unmarshal(resp, &devices)
	if err != nil {
		return Devices{}, err
	}
	return devices, nil
}

// 获取当前系统中所有终端设备信息-日志
func gsmSystemGetDevices2Log(tid int, d Devices) error {
	this := new(mysql_model.BeehiveLog)
	this.TaskId = tid
	this.Title = "尝试获取设备信息"
	this.Content = fmt.Sprintf("获取到%d台设备", len(d.Infos))
	//switch cid {
	//case 0:
	//	this.Content = strings.Join([]string{content, "未知失败,需要查阅日志"}, ":")
	//case 1:
	//	this.Content = strings.Join([]string{content, "读取成功"}, ":")
	//case 2:
	//	this.Content = strings.Join([]string{content, "系统未在运行，请先启动系统再执行"}, ":")
	//case 3:
	//	this.Content = strings.Join([]string{content, "未发现用户信息，请先连接终端设备"}, ":")
	//}
	return this.SetLog()
}

// 模拟短信时收件人下拉列表
func GsmSystemGetReceiveList(tid int) ([]mysql_model.BeehiveGsmSystemTty, error) {
	this := new(mysql_model.BeehiveGsmSystemTty)
	this.TaskId = tid
	list, err := this.GetTtyList()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 获取当前系统中所有终端设备信息-设备表
func gsmSystemGetDevices2Tty(tid int, devices Devices) error {
	// 正常逻辑每次获取设备信息都会删除上一次的结果 调试的时候需要注释掉 否则拿不到信息
	this := new(mysql_model.BeehiveGsmSystemTty)
	this.TaskId = tid
	_, err := this.DeleteByTaskId()
	if err != nil {
		return err
	}

	for _, info := range devices.Infos {
		fresh := new(mysql_model.BeehiveGsmSystemTty)
		fresh.TaskId = tid
		fresh.Imsi = info[0]
		fresh.Imei = info[1]
		fresh.Mobile = info[2]
		_, err := fresh.UpdateTty()
		if err != nil {
			continue
		}
	}
	return nil
}

// 获取当前系统中所有终端设备信息
func GsmSystemGetDevices(tid int) (Devices, error) {
	devices, err := gsmSystemGetDevices2Device()
	if err != nil {
		return Devices{}, err
	}
	err = gsmSystemGetDevices2Tty(tid, devices)
	if err != nil {
		return Devices{}, err
	}
	err = gsmSystemGetDevices2Log(tid, devices)
	if err != nil {
		return Devices{}, err
	}
	return devices, err
}

// 模糊搜索短信
func GsmSystemSearchSms(tid int, key, q string) map[string]interface{} {
	blog.Info("search", zap.Any("task_id:", tid), zap.Any("keyword", key), zap.Any("query", q))
	return mysql_model.GsmSystemSearchSms(tid, key, q)
}

// 短信选项卡 短信列表
func GsmSystemGetTaskAllSms(tid int) map[string]interface{} {
	return mysql_model.GetAllSms(tid)
}
func GmsSystemBulkDeleteSMS(ids []int) int {
	this := new(mysql_model.BeehiveGsmSystemSms)
	sms := this.DeleteSMS(ids)
	return sms
}

// 模拟短信
type Sms struct {
	TaskId  int      `json:"task_id"`
	Send    string   `json:"send"`
	Recv    []string `json:"recv"`
	Content string   `json:"content"`
}

func gsmSystemSimulateSms2Devices(imsi, sender, smsmessage string) (GsmSystemStandResp, error) {
	args := map[string]string{
		"imsi":       imsi,
		"sender":     sender,
		"smsmessage": smsmessage,
	}
	host, _ := GetHost(GsmSystem)
	resp, err := custom_util.HttpPostJson(nil, args, strings.Join([]string{host, common.SEND_SMS}, "/"))
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	var gsmSystemStandResp GsmSystemStandResp
	err = json.Unmarshal(resp, &gsmSystemStandResp)
	if err != nil {
		return GsmSystemStandResp{}, err
	}
	blog.Info(common.GET_DEVICES, zap.Any("request:", args), zap.Any("responds:", gsmSystemStandResp))
	return gsmSystemStandResp, nil
}
func gsmSystemSimulateSms2Task(tid int) (int64, error) {
	s := new(mysql_model.BeehiveGsmSystem)
	s.TaskId = tid
	return s.Update()
}
func gsmSystemSimulateSms2Log(tid int, recv, send, content string) error {
	this := new(mysql_model.BeehiveLog)
	this.TaskId = tid
	this.Content = fmt.Sprintf("发件人:%s 收件人:%s 短信内容:%s", send, recv, content)
	return this.SetLog()
}

// 模拟短信
func GsmSystemSimulateSms(sms Sms) int {
	var result []GsmSystemStandResp
	for _, recv := range sms.Recv {
		//m:=new(mysql_model.BeehiveGsmSystemTty)
		//m.Imsi=recv
		//mobile, _ :=m.FindMobileByImsi()
		device, err := gsmSystemSimulateSms2Devices(recv, sms.Send, sms.Content)
		if err != nil {
			continue
		}
		gsmSystemSimulateSms2Log(sms.TaskId, recv, sms.Send, sms.Content)
		result = append(result, device)
		gsmSystemSimulateSms2Task(sms.TaskId)

	}
	return len(result)
}

func Detail(tid int) (mysql_model.BeehiveGsmSystem, error) {
	this := new(mysql_model.BeehiveGsmSystem)
	return this.GetOne(tid)

}
