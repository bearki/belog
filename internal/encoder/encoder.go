package encoder

import (
	"time"

	"github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/level"
)

type Option struct {
	// 普通输入
	TimeKey     string
	TimeFormat  string
	LevelKey    string
	LevelFormat bool
	MsgKey      string
	FieldsKey   string

	// 调用栈输入
	StackSkip       uint
	StackKey        string
	StackFileKey    string
	StackFileFormat bool
	StackLineNoKey  string
	StackMethodKey  string

	// 调用栈输出
	StackFile   string
	StackLineNo int
	StackMethod string
}

// EncodeNormal 追加普通格式日志
func (opt *Option) EncodeNormal(dst []byte, t time.Time, l level.Level, m string, val ...field.Field) []byte {
	// 开始追加内容
	dst = appendTime(dst, t, opt.TimeFormat)
	dst = append(dst, ' ')
	dst = appendLevel(dst, l, opt.LevelFormat)
	// 是否追加调用栈
	if opt.StackSkip > 0 {
		opt.StackFile, opt.StackLineNo, opt.StackMethod = getCallStack(opt.StackSkip)
		dst = append(dst, ' ')
		dst = appendStack(dst, opt.StackFileFormat, opt.StackFile, opt.StackLineNo, opt.StackMethod)
	}
	dst = append(dst, ' ', ' ')
	dst = appendFieldAndMsg(dst, m, val...)
	dst = append(dst, "\r\n"...)
	// 追加完成
	return dst
}

// EncodeJSON 编码为JSON格式日志
func (opt *Option) EncodeJSON(dst []byte, t time.Time, l level.Level, m string, val ...field.Field) []byte {
	// 开始追加内容
	dst = append(dst, '{')
	dst = appendTimeJSON(dst, opt.TimeKey, t, opt.TimeFormat)
	dst = append(dst, `, `...)
	dst = appendLevelJSON(dst, opt.LevelKey, l, opt.LevelFormat)
	dst = append(dst, `, `...)
	// 是否追加调用栈
	if opt.StackSkip > 0 {
		opt.StackFile, opt.StackLineNo, opt.StackMethod = getCallStack(opt.StackSkip)
		dst = appendStackJSON(
			dst, opt.StackFileFormat, opt.StackKey,
			opt.StackFileKey, opt.StackFile,
			opt.StackLineNoKey, opt.StackLineNo,
			opt.StackMethodKey, opt.StackMethod,
		)
		dst = append(dst, `, `...)
	}
	dst = appendFieldAndMsgJSON(dst, opt.MsgKey, m, opt.FieldsKey, val...)
	dst = append(dst, "}\r\n"...)
	// 追加完成
	return dst
}
