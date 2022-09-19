package encoder

import (
	"github.com/bearki/belog/v2/level"
)

// AppendLevel 追加行格式的日志级别
//
// @params dst 目标切片
//
// @return 序列化后的日志级别字符串
//
// 返回示例，反引号内为实际内容:
// `[T]`
func AppendLevel(dst []byte, l level.Level) []byte {
	return append(dst, '[', l.GetLevelChar(), ']')
}

// AppendLevelJSON 追加行格式的日志级别
//
// @params dst 目标切片
//
// @params levelKey 日志级别JSON键名
//
// @return 序列化后的日志级别字符串
//
// 返回示例，反引号内为实际内容:
// `"level": "T"`
func AppendLevelJSON(dst []byte, levelKey string, l level.Level) []byte {
	dst = append(dst, '"')
	dst = append(dst, levelKey...)
	dst = append(dst, '"', ':', ' ', '"', l.GetLevelChar(), '"')
	return dst
}
