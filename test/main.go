package main

import "github.com/bearki/belog"

func main() {
	// 直接使用方式
	// belog.Trace("this is a trace log")
	// belog.Debug("this is a debug log")
	// belog.Info("this is a info log")
	// belog.Warn("this is a warn log")
	// belog.Error("this is a error log")
	// belog.Fatal("this is a fatal log")

	// 通过实例方式
	// 初始化一个实例(控制台引擎记录日志)
	var mylog = belog.New(belog.EngineFile, belog.FileEngineOption{
		LogPath: "./logs/app.log", // 日志储存路径
		MaxSize: 1,                // 日志单文件大小
		SaveDay: 1,                // 日志保存天数
	}).
		SetLevel(
			belog.LevelTrace,
			belog.LevelDebug,
			belog.LevelInfo,
			belog.LevelWarn,
			belog.LevelError,
			belog.LevelFatal,
		).
		OpenFileLine()
	mylog.Trace("this is a trace log")
	mylog.Debug("this is a debug log")
	mylog.Info("this is a info log")
	mylog.Warn("this is a warn log")
	mylog.Error("this is a error log")
	mylog.Fatal("this is a fatal log")

	for {
	}
}
