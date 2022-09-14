package test

import (
	"os"
	"testing"

	"github.com/bearki/belog"
	"github.com/bearki/belog/console"
	"github.com/bearki/belog/file"
	"github.com/bearki/belog/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestDefultBelog 默认方式输出日志
func TestDefultBelog(t *testing.T) {
	// 配置主引擎
	err := belog.SetEngine(file.New(), file.Options{
		LogPath: "./logs/test_default.log",
		MaxSize: 10,
		SaveDay: 7,
		Async:   false,
	})
	if err != nil {
		panic(err)
	}
	// 配置次引擎（次引擎配置失败不影响主引擎）
	err = belog.SetEngine(console.New(), nil)
	if err != nil {
		belog.Error("次引擎配置失败不影响主引擎: %s", err.Error())
	}
	// 开启行号记录
	belog.OpenFileLine()
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

// TestNewBelog 实例方式输出日志
func TestNewBelog(t *testing.T) {
	// 初始化一个实例(可实例化任意引擎)
	mylog, err := belog.New(
		file.New(), // 初始化文件引擎
		file.Options{
			LogPath:      "../logs/test_new.log", // 日志储存路径
			MaxSize:      100,                    // 日志单文件大小
			MaxLines:     1000000,                // 单文件最大行数
			SaveDay:      7,                      // 日志保存天数
			Async:        false,                  // 开启异步写入(main函数提前结束会导致日志未写入)
			AsyncChanCap: 100,                    // 异步缓存管道容量
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

	// 实例对象记录日志
	for i := 0; i < 1000001; i++ {
		mylog.Trace("this is a trace log")
	}
}

func TestConsoleLog(t *testing.T) {
	writeSyncer, _ := os.Create("../logs/zap.log")    //日志文件存放目录
	encoderConfig := zap.NewProductionEncoderConfig() //指定时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)               //获取编码器,NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	log := zap.New(core, zap.AddCaller())                             //AddCaller()为显示文件名和行号

	for i := 0; i < 1000000; i++ {
		log.Info("this is a info log")
	}
}
