package test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/bearki/belog/v2"
	"github.com/bearki/belog/v2/adapter/file"
	field2 "github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/logger"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestXxx(t *testing.T) {
	fmt.Println(0 << 0)
	fmt.Println(0 << 1)
	fmt.Println(1 << 0)
	fmt.Println(1 << 1)
}

func BenchmarkXxx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		zap.Int("index", i)
		p := field2.Int("index", i)
		p.Put()
	}
}

// BenchmarkNewFileBelog 实例方式输出文件日志
func BenchmarkNewFileBelog(b *testing.B) {
	// // 初始化文件日志适配器
	// fielAdapter, err := file.New(file.Options{
	// 	LogPath:      "../logs/test_new.log", // 日志储存路径
	// 	MaxSize:      200,                    // 日志单文件大小
	// 	MaxLines:     1000000,                // 单文件最大行数
	// 	SaveDay:      7,                      // 日志保存天数
	// 	Async:        true,                   // 开启异步写入(main函数提前结束会导致日志未写入)
	// 	AsyncChanCap: 100,                    // 异步缓存管道容量
	// })
	// if err != nil {
	// 	fmt.Printf("file adapter create failed, %s\r\n", err)
	// 	return
	// }
	// 初始化一个实例(可实例化任意适配器)
	mylog, err := belog.New(logger.Option{})
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}
	b.ReportAllocs()
	// defer mylog.Flush()

	// mylog.PrintCallStack()

	for i := 0; i < b.N; i++ {
		mylog.Trace(
			"this is a trace log",
			field2.Bool("key0", i%2 == 0),
			field2.Int8("key1", 1),
			field2.Bool("key0", i%2 == 0),
			field2.Int8("key1", 1),
			field2.Bool("key0", i%2 == 0),
			field2.Int8("key1", 1),
			field2.Bool("key0", i%2 == 0),
			field2.Int8("key1", 1),
			field2.Bool("key0", i%2 == 0),
			field2.Int8("key1", 1),
		)
		// mylog.Trace("this is a trace log")
	}
}

// TestNewFileBelog 实例方式输出文件日志
func TestNewFileBelog(t *testing.T) {
	// 初始化文件日志适配器
	fielAdapter, err := file.New(file.Options{
		LogPath:      "../logs/test_new.log", // 日志储存路径
		MaxSize:      200,                    // 日志单文件大小
		MaxLines:     1000000,                // 单文件最大行数
		SaveDay:      7,                      // 日志保存天数
		Async:        true,                   // 开启异步写入(main函数提前结束会导致日志未写入)
		AsyncChanCap: 1,                      // 异步缓存管道容量
	})
	if err != nil {
		fmt.Printf("file adapter create failed, %s\r\n", err)
		return
	}
	// 初始化一个实例(可实例化任意适配器)
	mylog, err := belog.New(logger.Option{}, fielAdapter)
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}
	defer mylog.Flush()
	tt := time.Now()
	for i := 0; i < 300000; i++ {
		mylog.Trace(
			"this is a trace log",
			// field.Int("index", i),
		)
		// mylog.Trace("this is a trace log, index: %v", i)
	}

	fmt.Println(time.Since(tt).Milliseconds())
}

func BenchmarkMyZap(b *testing.B) {
	file, _ := os.Create("../logs/zap.log")
	writeSyncer := zapcore.AddSync(file)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	for i := 0; i < b.N; i++ {
		logger.Error(
			"this is a trace log",
			zap.Bool("key0", i%2 == 0),
			// zap.Int8("key1", 1),
			// zap.Int32("key2", 2),
			// zap.Int64("key3", 3),
			// zap.Uint32("key4", 4),
			// zap.Uint64("key5", 5),
			// zap.Intp("key6", &i),
			// zap.String("key7", "test1"),
			// zap.String("key8", "test2"),
			// zap.String("key9", "test3"),
		)
	}
}

func TestZap(t *testing.T) {
	file, _ := os.Create("../logs/zap.log")
	writeSyncer := zapcore.AddSync(file)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	sugarLogger := logger.Sugar()
	defer sugarLogger.Sync()
	tt := time.Now()
	for i := 0; i < 300000; i++ {
		sugarLogger.Error(
			"this is a trace log",
			// zap.Int("index", i),
		)
	}
	fmt.Println(time.Since(tt).Milliseconds())
}

func BenchmarkLogrus(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(logrus.Fields{
			"url":     "http://foo.com",
			"attempt": 3,
			"backoff": time.Second,
		}).Info("failed to fetch URL")
	}
}

func BenchmarkZap(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	)
	logger := zap.New(core)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(
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
