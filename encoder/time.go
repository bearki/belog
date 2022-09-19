package encoder

import (
	"time"
)

// AppendTime 追加行格式的时间
//
// @params dst 目标切片
//
// @params t 实际时间
//
// @params format 序列化格式（与time.Format保持一致）
//
// @return 序列化后的行格式时间字符串
//
// 返回示例，反引号内为实际内容:
// `2006/01/02 15:04:05.000`
func AppendTime(dst []byte, t time.Time, format string) []byte {
	return t.AppendFormat(dst, format)
}

// AppendTimeJSON 追加JSON格式的时间
//
// @params dst 目标切片
//
// @params key 时间的JSON键名
//
// @params t 实际时间
//
// @params format 序列化格式（与time.Format保持一致）
//
// @return 序列化后的JSON格式时间字符串
//
// 返回示例，反引号内为实际内容:
// `"time": "2006/01/02 15:04:05.000"`
func AppendTimeJSON(dst []byte, key string, t time.Time, format string) []byte {
	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, `": "`...)
	dst = t.AppendFormat(dst, format)
	dst = append(dst, '"')
	return dst
}
