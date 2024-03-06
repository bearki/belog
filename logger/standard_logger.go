/**
 *@Title belog标准记录器
 *@Desc belog日志的标准实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"time"

	"github.com/bearki/belog/v3/field"
)

// StandardBelog 标准记录器
type StandardBelog struct {
	*belog
}

// 序列化格式
func (s *StandardBelog) format(t time.Time, l Level, msg string, val ...field.Field) {
	// 从对象池中获取一个日志字节流对象
	dst := logBytesPool.Get()

	// 是否需要调用栈
	if s.enabledStackPrint {
		fn, ln, mn := getCallStack(s.stackSkip)
		dst = s.encoder.EncodeStack(dst, t, l, fn, ln, mn, msg, val...)
		s.adapterPrintStack(t, l, dst, fn, ln, mn)
	} else {
		dst = s.encoder.Encode(dst, t, l, msg, val...)
		s.adapterPrint(t, l, dst)
	}

	// 避免使用defer，会有些许性能损耗
	// 回收切片
	logBytesPool.Put(dst)
}

// 高性能日志前置判断和序列化
func (s *StandardBelog) check(l Level, msg string, val ...field.Field) {
	// 判断当前级别日志是否需要记录
	if l < s.minLevel {
		return
	}

	// 编码器是否为空
	if s.encoder == nil {
		return
	}

	// 获取当前时间
	now := time.Now()
	// 执行格式化打印
	s.format(now, l, msg, val...)
}

// Trace 通知级别的日志（高性能序列化）
func (s *StandardBelog) Trace(msg string, val ...field.Field) {
	s.check(Trace, msg, val...)
}

// Debug 调试级别的日志（高性能序列化）
func (s *StandardBelog) Debug(msg string, val ...field.Field) {
	s.check(Debug, msg, val...)
}

// Info 普通级别的日志（高性能序列化）
func (s *StandardBelog) Info(msg string, val ...field.Field) {
	s.check(Info, msg, val...)
}

// Warn 警告级别的日志（高性能序列化）
func (s *StandardBelog) Warn(msg string, val ...field.Field) {
	s.check(Warn, msg, val...)
}

// Error 错误级别的日志（高性能序列化）
func (s *StandardBelog) Error(msg string, val ...field.Field) {
	s.check(Error, msg, val...)
}

// Fatal 致命级别的日志（高性能序列化）
func (s *StandardBelog) Fatal(msg string, val ...field.Field) {
	s.check(Fatal, msg, val...)
}

// GetSugarLogger 获取语法糖记录器
func (s *StandardBelog) GetSugarLogger() SugarLogger {
	return &SugarBelog{
		belog: s.belog,
	}
}
