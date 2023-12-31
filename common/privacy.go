package common

// 任务队列
const (
	PRIVACY_TASK_START = "privacy_task_start"
	PRIVACY_TASK_STOP  = "privacy_task_stop"
)

// 应用权限ID对应中文名映射
var PermissionIdZH = map[int]string{
	1:  "读取用户日历数据",
	2:  "写入用户日历数据",
	3:  "访问摄像头进行拍照",
	4:  "读取用户联系人数据",
	5:  "写入用户联系人数据",
	6:  "访问GMail账户列表",
	7:  "访问GPS位置",
	8:  "访问基站位置",
	9:  "录制音频",
	10: "读取设备的IMEI串号",
	11: "电话拨号",
	12: "读取通话记录",
	13: "编辑通话记录",
	14: "添加系统中的语音邮件",
	15: "使用SIP视频服务",
	16: "监视、修改播出的电话记录",
	17: "读取传感器数据",
	18: "读取外部存储器",
	19: "写入外部存储器",
	20: "发送SMS短信",
	21: "读取SMS短信",
	22: "监控接收到WAP PUSH信息",
	23: "监控收到的MMS彩信:将可对其进行编辑",
}

// 应用权限中文名映射
var Permission = map[string]string{
	"READ_CALENDAR":          "读取用户日历数据",
	"WRITE_CALENDAR":         "写入用户日历数据",
	"CAMERA":                 "访问摄像头进行拍照",
	"READ_CONTACTS":          "读取用户联系人数据",
	"WRITE_CONTACTS":         "写入用户联系人数据",
	"GET_ACCOUNTS":           "访问GMail账户列表",
	"ACCESS_FINE_LOCATION":   "访问GPS位置",
	"ACCESS_COARSE_LOCATION": "访问基站位置",
	"RECORD_AUDIO":           "录制音频",
	"READ_PHONE_STATE":       "读取设备的IMEI串号",
	"CALL_PHONE":             "电话拨号",
	"READ_CALL_LOG":          "读取通话记录",
	"WRITE_CALL_LOG":         "编辑通话记录",
	"ADD_VOICEMAIL":          "添加系统中的语音邮件",
	"USE_SIP":                "使用SIP视频服务",
	"PROCESS_OUTGOING_CALLS": "监视、修改播出的电话记录",
	"BODY_SENSORS":           "读取传感器数据",
	"READ_EXTERNAL_STORAGE":  "读取外部存储器",
	"WRITE_EXTERNAL_STORAGE": "写入外部存储器",
	"SEND_SMS":               "发送SMS短信",
	"READ_SMS":               "读取SMS短信",
	"RECEIVE_WAP_PUSH":       "监控接收到WAP PUSH信息",
	"RECEIVE_MMS":            "监控收到的MMS彩信:将可对其进行编辑",
}

// 授权状态
var Permission_state = map[string]string{
	"Access": "通过",
	"Reject": "拒绝",
}

var Permission_default = map[string]string{
	"allow":                                "允许",
	"default":                              "默认",
	"allow / switch COARSE_LOCATION=allow": "允许",
}

