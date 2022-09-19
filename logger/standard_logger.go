/**
 *@Title belog标准记录器
 *@Desc belog日志的标准实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"time"

	"github.com/bearki/belog/v2/logger/field"
)

// SugarBelog 标准记录器
type StandardBelog struct {
	*belog
}

// GetSugarLogger 获取语法糖记录器
func (b *StandardBelog) GetSugarLogger() SugarLogger {
	return &SugarBelog{
		belog: b.belog,
	}
}

// check 高性能日志前置判断和序列化
//
// @params level 日志级别
//
// @params message 日志消息
//
// @params val 字段信息
func (s *StandardBelog) check(level Level, message string, val ...field.Field) {
	// 判断当前级别日志是否需要记录
	if !s.levelIsExist(level) {
		// 当前级别日志不需要记录
		return
	}

	// 获取当前时间
	currTime := time.Now()
	// 声明日志容器
	var logBytes []byte

	// 是否禁用json序列化格式输出
	if s.disabledJsonFormat {
		// 结构化为:
		//
		// = 固定长度
		// + 动态长度
		//
		// k1: v1, k2: v2, ..., this is test message
		// +++++++++++++++++++++++++++++++++++++++++
		//
		// 全是动态长度

		// 计算byte大小
		size := len(message)
		for _, v := range val {
			// key: 10,  => "": => 4个字节
			// 每个键和值之间有一个冒号和一个空格 => 2个字节
			// 每个值后面有一个逗号和一个空格 => 2个字节
			size += 4 + len(v.KeyBytes) + len(v.ValBytes)
		}

		// 从对象池中获取一个日志字节流对象
		logBytes := logBytesPool.Get()
		// 用完后字节切片放回复用池
		defer func() {
			logBytesPool.Put(logBytes)
		}()

		// 判读对象容量是否足够
		if cap(logBytes) < size {
			logBytes = make([]byte, 0, size)
		}

		// 将字段和消息拼接为行格式
		logBytes = field.Append(logBytes, message, val...)
	} else {
		// 结构化为:
		//
		// = 固定长度
		// + 动态长度
		//
		// "$(fields)": {$("k1": v1, "k2": "v2", ...)}, "$(message)": "this is test message"
		// =+++++++++====++++++++++++++++++++++++++++====++++++++++====++++++++++++++++++++=
		// 1 len(fk)   4      sum(len(field)...)       4   len(mk)   4        len(m)       1
		//
		// const count = 1 + 4 + 4 + 4 + 1 = 14

		// 计算Byte大小
		size := 14 + len(s.fieldsJsonKey) + len(s.messageJsonKey) + len(message)
		for _, v := range val {
			// "key": 10 => "": => 4个字节
			// 每个键有两个双引号包裹 => 2个字节
			// 每个键和值之间有一个冒号和一个空格 => 2个字节
			size += 4 + len(v.KeyBytes) + len(v.ValPrefixBytes) + len(v.ValBytes) + len(v.ValSuffixBytes)
		}

		// 从对象池中获取一个日志字节流对象
		logBytes := logBytesPool.Get()
		// 用完后字节切片放回复用池
		defer func() {
			logBytesPool.Put(logBytes)
		}()

		// 判读对象容量是否足够
		if cap(logBytes) < size {
			logBytes = make([]byte, 0, size)
		}

		// 将字段和消息拼接为json格式
		logBytes = field.AppendJSON(logBytes, s.fieldsJsonKey, s.messageJsonKey, false, message, val...)
	}

	// 筛选合适的序列化方法执行序列化打印输出
	format := s.filterFormat()
	format(currTime, level, logBytes)
}

// Tracef 通知级别的日志（高性能序列化）
func (s *StandardBelog) Trace(message string, val ...field.Field) {
	s.check(LevelTrace, message, val...)
}

// Debugf 调试级别的日志（高性能序列化）
func (s *StandardBelog) Debug(message string, val ...field.Field) {
	s.check(LevelDebug, message, val...)
}

// Infof 普通级别的日志（高性能序列化）
func (s *StandardBelog) Info(message string, val ...field.Field) {
	s.check(LevelInfo, message, val...)
}

// Warnf 警告级别的日志（高性能序列化）
func (s *StandardBelog) Warn(message string, val ...field.Field) {
	s.check(LevelWarn, message, val...)
}

// Errorf 错误级别的日志（高性能序列化）
func (s *StandardBelog) Error(message string, val ...field.Field) {
	s.check(LevelError, message, val...)
}

// Fatalf 致命级别的日志（高性能序列化）
func (s *StandardBelog) Fatal(message string, val ...field.Field) {
	s.check(LevelFatal, message, val...)
}
