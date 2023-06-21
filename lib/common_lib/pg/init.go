package pg

import (
	"skygo_detection/service"

	"xorm.io/xorm"
)

var pgEngine *xorm.Engine

func InitPgEngine() {
	dbConfig := service.LoadConfig().DB
	config := NewConfig(dbConfig.HostName, dbConfig.Port, dbConfig.UserName, dbConfig.Password, dbConfig.DBName)

	if dbConfig.ShowSql == true {
		config.ShowSql = true
	}

	config.SetLevel(int(dbConfig.Log.Level))

	if dbConfig.Log.FilePath != "" {
		config.SetLogPath(dbConfig.Log.FilePath)
	}

	var err error
	pgEngine, err = NewPostgresEngine(config)
	if err != nil {
		panic(err)
	}
}

func GetEngine() *xorm.Engine {
	return pgEngine
}

func GetSession() *xorm.Session {
	return pgEngine.NewSession()
}
