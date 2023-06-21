package scanner

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/guardian/util"

	pkgerr "github.com/pkg/errors"

	"skygo_detection/common"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

type FirmwareScanner struct{}

func (this FirmwareScanner) Run(taskInfo qmap.QM) (int, error) {
	params := qmap.QM{
		"e_id": taskInfo.MustInt("scanner_id"),
	}
	if has, scannerTask := sys_service.NewSessionWithCond(params).GetOne(new(mysql_model.FirmwareTask)); has {
		var err error
		var sleepTime int
		// println("扫描任务信息:", scannerTask.ToString())
		switch scannerTask.MustInt("status") {
		case common.FIRMWARE_STATUS_PROJECT_CREATE:
			// 待创建固件扫描项目，开始创建固件扫描项目
			println("创建固件项目")
			sleepTime, err = this.createProject(scannerTask)
			if err == nil {
				// 项目创建成功，更新任务状态为:固件下载
				this.UpdateScannerTaskStatus(taskInfo, common.FIRMWARE_STATUS_FIRMWARE_DOWNLOAD, "更新任务状态为:固件下载", "")
				return sleepTime, err
			}
		case common.FIRMWARE_STATUS_FIRMWARE_DOWNLOAD:
			// 获取固件包下载进度
			println("获取固件包下载进度")
			sleepTime, err = this.getFirmwareDownloadProcess(scannerTask)
			if err == nil && sleepTime == 0 {
				// 固件下载成功，更新任务状态为:任务创建
				this.UpdateScannerTaskStatus(taskInfo, common.FIRMWARE_STATUS_TASK_CREATE, "更新任务状态为:任务创建", "")
			}
		case common.FIRMWARE_STATUS_TASK_CREATE:
			// 固件任务取消
			println("固件任务创建")
			sleepTime, err = this.createFirmwareTask(scannerTask)
			if err == nil && sleepTime == 0 {
				// 固件任务创建成功，更新任务状态为:任务创建
				this.UpdateScannerTaskStatus(taskInfo, common.FIRMWARE_STATUS_TASK_START, "更新任务状态为:任务启动", "")
			}
		case common.FIRMWARE_STATUS_TASK_START:
			// 固件任务启动
			println("固件任务启动")
			sleepTime, err = this.startFirmwareTask(scannerTask)
			if err == nil && sleepTime == 0 {
				// 固件任务启动成功，更新任务状态为:任务扫描中
				this.UpdateScannerTaskStatus(taskInfo, common.FIRMWARE_STATUS_TASK_SCANNING, "更新任务状态为:任务扫描中", "")
			}
		case common.FIRMWARE_STATUS_TASK_CANCEL:
			// 固件扫描任务取消
			println("固件任务取消")
			sleepTime, err = this.stopFirmwareTask(scannerTask)
			if err == nil && sleepTime == 0 {
				// 固件任务启动成功，更新任务状态为:扫描失败
				this.UpdateScannerTaskStatus(taskInfo, common.FIRMWARE_STATUS_SCAN_FAILURE, "更新任务状态为:扫描失败", "")
			}
			return -1, errors.New("固件扫描任务取消")
		case common.FIRMWARE_STATUS_TASK_SCANNING:
			// 固件扫描中，获取固件扫描进度
			println("检测扫描进度.......")
			sleepTime, err = this.scanStage(scannerTask)
			if err == nil && sleepTime == 0 {
				// 固件扫描成功，更新任务状态为:扫描完成(报告解析)
				this.UpdateScannerTaskStatus(taskInfo, common.FIRMWARE_STATUS_REPORT_ANALYSIS, "更新任务状态为:扫描完成(报告解析)", "")
			}
		case common.FIRMWARE_STATUS_REPORT_ANALYSIS:
			// 固件扫描成功，获取固件扫描报告
			println("获取固件扫描报告")
			sleepTime, err = this.fetchReport(scannerTask)
			if err == nil {
				// 固件报告解析完成，更新任务状态为:扫描成功
				this.UpdateScannerTaskStatus(taskInfo, common.FIRMWARE_STATUS_SCAN_SUCCESS, "更新任务状态为:扫描成功", "")
			}
		default:
			return -1, errors.New(fmt.Sprintf("未知固件扫描任务状态，状态码:%d", scannerTask.MustInt("status")))
		}
		if err != nil && sleepTime < 0 {
			println("更新固件扫描状态:失败")
			this.UpdateScannerTaskStatus(taskInfo, common.FIRMWARE_STATUS_SCAN_FAILURE, "更新任务状态为:扫描失败", err.Error())
		}
		return sleepTime, pkgerr.WithStack(err)
	} else {
		return -1, errors.New("未发现此任务")
	}
}

