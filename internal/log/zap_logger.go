package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	logLevelSeverity := map[zapcore.Level]string{
		zapcore.DebugLevel:  "DEBUG",
		zapcore.InfoLevel:   "INFO",
		zapcore.WarnLevel:   "WARNING",
		zapcore.ErrorLevel:  "ERROR",
		zapcore.DPanicLevel: "CRITICAL",
		zapcore.PanicLevel:  "ALERT",
		zapcore.FatalLevel:  "EMERGENCY",
	}
	conf := zap.NewProductionConfig()

	zap.LevelFlag("INFO", zap.InfoLevel, "info")
	conf.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "times"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(logLevelSeverity[level])
	}
	encoderConfig.StacktraceKey = ""
	conf.EncoderConfig = encoderConfig
	conf.OutputPaths = []string{
		"stderr",
	}

	l := zap.Must(conf.Build(zap.AddCallerSkip(2)))

	return &ZapLogger{
		sugar: l.Sugar(),
	}

}

func (p *ZapLogger) Debug(msg string, kv ...interface{}) {
	p.sugar.Debugw(msg, kv...)
}

func (p *ZapLogger) Info(msg string, kv ...interface{}) {
	p.sugar.Infow(msg, kv...)
}

func (p *ZapLogger) Warn(msg string, kv ...interface{}) {
	p.sugar.Warnw(msg, kv...)
}

func (p *ZapLogger) Error(msg string, err error, kv ...interface{}) {
	args := append([]interface{}{"error", err}, kv...)
	p.sugar.Errorw(msg, args...)
}

func (p *ZapLogger) Fatal(msg string, kv ...interface{}) {
	p.sugar.Fatalw(msg, kv...)
}
