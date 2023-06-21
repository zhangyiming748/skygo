package gps

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/gpslog"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

const (
	DEVICE_DEFAULT  = iota
	DEVICE_NO_START // 未启动
	DEVICE_RUNNING  // 启动成功，运行中
	DEVICE_FAIL     // 启动失败
	DEVICE_CLOSE    // 关闭
)

const (
	CHEAT_STATUS_ING  = 1 // 欺骗中
	CHEAT_STATUS_FAIL = 3 // 失败
)

const (
	URL_SDR        = "/api/sdr"
	URL_SDR_CONFIG = "/api/sdrConfig"
	URL_SDR_SAVE   = "/api/save-NMEA-on-server"
)

type Scl struct {
	Id             int     `json:"id"`
	AngleStart     int     `json:"angleStart"`
	AngleEnd       int     `json:"angleEnd"`
	SpeedDecayRate float32 `json:"speedDecayRate"`
	Acc            int     `json:"acc"`
	Jerk           int     `json:"jerk"`
}

func SetTime(taskId int, gpstime string) (int, error) {
	gt := mysql_model.GpsTask{}
	b, err := gt.FindGpsByTaskId(taskId)
	if err != nil {
		return 0, err
	}
	_ = SetGpsLog(taskId, "设置gps时间", gpstime)
	if !b {
		gt.TaskId = taskId
		gt.Time = gpstime
		gt.CreateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
		_, err := gt.Create()
		if err != nil {
			return 0, err
		}
	} else {
		gt.Time = gpstime
		gt.UpdateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
		gt.Update()
	}
	return gt.Id, nil
}

// 获取task
func GetTaskById(taskId int) (*mysql_model.Task, error) {
	task, has := mysql_model.TaskFindById(taskId)
	if !has {
		return task, errors.New("任务不存在")
	}
	if task.Status == TASK_COMPLETE {
		return task, errors.New("该任务已完成，不能再测试了")
	}
	return task, nil
}

// 获取gps_task
func GetGpsTaskByTaskId(taskId int) (mysql_model.GpsTask, error) {
	gt := mysql_model.GpsTask{}
	b, err := gt.FindGpsByTaskId(taskId)
	if !b {
		err = fmt.Errorf("没有找到任务")
	}
	return gt, err
}

func GetLatestCheatByTaskId(taskId int) (mysql_model.GpsCheat, error) {
	gc := mysql_model.GpsCheat{}
	_, err := gc.GetLatestByTaskId(taskId, 1)
	if err != nil {
		return gc, err
	}
	return gc, nil
}

func Cheat(taskId, searchId int) (bool, error) {
	gc, err := GetLatestCheatByTaskId(taskId)
	if err != nil {
		return false, err
	}
	if gc.Status == CHEAT_STATUS_ING {
		return false, errors.New("有正在欺骗中的任务，请先停止欺骗")
	}

	gs := mysql_model.GpsSearch{}
	b, err := gs.GetOne(searchId)
	if err != nil {
		return false, err
	}
	if !b {
		return false, errors.New("没有找到对应的运动轨迹线路")
	}
	if gs.Type != common.GPS_TYPE_MOTION {
		return false, errors.New("type数据异常")
	}
	if gs.TaskId != taskId {
		return false, errors.New("taskId数据异常")
	}

	gt := mysql_model.GpsSteerTemplate{}
	b, err = gt.GetOne(gs.TemplateId)
	if err != nil {
		return false, err
	}
	if !b {
		return false, errors.New("没找到模板")
	}

	gpsTask, err := GetGpsTaskByTaskId(taskId)
	if err != nil {
		return false, err
	}

	if gpsTask.Status != DEVICE_RUNNING {
		return false, errors.New("设备未启动")
	}

	name := uuid.NewV4().String()
	arg, err := MakeArg(gs.Req, name, gpsTask.Time, gt)
	if err != nil {
		return false, err
	}
	b, err = SdrCheat(arg)
	if err != nil {
		return false, err
	}
	if !b {
		return false, errors.New("SdrCheat false")
	}
	type Second struct {
		Name     string `json:"name"`
		RealTime bool   `json:"realtime"`
	}
	sc := Second{}
	sc.Name = name
	sc.RealTime = true
	// str := `{name: "请输入轨迹名", realtime: true}`
	b, err = SdrCheatSecond(sc)
	if err != nil {
		return false, err
	}
	gcCreate := mysql_model.GpsCheat{}
	gcCreate.Start = gs.Start
	gcCreate.Middle = gs.Middle
	gcCreate.End = gs.End
	gcCreate.Req = gs.Req
	gcCreate.Type = common.GPS_TYPE_MOTION
	gcCreate.SearchId = gs.Id
	gcCreate.TaskId = taskId
	gcCreate.Status = CHEAT_STATUS_ING
	gcCreate.CreateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	_, err = gcCreate.Create()
	_ = SetGpsLog(taskId, "开始运动轨迹欺骗", "成功")
	return b, err
}

