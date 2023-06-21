package service

import (
	"skygo_detection/guardian/src/config/watcher"
	"sync"
	"time"

	"xorm.io/core"
)

var (
	CONFIG_RELOAD_TAG = true
	onceWatch         = sync.Once{}
	config            *Config
)

type Config struct {
	Http           HttpConfig
	Https          HttpsConfig
	Log            LogConfig
	DB             DBConfig
	Redis          RedisConfig
	ES             ElasticSearchConfig
	Kafka          KafkaConfig
	Zookeeper      ZookeeperConfig
	JWT            JWTConfig
	Referer        RefererConfig
	EtlEngine      EtlEngineConfig
	SocPlat        SocPlatConfig
	MongoDB        MongoDBConfig
	S3             S3Config
	Firmware       FirmwareConfig
	ReportTemplate ReportTemplateConfig
	VehicleScreen  VehicleScreenConfig
	Download       DownloadConfig
	License        LicenseConfig
	Gps            GpsConfig
	Beehive        BeehiveConfig
	Ump            UmpConfig
	Hydra          HydraConfig
}

type ReportTemplateConfig struct {
	ReportTemplate string
	ReportHearder  string
	ReportFooter   string
	FontPath       string
	OutputImage    string
}

type FirmwareConfig struct {
	ScanHost  string
	AdminHost string
}

type S3Config struct {
	Bucket    string
	AccessKey string
	SecretKey string
	EndPoint  string
}

type HttpConfig struct {
	Host string
	Port int
}

type RedisConfig struct {
	Addr     []string
	Auth     string
	PoolSize int
	Timeout  time.Duration
}

type KafkaConfig struct {
	Brokers             []string
	ConsumerGroupPrefix string
	Version             string
}

type ZookeeperConfig struct {
	Brokers []string
}

type ElasticSearchConfig struct {
	User     string
	Password string
	Host     []string
}

type EtlEngineConfig struct {
	Addr                   string
	HeaderTokenAuthKeyName string
	HeaderTokenEncryptKey  string
}

type SocPlatConfig struct {
	UserAuthUrl string
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

type LogConfig struct {
	FilePath          string // 日志文件输出目录
	Level             string // 日志输出等级
	MaxSize           int    // 日志文件最大值
	MaxAge            int
	MaxBackups        int
	OutputProbability float32 // 日志输出概率(小数)
	ToStdout          bool    // 是否将日志输出到标准输出中
}

type JWTConfig struct {
	Algorithm  string
	SecretKey  string
	ExpireTime int
}

type RefererConfig struct {
	Url    []string
	Enable bool
}

type HttpsConfig struct {
	Enable bool
	Port   int
	Crt    string
	Key    string
}

// mongo连接配置
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

type VehicleScreenConfig struct {
	IsSqlite bool
	FilePath string
}

type DownloadConfig struct {
	DownloadPath string
}

type LicenseConfig struct {
	Path string
}

type GpsConfig struct {
	Url     string
	StartSh string
	StopSh  string
	Log     string
}

type BeehiveConfig struct {
	Host             string
	GsmSnifferPort   string
	LteSystemPort    string
	GsmSystemPort    string
	LetSystemNetwork string
}

type HydraConfig struct {
	Client string
	Server string
}

type UmpConfig struct {
	ClientID     string
	ClientSecret string
	ServerSso    string
	LoginUrl     string
}

// 全局获取
func LoadConfig() *Config {
	return loadConfig()
}

func LoadHttpConfig() *HttpConfig {
	return &loadConfig().Http
}

func LoadHttpsConfig() *HttpsConfig {
	return &loadConfig().Https
}

func LoadReportTemplateConfig() *ReportTemplateConfig {
	return &loadConfig().ReportTemplate
}

func LoadRedisConfig() *RedisConfig {
	return &loadConfig().Redis
}

func LoadKafkaConfig() *KafkaConfig {
	return &loadConfig().Kafka
}

func LoadMongoDBConfig() *MongoDBConfig {
	return &loadConfig().MongoDB
}

func LoadS3Config() *S3Config {
	return &loadConfig().S3
}

func LoadLogConfig() *LogConfig {
	return &loadConfig().Log
}

func LoadJWTConfig() *JWTConfig {
	return &loadConfig().JWT
}

func LoadFirmwareConfig() *FirmwareConfig {
	return &loadConfig().Firmware
}

func LoadLicenseConfig() *LicenseConfig {
	return &loadConfig().License
}

func LoadGpsConfig() *GpsConfig {
	return &loadConfig().Gps
}

func LoadBeehiveConfig() *BeehiveConfig {
	return &loadConfig().Beehive
}

func LoadUmpConfig() *UmpConfig {
	return &loadConfig().Ump
}

func LoadHydraConfig() *HydraConfig {
	return &loadConfig().Hydra
}

func loadConfig() *Config {
	if CONFIG_RELOAD_TAG {
		config = new(Config)
		if err := watcher.Get("config.tml").UnmarshalTOML(config); err != nil {
			if err != watcher.ErrNotExist {
				panic(err)
			}
		}
		CONFIG_RELOAD_TAG = false
		onceWatch.Do(func() { watcher.Watch("config.tml", new(configWatcher)) })
	}
	return config
}

type configWatcher struct{}

func (s *configWatcher) Set(val string) (err error) {
	CONFIG_RELOAD_TAG = true
	return
}
