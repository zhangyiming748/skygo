package beehive

import (
	"encoding/json"
	"errors"
	"fmt"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/blog"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/mysql_model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ChanelMap = map[int]string{1: "900", 2: "1800"}

const (
	SNIFF_CLOSE     = iota // 关闭系统
	SNIFF_SCANING          // 频点扫描中
	SNIFF_SCAN_STOP        // 频点扫描停止
	SNIFF_IMSI             // 嗅探imsi
	SNIFF_SMS              // 嗅探sms
	SNIFF_STOP             // 停止嗅探

	URL_SCANNER = "/scanner"
	URL_GETFREQ = "/getfreq"
	URL_SNIFFER = "/sniffer"
	URL_GETIMSI = "/getimsi"
	URL_GETSMS  = "/getsms"
	URL_STOP    = "/stop"
)

type ScanTimer struct {
	Stop   bool `json:"stop"`
	Minute int  `json:"minute"`
	Second int  `json:"second"`
}

type GSM struct {
	Status    bool
	MessageId int `json:"message_id"`
	Message   string
	Data      [][]string
}

type Sniffer struct{}

func (s Sniffer) Scan(ctx *gin.Context, taskId int, channel int) error {
	band, ok := ChanelMap[channel]
	if !ok {
		return errors.New("channel参数不正确")
	}
	bool, err := s.GSMSnifferScanner(band)
	if err != nil {
		return err
	}
	if !bool {
		return err
	}

	beehiveLog := mysql_model.BeehiveLog{}
	beehiveLog.TaskId = taskId
	beehiveLog.Title = "开始频点扫描"
	beehiveLog.Content = "扫描的频段" + band + "M"
	if err := beehiveLog.SetLog(); err != nil {
		return err
	}

	// 查找BeehiveGsmSniffer,如果没有则增加一条
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	b, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		return err
	}
	if !b {
		snifferModel.TaskId = taskId
		snifferModel.Status = SNIFF_SCANING
		snifferModel.Channel = channel
		snifferModel.ScanTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
		snifferModel.CreateTime = snifferModel.ScanTime
		_, err = snifferModel.Create()
	} else {
		snifferModel.Status = SNIFF_SCANING
		snifferModel.Channel = channel
		snifferModel.ScanTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
		snifferModel.UpdateTime = snifferModel.ScanTime
		_, err = snifferModel.Update()
	}
	if err != nil {
		return err
	}
	// 修改task的状态为测试中
	if err := s.UpdateTaskStatus(ctx, taskId, TASK_TESTING); err != nil {
		return err
	}

	return nil
}

func (s Sniffer) ReScan(ctx *gin.Context, taskId int) error {
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	b, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("没有sniffer任务，请检查参数")
	}

	band, ok := ChanelMap[snifferModel.Channel]
	if !ok {
		return errors.New("channel参数不正确")
	}
	bool, err := s.GSMSnifferScanner(band)
	if err != nil {
		return err
	}
	if !bool {
		return err
	}

	beehiveLog := mysql_model.BeehiveLog{}
	beehiveLog.TaskId = taskId
	beehiveLog.Title = "重新扫描频点"
	beehiveLog.Content = "扫描的频段" + band + "M"
	if err := beehiveLog.SetLog(); err != nil {
		return err
	}

	// 修改task的状态为测试中
	if err := s.UpdateTaskStatus(ctx, taskId, TASK_TESTING); err != nil {
		return err
	}

	snifferModel.Status = SNIFF_SCANING
	snifferModel.ScanTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	snifferModel.UpdateTime = snifferModel.ScanTime
	_, err = snifferModel.Update()
	if err != nil {
		return err
	}
	return nil
}

