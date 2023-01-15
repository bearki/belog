package encoder

import (
	"strings"
	"time"

	"github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/level"
)

// JsonEncoderOption JSON编码器参数
type JsonEncoderOption struct {
	BaseOption

	//日志记录时间JSON键名
	//
	// Default: "time"
	TimeKey string

	// 日志级别JSON键名
	//
	// Default: "level"
	LevelKey string

	// 日志消息JSON键名
	//
	// Default: "message"
	MsgKey string

	// 日志字段JSON键名
	//
	// Default: "fields"
	FieldsKey string

	// 调用栈JSON键名
	//
	// Default: "stack"
	StackKey string

	// 调用栈文件名JSON键名
	//
	// Default: "file"
	StackFileKey string

	// 调用栈行号JSON键名
	//
	// Default: "line"
	StackLineNoKey string

	// 调用栈函数名JSON键名
	//
	// Default: "method"
	StackMethodKey string
}

// DefaultJsonOption JSON编码器参数
var DefaultJsonOption = JsonEncoderOption{
	BaseOption:     DefaultBaseOption,
	TimeKey:        "time",
	LevelKey:       "level",
	MsgKey:         "message",
	FieldsKey:      "fields",
	StackKey:       "stack",
	StackFileKey:   "file",
	StackLineNoKey: "line",
	StackMethodKey: "method",
}

// JsonEncoder JSON编码器
type JsonEncoder struct {
	opt JsonEncoderOption
}

// checkJsonOptionValid 检查普通编码器参数有效性
func checkJsonOptionValid(opt JsonEncoderOption) JsonEncoderOption {
	// 检查基础参数有效性
	opt.BaseOption = checkBaseOptionValid(opt.BaseOption)
	// 构建剩余参数默认值
	if len(opt.TimeKey) == 0 {
		opt.TimeKey = DefaultJsonOption.TimeKey
	} else {
		opt.TimeKey = strings.ReplaceAll(opt.TimeKey, `"`, `\"`)
	}
	if len(opt.LevelKey) == 0 {
		opt.LevelKey = DefaultJsonOption.LevelKey
	} else {
		opt.LevelKey = strings.ReplaceAll(opt.LevelKey, `"`, `\"`)
	}
	if len(opt.MsgKey) == 0 {
		opt.MsgKey = DefaultJsonOption.MsgKey
	} else {
		opt.MsgKey = strings.ReplaceAll(opt.MsgKey, `"`, `\"`)
	}
	if len(opt.FieldsKey) == 0 {
		opt.FieldsKey = DefaultJsonOption.FieldsKey
	} else {
		opt.FieldsKey = strings.ReplaceAll(opt.FieldsKey, `"`, `\"`)
	}
	if len(opt.StackKey) == 0 {
		opt.StackKey = DefaultJsonOption.StackKey
	} else {
		opt.StackKey = strings.ReplaceAll(opt.StackKey, `"`, `\"`)
	}
	if len(opt.StackFileKey) == 0 {
		opt.StackFileKey = DefaultJsonOption.StackFileKey
	} else {
		opt.StackFileKey = strings.ReplaceAll(opt.StackFileKey, `"`, `\"`)
	}
	if len(opt.StackLineNoKey) == 0 {
		opt.StackLineNoKey = DefaultJsonOption.StackLineNoKey
	} else {
		opt.StackLineNoKey = strings.ReplaceAll(opt.StackLineNoKey, `"`, `\"`)
	}
	if len(opt.StackMethodKey) == 0 {
		opt.StackMethodKey = DefaultJsonOption.StackMethodKey
	} else {
		opt.StackMethodKey = strings.ReplaceAll(opt.StackMethodKey, `"`, `\"`)
	}
	// 检查完成
	return opt
}

// NewJsonEncoder 创建一个JSON格式编码器
func NewJsonEncoder(opt JsonEncoderOption) *JsonEncoder {
	// 检查参数有效性
	opt = checkJsonOptionValid(opt)
	// 创建编码器
	return &JsonEncoder{
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
func (e *JsonEncoder) Encode(dst []byte, t time.Time, l level.Level, msg string, val ...field.Field) []byte {
	// 开始追加内容
	dst = append(dst, '{')
	dst = appendTimeJSON(dst, e.opt.TimeKey, t, e.opt.TimeFormat)
	dst = append(dst, `, `...)
	dst = appendLevelJSON(dst, e.opt.LevelKey, l, e.opt.LevelFormat)
	dst = append(dst, `, `...)
	// 追加消息和字段内容
	dst = appendFieldAndMsgJSON(dst, e.opt.MsgKey, msg, e.opt.FieldsKey, val...)
	dst = append(dst, "}\r\n"...)
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
func (e *JsonEncoder) EncodeStack(dst []byte, t time.Time, l level.Level, fn string, ln int, mn string, msg string, val ...field.Field) []byte {
	// 开始追加内容
	dst = append(dst, '{')
	dst = appendTimeJSON(dst, e.opt.TimeKey, t, e.opt.TimeFormat)
	dst = append(dst, `, `...)
	dst = appendLevelJSON(dst, e.opt.LevelKey, l, e.opt.LevelFormat)
	dst = append(dst, `, `...)
	// 追加调用栈
	dst = appendStackJSON(
		dst, e.opt.StackFileFormat, e.opt.StackKey,
		e.opt.StackFileKey, fn,
		e.opt.StackLineNoKey, ln,
		e.opt.StackMethodKey, mn,
	)
	dst = append(dst, `, `...)
	// 追加消息和字段内容
	dst = appendFieldAndMsgJSON(dst, e.opt.MsgKey, msg, e.opt.FieldsKey, val...)
	dst = append(dst, "}\r\n"...)
	// 追加完成
	return dst
}