func SdrCheat(arg interface{}) (bool, error) {
	url := service.LoadGpsConfig().Url + URL_SDR_SAVE
	gpslog.Info("SdrCheat", zap.Any("url:", url))
	gpslog.Info("SdrCheat", zap.Any("arg:", arg))
	resp, err := custom_util.HttpPostJson(nil, arg, url)
	gpslog.Info("SdrCheat", zap.Any("resp:", string(resp)), zap.Any("err:", err))
	if err != nil {
		return false, err
	}
	return true, nil
}

func SdrCheatSecond(arg interface{}) (bool, error) {
	url := service.LoadGpsConfig().Url + URL_SDR
	gpslog.Info("SdrCheatSecond", zap.Any("url:", url))
	gpslog.Info("SdrCheatSecond", zap.Any("arg:", arg))
	resp, err := custom_util.HttpPostJson(nil, arg, url)
	gpslog.Info("SdrCheatSecond", zap.Any("resp:", string(resp)), zap.Any("err:", err))
	if err != nil {
		return false, err
	}
	type RespSecond struct {
		Status string
	}
	rs := RespSecond{}
	err = json.Unmarshal([]byte(resp), &rs)
	if err != nil {
		return false, err
	}
	if rs.Status != "ok" {
		return false, nil
	}
	return true, nil

}

func MakeArg(arg, name, gpsTime string, gt mysql_model.GpsSteerTemplate) (interface{}, error) {
	type Pos struct {
		Lng float32 `json:"lng"`
		Lat float32 `json:"lat"`
	}

	type MotionConf struct {
		MaxLongAcc              float32 `json:"maxLongAcc"`
		MaxLatAcc               float32 `json:"maxLatAcc"`
		MaxJerk                 int     `json:"maxJerk"`
		MaxSpeed                float32 `json:"maxSpeed"`
		StationaryPeriod        float32 `json:"stationaryPeriod"`
		StationaryPeriodEnd     float32 `json:"stationaryPeriodEnd"`
		PositionSmoothingFactor int     `json:"positionSmoothingFactor"`
		SpeedSmoothingFactor    int     `json:"speedSmoothingFactor"`
		SpeedChangeList         []Scl   `json:"speedChangeList"`
	}

	type Req struct {
		Order         int        `json:"order"`
		BeginTime     string     `json:"begin_time"`
		Distance      int        `json:"distance"`
		Duration      int        `json:"duration"`
		Speed         int        `json:"speed"`
		StartLocation Pos        `json:"start_location"`
		EndLocation   Pos        `json:"end_location"`
		MotionConf    MotionConf `json:"motion_conf"`
	}

	type CheatData struct {
		Name string `json:"name"`
		Data []Req  `json:"data"`
	}

	argReq := []Req{}
	if err := json.Unmarshal([]byte(arg), &argReq); err != nil {
		return "", err
	}
	mc := MotionConf{}
	mc.MaxLongAcc = gt.MaxLongAcc
	mc.MaxLatAcc = gt.MaxLatAcc
	mc.MaxJerk = gt.MaxJerk
	mc.MaxSpeed = gt.MaxSpeed
	mc.StationaryPeriod = gt.StationaryPeriod
	mc.StationaryPeriodEnd = gt.StationaryPeriodEnd
	mc.PositionSmoothingFactor = gt.PositionSmoothingFactor
	mc.SpeedSmoothingFactor = gt.SpeedSmoothingFactor
	defaltSpeed := SpeedChangeList(gt.Type)
	mc.SpeedChangeList = defaltSpeed

	i := 1
	for k := range argReq {
		argReq[k].Order = i
		i++
		argReq[k].BeginTime = gpsTime
		argReq[k].MotionConf = mc
	}

	cd := CheatData{}
	cd.Name = name
	cd.Data = argReq
	return cd, nil
}

