package gpslog

import (
	"go.uber.org/zap"

	"skygo_detection/lib/common_lib/log"
)

func Fatal(msg string, fields ...zap.Field) {
	log.GetGpsLogger().Fatal(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.GetGpsLogger().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.GetGpsLogger().Warn(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.GetGpsLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.GetGpsLogger().Debug(msg, fields...)
}
