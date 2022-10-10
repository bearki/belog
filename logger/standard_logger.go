/**
 *@Title belog标准记录器
 *@Desc belog日志的标准实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"time"

	"github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/internal/encoder"
	"github.com/bearki/belog/v2/level"
)

// StandardBelog 标准记录器
type StandardBelog struct {
	*belog
}

// format 序列化为行格式
//
// 2022/09/14 20:28:13.793 [T]  k1: v1, k2: v2, ..., this is test msg\r\n
func (s *StandardBelog) format(t time.Time, l level.Level, msg string, val ...field.Field) {
	// 从对象池中获取一个日志字节流对象
	c := logBytesPool.Get()

	// 声明空栈信息
	var (
		fn string
		ln int
		mn string
	)

	// 开始追加内容
	c = encoder.AppendTime(c, t, s.timeFormat)
	c = append(c, ' ')
	c = encoder.AppendLevel(c, l, s.levelFormat)
	if s.printCallStack {
		fn, ln, mn = encoder.GetCallStack(s.stackSkip)
		c = append(c, ' ')
		c = encoder.AppendStack(c, s.callStackFullPath, fn, ln, mn)
	}
	c = append(c, ' ', ' ')
	c = encoder.AppendFieldAndMsg(c, msg, val...)
	c = append(c, "\r\n"...)

	// 选择合适的适配器执行输出
	adapterPrint := s.filterAdapterPrint()
	adapterPrint(t, l, c, fn, ln, mn)

	// 避免使用defer，会有些许性能损耗
	// 回收切片
	logBytesPool.Put(c)
}

// formatJSON 序列化为JSON格式
//
// {"$(timeKey)": "$(2006/01/02 15:04:05.000)", "$(levelKey)": "D", "$(fieldsKey)": {$("k1": v1, "k2": "v2", ...)}, "$(msgKey)": "this is test msg"}\r\n
func (s *StandardBelog) formatJSON(t time.Time, l level.Level, msg string, val ...field.Field) {
	// 从对象池中获取一个日志字节流对象
	c := logBytesPool.Get()

	// 声明空栈信息
	var (
		fn string
		ln int
		mn string
	)

	// 开始追加内容
	c = append(c, '{')
	c = encoder.AppendTimeJSON(c, s.timeJsonKey, t, s.timeFormat)
	c = append(c, `, `...)
	c = encoder.AppendLevelJSON(c, s.levelJsonKey, l, s.levelFormat)
	c = append(c, `, `...)
	if s.printCallStack {
		fn, ln, mn = encoder.GetCallStack(s.stackSkip)
		c = encoder.AppendStackJSON(
			c, s.callStackFullPath, s.stackJsonKey,
			s.stackFileJsonKey, fn,
			s.stackLineNoJsonKey, ln,
			s.stackMethodJsonKey, mn,
		)
		c = append(c, `, `...)
	}
	c = encoder.AppendFieldAndMsgJSON(c, s.messageJsonKey, msg, s.fieldsJsonKey, val...)
	c = append(c, "}\r\n"...)

	// 选择合适的适配器执行输出
	adapterPrint := s.filterAdapterPrint()
	adapterPrint(t, l, c, fn, ln, mn)

	// 避免使用defer，会有些许性能损耗
	// 回收切片
	logBytesPool.Put(c)
}

// check 高性能日志前置判断和序列化
func (s *StandardBelog) check(l level.Level, msg string, val ...field.Field) {
	// 判断当前级别日志是否需要记录
	if !s.levelIsExist(l) {
		// 当前级别日志不需要记录
		return
	}

	// 获取当前时间
	now := time.Now()

	// 是否禁用json序列化格式输出
	if s.disabledJsonFormat {
		// 执行行格式打印
		s.format(now, l, msg, val...)
	} else {
		// 执行JSON格式打印
		s.formatJSON(now, l, msg, val...)
	}
}

// Trace 通知级别的日志（高性能序列化）
func (s *StandardBelog) Trace(msg string, val ...field.Field) {
	s.check(level.Trace, msg, val...)
}

// Debug 调试级别的日志（高性能序列化）
func (s *StandardBelog) Debug(msg string, val ...field.Field) {
	s.check(level.Debug, msg, val...)
}

// Info 普通级别的日志（高性能序列化）
func (s *StandardBelog) Info(msg string, val ...field.Field) {
	s.check(level.Info, msg, val...)
}

// Warn 警告级别的日志（高性能序列化）
func (s *StandardBelog) Warn(msg string, val ...field.Field) {
	s.check(level.Warn, msg, val...)
}

// Error 错误级别的日志（高性能序列化）
func (s *StandardBelog) Error(msg string, val ...field.Field) {
	s.check(level.Error, msg, val...)
}

// Fatal 致命级别的日志（高性能序列化）
func (s *StandardBelog) Fatal(msg string, val ...field.Field) {
	s.check(level.Fatal, msg, val...)
}

// GetSugarLogger 获取语法糖记录器
func (s *StandardBelog) GetSugarLogger() SugarLogger {
	return &SugarBelog{
		belog: s.belog,
	}
}
