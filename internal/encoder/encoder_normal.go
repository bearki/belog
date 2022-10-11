package encoder

import (
	"time"

	"github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/level"
)

// NormalEncoderOption 普通编码器参数
type NormalEncoderOption struct {
	BaseOption
}

// DefaultNormalEncoderOption 普通编码器默认参数
var DefaultNormalOption = NormalEncoderOption{
	BaseOption: DefaultBaseOption,
}

// NormalEncoder 普通编码器
type NormalEncoder struct {
	opt NormalEncoderOption
}

// checkNormalOptionValid 检查普通编码器参数有效性
func checkNormalOptionValid(opt NormalEncoderOption) NormalEncoderOption {
	// 检查基础参数有效性
	opt.BaseOption = checkBaseOptionValid(opt.BaseOption)
	// 检查完成
	return opt
}

// NewNormalEncoder 创建一个普通格式编码器
func NewNormalEncoder(opt NormalEncoderOption) *NormalEncoder {
	// 检查参数有效性
	opt = checkNormalOptionValid(opt)
	// 创建编码器
	return &NormalEncoder{
		opt: opt,
	}
}

// Encode 编码输出方法
//
// @params dst 填充目标
//
// @params t 日志记录时间
//
// @params l 日志级别
//
// @params msg 日志描述
//
// @params val 日志内容字段
//
// @return 填充后的内容
func (e *NormalEncoder) Encode(dst []byte, t time.Time, l level.Level, msg string, val ...field.Field) []byte {
	// 开始追加内容
	dst = appendTime(dst, t, e.opt.TimeFormat)
	dst = append(dst, ' ')
	dst = appendLevel(dst, l, e.opt.LevelFormat)
	// 追加消息和字段内容
	dst = append(dst, ' ', ' ')
	dst = appendFieldAndMsg(dst, msg, val...)
	dst = append(dst, "\r\n"...)
	// 追加完成
	return dst
}

// EncodeStack 含调用栈编码输出方法
//
// @params dst 填充目标
//
// @params t 日志记录时间
//
// @params l 日志级别
//
// @params fn 调用栈文件名
//
// @params ln 调用栈行号
//
// @params mn 调用栈函数名
//
// @params msg 日志描述
//
// @params val 日志内容字段
//
// @return 填充后的内容
func (e *NormalEncoder) EncodeStack(dst []byte, t time.Time, l level.Level, fn string, ln int, mn string, msg string, val ...field.Field) []byte {
	// 开始追加内容
	dst = appendTime(dst, t, e.opt.TimeFormat)
	dst = append(dst, ' ')
	dst = appendLevel(dst, l, e.opt.LevelFormat)
	// 追加调用栈
	dst = append(dst, ' ')
	dst = appendStack(dst, e.opt.StackFileFormat, fn, ln, mn)
	// 追加消息和字段内容
	dst = append(dst, ' ', ' ')
	dst = appendFieldAndMsg(dst, msg, val...)
	dst = append(dst, "\r\n"...)
	// 追加完成
	return dst
}