// 创建固件扫描项目
func (this FirmwareScanner) createProject(firmWareData *qmap.QM) (int, error) {
	body := qmap.QM{
		"project_name":     firmWareData.MustString("name"),
		"device_name":      firmWareData.MustString("device_name"),
		"device_model":     firmWareData.MustString("device_model"),
		"firmware_version": firmWareData.MustString("firmware_version"),
		"device_type":      firmWareData.MustString("device_type"),
	}
	req, err := http.NewRequest("POST", service.LoadFirmwareConfig().ScanHost+"/api/projects", bytes.NewBuffer([]byte(body.ToString())))
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		println("项目创建失败")
		return -1, errors.New(fmt.Sprintf("项目创建失败，%s", err.Error()))
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	respQM, err := qmap.NewWithString(string(respBody))
	if err != nil {
		println("项目创建失败")
		return -1, errors.New(fmt.Sprintf("项目创建失败，%s", err.Error()))
	}
	code := respQM.Int("code")
	if code != 1000 {
		return -1, errors.New(fmt.Sprintf("项目创建失败，错误码:%d", code))
	}
	projectId := respQM.Map("data").Int("id")
	// 创建项目完成后，向固件扫描服务提交固件下载url
	if err := this.postFirmwareDownloadUrl(firmWareData.MustString("file_id"), projectId); err != nil {
		println("推送url失败")
		return -1, errors.New("推送url失败")
	}
	// 上传完url,更新项目id
	updateTask := qmap.QM{
		"yafaf_project_id": projectId,
	}
	if err := new(mysql_model.FirmwareTask).Update(firmWareData.Int("id"), updateTask); err != nil {
		return -1, err
	}
	return 1, nil
}

// 向附件扫描服务提交固件下载url
func (this FirmwareScanner) postFirmwareDownloadUrl(fileId string, projectId int) error {
	downloadUrl := fmt.Sprintf("%s/api/v1/firmware/download?name=%s", service.LoadFirmwareConfig().AdminHost, fileId)
	request := qmap.QM{
		"pid": fmt.Sprintf("%d", projectId),
		"url": downloadUrl,
	}
	url := service.LoadFirmwareConfig().ScanHost + "/api/remotedownload"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(request.ToString())))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respStr, _ := ioutil.ReadAll(resp.Body)
	if respQM, err := qmap.NewWithString(string(respStr)); err == nil {
		if code := respQM.Int("code"); code != 1000 {
			return errors.New(fmt.Sprintf("post firmware url error, error code : %d", code))
		}
	} else {
		return err
	}
	return nil
}

// 获取下载状态
func (this FirmwareScanner) getFirmwareDownloadProcess(firmwareTask *qmap.QM) (int, error) {
	yafafProjectId := firmwareTask.MustInt("yafaf_project_id")
	url := fmt.Sprintf("%s/api/remotedownload?pid=%d", service.LoadFirmwareConfig().ScanHost, yafafProjectId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respStr, _ := ioutil.ReadAll(resp.Body)
	respQM, err := qmap.NewWithString(string(respStr))
	if err != nil {
		return -1, err
	}
	if code := respQM.Int("code"); code != 1000 {
		return -1, errors.New(fmt.Sprintf("固件下载失败，错误码为:%d", code))
	}
	status := respQM.Map("data").String("status")
	// 如果解析引擎下载完成
	if status == "finished" {
		println("检测下载任务为finished")
		updateTask := qmap.QM{
			"yafaf_download_path": respQM.Map("data").String("path"),
		}
		if err := new(mysql_model.FirmwareTask).Update(firmwareTask.Int("id"), updateTask); err != nil {
			return -1, err
		}
		return 0, nil
	} else if status == "none" {
		println("检测下载状态为none")
		return -1, errors.New("检测下载状态为:none")
	} else if status == "downloading" {
		println("检测下载状态为:downloading")
		return 3, nil
	} else {
		println("检测下载状态为:error")
		return -1, errors.New("固件下载失败")
	}
}

// 创建固件任务
func (this FirmwareScanner) createFirmwareTask(firmwareTask *qmap.QM) (int, error) {
	templateId := firmwareTask.Int("template_id")
	templateName := this.getTemplateName(templateId)
	yafafId, err := this.createTask(firmwareTask.MustString("name"), firmwareTask.String("yafaf_download_path"), firmwareTask.MustInt("yafaf_project_id"), templateId, templateName)
	if err != nil {
		return -1, err
	}
	if yafafId <= 0 {
		println(fmt.Sprintf("任务创建失败，yafaf_id:%d", yafafId))
		return -1, errors.New(fmt.Sprintf("任务创建失败，yafaf_id:%d", yafafId))
	} else {
		println(fmt.Sprintf("任务创建成功，yafaf_id:%d", yafafId))
	}
	updateTask := qmap.QM{
		"yafaf_id": yafafId,
	}
	if err := new(mysql_model.FirmwareTask).Update(firmwareTask.Int("id"), updateTask); err != nil {
		return -1, err
	}
	return 0, nil
}

/*
 * 创建扫描任务
 */
func (this FirmwareScanner) createTask(projectName, downloadPath string, projectId int, templateId int, templateName string) (int, error) {
	requestBody := qmap.QM{
		"task_name":     fmt.Sprintf("%s_%s", projectName, time.Now().Format("2006-01-02")),
		"project_name":  projectName,
		"project_id":    fmt.Sprintf("%d", projectId),
		"status":        "ready",
		"file_type":     "multifile",
		"template_id":   templateId,
		"template_name": templateName,
		"tempfile":      downloadPath,
		"workpath":      "",
	}
	url := service.LoadFirmwareConfig().ScanHost + "/api/tasks"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody.ToString())))
	req.Header.Set("Content-Type", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respStr, _ := ioutil.ReadAll(resp.Body)
	if respQM, err := qmap.NewWithString(string(respStr)); err == nil {
		if code := respQM.Int("code"); code == 1000 {
			return respQM.Map("data").Int("id"), nil
		} else {
			return -1, errors.New(fmt.Sprintf("任务创建创建失败，错误码:%d", code))
		}
	} else {
		panic(err)
	}
}

