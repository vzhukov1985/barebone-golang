package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	logger *zap.SugaredLogger
	atom   zap.AtomicLevel
)

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("02/01/2006 15:04:05"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.ConsoleSeparator = " "
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	atom = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}
