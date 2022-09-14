/**
 *@Title belog默认实例
 *@Desc 该实例可通过控制台打印日子
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"sync"

	"github.com/bearki/belog/v2/adapter/console"
	"github.com/bearki/belog/v2/logger"
)

var (
	// 默认实例(控制台适配器记录日志)
	belogDefault logger.Logger
	// 是否已经重新配置过默认适配器
	isSetBelogDefault = false
	// 仅初始化一次的适配器
	initOnce sync.Once
)

// New 初始化一个日志记录器实例
//
// @params adapter 日志适配器
//
// @return 日志记录器实例
func New(Adapter logger.Adapter) (logger.Logger, error) {
	// 返回日志示例指针
	return logger.New(Adapter)
}

// SetAdapter 配置默认实例适配器
//
// @params adapter 适配器实例
//
// @return error 错误信息
func SetAdapter(Adapter logger.Adapter) error {
	// 判断默认适配器是否已经重新配置过
	if !isSetBelogDefault {
		// 未配置过适配器，清空自带的console适配器
		belogDefault = nil
	}
	// 判断适配器是否为空
	if belogDefault == nil { // 空适配器需要先初始化一个适配器
		var err error
		belogDefault, err = New(Adapter)
		if err != nil {
			return err
		}
		belogDefault.SetSkip(1)  // 固定的函数栈层数
		isSetBelogDefault = true // 赋值适配器已经配置
		return nil
	} else { // 已有实例可直接进行增加
		return belogDefault.SetAdapter(Adapter)
	}
}

// SetLevel 默认实例日志级别设置
//
// @params val 日志记录级别（会覆盖上一次的级别配置）
func SetLevel(val ...logger.Level) logger.Logger {
	initDefaultAdapter()
	return belogDefault.SetLevel(val...)
}

// OpenFileLine 默认实例是否记录调用栈
//
// @return 日志记录器实例
func PrintCallStack() logger.Logger {
	initDefaultAdapter()
	return belogDefault.PrintCallStack()
}

// initDefaultAdapter 初始化默认适配器
func initDefaultAdapter() {
	// 判断适配器是否初始化
	if belogDefault == nil {
		initOnce.Do(func() {
			belogDefault, _ = New(console.New())
			belogDefault.SetSkip(1) // 固定的函数栈层数
		})
	}
}

// Trace 通知级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func Trace(format string, val ...interface{}) {
	// 默认适配器为空时初始化一次默认适配器
	initDefaultAdapter()
	belogDefault.Trace(format, val...)
}

// Debug 调试级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func Debug(format string, val ...interface{}) {
	// 默认适配器为空时初始化一次默认适配器
	initDefaultAdapter()
	belogDefault.Debug(format, val...)
}

// Info 普通级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func Info(format string, val ...interface{}) {
	// 默认适配器为空时初始化一次默认适配器
	initDefaultAdapter()
	belogDefault.Info(format, val...)
}

// Warn 警告级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func Warn(format string, val ...interface{}) {
	// 默认适配器为空时初始化一次默认适配器
	initDefaultAdapter()
	belogDefault.Warn(format, val...)
}

// Error 错误级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func Error(format string, val ...interface{}) {
	// 默认适配器为空时初始化一次默认适配器
	initDefaultAdapter()
	belogDefault.Error(format, val...)
}

// Fatal 致命级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func Fatal(format string, val ...interface{}) {
	// 默认适配器为空时初始化一次默认适配器
	initDefaultAdapter()
	belogDefault.Fatal(format, val...)
}