func (s Sniffer) GetFreq(ctx *gin.Context, taskId int) ([]string, error) {
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	has, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("sniffer任务不存在")
	}
	gsm, err := s.GSMSnifferGetfreq()
	if err != nil {
		return nil, err
	}
	data := make([]string, 0)
	if len(gsm.Data) > 0 {
		allFreq := map[string]string{}
		if snifferModel.Frequency != "" {
			// 老的放map里，方便去重
			oldFreqs := strings.Split(snifferModel.Frequency, ",")
			fmt.Println("oldFreqs:", oldFreqs)
			if len(oldFreqs) > 0 {
				for _, v := range oldFreqs {
					if _, ok := allFreq[v]; !ok {
						allFreq[v] = v
					}
				}
			}
		}

		// 新的也放map里
		newFreq := []string{}
		for _, v := range gsm.Data {
			newFreq = append(newFreq, v[1])
			if _, ok := allFreq[v[1]]; !ok {
				allFreq[v[1]] = v[1]
			}
		}

		for _, v := range allFreq {
			if v != "" {
				data = append(data, v)
			}
		}

		// 保存频点
		frequency := strings.Join(data, ",")
		snifferModel.Frequency = frequency
		snifferModel.UpdateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
		snifferModel.Update()

		beehiveLog := mysql_model.BeehiveLog{}
		beehiveLog.TaskId = taskId
		beehiveLog.Title = "扫描频点"
		beehiveLog.Content = "扫描到频点是" + strings.Join(newFreq, ",")
		if err := beehiveLog.SetLog(); err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (s Sniffer) GetTimer(ctx *gin.Context, taskId int) (ScanTimer, error) {
	st := ScanTimer{}
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	has, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		return st, err
	}
	if !has {
		return st, errors.New("sniffer任务不存在")
	}

	end, _ := time.ParseInLocation("2006-01-02 15:04:05", snifferModel.ScanTime, time.Local)
	end = end.Add(time.Duration(300) * time.Second)
	now := time.Now()
	if end.Before(now) {
		st.Stop = true
		return st, nil
	}
	d := end.Sub(now)
	st.Minute = int(d.Minutes())
	st.Second = int(d.Seconds()) % 60
	return st, nil
}

func (s Sniffer) StartSniff(ctx *gin.Context, taskId int, freq, mode string) (bool, error) {
	// 查找BeehiveGsmSniffer,如果没有则增加一条
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	b, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		return false, err
	}
	if !b {
		return false, errors.New("没有sniffer任务，请检查参数")
	}

	err = s.GSMSnifferStop()
	if err != nil {
		return false, err
	}
	gsm, err := s.GSMSniffer(freq, mode)
	if err != nil {
		return false, err
	}

	if !gsm.Status {
		return false, errors.New(gsm.Message)
	}

	beehiveLog := mysql_model.BeehiveLog{}
	beehiveLog.TaskId = taskId
	beehiveLog.Title = "开始嗅探" + mode
	if err := beehiveLog.SetLog(); err != nil {
		return false, err
	}

	if mode == "imsi" {
		snifferModel.Status = SNIFF_IMSI
	} else if mode == "sms" {
		snifferModel.Status = SNIFF_SMS
	} else {
		return false, errors.New("请检查mode参数")
	}

	snifferModel.UpdateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	snifferModel.SniffFreq = freq
	_, err = snifferModel.Update()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s Sniffer) GetImsi(ctx *gin.Context, taskId int) (int, error) {
	data := []mysql_model.BeehiveGsmSnifferImsi{}
	resp, err := s.GSMSnifferGetimsi()
	if err != nil {
		return 0, err
	}

	if len(resp.Data) > 0 {
		for _, v := range resp.Data {
			tmp := mysql_model.BeehiveGsmSnifferImsi{}
			tmp.TaskId = taskId
			tmp.Content = v[0]
			tmp.Date = custom_util.TimeTransLocal(v[1])
			tmp.CreateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
			data = append(data, tmp)
		}
		mysql.GetSession().Insert(data)

		beehiveLog := mysql_model.BeehiveLog{}
		beehiveLog.TaskId = taskId
		beehiveLog.Title = "嗅探imsi"
		beehiveLog.Content = "嗅探到" + strconv.Itoa(len(resp.Data)) + "条imsi信息"
		if err := beehiveLog.SetLog(); err != nil {
			return 0, err
		}

	}
	return len(resp.Data), nil
}

