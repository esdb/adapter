package plz_zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/v2pro/plz/logging"
)

func Adapt(zapLogger *zap.Logger) logging.Logger {
	return &loggerAdapter{zapLogger}
}

type loggerAdapter struct {
	zapLogger *zap.Logger
}

func (adapter *loggerAdapter) Log(level logging.Level, msg string, kv ...interface{}) {
	if ce := adapter.zapLogger.Check(translateLevel(level), msg); ce != nil {
		zapFields := []zapcore.Field{}
		for i := 0; i < len(kv); i += 2 {
			zapFields = append(zapFields, zap.Reflect(kv[i].(string), kv[i+1]))
		}
		ce.Write(zapFields...)
	}
}

func (adapter *loggerAdapter) Error(msg string, kv ...interface{}) {
	adapter.Log(logging.LEVEL_ERROR, msg, kv...)
}
func (adapter *loggerAdapter) Info(msg string, kv ...interface{}) {
	adapter.Log(logging.LEVEL_INFO, msg, kv...)
}
func (adapter *loggerAdapter) Debug(msg string, kv ...interface{}) {
	adapter.Log(logging.LEVEL_DEBUG, msg, kv...)
}
func (adapter *loggerAdapter) ShouldLog(level logging.Level) bool {
	return adapter.zapLogger.Check(translateLevel(level), "") != nil
}

func translateLevel(level logging.Level) zapcore.Level {
	if level.Severity >= logging.LEVEL_FATAL.Severity {
		return zapcore.FatalLevel
	} else if level.Severity >= logging.LEVEL_ERROR.Severity {
		return zapcore.ErrorLevel
	} else if level.Severity >= logging.LEVEL_WARNING.Severity {
		return zapcore.WarnLevel
	} else if level.Severity >= logging.LEVEL_INFO.Severity {
		return zapcore.InfoLevel
	} else {
		return zapcore.DebugLevel
	}
}
