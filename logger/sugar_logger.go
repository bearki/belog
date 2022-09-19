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

	field2 "github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/internal/convert"
	"github.com/bearki/belog/v2/level"
)

// SugarBelog 语法糖记录器
type SugarBelog struct {
	*belog
}

// check 常规日志前置判断和序列化
func (s *SugarBelog) check(l level.Level, msg string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if !s.levelIsExist(l) {
		// 当前级别日志不需要记录
		return
	}

	// 获取当前时间
	now := time.Now()
	// 数据容器
	var logBytes []byte

	// 是否禁用JSON格式输出
	if s.disabledJsonFormat {
		// 使用行格式输出

		// 静态字符串时开启优化打印
		if len(val) == 0 {
			logBytes = convert.StringToBytes(msg)
		} else {
			logBytes = convert.StringToBytes(fmt.Sprintf(msg, val...))
		}
	} else {
		// 获取标准记录器
		// stdLogger := s.GetLogger()

		// 使用json格式输出

		// 静态字符串时开启优化打印
		if len(val) == 0 {
			// stdLogger.preCheck(l, msg)
			return
		}

		// 是否满足(msg string, val1 field.Field, val2 field.Field...)
		var fields []field2.Field
		for _, v := range val {
			// 是否为Field类型
			if item, ok := v.(field2.Field); ok {
				fields = append(fields, item)
			}
		}
		if len(fields) == len(val) {
			// 使用高性能格式化输出
			// stdLogger.preCheck(l, msg, fields...)
			return
		}

		// 是否满足(msg string, key1 string, val1 interface{}, key2 string, val2 interface{}...)
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
				fields = append(fields, field2.Interface(key, v))
			}
		}
		if len(fields) == len(val) {
			// 使用高性能格式化输出
			// s.preCheck(l, msg, fields...)
			return
		}

		// 格式化为json数据

		// 使用标准库的JSON输出
	}

	// 筛选合适的序列化方法执行序列化打印输出
	format := s.filterFormat()
	format(now, l, logBytes)
}

// Trace 通知级别的日志
func (s *SugarBelog) Trace(msg string, val ...interface{}) {
	s.check(level.Trace, msg, val...)
}

// Debug 调试级别的日志
func (s *SugarBelog) Debug(msg string, val ...interface{}) {
	s.check(level.Debug, msg, val...)
}

// Info 普通级别的日志
func (s *SugarBelog) Info(msg string, val ...interface{}) {
	s.check(level.Info, msg, val...)
}

// Warn 警告级别的日志
func (s *SugarBelog) Warn(msg string, val ...interface{}) {
	s.check(level.Warn, msg, val...)
}

// Error 错误级别的日志
func (s *SugarBelog) Error(msg string, val ...interface{}) {
	s.check(level.Error, msg, val...)
}

// Fatal 致命级别的日志
func (s *SugarBelog) Fatal(msg string, val ...interface{}) {
	s.check(level.Fatal, msg, val...)
}