// 应用中文名对应
var AppsName = map[string]string{
	"android":                                                "Android 系统",
	"android.auto_generated_rro__":                           "android.auto_generated_rro__",
	"android.ext.services":                                   "Android Services Library",
	"android.ext.shared":                                     "Android Shared Library",
	"android.telephony.overlay.cmcc":                         "android.telephony.overlay.cmcc",
	"com.android.apps.tag":                                   "Tags",
	"com.android.backupconfirm":                              "com.android.backupconfirm",
	"com.android.bips":                                       "默认打印服务",
	"com.android.bluetooth":                                  "蓝牙",
	"com.android.bluetoothmidiservice":                       "Bluetooth MIDI Service",
	"com.android.bookmarkprovider":                           "Bookmark Provider",
	"com.android.calculator2":                                "计算器",
	"com.android.calendar":                                   "日历",
	"com.android.calllogbackup":                              "Call Log Backup/Restore",
	"com.android.camera2":                                    "相机",
	"com.android.captiveportallogin":                         "CaptivePortalLogin",
	"com.android.carrierconfig":                              "com.android.carrierconfig",
	"com.android.carrierdefaultapp":                          "运营商默认应用",
	"com.android.cellbroadcastreceiver":                      "小区广播",
	"com.android.certinstaller":                              "证书安装程序",
	"com.android.companiondevicemanager":                     "配套设备管理器",
	"com.android.contacts":                                   "通讯录",
	"com.android.cts.ctsshim":                                "com.android.cts.ctsshim",
	"com.android.cts.priv.ctsshim":                           "com.android.cts.priv.ctsshim",
	"com.android.defcontainer":                               "软件包权限帮助程序",
	"com.android.deskclock":                                  "时钟",
	"com.android.dialer":                                     "电话",
	"com.android.documentsui":                                "文件",
	"com.android.dreams.basic":                               "基本互动屏保",
	"com.android.dreams.phototable":                          "照片屏幕保护程序",
	"com.android.egg":                                        "PAINT.APK",
	"com.android.email":                                      "电子邮件",
	"com.android.emergency":                                  "急救信息",
	"com.android.externalstorage":                            "外部存储设备",
	"com.android.gallery3d":                                  "图库",
	"com.android.htmlviewer":                                 "HTML 查看程序",
	"com.android.inputdevices":                               "输入设备",
	"com.android.inputmethod.latin":                          "Android 键盘 (AOSP)",
	"com.android.internal.display.cutout.emulation.corner":   "边角刘海屏",
	"com.android.internal.display.cutout.emulation.double":   "双刘海屏",
	"com.android.internal.display.cutout.emulation.noCutout": "隐藏",
	"com.android.internal.display.cutout.emulation.tall":     "长型刘海屏",
	"com.android.keychain":                                   "密钥链",
	"com.android.launcher3":                                  "Quickstep",
	"com.android.location.fused":                             "一体化位置信息",
	"com.android.managedprovisioning":                        "工作资料设置",
	"com.android.messaging":                                  "短信",
	"com.android.mms.service":                                "MmsService",
	"com.android.mtp":                                        "MTP 主机",
	"com.android.music":                                      "音乐",
	"com.android.musicfx":                                    "MusicFX",
	"com.android.nfc":                                        "NFC服务",
	"com.android.onetimeinitializer":                         "One Time Init",
	"com.android.packageinstaller":                           "软件包安装程序",
	"com.android.pacprocessor":                               "PacProcessor",
	"com.android.phone":                                      "电话服务",
	"com.android.printservice.recommendation":                "Print Service Recommendation Service",
	"com.android.printspooler":                               "打印处理服务",
	"com.android.providers.blockednumber":                    "存储已屏蔽的号码",
	"com.android.providers.calendar":                         "日历存储",
	"com.android.providers.contacts":                         "联系人存储",
	"com.android.providers.downloads":                        "内容下载管理器",
	"com.android.providers.downloads.ui":                     "下载",
	"com.android.providers.media":                            "媒体存储设备",
	"com.android.providers.settings":                         "设置存储",
	"com.android.providers.telephony":                        "电话和短信存储",
	"com.android.providers.userdictionary":                   "用户字典",
	"com.android.provision":                                  "com.android.provision",
	"com.android.proxyhandler":                               "ProxyHandler",
	"com.android.quicksearchbox":                             "Search",
	"com.android.se":                                         "SecureElementApplication",
	"com.android.server.telecom":                             "通话管理",
	"com.android.settings":                                   "设置",
	"com.android.settings.intelligence":                      "Settings Suggestions",
	"com.android.sharedstoragebackup":                        "com.android.sharedstoragebackup",
	"com.android.shell":                                      "Shell",
	"com.android.simappdialog":                               "Sim App Dialog",
	"com.android.smspush":                                    "com.android.smspush",
	"com.android.statementservice":                           "Intent Filter Verification Service",
	"com.android.storagemanager":                             "存储空间管理器",
	"com.android.systemui":                                   "系统界面",
	"com.android.systemui.theme.dark":                        "深色",
	"com.android.traceur":                                    "系统跟踪",
	"com.android.vpndialogs":                                 "VpnDialogs",
	"com.android.wallpaper.livepicker":                       "Live Wallpaper Picker",
	"com.android.wallpaperbackup":                            "com.android.wallpaperbackup",
	"com.android.wallpapercropper":                           "com.android.wallpapercropper",
	"com.android.wallpaperpicker":                            "com.android.wallpaperpicker",
	"com.android.webview":                                    "Android System WebView",
	"com.google.android.theme.pixel":                         "Pixel",
	"com.huawei.works.videomeetinglite":                      "WeMeeting",
	"com.qihoo.policysdk":                                    "IDPSDemo",
	"com.qihoo360.applicationmanagerdemo":                    "ApplicationManagerDemo",
	"com.qihoo360.caridps.filemonitordemo":                   "FileMonitorDemo",
	"com.qihoo360.caridps.policymanager":                     "idps-android-framework-policymanager",
	"com.qihoo360.carkeeper":                                 "汽车管家",
	"com.qihoo360.networkclient":                             "NetworkClient",
	"com.qihoo360.vcaridps":                                  "系统检测与清理",
	"com.qihoo360.vprivacy":                                  "VPrivacyScanner",
	"com.qihoo360.vulscan":                                   "SystemVulScan",
	"com.qualcomm.timeservice":                               "com.qualcomm.timeservice",
	"com.taobao.taobao":                                      "淘宝",
	"org.chromium.webview_shell":                             "WebView Shell",
}
