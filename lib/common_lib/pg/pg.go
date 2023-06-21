package pg

import (
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type DbName = string

func NewPostgresEngine(conf *Config) (*xorm.Engine, error) {
	return generatePostgresEngine(conf)
}

var PostgresEngines = map[DbName]*xorm.Engine{}

func generatePostgresEngine(conf *Config) (*xorm.Engine, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname)

	engine, err := xorm.NewEngine("postgres", dsn)
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
