/**
 *@Title belog默认实例
 *@Desc 该实例可通过控制台打印日子
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"github.com/bearki/belog/console"
	"github.com/bearki/belog/logger"
)

// 默认实例(控制台引擎记录日志)
var belogDefault logger.Logger

// 初始化一个默认实例
func init() {
	belogDefault, _ = New(new(console.Engine), nil) // 初始化
	belogDefault.OpenFileLine()                     // 开启文件行号打印
	belogDefault.SetSkip(1)                         // 因为又封装了一层，故需要跳过一层函数栈
}

// New 初始化一个日志记录器实例
// @params engine  belogEngine 必选的基础日志引擎
// @params options interface{} 引擎的配置参数
// @return         *BeLog      日志记录器实例指针
// @return         error       错误信息
func New(engine logger.Engine, options interface{}) (logger.Logger, error) {
	// 返回日志示例指针
	return logger.New(engine, options)
}

// Trace 通知级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Trace(format string, val ...interface{}) {
	belogDefault.Trace(format, val...)
}

// Debug 调试级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Debug(format string, val ...interface{}) {
	belogDefault.Debug(format, val...)
}

// Info 普通级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Info(format string, val ...interface{}) {
	belogDefault.Info(format, val...)
}

// Warn 警告级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Warn(format string, val ...interface{}) {
	belogDefault.Warn(format, val...)
}

// Error 错误级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Error(format string, val ...interface{}) {
	belogDefault.Error(format, val...)
}

// Fatal 致命级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func Fatal(format string, val ...interface{}) {
	belogDefault.Fatal(format, val...)
}
