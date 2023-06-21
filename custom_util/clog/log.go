package clog

import (
	"go.uber.org/zap"

	"skygo_detection/lib/common_lib/log"
)

func Fatal(msg string, fields ...zap.Field) {
	log.GetEtlLogger().Fatal(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.GetEtlLogger().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.GetEtlLogger().Warn(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.GetEtlLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.GetEtlLogger().Debug(msg, fields...)
}
