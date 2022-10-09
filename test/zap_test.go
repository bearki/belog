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
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	)
	l := zap.New(core)

	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		l.Info(
			"this is a trace log",
			zap.Bool("key0", i%2 == 0),
			zap.Int8("key1", 1),
			zap.Bool("key2", i%2 == 0),
			zap.Int8("key3", 1),
			zap.Bool("key4", i%2 == 0),
			zap.Int8("key5", 1),
			zap.Bool("key6", i%2 == 0),
			zap.Int8("key7", 1),
			zap.Bool("key8", i%2 == 0),
			zap.String("key9", "不是编号编号分别为发布会我"),
		)
	}
}
