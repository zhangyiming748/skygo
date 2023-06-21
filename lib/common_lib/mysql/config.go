package mysql

import (
	"time"
)

func NewConfig(host string, port int, user, password, dbName, charset string) *Config {
	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Dbname:   dbName,
		Charset:  charset,
	}
}

type Config struct {
	Host        string
	Port        int
	User        string
	Password    string
	Dbname      string
	MaxOpenCons int    // 连接池中最大连接数
	MaxIdleCons int    // 连接池中最大空闲连接数
	MaxLifeTime int    // 单个连接最大存活时间(单位:秒)
	ShowSql     bool   // 是否在控制台打印出生成的SQL语句
	LogPath     string // 日志输出路径
	Level       int
	Charset     string
}

func (x *Config) SetLogPath(logPath string) *Config {
	x.LogPath = logPath
	return x
}

func (x *Config) SetLevel(level int) *Config {
	x.Level = level
	return x
}

func (x *Config) SetMaxOpenCons(cons int) *Config {
	x.MaxOpenCons = cons
	return x
}

func (x *Config) SetMaxIdleCons(cons int) *Config {
	x.MaxIdleCons = cons
	return x
}

func (x *Config) SetMaxLifeTime(cons int) *Config {
	x.MaxIdleCons = cons
	return x
}

const (
	DefaultMaxOpenConns = 500              // 连接池中最大连接数 默认值
	DefaultMaxIdleConns = 5                // 连接池中最大空闲连接数 默认值
	DefaultMaxLifeTime  = 30 * time.Second // 单个连接最大存活时间(单位:秒) 默认值
)
