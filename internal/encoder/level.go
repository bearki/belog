package encoder

import (
	"github.com/bearki/belog/v2/level"
)

// appendLevel 追加行格式的日志级别
//
//	@var dst 目标切片
//	@var l 级别
//	@var useFullString 使用级别的完整字符串
//	@return 序列化后的日志级别字符串
//
// 返回示例，反引号内为实际内容:
// `[T]`
func appendLevel(dst []byte, l level.Level, useFullString bool) []byte {
	dst = append(dst, '[')
	if useFullString {
		dst = append(dst, l.String()...)
	} else {
		dst = append(dst, l.Byte())
	}
	dst = append(dst, ']')
	return dst
}

// appendLevelJSON 追加行格式的日志级别
//
//	@var dst 目标切片
//	@var levelKey 日志级别JSON键名
//	@var l 级别
//	@var useFullString 使用级别的完整字符串
//	@return 序列化后的日志级别字符串
//
// 返回示例，反引号内为实际内容:
// `"level": "T"`
func appendLevelJSON(dst []byte, levelKey string, l level.Level, useFullString bool) []byte {
	dst = append(dst, '"')
	dst = append(dst, levelKey...)
	dst = append(dst, '"', ':', ' ', '"')
	if useFullString {
		dst = append(dst, l.String()...)
	} else {
		dst = append(dst, l.Byte())
	}
	dst = append(dst, '"')
	return dst
}