// 启动固件扫描任务
func (this FirmwareScanner) startFirmwareTask(firmwareTask *qmap.QM) (int, error) {
	return this.changeTaskStatus(firmwareTask.Int("yafaf_id"), "start")
}

// 停止固件扫描任务
func (this FirmwareScanner) stopFirmwareTask(firmwareTask *qmap.QM) (int, error) {
	return this.changeTaskStatus(firmwareTask.Int("yafaf_id"), "stop")
}

/*
 * 更改固件扫描任务状态 start stop
 */
func (this FirmwareScanner) changeTaskStatus(yafafId int, status string) (int, error) {
	requestBody := qmap.QM{
		"id":        yafafId,
		"operation": status,
	}
	url := service.LoadFirmwareConfig().ScanHost + "/api/tasks"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(requestBody.ToString())))
	req.Header.Set("Content-Type", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	respBodyQM, err := qmap.NewWithString(string(respBody))
	if err != nil {
		panic(err)
	}
	if respBodyQM.Int("code") == 1000 {
		return 0, nil
	} else {
		return -1, errors.New("firmware task start error")
	}
}

// 查询固件包扫描进度
func (this FirmwareScanner) scanStage(firmWareData *qmap.QM) (int, error) {
	yafafId := firmWareData.MustInt("yafaf_id")
	// 获取任务状态
	url := fmt.Sprintf("%s/api/tasks/status?task_id=%d", service.LoadFirmwareConfig().ScanHost, yafafId)
	result, err := getResultFromUrl(url)
	if err != nil {
		println("解析引擎API api/tasks/status 异常")
		return -1, errors.New("解析引擎API api/tasks/status 异常")
	}
	respQM, err := qmap.NewWithString(string(result))
	if err != nil {
		panic(err)
	}
	status := respQM.Map("data").String("status")
	if status != "completed" {
		println("任务检测进行中，检测状态为:" + status)
		return 3, nil
	} else {
		println("任务检测完毕")
		return 0, nil
	}
}

func getResultFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (this FirmwareScanner) fetchReport(firmwareTask *qmap.QM) (int, error) {
	yafafId := firmwareTask.MustInt("yafaf_id")
	resultUrl := fmt.Sprintf("%s/api/tasks/result?task_id=%d&format=json&plugin_name=", service.LoadFirmwareConfig().ScanHost, yafafId)
	result, err := getResultFromUrl(resultUrl)
	if err != nil {
		return -1, errors.New("报告获取失败")
	}
	resultQM, err := qmap.NewWithString(string(result))
	if err != nil {
		return -1, err
	}
	updateTask := qmap.QM{
		"source_report": string(result),
	}
	if err := new(mysql_model.FirmwareTask).Update(firmwareTask.Int("id"), updateTask); err != nil {
		return -1, err
	}
	err = this.analysisReport(firmwareTask, &resultQM)
	return -1, err
}

func (this FirmwareScanner) analysisReport(taskInfo, report *qmap.QM) error {
	cveCount := 0
	binaryCount := 0
	linuxCount := 0
	passRiskCount := 0
	overCertCount := 0
	leaksCount := 0
	passSuspectCount := 0
	apkRiskCount := 0
	rtsCategory := new(mysql_model.FirmwareReportRtsCategory)
	rtsCategory.ScannerId = taskInfo.Int("id")
	println("开始报告解析:")
	reportData := report.Map("data")
	for k, _ := range reportData {
		switch k {
		case "init_files":
			println("报告解析:init_files")
			rtsCategory.InitFiles = 1
			this.initFilesAnalysisReport(*taskInfo, reportData.Map("init_files"))
		case "elf_scanner":
			println("报告解析:elf_scanner")
			binaryCount = this.elfScannerAnalysisReport(*taskInfo, reportData.Map("elf_scanner"))
			if binaryCount > 0 {
				rtsCategory.ElfScanner = 1
				rtsCategory.BinaryHardening = 1
			}
		case "leaks_scanner":
			println("报告解析:leaks_scanner")
			leaksCount = this.leaksScannerAnalysisReport(*taskInfo, reportData.Map("leaks_scanner"))
			if leaksCount > 0 {
				rtsCategory.LeaksScanner = 1
			}
		case "certificates_scanner":
			println("报告解析:certificates_scanner")
			overCertCount = this.certificatesScannerAnalysisReport(*taskInfo, reportData.Map("certificates_scanner"))
			if overCertCount > 0 {
				rtsCategory.CertificatesScanner = 1
			}
		case "password_scanner":
			println("报告解析:password_scanner")
			passRiskCount, passSuspectCount = this.passwordScannerAnalysisReport(*taskInfo, reportData.Map("password_scanner"))
			if (passRiskCount + passSuspectCount) > 0 {
				rtsCategory.PasswordScanner = 1
			}
		case "linux_basic_audit":
			println("报告解析:linux_basic_audit")
			linuxCount = this.linuxBasicAuditScannerAnalysisReport(*taskInfo, reportData.Map("linux_basic_audit"))
			if linuxCount > 0 {
				rtsCategory.LinuxBasicAudit = 1
			}
		case "version_scanner":
			println("报告解析:version_scanner")
			cveCount = this.versionScannerAnalysisReport(*taskInfo, reportData.Interface("version_scanner").([]interface{}))
			if cveCount > 0 {
				rtsCategory.VersionScanner = 1
			}
		case "apk_info":
			println("报告解析:apk_info")
			apkRiskCount = this.apkInfoScannerAnalysisReport(*taskInfo, reportData)
			count := 0
			for _, v := range reportData.Map("apk_info") {
				count += len(v.(map[string]interface{}))
			}
			if count > 0 {
				rtsCategory.ApkInfo = 1
			}
		case "apk_sensitive_info":
			println("报告解析:apk_sensitive_info")
			this.apkSensitiveScannerAnalysisReport(*taskInfo, reportData.Map("apk_sensitive_info"))
			count := 0
			for _, v := range reportData.Map("apk_sensitive_info") {
				count += len(v.(map[string]interface{}))
			}
			if count > 0 {
				rtsCategory.ApkSensitiveInfo = 1
			}
		case "apk_common_vul":
			println("报告解析:apk_common_vul")
			this.apkCommonVulAnalysisReport(*taskInfo, reportData.Map("apk_common_vul"))
			count := 0
			for _, v := range reportData.Map("apk_sensitive_info") {
				count += len(v.(map[string]interface{}))
			}
			if count > 0 {
				rtsCategory.ApkCommonVul = 1
			}
		}
	}
	if _, err := sys_service.NewOrm().InsertOne(rtsCategory); err != nil {
		panic(err)
	}
	rtsRisk := new(mysql_model.FirmwareReportRtsRisk)
	rtsRisk.ScannerId = taskInfo.Int("id")
	rtsRisk.RiskCount = cveCount + binaryCount + linuxCount + passRiskCount + overCertCount
	rtsRisk.PassRiskCount = passRiskCount
	rtsRisk.BinaryCount = binaryCount
	rtsRisk.LinuxCount = linuxCount
	rtsRisk.OverCertCount = overCertCount
	rtsRisk.RiskSuspectCount = leaksCount + passSuspectCount
	if apkRiskCount > 0 {
		rtsRisk.RiskCount = apkRiskCount
	}
	if _, err := sys_service.NewOrm().InsertOne(rtsRisk); err != nil {
		panic(err)
	}
	println("报告解析完毕")
	return nil
}

