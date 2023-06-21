package mysql

import (
	"skygo_detection/service"

	"xorm.io/xorm"
)

var mysqlEngine *xorm.Engine

func InitMysqlEngine() {
	dbConfig := service.LoadConfig().DB
	config := NewConfig(dbConfig.HostName, dbConfig.Port, dbConfig.UserName, dbConfig.Password, dbConfig.DBName, dbConfig.Charset)

	if dbConfig.ShowSql == true {
		config.ShowSql = true
	}

	config.SetLevel(int(dbConfig.Log.Level))

	if dbConfig.Log.FilePath != "" {
		config.SetLogPath(dbConfig.Log.FilePath)
	}

	var err error
	mysqlEngine, err = NewMysqlEngine(config)
	if err != nil {
		panic(err)
	}
	err2 := mysqlEngine.Ping()
	if err2 != nil {
		panic(err2)
	}
}

func GetEngine() *xorm.Engine {
	return mysqlEngine
}

func GetSession() *xorm.Session {
	return mysqlEngine.NewSession()
}
