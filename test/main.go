package main

import (
	"github.com/bearki/belog"
	"github.com/bearki/belog/file"
	"github.com/bearki/belog/logger"
)

// 全局日志对象
var Log logger.Logger

func main() {
	// 直接使用方式（仅能输出到控制台）
	belog.Trace("this is a trace log")
	belog.Debug("this is a debug log")
	belog.Info("this is a info log")
	belog.Warn("this is a warn log")
	belog.Error("this is a error log")
	belog.Fatal("this is a fatal log")

	// 实例方式
	// 初始化一个实例(可实例化任意引擎)
	mylog, err := belog.New(
		new(file.Engine), // 初始化文件引擎
		file.Options{
			LogPath:      "./logs/app.log", // 日志储存路径
			MaxSize:      128,              // 日志单文件大小
			SaveDay:      7,                // 日志保存天数
			Async:        true,             // 开启异步写入
			AsyncChanCap: 20,               // 异步缓存管道容量
		},
	)
	if err != nil { // 初始化失败将不能执行任何后续操作，否则会引起恐慌
		panic("belog init error: " + err.Error())
	}

	// 配置日志记录级别
	mylog.SetLevel(
		logger.LevelTrace,
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	).OpenFileLine() // 开启行号记录

	/**********支持多引擎同时输出**********/
	/*************************************/
	// err = mylog.SetEngine( // 增加控制台引擎
	// 	new(console.Engine),
	// 	nil,
	// )
	// if err != nil { // 该错误不影响已添加的引擎
	// 	mylog.Error("add console engine to log error: %s", err.Error())
	// }
	/*************************************/

	// 实例对象记录日志
	mylog.Trace("this is a trace log")
	mylog.Debug("this is a debug log")
	mylog.Info("this is a info log")
	mylog.Warn("this is a warn log")
	mylog.Error("this is a error log")
	mylog.Fatal("this is a fatal log")
}
