package encoder

import (
	"github.com/bearki/belog/v3/logger"
)

// appendLevel 追加行格式的日志级别
//
//	@param	dst				目标切片
//	@param	l				级别
//	@param	useFullString	是否使用全称
//	@return	序列化后的日志级别字符串
//
// 返回示例: [T]
func appendLevel(dst []byte, l logger.Level, useFullString bool) []byte {
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
//	@param	dst				目标切片
//	@param	levelKey		日志级别JSON键名
//	@param	l				级别
//	@param	useFullString	是否使用全称
//	@return	序列化后的日志级别字符串
//
// 返回示例: "level": "T"
func appendLevelJSON(dst []byte, levelKey string, l logger.Level, useFullString bool) []byte {
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