func Start(taskId int, gpsTime string) {
	go SdrStart(taskId, gpsTime)
}

func Close(taskId int) {
	go SdrClose(taskId)
}

func CreateLine(gs mysql_model.GpsSearch) (int, error) {
	gt := mysql_model.GpsSteerTemplate{}
	b, err := gt.GetOne(gs.TemplateId)
	if err != nil {
		return 0, err
	}
	if !b {
		return 0, errors.New("没找到模板")
	}

	gs.CreateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	gs.Type = common.GPS_TYPE_MOTION
	_, err = gs.Create()
	_ = SetGpsLog(gs.TaskId, "生成运动轨迹路线", "设置行驶模板为:"+gt.Name)
	_ = SetGpsLog(gs.TaskId, "生成运动轨迹路线", "设置的路线为:"+gs.Start+gs.Middle+gs.End)
	return gs.Id, err
}

func LineHistory(taskId, templateId int) ([]map[string]interface{}, error) {
	gs := mysql_model.GpsSearch{}
	rs, err := gs.FindMotion(taskId, templateId, common.GPS_SEARCH_MOTION_LIMIT)
	if err != nil {
		return nil, err
	}
	data := make([]map[string]interface{}, 0)
	for _, v := range rs {
		m := map[string]interface{}{}
		m["id"] = v.Id
		m["start"] = v.Start
		m["middle"] = v.Middle
		m["end"] = v.End
		data = append(data, m)
	}
	return data, nil
}

func Online() (bool, error) {
	url := service.LoadGpsConfig().Url + URL_SDR
	gpslog.Info("Online", zap.Any("url:", url))
	resp, err := custom_util.HttpGet(nil, nil, url)
	gpslog.Info("Online", zap.Any("resp:", string(resp)), zap.Any("err:", err))
	if err != nil {
		return false, err
	}
	result := map[string]bool{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return false, err
	}
	return result["online"], nil
}

func SdrConfig(gpsTime string) error {
	url := service.LoadGpsConfig().Url + URL_SDR_CONFIG
	gpslog.Info("SdrConfig", zap.Any("url:", url))
	args := map[string]interface{}{"gps_time": gpsTime, "sampleRate": 2600000, "bandwidth": 2500000}
	gpslog.Info("SdrConfig", zap.Any("args:", args))
	resp, err := custom_util.HttpPostJson(nil, args, url)
	gpslog.Info("SdrConfig", zap.Any("resp:", string(resp)), zap.Any("err:", err))
	if err != nil {
		return err
	}
	return nil
}

