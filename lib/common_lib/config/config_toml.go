package config

//
// import (
// 	"path/filepath"
// 	"time"
//
// 	"github.com/BurntSushi/toml"
// 	"xorm.io/core"
// )
//
// // 初始化
// func InitConfig(path string) {
// 	if path == "" {
// 		panic("config path must not be empty")
// 	}
//
// 	path, _ = filepath.Abs(path)
// 	if c, err := setConfig(path); err != nil {
// 		panic(err)
// 	} else {
// 		config = c
// 	}
// }
//
// var config *Config
//
// func setConfig(path string) (*Config, error) {
// 	config := new(Config)
// 	_, err := toml.DecodeFile(path, config)
// 	return config, err
// }
//
// // --------------基础配置---------------------------
// type Config struct {
// 	Http           HttpConfig
// 	Https          HttpsConfig
// 	Log            LogConfig
// 	DB             DBConfig
// 	Redis          RedisConfig
// 	ES             ElasticSearchConfig
// 	Kafka          KafkaConfig
// 	Zookeeper      ZookeeperConfig
// 	JWT            JWTConfig
// 	Referer        RefererConfig
// 	EtlEngine      EtlEngineConfig
// 	SocPlat        SocPlatConfig
// 	MongoDB        MongoDBConfig
// 	S3             S3Config
// 	Firmware       FirmwareConfig
// 	ReportTemplate ReportTemplateConfig
// 	VehicleScreen  VehicleScreenConfig
// }
//
// type ReportTemplateConfig struct {
// 	ReportTemplate string
// 	ReportHearder  string
// 	ReportFooter   string
// 	FontPath       string
// 	OutputImage    string
// }
//
// type FirmwareConfig struct {
// 	ScanHost  string
// 	AdminHost string
// }
//
// type S3Config struct {
// 	Bucket    string
// 	AccessKey string
// 	SecretKey string
// 	EndPoint  string
// }
//
// type HttpConfig struct {
// 	Host string
// 	Port int
// }
//
// type RedisConfig struct {
// 	Addr     []string
// 	Auth     string
// 	PoolSize int
// 	Timeout  time.Duration
// }
//
// type KafkaConfig struct {
// 	Brokers             []string
// 	ConsumerGroupPrefix string
// 	Version             string
// }
//
// type ZookeeperConfig struct {
// 	Brokers []string
// }
//
// type ElasticSearchConfig struct {
// 	User     string
// 	Password string
// 	Host     []string
// }
//
// type EtlEngineConfig struct {
// 	Addr                   string
// 	HeaderTokenAuthKeyName string
// 	HeaderTokenEncryptKey  string
// }
//
// type SocPlatConfig struct {
// 	UserAuthUrl string
// }
//
// type DBConfig struct {
// 	DBName            string
// 	HostName          string
// 	Port              int
// 	UserName          string
// 	Password          string
// 	Charset           string
// 	Log               DBLogConfig
// 	MaxLifeTime       time.Duration
// 	MaxIdleConnection int
// 	MaxOpenConnection int
// 	ShowSql           bool
// }
//
// type DBLogConfig struct {
// 	FilePath string
// 	Level    core.LogLevel
// }
//
// type LogConfig struct {
// 	FilePath          string // 日志文件输出目录
// 	Level             string // 日志输出等级
// 	MaxSize           int    // 日志文件最大值
// 	MaxAge            int
// 	MaxBackups        int
// 	OutputProbability float32 // 日志输出概率(小数)
// 	ToStdout          bool    // 是否将日志输出到标准输出中
// }
//
// type JWTConfig struct {
// 	Algorithm  string
// 	SecretKey  string
// 	ExpireTime int
// }
//
// type RefererConfig struct {
// 	Url    []string
// 	Enable bool
// }
//
// type HttpsConfig struct {
// 	Enable bool
// 	Port   int
// 	Crt    string
// 	Key    string
// }
//
// type MongoDBConfig struct {
// 	Host       string
// 	Port       int32
// 	Username   string
// 	Password   string
// 	AuthSource string
// 	ReplicaSet string
// 	ExtraUrl   string // 网络连接额外参数（亿咖通会用到"/test"）
// 	DBName     string
// 	QcmUrl     string
// 	QcmKey     string
// }
//
// type VehicleScreenConfig struct {
// 	IsSqlite bool
// 	FilePath int32
// }
//
// // 全局获取
// func GetConfig() *Config {
// 	return config
// }
//
// func GetHttpConfig() *HttpConfig {
// 	return &config.Http
// }
//
// func GetHttpsConfig() *HttpsConfig {
// 	return &config.Https
// }
//
// func GetReportTemplateConfig() *ReportTemplateConfig {
// 	return &config.ReportTemplate
// }
//
// func GetLogConfig() *LogConfig {
// 	return &config.Log
// }
//
// func GetMongoDBConfig() *MongoDBConfig {
// 	mgoConfig := &config.MongoDB
// 	if mgoConfig.Host == "" && mgoConfig.QcmUrl != "" && mgoConfig.QcmKey != "" {
// 		if qcmConn, err := GetQCMConfig(mgoConfig.QcmUrl, mgoConfig.QcmKey); err == nil {
// 			if len(qcmConn) > 0 {
// 				mgoConfig.Host = qcmConn[0].Ip
// 				mgoConfig.Port = int32(qcmConn[0].Port)
// 			}
// 		} else {
// 			panic(err)
// 		}
// 	}
// 	return mgoConfig
// }