func (this FirmwareScanner) initFilesAnalysisReport(firmwareTask, initFile qmap.QM) {
	reportInit := new(mysql_model.FirmwareReportInit)
	reportInit.ScannerId = firmwareTask.Int("id")
	reportInit.TotalSize = initFile.String("total_size")
	reportInit.CreateTime = int(time.Now().Unix())
	if fileInfo, has := initFile.TryMap("file_info"); has {
		reportInit.DirNum = fileInfo.Int("dir_num")
		reportInit.FileNum = fileInfo.Int("file_num")
		reportInit.LinkNum = fileInfo.Int("link_num")
		reportInit.NodeNum = fileInfo.Int("node_num")
	}
	if systemGuss, has := initFile.TryMap("system_guss"); has {
		if printable, has := systemGuss.TryMap("printable"); has {
			if arch, has := printable.TryMap("arch"); has {
				reportInit.Arch = arch.ToString()
			}
			if system, has := printable.TryMap("system"); has {
				reportInit.System = system.ToString()
			}
		}
	}
	for k, v := range initFile {
		switch k {
		case "compressed":
			reportInit.Compressed = qmap.QM(v.(map[string]interface{})).ToString()
		case "disk_usage":
			reportInit.DiskUsage = util.SliceToString(v.([]interface{}))
		case "dulplicated":
			reportInit.Dulplicated = qmap.QM(v.(map[string]interface{})).ToString()
		case "filesystem":
			reportInit.Filesystem = qmap.QM(v.(map[string]interface{})).ToString()
		case "firmware":
			reportInit.Firmware = qmap.QM(v.(map[string]interface{})).ToString()
		case "system_guss":
			reportInit.SystemGuess = qmap.QM(v.(map[string]interface{})).ToString()
		}
	}
	if _, err := sys_service.NewOrm().InsertOne(reportInit); err != nil {
		panic(err)
	}
}

