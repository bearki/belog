/**
 *@Title belog默认实例
 *@Desc 该实例可通过控制台打印日志
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"github.com/bearki/belog/v3/adapter/console"
	"github.com/bearki/belog/v3/encoder"
	"github.com/bearki/belog/v3/field"
	"github.com/bearki/belog/v3/logger"
)

// DefaultLog 默认实例(控制台适配器记录日志)
var DefaultLog, _ = logger.New(
	logger.Option{
		EnabledStackPrint: false,
		Encoder:           encoder.NewJsonEncoder(encoder.DefaultJsonOption),
	},
	console.New(console.Option{
		DisabledBuffer: true,
	}),
)

func init() {
	DefaultLog.SetSkip(1)
}

// Trace 通知级别的日志（默认实例）
func Trace(msg string, val ...field.Field) {
	DefaultLog.Trace(msg, val...)
}

// Debug 调试级别的日志（默认实例）
func Debug(msg string, val ...field.Field) {
	DefaultLog.Debug(msg, val...)
}

// Info 普通级别的日志（默认实例）
func Info(msg string, val ...field.Field) {
	DefaultLog.Info(msg, val...)
}

// Warn 警告级别的日志（默认实例）
func Warn(msg string, val ...field.Field) {
	DefaultLog.Warn(msg, val...)
}

// Error 错误级别的日志（默认实例）
func Error(msg string, val ...field.Field) {
	DefaultLog.Error(msg, val...)
}

// Fatal 致命级别的日志（默认实例）
func Fatal(msg string, val ...field.Field) {
	DefaultLog.Fatal(msg, val...)
}

// New 初始化一个日志记录器实例
//
//	@var option 日志记录器初始化参数
//	@var adapter 日志适配器
//	@return 日志记录器实例
func New(option logger.Option, adapter ...logger.Adapter) (logger.Logger, error) {
	// 检查编码器
	if option.Encoder == nil {
		// 赋值默认编码器
		option.Encoder = encoder.NewJsonEncoder(encoder.DefaultJsonOption)
	}
	// 返回日志示例指针
	return logger.New(option, adapter...)
}
