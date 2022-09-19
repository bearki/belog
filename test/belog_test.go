package test

import (
	"fmt"
	"testing"

	"github.com/bearki/belog/v2"
	"github.com/bearki/belog/v2/adapter/file"
	field2 "github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/logger"
)

// BenchmarkBelogLoggerFileWrite 测试belog标准记录器文件写入
func BenchmarkBelogLoggerFileWrite(b *testing.B) {
	// 初始化文件日志适配器
	fileAdapter, err := file.New(file.Options{
		LogPath:      "../logs/test_new.log", // 日志储存路径
		MaxSize:      200,                    // 日志单文件大小
		MaxLines:     1000000,                // 单文件最大行数
		SaveDay:      7,                      // 日志保存天数
		Async:        true,                   // 开启异步写入(main函数提前结束会导致日志未写入)
		AsyncChanCap: 100,                    // 异步缓存管道容量
	})
	if err != nil {
		b.Fatalf("file adapter create failed, %s\r\n", err)
	}
	// 初始化一个实例(可实例化任意适配器)
	l, err := belog.New(logger.Option{}, fileAdapter)
	if err != nil {
		b.Fatalf("belog logger create failed, %s\r\n", err)
	}
	b.ReportAllocs()
	defer l.Flush()

	// l.PrintCallStack()

	for i := 0; i < b.N; i++ {
		l.Trace(
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
	}
}

// BenchmarkBelogLoggerFormat 测试belog标准记录器序列化
func BenchmarkBelogLoggerFormat(b *testing.B) {
	// 初始化一个实例(无适配器)
	l, err := belog.New(logger.Option{})
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}
	b.ReportAllocs()

	// l.PrintCallStack()

	for i := 0; i < b.N; i++ {
		l.Trace(
			"this is a trace log",
			// field2.Bool("key0", i%2 == 0),
			// field2.Int8("key1", 1),
			// field2.Bool("key0", i%2 == 0),
			// field2.Int8("key1", 1),
			// field2.Bool("key0", i%2 == 0),
			// field2.Int8("key1", 1),
			// field2.Bool("key0", i%2 == 0),
			// field2.Int8("key1", 1),
			// field2.Bool("key0", i%2 == 0),
			// field2.Int8("key1", 1),
		)
	}
}