func (this FirmwareScanner) elfScannerAnalysisReport(firmwareTask, elfScanner qmap.QM) int {
	binaryCount := 0
	elf := new(mysql_model.FirmwareReportElf)
	elf.ScannerId = firmwareTask.Int("id")
	elf.CreateTime = int(time.Now().Unix())
	for k, v := range elfScanner {
		switch k {
		case "executable":
			elf.Executable = len(v.([]interface{}))
		case "kernel_module":
			elf.KernelModule = len(v.([]interface{}))
		case "shared_lib":
			elf.SharedLib = len(v.([]interface{}))
		}
		if k == "binary_hardening" {
			binaryCount = len(v.([]interface{}))
			var nx, pie, relro, canary, stripped int
			for _, item := range v.([]interface{}) {
				itemQM := qmap.QM(item.(map[string]interface{}))
				rtsBinary := new(mysql_model.FirmwareReportRtsBinary)
				rtsBinary.Type = "binary_hardening"
				rtsBinary.ScannerId = firmwareTask.Int("id")
				rtsBinary.IsElf = 1
				if itemQM.Bool("hardenable") == true {
					rtsBinary.IsDoubt = 1
				} else {
					rtsBinary.IsDoubt = 0
				}
				rtsBinary.FileName = itemQM.String("file_name")
				rtsBinary.FullPath = itemQM.String("full_path")
				rtsBinary.MagicInfo = itemQM.String("magic_info")
				rtsBinary.RelaPath = itemQM.String("rela_path")
				rtsBinary.Result = itemQM.Map("result").ToString()
				rtsBinary.CreateTime = int(time.Now().Unix())
				{
					binaryResult := itemQM.Map("result")
					if binaryResult.String("nx") == "yes" {
						nx++
					}
					if binaryResult.String("pie") == "yes" {
						pie++
					}
					if binaryResult.String("relro") == "yes" {
						relro++
					}
					if binaryResult.String("canary") == "yes" {
						canary++
					}
					if binaryResult.String("stripped") == "yes" {
						stripped++
					}
				}
				if _, err := sys_service.NewOrm().InsertOne(rtsBinary); err != nil {
					panic(err)
				}
			}
			rtsBinaryTotal := new(mysql_model.FirmwareReportRtsBinaryTotal)
			rtsBinaryTotal.ScannerId = firmwareTask.Int("id")
			rtsBinaryTotal.CreateTime = int(time.Now().Unix())
			rtsBinaryTotal.Count = binaryCount
			rtsBinaryTotal.Canary = canary
			rtsBinaryTotal.Nx = nx
			rtsBinaryTotal.Pie = pie
			rtsBinaryTotal.Relro = relro
			rtsBinaryTotal.Stripped = stripped
			if _, err := sys_service.NewOrm().InsertOne(rtsBinaryTotal); err != nil {
				panic(err)
			}
		} else if k == "executable" || k == "shared_lib" || k == "kernel_module" {
			for _, item := range v.([]interface{}) {
				rtsElf := new(mysql_model.FirmwareReportRtsElf)
				rtsElf.ScannerId = firmwareTask.Int("id")
				rtsElf.Type = k
				rtsElf.Executable = item.(string)
				rtsElf.CreateTime = int(time.Now().Unix())
				if _, err := sys_service.NewOrm().InsertOne(rtsElf); err != nil {
					panic(err)
				}
			}
		}
	}
	if _, err := sys_service.NewOrm().InsertOne(elf); err != nil {
		panic(err)
	}
	return binaryCount
}

func (this FirmwareScanner) leaksScannerAnalysisReport(firmwareTask, leaksScanner qmap.QM) int {
	leaksCount := 0
	for k, v := range leaksScanner {
		leaksCount += len(v.([]interface{}))
		for _, item := range v.([]interface{}) {
			itemQM := qmap.QM(item.(map[string]interface{}))
			rtsLeaks := new(mysql_model.FirmwareReportRtsLeaks)
			rtsLeaks.Type = k
			rtsLeaks.ScannerId = firmwareTask.Int("id")
			rtsLeaks.Info = itemQM.String("info")
			rtsLeaks.Path = itemQM.String("path")
			rtsLeaks.FullPath = itemQM.String("full_path")
			rtsLeaks.Origin = itemQM.String("origin")
			rtsLeaks.CreateTime = int(time.Now().Unix())
			if _, err := sys_service.NewOrm().InsertOne(rtsLeaks); err != nil {
				panic(err)
			}
		}
	}
	return leaksCount
}

func (this FirmwareScanner) certificatesScannerAnalysisReport(firmwareTask, certificatesScanner qmap.QM) int {
	overCertCount := 0
	certTotal := new(mysql_model.FirmwareReportRtsCertTotal)
	certTotal.ScannerId = firmwareTask.Int("id")
	certTotal.CreateTime = int(time.Now().Unix())
	if pk, has := certificatesScanner.TrySlice("private_key"); has {
		certTotal.PrivateKeyCount = len(pk)
	} else {
		certTotal.PrivateKeyCount = 0
	}
	if certs, has := certificatesScanner.TrySlice("certificate"); has {
		certTotal.CertificateCount = len(certs)
		for _, cert := range certs {
			certQM := qmap.QM(cert.(map[string]interface{}))
			if info, has := certQM.TryMap("info"); has {
				if jsonQM, has := info.TryMap("json"); has {
					expireTime := jsonQM.Map("Validity").String("Not After")
					stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", expireTime, time.Local)
					if stamp.Unix() < time.Now().Unix() {
						overCertCount++
					}
				}
			}
		}
		if overCertCount > 0 {
			certTotal.Part = "Cert"
			certTotal.Level = 3
		}
		certTotal.CertificateOverdateCount = overCertCount
		if _, err := sys_service.NewOrm().InsertOne(certTotal); err != nil {
			panic(err)
		}
	}
	for k, v := range certificatesScanner {
		if v == nil {
			continue
		}
		for _, item := range v.([]interface{}) {
			itemQM := qmap.QM(item.(map[string]interface{}))
			rtsCert := new(mysql_model.FirmwareReportRtsCert)
			rtsCert.Type = k
			rtsCert.ScannerId = firmwareTask.Int("id")
			rtsCert.Path = itemQM.String("path")
			rtsCert.FileName = itemQM.String("file_name")
			rtsCert.Content = itemQM.String("content")
			rtsCert.Info = `{}`
			if info, has := itemQM.TryInterface("info"); has {
				switch info.(type) {
				case map[string]interface{}, qmap.QM:
					rtsCert.Info = itemQM.Map("info").ToString()
				}
			}
			rtsCert.CreateTime = int(time.Now().Unix())
			if _, err := sys_service.NewOrm().InsertOne(rtsCert); err != nil {
				panic(err)
			}
		}
	}
	return overCertCount
}