func (s Sniffer) GetSms(ctx *gin.Context, taskId int) (int, error) {
	data := []mysql_model.BeehiveGsmSnifferSms{}
	resp, err := s.GSMSnifferGetsms()
	if err != nil {
		return 0, err
	}

	if len(resp.Data) > 0 {
		for _, v := range resp.Data {
			tmp := mysql_model.BeehiveGsmSnifferSms{}
			tmp.TaskId = taskId
			tmp.Content = v[0]
			tmp.Date = custom_util.TimeTransLocal(v[1])
			tmp.CreateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
			data = append(data, tmp)
		}
		mysql.GetSession().Insert(data)

		beehiveLog := mysql_model.BeehiveLog{}
		beehiveLog.TaskId = taskId
		beehiveLog.Title = "嗅探sms"
		beehiveLog.Content = "嗅探到" + strconv.Itoa(len(resp.Data)) + "条sms信息"
		if err := beehiveLog.SetLog(); err != nil {
			return 0, err
		}
	}
	return len(resp.Data), nil
}

func (s Sniffer) GetImsiCount(ctx *gin.Context, taskId int) (int, error) {
	c, err := mysql.GetSession().Table(mysql_model.BeehiveGsmSnifferImsi{}).Where("task_id=?", taskId).Where("delete_time=0").Count()
	if err != nil {
		return 0, err
	}
	return int(c), nil
}

func (s Sniffer) GetSmsCount(ctx *gin.Context, taskId int) (int, error) {
	c, err := mysql.GetSession().Table(mysql_model.BeehiveGsmSnifferSms{}).Where("task_id=?", taskId).Where("delete_time=0").Count()
	if err != nil {
		return 0, err
	}
	return int(c), nil
}

func (s Sniffer) DelImsi(ctx *gin.Context, taskId int, ids []interface{}) (int64, error) {
	data := map[string]interface{}{}
	data["delete_time"] = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	return mysql.GetSession().Table(mysql_model.BeehiveGsmSnifferImsi{}).Where("task_id=?", taskId).In("id", ids).Update(data)
}

func (s Sniffer) DelSms(ctx *gin.Context, taskId int, ids []interface{}) (int64, error) {
	data := map[string]interface{}{}
	data["delete_time"] = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	return mysql.GetSession().Table(mysql_model.BeehiveGsmSnifferSms{}).Where("task_id=?", taskId).In("id", ids).Update(data)
}

func (s Sniffer) Close(ctx *gin.Context, taskId int) error {
	err := s.GSMSnifferStop()
	if err != nil {
		return err
	}

	snifferModel := mysql_model.BeehiveGsmSniffer{}
	has, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("sniffer任务不存在")
	}
	snifferModel.Status = SNIFF_CLOSE
	snifferModel.UpdateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	snifferModel.Frequency = ""
	// 修改sniffer的状态为初始化0
	if _, err := snifferModel.Update("status", "update_time", "frequency"); err != nil {
		return err
	}

	// 修改task的状态为暂停
	if err := s.UpdateTaskStatus(ctx, taskId, TASK_PAUSE); err != nil {
		return err
	}

	beehiveLog := mysql_model.BeehiveLog{}
	beehiveLog.TaskId = taskId
	beehiveLog.Title = "关闭系统"
	if err := beehiveLog.SetLog(); err != nil {
		return err
	}
	return nil
}