func SdrStart(taskId int, gpsTime string) {
	defer func() {
		if err := recover(); err != nil {
			gpslog.Error("SdrStart:", zap.Any("taskId:", taskId), zap.Any("err:", err))
			_ = SetGpsLog(taskId, "启动系统", "失败")
		} else {
			_ = SetGpsLog(taskId, "启动系统", "成功")
		}
	}()

	// docker run
	// docker start
	sh := service.LoadGpsConfig().StartSh
	err := Sh(sh)
	if err != nil {
		panic(err)
	}

	// set sdrConfig
	err = SdrConfig(gpsTime)
	if err != nil {
		panic(err)
	}

	// 修改gps_task的status
	gt := mysql_model.GpsTask{}
	_, err = gt.UpdateStatusByTaskId(taskId, DEVICE_RUNNING)
	if err != nil {
		panic(err)
	}

	// 修改task状态为测试中
	taskLogic := new(Task)
	_, err = taskLogic.Start(taskId)
	if err != nil {
		panic(err)
	}
}

func SdrClose(taskId int) {
	defer func() {
		if err := recover(); err != nil {
			gpslog.Error("SdrClose:", zap.Any("taskId:", taskId), zap.Any("err:", err))
			_ = SetGpsLog(taskId, "关闭系统", "失败")
		} else {
			_ = SetGpsLog(taskId, "关闭系统", "成功")
		}
	}()

	sh := service.LoadGpsConfig().StopSh
	err := Sh(sh)
	if err != nil {
		panic(err)
	}

	gt := mysql_model.GpsTask{}
	_, err = gt.UpdateStatusByTaskId(taskId, DEVICE_CLOSE)
	if err != nil {
		panic(err)
	}

	// 修改task状态为暂停
	taskLogic := new(Task)
	_, err = taskLogic.Close(taskId)
	if err != nil {
		panic(err)
	}
}