func (this FirmwareScanner) passwordScannerAnalysisReport(firmwareTask, passwordScanner qmap.QM) (int, int) {
	passRiskCount := 0
	passSuspectCount := 0
	if printable, has := passwordScanner.TryInterface("printable"); has {
		for _, item := range printable.([]interface{}) {
			itemQM := qmap.QM(item.(map[string]interface{}))
			rtsPwd := new(mysql_model.FirmwareReportRtsPwd)
			rtsPwd.Type = "printable"
			rtsPwd.ScannerId = firmwareTask.Int("id")
			rtsPwd.Path = itemQM.String("path")
			rtsPwd.FullPath = itemQM.String("full_path")
			rtsPwd.Flag = itemQM.Map("flag").ToString()
			rtsPwd.Content = itemQM.String("content")
			rtsPwd.CreateTime = int(time.Now().Unix())
			if _, err := sys_service.NewOrm().InsertOne(rtsPwd); err != nil {
				panic(err)
			}
			if itemQM.Map("flag").String("level") == "warning" {
				passRiskCount++
			} else {
				passSuspectCount++
			}
		}
	}
	return passRiskCount, passSuspectCount
}

func (this FirmwareScanner) linuxBasicAuditScannerAnalysisReport(firmwareTask, auditScanner qmap.QM) int {
	linuxCount := 0
	linuxTotal := new(mysql_model.FirmwareReportRtsLinuxTotal)
	linuxTotal.ScannerId = firmwareTask.Int("id")
	linuxTotal.CreateTime = int(time.Now().Unix())
	for k, v := range auditScanner {
		if v == nil || len(v.([]interface{})) == 0 {
			continue
		}
		for _, vv := range v.([]interface{}) {
			var vvQM qmap.QM = vv.(map[string]interface{})
			fullPath := vvQM.String("fullpath")
			if list := vvQM.Slice("list"); len(list) > 0 {
				linuxCount += len(list)
				for _, item := range list {
					itemQM := qmap.QM(item.(map[string]interface{}))
					rtsLinux := new(mysql_model.FirmwareReportRtsLinux)
					rtsLinux.Type = k
					rtsLinux.ScannerId = firmwareTask.Int("id")
					rtsLinux.FullPath = fullPath
					rtsLinux.Detail = itemQM.ToString()
					rtsLinux.CreateTime = int(time.Now().Unix())
					if _, err := sys_service.NewOrm().InsertOne(rtsLinux); err != nil {
						panic(err)
					}
				}
			}
		}
	}
	linuxTotal.AbnormalCount = linuxCount
	if _, err := sys_service.NewOrm().InsertOne(linuxTotal); err != nil {
		panic(err)
	}
	return linuxCount
}

func (this FirmwareScanner) versionScannerAnalysisReport(firmwareTask qmap.QM, versionScanner []interface{}) int {
	cveCount := 0
	for _, scan := range versionScanner {
		scanQM := qmap.QM(scan.(map[string]interface{}))
		fileName := scanQM.String("file_name")
		path := scanQM.String("path")
		version := scanQM.String("version")
		var no, low, mid, high, heavy int
		for _, cve := range scanQM.Interface("cve_list").(map[string]interface{}) {
			cveCount++
			cveQM := qmap.QM(cve.(map[string]interface{}))
			cveItem := new(mysql_model.FirmwareReportRtsCve)
			cveItem.ScannerId = firmwareTask.Int("id")
			cveItem.Type = "version_scanner"
			cveItem.CreateTime = int(time.Now().Unix())
			cveItem.Path = path
			cveItem.FileName = fileName
			cveItem.Version = version
			cveItem.Cve = cveQM.String("cve")
			cveItem.Vendor = cveQM.String("vendor")
			cveItem.VersionEndExcluding = cveQM.String("version_end_excluding")
			cveItem.VersionEndIncluding = cveQM.String("version_end_including")
			cveItem.VersionStartExcluding = cveQM.String("version_start_excluding")
			cveItem.VersionStartIncluding = cveQM.String("version_start_including")
			cvssv3 := cveQM.Map("cvssv3")
			cvssv2 := cveQM.Map("cvssv2")
			if cvssv3 != nil {
				cveItem.Cvssv3 = cvssv3.ToString()
			}
			if cvssv2 != nil {
				cveItem.Cvssv2 = cvssv2.ToString()
			}
			cveItem.Cvssv2score = 0.1
			if cvssv3 != nil {
				if cvssV3, has := cvssv3.TryMap("cvssV3"); has {
					cveItem.Cvssv2score = cvssV3.Float32("baseScore")
					cveItem.Vector = cvssV3.String("attackVector")
				}
			} else if cvssv2 != nil {
				if cvssV2, has := cvssv2.TryMap("cvssV2"); has {
					cveItem.Cvssv2score = cvssV2.Float32("baseScore")
					cveItem.Vector = cvssV2.String("attackVector")
				}
			}
			score := cveItem.Cvssv2score * 10
			if score < 1 {
				no++
			} else if score >= 1 && score <= 39 {
				low++
			} else if score >= 40 && score <= 69 {
				mid++
			} else if score >= 70 && score <= 89 {
				high++
			} else if score >= 90 && score <= 100 {
				heavy++
			}
			cveItem.Level = this.culLevel(cveItem.Cvssv2score)
			if desc, has := cveQM.TryInterface("description"); has && len(desc.([]interface{})) > 0 {
				cveItem.Description = util.SliceToString(desc.([]interface{}))
			}
			if _, err := sys_service.NewOrm().InsertOne(cveItem); err != nil {
				panic(err)
			}
		}
		cveTotal := new(mysql_model.FirmwareReportRtsCveTotal)
		cveTotal.ScannerId = firmwareTask.Int("id")
		cveTotal.CreateTime = int(time.Now().Unix())
		cveTotal.No = no
		cveTotal.Low = low
		cveTotal.Mid = mid
		cveTotal.High = high
		cveTotal.Heavy = heavy
		if _, err := sys_service.NewOrm().InsertOne(cveTotal); err != nil {
			panic(err)
		}
	}
	return cveCount
}

