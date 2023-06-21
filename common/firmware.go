package common

// 固件扫描状态(firmware Status)
// 状态 	1 待上传 2 上传完成 3 上传失败 4 取消上传 5 (下载完成) 扫描中 6 (创建任务) 扫描中 7 取消扫描 8 扫描完成 9 扫描失败 10 已解析 0 已删除
const (
	FW_FINISHED = 10 //固件扫描状态:扫描完成，解析成功
	FW_ABNORMAL = 9  //固件扫描状态:扫描异常，解析失败
	FW_SCANNING = 8  //固件扫描状态:扫描异常，解析失败
)
const (
	FIRM_WARE_API = "http://10.220.171.247:8080"
	//FIRM_WARE_DOWNLOAD_MOMAIN = "http://beta.vadmin.car.qihoo.net"
	//FIRM_WARE_API  = "http://10.19.42.125:8080"
	FIRM_WARE_DOWNLOAD_MOMAIN = "http://qa.vadmin.car.qihoo.net:7100"
)
