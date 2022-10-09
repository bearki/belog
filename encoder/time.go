package encoder

import (
	"strconv"
	"time"
)

const (
	TimeFormat1 = "2006/01/02 15:04:05.000"
	TimeFormat2 = "2006-01-02 15:04:05.000"
	TimeFormat3 = "2006/01/02 15:04:05"
	TimeFormat4 = "2006-01-02 15:04:05"
	TimeFormat5 = "2006/01/02 15:04"
	TimeFormat6 = "2006-01-02 15:04"
	TimeFormat7 = "2006/01/02"
	TimeFormat8 = "2006-01-02"

	TimeFormatUnix      = "Unix"      // 秒级时间戳
	TimeFormatUnixMilli = "UnixMilli" // 毫秒级时间戳
	TimeFormatUnixMicro = "UnixMicro" // 微秒级时间戳
	TimeFormatUnixNano  = "UnixNano"  // 纳秒级时间戳
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
	// 是否使用时间戳
	switch format {

	case TimeFormatUnix:
		return strconv.AppendInt(dst, t.Unix(), 10)

	case TimeFormatUnixMilli:
		return strconv.AppendInt(dst, t.UnixMilli(), 10)

	case TimeFormatUnixMicro:
		return strconv.AppendInt(dst, t.UnixMicro(), 10)

	case TimeFormatUnixNano:
		return strconv.AppendInt(dst, t.UnixNano(), 10)

	default:
		return t.AppendFormat(dst, format)

	}
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
	dst = append(dst, `": `...)

	// 是否使用时间戳
	switch format {

	case TimeFormatUnix:
		return strconv.AppendInt(dst, t.Unix(), 10)

	case TimeFormatUnixMilli:
		return strconv.AppendInt(dst, t.UnixMilli(), 10)

	case TimeFormatUnixMicro:
		return strconv.AppendInt(dst, t.UnixMicro(), 10)

	case TimeFormatUnixNano:
		return strconv.AppendInt(dst, t.UnixNano(), 10)

	default:
		dst = append(dst, '"')
		dst = t.AppendFormat(dst, format)
		dst = append(dst, '"')
		return dst

	}
}