func (this FirmwareScanner) apkInfoScannerAnalysisReport(firmwareTask, apkScanner qmap.QM) int {
	var no, low, mid, high, heavy, level, total, soPie, soDebug, soCanary, urlCount, ipCount, tokenCount, akCount, emailCount, certCount int
	vQM := qmap.QM{}
	if apkScanner["apk_info"] != nil {
		for apkName, v := range apkScanner["apk_info"].(map[string]interface{}) {
			vQM = v.(map[string]interface{})
			vQM["pkg_name"] = apkName
			for k, val := range v.(map[string]interface{}) {
				if valData, ok := val.(map[string]interface{}); ok {
					if valData["risk_level"] != nil && valData["detail"] != nil {
						if _, ok := valData["detail"].([]interface{}); !ok {
							continue
						}
						level, count := this.calLevel(valData["risk_level"].(string), valData["detail"].([]interface{}))
						switch level {
						case 0:
							no += count
						case 1:
							low += count
							total += count
						case 2:
							mid += count
							total += count
						case 3:
							high += count
							total += count
						case 4:
							heavy += count
							total += count
						}
						if k == "so_pie" {
							soPie = count
						} else if k == "so_debug" {
							soDebug = count
						} else if k == "so_canary" {
							soCanary = count
						}
					}
				}
			}
		}
	}

	if apkScanner["apk_common_vul"] != nil {
		for _, v := range apkScanner["apk_common_vul"].(map[string]interface{}) {
			for _, val := range v.(map[string]interface{}) {
				for kk, vv := range val.(map[string]interface{}) {
					if kk == "webview_remote_excute" {
						continue
					}
					if valData, ok := vv.(map[string]interface{}); ok {
						if valData["risk_level"] != nil && valData["detail"] != nil {
							level, count := this.calLevel(valData["risk_level"].(string), valData["detail"].([]interface{}))
							if valData["vul_desc"].(string) != "存在风险" {
								continue
							}
							switch level {
							case 0:
								no += count
							case 1:
								low += count
								total += count
							case 2:
								mid += count
								total += count
							case 3:
								high += count
								total += count
							case 4:
								heavy += count
								total += count
							}
						}
					}
				}
			}
		}
	}

	if apkScanner["apk_sensitive_info"] != nil {
		for _, v := range apkScanner["apk_sensitive_info"].(map[string]interface{}) {
			for k, val := range v.(map[string]interface{}) {
				if k == "sensitive_info" {
					for _, item := range val.([]interface{}) {
						for t, typeDatas := range item.(map[string]interface{}) {
							count := len(typeDatas.([]interface{}))
							switch t {
							case "urls":
								urlCount = count
								break
							case "ip":
								ipCount = count
								break
							case "token":
								tokenCount = count
								break
							case "access_key":
								akCount = count
								break
							case "email":
								emailCount = count
								break
							}
						}
					}
				}
				if k == "certs_list" {
					certCount = len(val.([]interface{}))
				}
			}
		}
	}

	so := qmap.QM{
		"so_pie":    soPie,
		"so_debug":  soDebug,
		"so_canary": soCanary,
	}
	sensitiveCount := qmap.QM{
		"url":        urlCount,
		"ip":         ipCount,
		"token":      tokenCount,
		"access_key": akCount,
		"email":      emailCount,
		"cert":       certCount,
	}
	if heavy > 0 || high > 0 {
		level = 3
	} else if mid > 0 {
		level = 2
	} else if low > 0 {
		level = 1
	} else {
		level = 0
	}

	rtsApk := new(mysql_model.FirmwareReportRtsApkLevel)
	rtsApk.Type = "apk_info"
	rtsApk.ScannerId = firmwareTask.Int("id")
	rtsApk.OriginalContent = vQM.ToString()
	rtsApk.CreateTime = int(time.Now().Unix())
	rtsApk.No = no
	rtsApk.Low = low
	rtsApk.Mid = mid
	rtsApk.High = high
	rtsApk.Heavy = heavy
	rtsApk.Total = total
	rtsApk.So = so.ToString()
	rtsApk.SensitiveCount = sensitiveCount.ToString()
	rtsApk.Level = level
	if _, err := sys_service.NewOrm().InsertOne(rtsApk); err != nil {
		panic(err)
	}
	return total
}

func (this FirmwareScanner) calLevel(level string, detail []interface{}) (int, int) {
	l := 0
	switch level {
	case "低危":
		l = 1
	case "中危":
		l = 2
	case "高危":
		l = 3
	case "严重":
		l = 4
	default:
		l = 0
	}

	return l, len(detail)
}

