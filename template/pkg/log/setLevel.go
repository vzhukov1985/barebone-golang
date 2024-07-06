package log

import "go.uber.org/zap/zapcore"

type Level = zapcore.Level

const (
	DebugLevel Level = zapcore.DebugLevel
	InfoLevel  Level = zapcore.InfoLevel
	WarnLevel  Level = zapcore.WarnLevel
	ErrorLevel Level = zapcore.ErrorLevel
	PanicLevel Level = zapcore.PanicLevel
	FatalLevel Level = zapcore.FatalLevel
)

func SetLevel(l Level) {
	atom.SetLevel(l)
}
