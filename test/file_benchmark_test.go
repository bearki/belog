package test

import (
	"testing"

	"github.com/bearki/belog/v3"
	"github.com/bearki/belog/v3/adapter/file"
	"github.com/bearki/belog/v3/encoder"
	"github.com/bearki/belog/v3/field"
	"github.com/bearki/belog/v3/logger"
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

	// 配置日志记录器参数
	opt := logger.Option{
		EnabledStackPrint: true,
		Encoder:           encoder.NewJsonEncoder(encoder.DefaultJsonOption),
	}

	// 初始化一个实例(可实例化任意适配器)
	l, err := belog.New(opt, fileAdapter)
	if err != nil {
		b.Fatalf("belog logger create failed, %s\r\n", err)
	}

	// 程序结束需要刷新缓冲区
	defer l.Flush()

	// 重置测试参数
	b.ReportAllocs()
	b.ResetTimer()

	// 开始测试
	for i := 0; i < b.N; i++ {
		l.Trace(
			"this is a trace log",
			field.Bool("key0", i%2 == 0),
			field.Int8("key1", 1),
			field.Ints("int10", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}),
			field.Interface("any", map[string]int{"1": 1, "2": 2, "3": 3}),
		)
	}
}
