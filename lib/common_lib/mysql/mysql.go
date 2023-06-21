package mysql

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type DbName = string

func NewMysqlEngine(conf *Config) (*xorm.Engine, error) {
	return generateMysqlEngine(conf)
}

var MysqlEngines = map[DbName]*xorm.Engine{}

func GetDBConnectionStr(dbConfig *Config) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s", dbConfig.User,
		dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Dbname, dbConfig.Charset)
}

func generateMysqlEngine(conf *Config) (*xorm.Engine, error) {
	dsn := GetDBConnectionStr(conf)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if conf.MaxOpenCons == 0 {
		engine.SetMaxOpenConns(DefaultMaxOpenConns)
	} else {
		engine.SetMaxOpenConns(conf.MaxOpenCons)
	}

	if conf.MaxIdleCons == 0 {
		engine.SetMaxIdleConns(DefaultMaxIdleConns)
	} else {
		engine.SetMaxIdleConns(conf.MaxIdleCons)
	}

	if conf.MaxLifeTime == 0 {
		engine.SetConnMaxLifetime(DefaultMaxLifeTime)
	} else {
		engine.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifeTime))
	}

	// 是否打印sql.default false
	engine.ShowSQL(conf.ShowSql)

	// 日志
	if conf.LogPath != "" {
		f, err := os.Create(conf.LogPath)
		if err != nil {
			println(err.Error())
		} else {
			logger := log.NewSimpleLogger(f)
			logger.ShowSQL(conf.ShowSql)              // 是否打印sql.default false
			logger.SetLevel(log.LogLevel(conf.Level)) // 日志等级
			engine.SetLogger(logger)
		}
	}

	return engine, nil
}