func Sh(sh string) error {
	cmd := exec.Command("bash", "-c", sh)
	//打开一个文件
	f, _ := os.OpenFile(service.LoadGpsConfig().Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	defer f.Close()
	cmd.Stderr = f
	cmd.Stdout = f
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

// 返回t类型的speedChangeList数据
func SpeedChangeList(t int) []Scl {
	m := []Scl{}
	switch t {
	case 1:
		m = append(m,
			Scl{Id: 1, AngleStart: 180, AngleEnd: 150, SpeedDecayRate: 1, Acc: 10, Jerk: 10},
			Scl{Id: 2, AngleStart: 150, AngleEnd: 130, SpeedDecayRate: 0.5, Acc: 10, Jerk: 10},
			Scl{Id: 3, AngleStart: 130, AngleEnd: 120, SpeedDecayRate: 0.4, Acc: 10, Jerk: 10},
			Scl{Id: 4, AngleStart: 120, AngleEnd: 110, SpeedDecayRate: 0.2, Acc: 10, Jerk: 10},
			Scl{Id: 5, AngleStart: 110, AngleEnd: 100, SpeedDecayRate: 0.15, Acc: 10, Jerk: 10},
			Scl{Id: 6, AngleStart: 100, AngleEnd: 90, SpeedDecayRate: 0.1, Acc: 10, Jerk: 10},
			Scl{Id: 7, AngleStart: 90, AngleEnd: 0, SpeedDecayRate: 0, Acc: 10, Jerk: 10},
		)
	case 2:
		m = append(m,
			Scl{Id: 1, AngleStart: 180, AngleEnd: 150, SpeedDecayRate: 1, Acc: 10, Jerk: 10},
			Scl{Id: 2, AngleStart: 150, AngleEnd: 130, SpeedDecayRate: 0.5, Acc: 10, Jerk: 10},
			Scl{Id: 3, AngleStart: 130, AngleEnd: 120, SpeedDecayRate: 0.4, Acc: 10, Jerk: 10},
			Scl{Id: 4, AngleStart: 120, AngleEnd: 110, SpeedDecayRate: 0.2, Acc: 10, Jerk: 10},
			Scl{Id: 5, AngleStart: 110, AngleEnd: 100, SpeedDecayRate: 0.15, Acc: 10, Jerk: 10},
			Scl{Id: 6, AngleStart: 100, AngleEnd: 90, SpeedDecayRate: 0.1, Acc: 10, Jerk: 10},
			Scl{Id: 7, AngleStart: 90, AngleEnd: 0, SpeedDecayRate: 0, Acc: 10, Jerk: 10},
		)
	case 3:
		m = append(m,
			Scl{Id: 1, AngleStart: 180, AngleEnd: 150, SpeedDecayRate: 1, Acc: 10, Jerk: 10},
			Scl{Id: 2, AngleStart: 150, AngleEnd: 130, SpeedDecayRate: 0.5, Acc: 10, Jerk: 10},
			Scl{Id: 3, AngleStart: 130, AngleEnd: 120, SpeedDecayRate: 0.4, Acc: 10, Jerk: 10},
			Scl{Id: 4, AngleStart: 120, AngleEnd: 110, SpeedDecayRate: 0.2, Acc: 10, Jerk: 10},
			Scl{Id: 5, AngleStart: 110, AngleEnd: 100, SpeedDecayRate: 0.15, Acc: 10, Jerk: 10},
			Scl{Id: 6, AngleStart: 100, AngleEnd: 90, SpeedDecayRate: 0.1, Acc: 10, Jerk: 10},
			Scl{Id: 7, AngleStart: 90, AngleEnd: 0, SpeedDecayRate: 0, Acc: 10, Jerk: 10},
		)
	case 4:
		m = append(m,
			Scl{Id: 1, AngleStart: 180, AngleEnd: 150, SpeedDecayRate: 1, Acc: 7, Jerk: 7},
			Scl{Id: 2, AngleStart: 150, AngleEnd: 130, SpeedDecayRate: 0.5, Acc: 7, Jerk: 7},
			Scl{Id: 3, AngleStart: 130, AngleEnd: 120, SpeedDecayRate: 0.4, Acc: 7, Jerk: 7},
			Scl{Id: 4, AngleStart: 120, AngleEnd: 110, SpeedDecayRate: 0.2, Acc: 7, Jerk: 7},
			Scl{Id: 5, AngleStart: 110, AngleEnd: 100, SpeedDecayRate: 0.15, Acc: 7, Jerk: 7},
			Scl{Id: 6, AngleStart: 100, AngleEnd: 90, SpeedDecayRate: 0.1, Acc: 7, Jerk: 7},
			Scl{Id: 7, AngleStart: 90, AngleEnd: 0, SpeedDecayRate: 0, Acc: 7, Jerk: 7},
		)
	case 5:
		m = append(m,
			Scl{Id: 1, AngleStart: 180, AngleEnd: 150, SpeedDecayRate: 1, Acc: 10, Jerk: 10},
			Scl{Id: 2, AngleStart: 150, AngleEnd: 130, SpeedDecayRate: 0.5, Acc: 10, Jerk: 10},
			Scl{Id: 3, AngleStart: 130, AngleEnd: 120, SpeedDecayRate: 0.4, Acc: 10, Jerk: 10},
			Scl{Id: 4, AngleStart: 120, AngleEnd: 110, SpeedDecayRate: 0.2, Acc: 10, Jerk: 10},
			Scl{Id: 5, AngleStart: 110, AngleEnd: 100, SpeedDecayRate: 0.15, Acc: 10, Jerk: 10},
			Scl{Id: 6, AngleStart: 100, AngleEnd: 90, SpeedDecayRate: 0.1, Acc: 10, Jerk: 10},
			Scl{Id: 7, AngleStart: 90, AngleEnd: 0, SpeedDecayRate: 0, Acc: 10, Jerk: 10},
		)
	}
	return m
}

func SetGpsLog(taskId int, title, content string) error {
	beehiveLog := mysql_model.BeehiveLog{}
	beehiveLog.TaskId = taskId
	beehiveLog.Title = title
	beehiveLog.Content = content
	if err := beehiveLog.SetLog(); err != nil {
		gpslog.Error("SetGpsLog:", zap.Any("err:", err))
		return err
	}
	return nil
}
