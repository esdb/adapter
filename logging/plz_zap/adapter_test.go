package plz_zap

import (
	"testing"
	"github.com/v2pro/plz"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
	"github.com/v2pro/plz/logging"
)

func Test_zap(t *testing.T) {
	logger, _ := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   nil,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	logging.Providers = append(logging.Providers, func(loggerKv []interface{}) logging.Logger {
		return Adapt(logger)
	})
	plz.Log().Info("Failed to fetch URL.", "hello", "world")
}