// 停止扫描，记录嗅探时间，前端得停定时器
func (s Sniffer) StopScan(ctx *gin.Context, taskId int) error {
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	b, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("没有sniffer任务，请检查参数")
	}

	snifferModel.Status = SNIFF_SCAN_STOP
	snifferModel.SniffTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	snifferModel.UpdateTime = snifferModel.SniffTime
	_, err = snifferModel.Update()
	if err != nil {
		return err
	}

	// 修改task的状态为暂停
	if err := s.UpdateTaskStatus(ctx, taskId, TASK_PAUSE); err != nil {
		return err
	}

	beehiveLog := mysql_model.BeehiveLog{}
	beehiveLog.TaskId = taskId
	beehiveLog.Title = "停止频点扫描"
	if err := beehiveLog.SetLog(); err != nil {
		return err
	}
	return nil
}

// 停止嗅探，前端也得停定时器
func (s Sniffer) StopSniff(ctx *gin.Context, taskId int) error {
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	b, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("没有sniffer任务，请检查参数")
	}

	snifferModel.Status = SNIFF_STOP
	snifferModel.UpdateTime = snifferModel.ScanTime
	_, err = snifferModel.Update()
	if err != nil {
		return err
	}

	if err = s.UpdateTaskStatus(ctx, taskId, TASK_PAUSE); err != nil {
		return err
	}

	beehiveLog := mysql_model.BeehiveLog{}
	beehiveLog.TaskId = taskId
	beehiveLog.Title = "停止嗅探"
	if err := beehiveLog.SetLog(); err != nil {
		return err
	}
	return nil
}

// 扫描停止4个小时后，要提示要不要重新扫描
func (s Sniffer) IfStopScan(ctx *gin.Context, taskId int) bool {
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	_, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		blog.Error("IfStopScan", zap.Any("err:", err.Error()))
		return false
	}
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", snifferModel.SniffTime, time.Local)
	end = end.Add(time.Duration(4) * time.Hour)
	now := time.Now()
	if end.Before(now) {
		snifferModel.SniffTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
		snifferModel.Update()
		return true
	}
	return false
}

// 获取task
func (s Sniffer) GetTaskById(ctx *gin.Context, taskId int) (*mysql_model.Task, error) {
	task, has := mysql_model.TaskFindById(taskId)
	if !has {
		return task, errors.New("任务不存在")
	}
	if task.Status == TASK_COMPLETE {
		return task, errors.New("该任务已完成，不能再测试了")
	}
	return task, nil
}

func (s Sniffer) UpdateTaskStatus(ctx *gin.Context, taskId, status int) error {
	_, err := mysql_model.UpdateStatus(taskId, status)
	return err
}

func (s Sniffer) GSMSnifferScanner(band string) (bool, error) {
	gsm := GSM{}
	args := map[string]string{"band": band}
	url, err := GetHost(GsmSniffer)
	if err != nil {
		return false, err
	}
	url += URL_SCANNER
	blog.Info(URL_SCANNER, zap.Any("url:", url), zap.Any("args:", args))
	resp, err := custom_util.HttpPostJson(nil, args, url)
	blog.Info(URL_SCANNER, zap.Any("resp:", string(resp)))
	if err != nil {
		return false, err
	}
	// resp = []byte(`{"status": true, "message_id": 1, "message": "Start scanner successfully"}`)
	err = json.Unmarshal(resp, &gsm)
	blog.Info(URL_SCANNER, zap.Any("gsm:", gsm))
	if err != nil {
		return false, err
	}
	if gsm.MessageId == 1 || gsm.MessageId == 2 {
		return true, nil
	}
	return false, errors.New(gsm.Message)
}

func (s Sniffer) GSMSnifferGetfreq() (GSM, error) {
	gsm := GSM{}
	url, err := GetHost(GsmSniffer)
	if err != nil {
		return gsm, err
	}
	url += URL_GETFREQ
	blog.Info(URL_GETFREQ, zap.Any("url:", url))
	resp, err := custom_util.HttpPostJson(nil, nil, url)
	blog.Info(URL_GETFREQ, zap.Any("resp:", string(resp)))
	if err != nil {
		return gsm, err
	}
	// resp = []byte(`{"status":true,"message_id":1,"message":"Get frequency success.","data":[["40","943.0M","41285","4421","460","0","-58"],["50","953.0M","41285","4421","460","0","-58"]]}`)
	err = json.Unmarshal(resp, &gsm)
	blog.Info(URL_GETFREQ, zap.Any("gsm:", gsm))
	if err != nil {
		return gsm, err
	}
	return gsm, nil
}

