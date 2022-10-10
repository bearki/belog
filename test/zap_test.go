package test

import (
	"io"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BenchmarkZapLoggerFormat 测试zap标准记录器序列化
func BenchmarkZapLoggerFormat(b *testing.B) {
	cfg := zap.NewProductionConfig()
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
			zap.Ints("int10", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}),
		)
	}
}
