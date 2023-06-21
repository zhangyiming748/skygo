package bootstrap

import (
	"skygo_detection/lib/common_lib/mysql"
)

func InitService() {
	// mysql数据库初始
	mysql.InitMysqlEngine()

	// 验证器
	// InitValidator()
}
