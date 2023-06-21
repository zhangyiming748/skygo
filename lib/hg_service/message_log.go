package hg_service

import (
	"skygo_detection/service"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func AnalysisResultLog(messageType, taskId, rawMsg, errMsg string) {
	logger := service.GetDefaultLogger("analysis_message_result")
	defer logger.Sync()
	errLevel := zapcore.InfoLevel
	if errMsg != "" {
		errLevel = zapcore.ErrorLevel
	}
	if logger.Core().Enabled(errLevel) {
		logger.Check(errLevel, errMsg).Write(
			zap.String("task_id", taskId),
			zap.String("message_type", messageType),
			zap.String("raw_message", rawMsg),
			zap.String("create_time", time.Now().Format(time.RFC3339)),
		)
	}
}

func AnalysisRawLog(messageType, taskId, rawMsg, errMsg string) {
	logger := service.GetDefaultLogger("analysis_message_raw")
	defer logger.Sync()
	errLevel := zapcore.InfoLevel
	if errMsg != "" {
		errLevel = zapcore.ErrorLevel
	}
	if logger.Core().Enabled(errLevel) {
		logger.Check(errLevel, errMsg).Write(
			zap.String("task_id", taskId),
			zap.String("message_type", messageType),
			zap.String("raw_message", rawMsg),
			zap.String("create_time", time.Now().Format(time.RFC3339)),
		)
	}
}

func TerminalReceiveLog(messageType, taskId, rawMsg, errMsg string) {
	logger := service.GetDefaultLogger("terminal_message_receive")
	defer logger.Sync()
	errLevel := zapcore.InfoLevel
	if errMsg != "" {
		errLevel = zapcore.ErrorLevel
	}
	if logger.Core().Enabled(errLevel) {
		logger.Check(errLevel, errMsg).Write(
			zap.String("task_id", taskId),
			zap.String("message_type", messageType),
			zap.String("raw_message", rawMsg),
			zap.String("create_time", time.Now().Format(time.RFC3339)),
		)
	}
}

func TerminalSendLog(messageType, taskId, rawMsg, errMsg string) {
	logger := service.GetDefaultLogger("terminal_message_send")
	defer logger.Sync()
	errLevel := zapcore.InfoLevel
	if errMsg != "" {
		errLevel = zapcore.ErrorLevel
	}
	if logger.Core().Enabled(errLevel) {
		logger.Check(errLevel, errMsg).Write(
			zap.String("task_id", taskId),
			zap.String("message_type", messageType),
			zap.String("raw_message", rawMsg),
			zap.String("create_time", time.Now().Format(time.RFC3339)),
		)
	}
}
