package logger

import (
	"time"

	"github.com/bearki/belog/v2/logger/field"
)

// Adapter 适配器接口
type Adapter interface {
	// Name 用于获取适配器名称
	//
	// 注意：请确保适配器名称不与其他适配器名称冲突
	Name() string

	// Print 普通日志打印方法
	//
	// @params logTime 日记记录时间
	//
	// @params level 日志级别
	//
	// @params content 日志内容
	Print(logTime time.Time, level Level, content []byte)

	// PrintStack 调用栈日志打印方法
	//
	// @params logTime 日记记录时间
	//
	// @params level 日志级别
	//
	// @params content 日志内容
	//
	// @params fileName 日志记录调用文件路径
	//
	// @params lineNo 日志记录调用文件行号
	//
	// @params methodName 日志记录调用函数名
	PrintStack(logTime time.Time, level Level, content []byte, fileName []byte, lineNo int, methodName []byte)

	// Flush 日志缓存刷新
	//
	// 用于日志缓冲区刷新,
	// 接收到该通知后需要立即将缓冲区中的日志持久化,
	// 因为程序很有可能将在短时间内退出
	Flush()
}

// BaseLogger 基础日志接口
type BaseLogger interface {
	SetAdapter(Adapter) error // 适配器设置
	SetLevel(...Level)        // 日志级别设置
	SetSkip(uint)             // 函数栈配置
	Flush()                   // 日志缓存刷新
}

// Logger 标准日志接口
type Logger interface {
	BaseLogger
	GetSugarLogger() SugarLogger  // 获取语法糖记录器
	Trace(string, ...field.Field) // 通知级别的日志（高性能序列化）
	Debug(string, ...field.Field) // 调试级别的日志（高性能序列化）
	Info(string, ...field.Field)  // 普通级别的日志（高性能序列化）
	Warn(string, ...field.Field)  // 警告级别的日志（高性能序列化）
	Error(string, ...field.Field) // 错误级别的日志（高性能序列化）
	Fatal(string, ...field.Field) // 致命级别的日志（高性能序列化）
}

// SugarLogger 语法糖日志接口
type SugarLogger interface {
	BaseLogger
	Trace(string, ...interface{}) // 通知级别的日志
	Debug(string, ...interface{}) // 调试级别的日志
	Info(string, ...interface{})  // 普通级别的日志
	Warn(string, ...interface{})  // 警告级别的日志
	Error(string, ...interface{}) // 错误级别的日志
	Fatal(string, ...interface{}) // 致命级别的日志
}
