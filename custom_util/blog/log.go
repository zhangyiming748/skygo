package blog

import (
	"go.uber.org/zap"

	"skygo_detection/lib/common_lib/log"
)

func Fatal(msg string, fields ...zap.Field) {
	log.GetBeehiveLogger().Fatal(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.GetBeehiveLogger().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.GetBeehiveLogger().Warn(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.GetBeehiveLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.GetBeehiveLogger().Debug(msg, fields...)
}
