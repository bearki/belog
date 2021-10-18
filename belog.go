/**
 *@Title belog默认实例
 *@Desc 该实例可通过控制台打印日子
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"sync"

	"github.com/bearki/belog/console"
	"github.com/bearki/belog/logger"
)

// 默认实例(控制台引擎记录日志)
var belogDefault logger.Logger

// 仅初始化一次的引擎
var initOnce sync.Once

// New 初始化一个日志记录器实例
// @params engine  belogEngine 必选的基础日志引擎
// @params options interface{} 引擎的配置参数
// @return         *BeLog      日志记录器实例指针
// @return         error       错误信息
func New(engine logger.Engine, options interface{}) (logger.Logger, error) {
	// 返回日志示例指针
	return logger.New(engine, options)
}

// SetEngine 配置默认实现的引擎
// @params engine  Engine      引擎对象
// @params options interface{} 引擎参数
// @return         error       错误信息
func SetEngine(engine logger.Engine, options interface{}) error {
	// 判断引擎是否为空
	if belogDefault == nil { // 空引擎需要先初始化一个引擎
		var err error
		belogDefault, err = New(engine, options)
		belogDefault.SetSkip(1) // 固定的函数栈层数
		return err
	} else { // 已有实例可直接进行增加
		return belogDefault.SetEngine(engine, options)
	}
}

// SetLevel 默认实现的日志级别设置
func SetLevel(val ...logger.BeLevel) logger.Logger {
	initDefaultEngine()
	return belogDefault.SetLevel(val...)
}

// OpenFileLine 默认实现的行号打印开启
func OpenFileLine() logger.Logger {
	initDefaultEngine()
	return belogDefault.OpenFileLine()
}

// initDefaultEngine 初始化默认引擎
func initDefaultEngine() {
	// 判断引擎是否初始化
	if belogDefault == nil {
		initOnce.Do(func() {
			belogDefault, _ = New(new(console.Engine), nil)
			belogDefault.SetSkip(1) // 固定的函数栈层数
		})
	}
}

// Trace 通知级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Trace(format string, val ...interface{}) {
	// 默认引擎为空时初始化一次默认引擎
	initDefaultEngine()
	belogDefault.Trace(format, val...)
}

// Debug 调试级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Debug(format string, val ...interface{}) {
	// 默认引擎为空时初始化一次默认引擎
	initDefaultEngine()
	belogDefault.Debug(format, val...)
}

// Info 普通级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Info(format string, val ...interface{}) {
	// 默认引擎为空时初始化一次默认引擎
	initDefaultEngine()
	belogDefault.Info(format, val...)
}

// Warn 警告级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Warn(format string, val ...interface{}) {
	// 默认引擎为空时初始化一次默认引擎
	initDefaultEngine()
	belogDefault.Warn(format, val...)
}

// Error 错误级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Error(format string, val ...interface{}) {
	// 默认引擎为空时初始化一次默认引擎
	initDefaultEngine()
	belogDefault.Error(format, val...)
}

// Fatal 致命级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Fatal(format string, val ...interface{}) {
	// 默认引擎为空时初始化一次默认引擎
	initDefaultEngine()
	belogDefault.Fatal(format, val...)
}
