package sys_service

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var (
	httpLogger    *zap.Logger
	rpcLogger     *zap.Logger
	consoleLogger *zap.Logger
	output        = map[string]*os.File{}
	loggerMu      sync.Mutex
)

func GetHttpLoggerInstance() *zap.Logger {
	if httpLogger == nil {
		loggerMu.Lock()
		logConfig := GetHttpLogConfig()
		output := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logConfig.FilePath,
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
		httpLogger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
		defer loggerMu.Unlock()
	}
	return httpLogger
}

func GetRpcLoggerInstance() *zap.Logger {
	if rpcLogger == nil {
		logConfig := GetRpcLogConfig()
		output := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logConfig.FilePath,
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
		rpcLogger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	return rpcLogger
}

func GetConsoleLoggerInstance() *zap.Logger {
	if consoleLogger == nil {
		logConfig := GetConsoleLogConfig()
		output := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logConfig.FilePath,
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
		consoleLogger = zap.New(core, zap.AddCaller())
	}
	return consoleLogger
}

func GetLogFile(fileName string) *os.File {
	if _, ok := output[fileName]; !ok {
		outputFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		output[fileName] = outputFile
	}
	return output[fileName]
}

func ErrorToLevel(err error) zapcore.Level {
	if err != nil {
		return zap.ErrorLevel
	} else {
		return zap.InfoLevel
	}
}
