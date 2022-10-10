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

// format 序列化格式
func (s *StandardBelog) format(t time.Time, l level.Level, msg string, val ...field.Field) {
	// 从对象池中获取一个日志字节流对象
	dst := logBytesPool.Get()

	// 构建拼接参数
	opt := &encoder.Option{
		TimeKey:     s.timeJsonKey,
		TimeFormat:  s.timeFormat,
		LevelKey:    s.levelJsonKey,
		LevelFormat: s.levelFormat,
		MsgKey:      s.messageJsonKey,
		FieldsKey:   s.fieldsJsonKey,
	}

	// 是否打印调用栈
	if s.printCallStack {
		opt.StackSkip = s.stackSkip
		opt.StackKey = s.stackJsonKey
		opt.StackFileKey = s.stackFileJsonKey
		opt.StackFileFormat = s.stackFileFormat
		opt.StackLineNoKey = s.stackLineNoJsonKey
		opt.StackMethodKey = s.stackMethodJsonKey
		opt.StackFile = ""
		opt.StackLineNo = 0
		opt.StackMethod = ""
	}

	// 是否禁用JSON编码
	if s.disabledJsonFormat {
		// 执行日志普通编码
		dst = opt.EncodeNormal(dst, t, l, msg, val...)
	} else {
		// 执行日志JSON编码
		dst = opt.EncodeJSON(dst, t, l, msg, val...)
	}

	// 选择合适的适配器执行输出
	adapterPrint := s.filterAdapterPrint()
	adapterPrint(t, l, dst, opt.StackFile, opt.StackLineNo, opt.StackMethod)

	// 避免使用defer，会有些许性能损耗
	// 回收切片
	logBytesPool.Put(dst)
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
	// 执行格式化打印
	s.format(now, l, msg, val...)
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
