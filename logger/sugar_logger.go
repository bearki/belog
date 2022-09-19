/**
 *@Title belog语法糖代码
 *@Desc belog日志的语法糖实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"fmt"
	"time"

	"github.com/bearki/belog/v2/internal/convert"
	"github.com/bearki/belog/v2/logger/field"
)

// SugarBelog 语法糖记录器
type SugarBelog struct {
	*belog
}

// check 常规日志前置判断和序列化
//
// @params level 日志级别
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (s *SugarBelog) check(level Level, message string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if !s.levelIsExist(level) {
		// 当前级别日志不需要记录
		return
	}

	// 获取当前时间
	currTime := time.Now()
	// 数据容器
	var logBytes []byte

	// 是否禁用JSON格式输出
	if s.disabledJsonFormat {
		// 使用行格式输出

		// 静态字符串时开启优化打印
		if len(val) == 0 {
			logBytes = convert.StringToBytes(message)
		} else {
			logBytes = convert.StringToBytes(fmt.Sprintf(message, val...))
		}
	} else {
		// 获取标准记录器
		// stdLogger := s.GetLogger()

		// 使用json格式输出

		// 静态字符串时开启优化打印
		if len(val) == 0 {
			// stdLogger.preCheck(level, message)
			return
		}

		// 是否满足(message string, val1 field.Field, val2 field.Field...)
		var fields []field.Field
		for _, v := range val {
			// 是否为Field类型
			if item, ok := v.(field.Field); ok {
				fields = append(fields, item)
			}
		}
		if len(fields) == len(val) {
			// 使用高性能格式化输出
			// stdLogger.preCheck(level, message, fields...)
			return
		}

		// 是否满足(message string, key1 string, val1 interface{}, key2 string, val2 interface{}...)
		fields = fields[:]
		key := ""
		for i, v := range val {
			if i%2 == 0 {
				if item, ok := v.(string); ok {
					key = item
				} else {
					// 一旦匹配到偶数项不是字符串就立即跳出
					break
				}
			} else {
				fields = append(fields, field.Interface(key, v))
			}
		}
		if len(fields) == len(val) {
			// 使用高性能格式化输出
			// s.preCheck(level, message, fields...)
			return
		}

		// 格式化为json数据

		// 使用标准库的JSON输出
	}

	// 筛选合适的序列化方法执行序列化打印输出
	format := s.filterFormat()
	format(currTime, level, logBytes)
}

// Trace 通知级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (s *SugarBelog) Trace(format string, val ...interface{}) {
	s.check(LevelTrace, format, val...)
}

// Debug 调试级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (s *SugarBelog) Debug(format string, val ...interface{}) {
	s.check(LevelDebug, format, val...)
}

// Info 普通级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (s *SugarBelog) Info(format string, val ...interface{}) {
	s.check(LevelInfo, format, val...)
}

// Warn 警告级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (s *SugarBelog) Warn(format string, val ...interface{}) {
	s.check(LevelWarn, format, val...)
}

// Error 错误级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (s *SugarBelog) Error(format string, val ...interface{}) {
	s.check(LevelError, format, val...)
}

// Fatal 致命级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (s *SugarBelog) Fatal(format string, val ...interface{}) {
	s.check(LevelFatal, format, val...)
}
