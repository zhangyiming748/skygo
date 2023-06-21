package common

const (
	START_FAILED       = iota //未知启动失败,需要查阅日志
	START_SUCCESS             //启动成功
	ISRUNNING                 //设备正在运行中
	DEVICE_NOT_CONNECT        //未连接USRP 请接入设备后再尝试
)
const (
	STOP_FAILED  = iota //关闭设备失败,需要查阅日志
	STOP_SUCCESS        // 关闭设备成功
	NOT_RUNNING         //程序未在运行在,无需关闭
)

const (
	LteUrlCrackApn       = "/crackapn"
	LteUrlGetCrackResult = "/getcrackresult"
	LteUrlPasswordUpload = "/passwordupload"
)

const (
	START_2G_SYSTEM = "start"    // 启动2G系统
	STOP_2G_SYSTEM  = "stop"     // 停止2G系统
	CONFIG_SYSTEM   = "config"   // 对整套系统进行基础配置
	BEFORE_START    = "iptables" // 配置系统网络数据转发
	GET_SMS         = "smsinfo"  // 获取当前系统中短信相关信息
	GET_DEVICES     = "ueinfo"   // 获取当前系统中所有终端设备信息
	SEND_SMS        = "sendsms"  // 向指定imsi设备发送短信
)
const IFACE = "wlo1" // 设置配置前或开启设备前需要配置防火墙

const (
	// LteSystem
	SET_IMSL         = "ipaddress:8081/writesim"  // 写卡
	START_lTE_SYSTEM = "ipaddress:8081/start"     // 启动LTE设备
	BASICINFO_INFO   = "ipaddress:8081/basicinfo" //获取设备信息
	STOP_lTE_SYSTEM  = "ipaddress:8081/stop"      // 关闭系统
	GET_PACKAGE      = "ipaddress:8081/getfile"   // 抓包

)