func (this FirmwareScanner) apkSensitiveScannerAnalysisReport(firmwareTask, apkSensitiveScanner qmap.QM) {
	for k, v := range apkSensitiveScanner {
		vQM := qmap.QM(v.(map[string]interface{}))
		if certsList, has := vQM.TrySlice("certs_list"); has && len(certsList) > 0 {
			for _, item := range certsList {
				rtsApk := new(mysql_model.FirmwareReportRtsApkSensitive)
				rtsApk.Type = "cert"
				rtsApk.PkgName = k
				rtsApk.Content = item.([]interface{})[0].(string)
				rtsApk.ScannerId = firmwareTask.Int("id")
				rtsApk.CreateTime = int(time.Now().Unix())
				if _, err := sys_service.NewOrm().InsertOne(rtsApk); err != nil {
					panic(err)
				}
			}
		}
		if sensitiveInfo, has := vQM.TrySlice("sensitive_info"); has && len(sensitiveInfo) > 0 {
			for _, item := range sensitiveInfo {
				itemQM := qmap.QM(item.(map[string]interface{}))
				for t, typeDatas := range itemQM {
					for _, typeData := range typeDatas.([]interface{}) {
						rtsApk := new(mysql_model.FirmwareReportRtsApkSensitive)
						rtsApk.Type = t
						rtsApk.PkgName = k
						rtsApk.Content = typeData.(string)
						rtsApk.ScannerId = firmwareTask.Int("id")
						rtsApk.CreateTime = int(time.Now().Unix())
						if _, err := sys_service.NewOrm().InsertOne(rtsApk); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}
}

func (this FirmwareScanner) apkCommonVulAnalysisReport(firmwareTask, apkCommonVul qmap.QM) {
	for k, v := range apkCommonVul {
		vQM := qmap.QM(v.(map[string]interface{}))
		for parentType, item := range vQM {
			for t, typeData := range item.(map[string]interface{}) {
				typeDataQM := qmap.QM(typeData.(map[string]interface{}))
				apkVul := new(mysql_model.FirmwareReportRtsApkVul)
				apkVul.ScannerId = firmwareTask.Int("id")
				apkVul.Type = t
				apkVul.PkgName = k
				apkVul.CreateTime = int(time.Now().Unix())
				apkVul.VulDesc = typeDataQM.String("vul_desc")
				apkVul.Reference = typeDataQM.String("reference")
				apkVul.ParentType = parentType
				apkVul.Objective = typeDataQM.String("objective")
				apkVul.RiskLevel = typeDataQM.String("risk_level")
				apkVul.CheckResult = typeDataQM.String("check_result")
				apkVul.Solution = typeDataQM.String("solution")
				apkVul.Name = typeDataQM.String("name")
				apkVul.VulEffect = typeDataQM.String("vul_effect")
				apkVul.Detail = util.SliceToString(typeDataQM.Slice("detail"))
				if _, err := sys_service.NewOrm().InsertOne(apkVul); err != nil {
					panic(err)
				}
			}
		}
	}
}

/**
 *  旧方法通过分值计算Level的方法
 *  1低，2中，3高，4严重，-1默认
 */
func (this FirmwareScanner) culLevel(score float32) int {
	score = score * 10
	level := 1
	if score >= 1 && score <= 39 {
		level = 1
	} else if score >= 40 && score <= 69 {
		level = 2
	} else if score >= 70 && score <= 89 {
		level = 3
	} else if score >= 90 && score <= 100 {
		level = 4
	}
	return level
}

func (this *FirmwareScanner) UpdateScannerTaskStatus(taskInfo qmap.QM, status int, msg, errMsg string) error {
	updateTask := qmap.QM{
		"status": status,
	}
	if err := new(mysql_model.FirmwareTask).Update(taskInfo.MustInt("scanner_id"), updateTask); err == nil {
		new(mysql_model.ScannerTaskLog).Insert(taskInfo, msg, errMsg)
	} else {
		panic(err)
	}
	return nil
}

// 查询固件扫描模板信息
// template_id:100, template_name:通用IoT固件检测模板
// template_id:101, template_name:fs
// template_id:102, template_name:APK扫描
// template_id:103, template_name:固件检测_附带二进制安全检测
// template_id:104, template_name:二进制单文件扫描new
func (this FirmwareScanner) initTemplateMap() {
	url := fmt.Sprintf("%s/api/templates?page=1&num=10000&project_id=0", service.LoadFirmwareConfig().ScanHost)
	result, err := getResultFromUrl(url)
	if err != nil {
		panic(err)
	}
	resultQM, err := qmap.NewWithString(string(result))
	if err != nil {
		panic(err)
	}
	for _, item := range resultQM.Map("data").Slice("list") {
		var itemQM qmap.QM = item.(map[string]interface{})
		println(fmt.Sprintf("template_id:%d, template_name:%s", itemQM.Int("id"), itemQM.String("template_name")))
	}
}

func (this FirmwareScanner) getTemplateName(templateId int) string {
	switch templateId {
	case common.FIRMWARE_TEMPLATE_ID_100:
		return "通用IoT固件检测模板"
	case common.FIRMWARE_TEMPLATE_ID_101:
		return "fs"
	case common.FIRMWARE_TEMPLATE_ID_102:
		return "APK扫描"
	case common.FIRMWARE_TEMPLATE_ID_103:
		return "固件检测_附带二进制安全检测"
	case common.FIRMWARE_TEMPLATE_ID_104:
		return "二进制单文件扫描new"
	}
	return ""
}
