package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bearki/belog/v2"
	"github.com/bearki/belog/v2/adapter/console"
	"github.com/bearki/belog/v2/adapter/file"
	"github.com/bearki/belog/v2/logger"
)

// TestDefultBelog 默认方式输出日志
func TestDefultBelog(t *testing.T) {
	// 记录调用栈
	belog.PrintCallStack()
	// 指定需要记录的日志级别（默认全部级别）
	belog.SetLevel(
		logger.LevelTrace,
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	)
	// 打印日志
	belog.Trace("this is a trace log")
	belog.Debug("this is a debug log")
	belog.Info("this is a info log")
	belog.Warn("this is a warn log")
	belog.Error("this is a error log")
	belog.Fatal("this is a fatal log")
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
		AsyncChanCap: 100,                    // 异步缓存管道容量
	})
	if err != nil {
		panic("file adapter create failed, " + err.Error())
	}

	// 初始化一个实例(可实例化任意适配器)
	mylog, err := belog.New(fielAdapter)
	if err != nil {
		panic("belog logger create failed, " + err.Error())
	}

	// 配置日志记录级别
	mylog.SetLevel(
		logger.LevelTrace,
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	)

	// 开启调用栈打印
	mylog.PrintCallStack()

	tt := time.Now()
	// 实例对象记录日志
	for i := 0; i < 100000; i++ {
		mylog.Tracef(
			"this is a trace log",
			logger.Intf("value", i),
			logger.Intf("index", i),
			logger.Intf("bba", i),
		)
	}
	mylog.Flush()
	fmt.Println(time.Since(tt).Milliseconds())
}

// TestNewFileBelog 实例方式输出文件和控制台日志
func TestNewFileAndConsoleBelog(t *testing.T) {
	// 初始化文件日志适配器
	fielAdapter, err := file.New(file.Options{
		LogPath:      "../logs/test_new.log", // 日志储存路径
		MaxSize:      100,                    // 日志单文件大小
		MaxLines:     1000000,                // 单文件最大行数
		SaveDay:      7,                      // 日志保存天数
		Async:        true,                   // 开启异步写入(main函数提前结束会导致日志未写入)
		AsyncChanCap: 100,                    // 异步缓存管道容量
	})
	if err != nil {
		panic("file adapter create failed, " + err.Error())
	}

	// 初始化一个实例(可实例化任意适配器)
	mylog, err := belog.New(fielAdapter)
	if err != nil {
		panic("belog logger create failed, " + err.Error())
	}

	// 配置日志记录级别
	mylog.SetLevel(
		logger.LevelTrace,
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	)

	// 开启调用栈打印
	mylog.PrintCallStack()

	// 添加控制台适配器
	err = mylog.SetAdapter(console.New())
	if err != nil {
		panic("add belog logger adapter failed, " + err.Error())
	}

	// 实例对象记录日志
	for i := 0; i < 1000000; i++ {
		mylog.Trace("this is a trace log")
	}
}
