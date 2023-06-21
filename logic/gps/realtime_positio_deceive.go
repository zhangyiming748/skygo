package gps

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/clog"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
	"time"
)

const (
	TASK_TYPE_REALTIME = 1
	CHEAT_STATUS_START = 1
	CHEAT_STATUS_STOP  = 2
)

type Result struct {
	Status string
}

func SetLogs(taskId int, content string) {
	log := new(mysql_model.BeehiveLog)
	log.TaskId = taskId
	log.Content = content
	err := log.SetLog()
	if err != nil {
		return
	}
	return
}

func StartRealtimePositionDeceive(taskId int, start string, lng float32, lat float32) (Result, error) {
	var r Result

	gpsTask, err := GetGpsTaskByTaskId(taskId)
	if err != nil {
		return r, err
	}
	if gpsTask.Status != DEVICE_RUNNING {
		return r, errors.New("设备未启动")
	}

	args := make(map[string]interface{})
	args["lng"] = lng
	args["lat"] = lat
	args["h"] = 100
	url := service.LoadConfig().Gps.Url
	url += "/api/sdr"
	resp, err := custom_util.HttpPostJson(nil, args, url)
	clog.Info("StartRealtimePositionDeceive", zap.Any("url: ", url),
		zap.Any("responds: ", string(resp)), zap.Any("err: ", err))
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return r, err
	}

	module := new(mysql_model.GpsCheat)
	module.TaskId = taskId
	module.Lng = lng
	module.Lat = lat
	module.Start = start
	module.Type = TASK_TYPE_REALTIME
	module.Status = CHEAT_STATUS_START
	module.Resp = r.Status
	module.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	_, err = module.Create()
	if err != nil {
		return r, err
	}
	SetLogs(taskId, "开始实时位置欺骗")
	return r, nil
}

func StopRealtimePositionDeceive(taskId int) (string, error) {
	url := service.LoadConfig().Gps.Url
	url += "/api/sdr"
	resp, err := custom_util.HttpPostJsonPut(nil, nil, url)
	clog.Info("StartRealtimePositionDeceive", zap.Any("url: ", url),
		zap.Any("responds: ", string(resp)), zap.Any("err: ", err))
	if err != nil {
		return "", err
	}

	module := new(mysql_model.GpsCheat)
	module.Status = CHEAT_STATUS_STOP
	err = module.Update(taskId)
	if err != nil {
		return "", err
	}
	SetLogs(taskId, "停止实时位置欺骗")
	return "success", nil
}