// 嗅探
func (s Sniffer) GSMSniffer(freq, mode string) (GSM, error) {
	gsm := GSM{}
	args := map[string]string{"freq": freq, "mode": mode}
	url, err := GetHost(GsmSniffer)
	if err != nil {
		return gsm, err
	}
	url += URL_SNIFFER
	blog.Info(URL_SNIFFER, zap.Any("url:", url), zap.Any("args:", args))
	resp, err := custom_util.HttpPostJson(nil, args, url)
	blog.Info(URL_SNIFFER, zap.Any("resp:", string(resp)))
	if err != nil {
		return gsm, err
	}
	// resp = []byte(`{"status": true, "message_id": 1, "message": "Start scanner successfully"}`)
	err = json.Unmarshal(resp, &gsm)
	blog.Info(URL_SNIFFER, zap.Any("gsm:", gsm))
	if err != nil {
		return gsm, err
	}
	return gsm, nil
}

func (s Sniffer) GSMSnifferGetimsi() (GSM, error) {
	gsm := GSM{}
	url, err := GetHost(GsmSniffer)
	if err != nil {
		return gsm, err
	}
	url += URL_GETIMSI
	blog.Info(URL_GETIMSI, zap.Any("url:", url))
	resp, err := custom_util.HttpPostJson(nil, nil, url)
	blog.Info(URL_GETIMSI, zap.Any("resp:", string(resp)))
	if err != nil {
		return gsm, err
	}
	// resp = []byte(`{"status":true,"message_id":1,"message":"Get frequency success.","data":[["460028106055249","Sep 23, 2021 06:42:47.778473844 UTC"],["160028106055248","Sep 23, 2022 06:42:47.778473844 UTC"]]}`)
	err = json.Unmarshal(resp, &gsm)
	blog.Info(URL_GETIMSI, zap.Any("gsm:", gsm))
	if err != nil {
		return gsm, err
	}
	return gsm, nil
}

func (s Sniffer) GSMSnifferGetsms() (GSM, error) {
	gsm := GSM{}
	url, err := GetHost(GsmSniffer)
	if err != nil {
		return gsm, err
	}
	url += URL_GETSMS
	blog.Info(URL_GETSMS, zap.Any("url:", url))
	resp, err := custom_util.HttpPostJson(nil, nil, url)
	blog.Info(URL_GETSMS, zap.Any("resp:", string(resp)))
	if err != nil {
		return gsm, err
	}
	// resp = []byte(`{"status":true,"message_id":1,"message":"Get frequency success.","data":[["我是短信内容","Sep 23, 2021 06:42:47.778473844 UTC"],["北京欢迎您","Sep 23, 2022 06:42:47.778473844 UTC"]]}`)
	err = json.Unmarshal(resp, &gsm)
	blog.Info(URL_GETSMS, zap.Any("gsm:", gsm))
	if err != nil {
		return gsm, err
	}
	return gsm, nil
}

func (s Sniffer) GSMSnifferStop() error {
	gsm := GSM{}
	url, err := GetHost(GsmSniffer)
	if err != nil {
		return err
	}
	url += URL_STOP
	blog.Info(URL_STOP, zap.Any("url:", url))
	resp, err := custom_util.HttpPostJson(nil, nil, url)
	blog.Info(URL_STOP, zap.Any("resp:", string(resp)))
	err = json.Unmarshal(resp, &gsm)
	blog.Info(URL_STOP, zap.Any("gsm:", gsm))
	if err != nil {
		return err
	}
	if gsm.MessageId == 0 {
		return errors.New(gsm.Message)
	}
	return nil
}
