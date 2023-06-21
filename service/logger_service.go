package service

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggerInstance = map[string]*zap.Logger{}
var loggerMu = &sync.Mutex{}

func GetDefaultLogger(loggerName string) *zap.Logger {
	if logger, has := loggerInstance[loggerName]; has {
		return logger
	} else {
		loggerMu.Lock()
		newLogger := newLoggerInstance(loggerName)
		loggerInstance[loggerName] = newLogger
		loggerMu.Unlock()
		return newLogger
	}
}

func newLoggerInstance(loggerName string) *zap.Logger {
	logConfig := LoadLogConfig()
	logPath := logConfig.FilePath
	if strings.HasSuffix(logPath, "/") {
		logPath = logPath[0 : len(logPath)-1]
	}
	logFile := fmt.Sprintf("%s/%s.log", logPath, loggerName)
	output := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
	})
	if logConfig.ToStdout {
		output = zapcore.NewMultiWriteSyncer(output, os.Stdout)
	}
	level := zapcore.DebugLevel
	level.Set(logConfig.Level)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		output,
		zap.NewAtomicLevelAt(level),
	)
	return zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
}
