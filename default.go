/**
 *@Title belog默认实例
 *@Desc 该实例可通过控制台打印日子
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

// 初始化一个默认实例(控制台引擎记录日志)
var belogDefault = New().
	SetEngine(EngineConsole).
	SetLevel(LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal).
	OpenFileLine().
	SetSkip(1) // 因为又封装了一层，故需要跳过一层函数栈

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
