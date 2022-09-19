package test

import (
	"io"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BenchmarkZapLoggerFormat 测试zap标准记录器序列化
func BenchmarkZapLoggerFormat(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	)
	l := zap.New(core)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Info(
			"this is a trace log",
			zap.Bool("key0", i%2 == 0),
			zap.Int8("key1", 1),
			zap.Bool("key0", i%2 == 0),
			zap.Int8("key1", 1),
			zap.Bool("key0", i%2 == 0),
			zap.Int8("key1", 1),
			zap.Bool("key0", i%2 == 0),
			zap.Int8("key1", 1),
			zap.Bool("key0", i%2 == 0),
			zap.Int8("key1", 1),
		)
	}
}
