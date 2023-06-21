package beehive

import (
	"encoding/json"
	"math/rand"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/lib/common_lib/mysql"

	"testing"
	"time"
)

func TestGetJson(t *testing.T) {
	resp := []byte(`{"status": true, "message_id": 1, "message": "Start scanner successfully"}`)
	result := new(GsmSystemStandResp)
	json.Unmarshal(resp, &result)
	t.Logf("value is %+v\n", result.Status)
	t.Logf("value is %+v\n", result.Message)
	t.Logf("value is %+v\n", result.MessageId)
	smss := []byte(`{"status":true,"message_id":1,"message":"Success","infos":[["2021-08-13T17:31:06.2","Hhh","10000001","001010123456780",null,"233"],["2021-08-13T17:32:26.9","Fhhhj","10000001","001010123456780","10000000","001012333333333"]]}`)
	sms_list := new(GsmSystemSmsInfo)
	json.Unmarshal(smss, &sms_list)
	t.Logf("value is %+v\n", sms_list.Status)
	t.Logf("value is %+v\n", sms_list.MessageId)
	t.Logf("value is %+v\n", sms_list.Message)
	for i, sms := range sms_list.Infos {
		t.Logf("第%v条信息\n", i)
		t.Logf("时间:%v\n短信内容:%v\n发送者电话号码:%v\n发送者imei:%v\n接收者电话号:%v\n码接收者imei:%v\n", sms[0], sms[1], sms[2], sms[3], sms[4], sms[5])
	}
	t.Logf("value is %+v\n", sms_list.Status)

}
func TestingInit() {
	sys_service.InitConfigWatcher("qa", "../../config/qa/config.tml")
	mysql.InitMysqlEngine()
}

func TestSetConfigWithDevice(t *testing.T) {
	TestingInit()

}
func TestSetConfigWithTask(t *testing.T) {
	TestingInit()
	_, err := GsmSystemSetConfig2Task(1, 1)
	if err != nil {
		return
	}
	t.Logf("配置设备记录到任务详情表")
}
func TestSetConfigWithLog(t *testing.T) {
	TestingInit()
	err := gsmSystemSetConfig2Log(5, 3, 2)
	if err != nil {
		return
	}
	t.Logf("配置设备记录到日志表")
}
func BenchmarkSetConfigWithTask(b *testing.B) {
	TestingInit()
	rand.Seed(time.Now().Unix())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GsmSystemSetConfig2Task(rand.Intn(100), rand.Intn(3))
		b.Logf("配置设备记录到任务详情表")
	}
	b.StopTimer()
}

// 测试获取设备信息
// []byte(`{"status":true,"message_id":1,"message":"Success","infos":[["001010123456780","351615087961130","10000001","none"],["001012333333333","355754071347990","10000002","192.168.99.2"]]}`)
func TestGetDevices(t *testing.T) {
	TestingInit()
	tid := 7
	devices, err := GsmSystemGetDevices(tid)
	if err != nil {
		return
	}
	for _, device := range devices.Infos {
		t.Logf("获取到的设备%v\n", device)
	}
	t.Logf("")
}
func TestGsmSystemGetDevices2Tty(t *testing.T) {
	TestingInit()
	tid := 7
	resp := [][]string{{"001010123456780", "351615087961130", "10000001", "none"}, {"001012333333333", "355754071347990", "10000002", "192.168.99.2"}}
	devices := Devices{
		Status:    true,
		MessageId: 1,
		Message:   "假信息",
		Infos:     resp,
	}
	tty := gsmSystemGetDevices2Tty(tid, devices)

	t.Logf("%+v\n", tty)
}
func TestGsmSystemGetReceiveList(t *testing.T) {
	TestingInit()
	tid := 7
	list, err := GsmSystemGetReceiveList(tid)
	if err != nil {
		return
	}
	for i, v := range list {
		t.Logf("fetch %d is %v\n", i, v.Mobile)
	}
}

func TestSimulator(t *testing.T) {
	TestingInit()
	sms := Sms{
		TaskId:  7,
		Send:    "1415926535897",
		Recv:    []string{"001012333333333"},
		Content: "这是一条假信息",
	}
	GsmSystemSimulateSms(sms)
}

// 模拟短信发送后过更新任务时间
func TestSimulator2Task(t *testing.T) {
	TestingInit()
	gsmSystemSimulateSms2Task(7)
}
func TestSimulator2log(t *testing.T) {
	TestingInit()
	gsmSystemSimulateSms2Log(7, "1", "", "")
}
