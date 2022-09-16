package test

import (
	"testing"

	"github.com/bearki/belog/v2"
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
