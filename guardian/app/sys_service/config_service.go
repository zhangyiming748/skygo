package sys_service

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"xorm.io/core"

	"skygo_detection/guardian/src/config/watcher"
)

var (
	// 设置系统环境变量 dev docker online
	ENV              = "dev"
	CONFIG_PATH      = "./config/dev/"
	CONFIG_FILE_NAME = "sys.tml"
	G_RELOAD_TAG     = true
	ROOT_DIR         = "."
	onceWatch        = sync.Once{}
	guardianConfig   *GuardianConfig
)

// 设置系统环境变量 dev qa ecarx_qa online
func InitConfigWatcher(env, configPath string) {
	CONFIG_FILE_NAME = "sys.tml"
	if configPath == "" {
		// 非容器部署环境变量, 采用正则匹配
		reg, _ := regexp.Compile(`^(dev|op_qa)(_[a-zA-Z0-9]+)?$`)
		switch {
		case reg.MatchString(env) == true:
			ENV = env
			CONFIG_PATH = "./config/" + ENV + "/"
		case env == "qa" || env == "online" || env == "ecarx_qa":
			ENV = env
			CONFIG_PATH = "/src/config/" + ENV + "/"
		default:
			panic("unknown environment param!")
		}
	} else {
		path, fileName := path.Split(configPath)
		CONFIG_PATH = path
		if fileName != "" {
			CONFIG_FILE_NAME = fileName
		}

	}
	// 获取程序执行的根目录
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		// 如果程序不是通过 go run * 来执行的，则当前目录就是程序执行的根目录
		if !strings.HasPrefix(dir, "/tmp/go-build") {
			ROOT_DIR = dir
		}
	}

	if err := watcher.InitWatch(CONFIG_PATH); err != nil {
		panic(err)
	}
}

type GuardianConfig struct {
	Service       ServiceConfig
	Rpc           RpcConfig
	RpcLog        LogConfig
	Http          HttpConfig
	Https         HttpsConfig
	HttpLog       LogConfig
	DB            DBConfig
	MongoDB       MongoDBConfig
	ConsoleLog    LogConfig
	RpcServiceMap []*RpcService
}

type ServiceConfig struct {
	Name string
}

type RpcConfig struct {
	Port          int
	Tag           string
	CheckTimeout  int
	CheckInterval int
}

type HttpConfig struct {
	Host string
	Port int
}

type HttpsConfig struct {
	Enable bool
	Port   int
	Crt    string
	Key    string
}

type DBConfig struct {
	DBName            string
	HostName          string
	Port              int
	UserName          string
	Password          string
	Charset           string
	Log               DBLogConfig
	MaxLifeTime       time.Duration
	MaxIdleConnection int
	MaxOpenConnection int
	ShowSql           bool
}

type DBLogConfig struct {
	FilePath string
	Level    core.LogLevel
}

type MongoDBConfig struct {
	Host       string
	Port       int32
	Username   string
	Password   string
	AuthSource string
	ReplicaSet string
	ExtraUrl   string // 网络连接额外参数（亿咖通会用到"/test"）
	DBName     string
	QcmUrl     string
	QcmKey     string
}

type LogConfig struct {
	FilePath          string
	Level             string
	MaxSize           int
	MaxAge            int
	MaxBackups        int
	OutputProbability float32 // 日志输出概率(小数)
	ToStdout          bool    // 是否将日志输出到标准输出中
}

type RpcService struct {
	Name string
	Host string
	Port int
}

func GetServiceConfig() *ServiceConfig {
	return &loadConfig().Service
}

func GetRpcConfig() *RpcConfig {
	return &loadConfig().Rpc
}

func GetRpcLogConfig() *LogConfig {
	return &loadConfig().RpcLog
}

func GetHttpConfig() *HttpConfig {
	return &loadConfig().Http
}

func GetHttpsConfig() *HttpsConfig {
	return &loadConfig().Https
}

func GetHttpLogConfig() *LogConfig {
	return &loadConfig().HttpLog
}

func GetDefaultDBConfig() *DBConfig {
	return &loadConfig().DB
}

func LoadMongoDBConfig() *MongoDBConfig {
	mgoConfig := &loadConfig().MongoDB
	if mgoConfig.Host == "" && mgoConfig.QcmUrl != "" && mgoConfig.QcmKey != "" {
		if qcmConn, err := GetQCMConfig(mgoConfig.QcmUrl, mgoConfig.QcmKey); err == nil {
			if len(qcmConn) > 0 {
				mgoConfig.Host = qcmConn[0].Ip
				mgoConfig.Port = int32(qcmConn[0].Port)
			}
		} else {
			panic(err)
		}
	}
	return mgoConfig
}

func GetConsoleLogConfig() *LogConfig {
	return &loadConfig().ConsoleLog
}

func GetRpcServiceMapConfig() *[]*RpcService {
	return &loadConfig().RpcServiceMap
}

func GetServiceName() string {
	return GetServiceConfig().Name
}

var loadConfigMutex sync.Mutex // 锁机制，确保调用loadConfig()方法是协程安全的

func loadConfig() *GuardianConfig {
	loadConfigMutex.Lock()
	if G_RELOAD_TAG {
		guardianConfig = new(GuardianConfig)
		if err := watcher.Get(CONFIG_FILE_NAME).UnmarshalTOML(guardianConfig); err != nil {
			if err != watcher.ErrNotExist {
				panic(err)
			}
		}
		G_RELOAD_TAG = false
		onceWatch.Do(func() { watcher.Watch(CONFIG_FILE_NAME, new(sysWatcher)) })
	}

	loadConfigMutex.Unlock()

	return guardianConfig
}

type sysWatcher struct{}

func (s *sysWatcher) Set(val string) (err error) {
	G_RELOAD_TAG = true
	return
}
