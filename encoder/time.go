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

// 追加时间值
func appendTimeValue(isJson bool, dst []byte, t time.Time, format string) []byte {
	// 赋值默认格式
	if len(format) == 0 {
		format = TimeFormat1
	}

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
		// 格式化为字符串时判断是否需要追加双引号
		if isJson {
			dst = append(dst, '"')
			dst = t.AppendFormat(dst, format)
			dst = append(dst, '"')
			return dst
		}

		return t.AppendFormat(dst, format)
	}
}

// 追加行格式的时间
//
//	@param	dst		目标切片
//	@param	key		时间的键名
//	@param	t		实际时间
//	@param	format	序列化格式（与time.Format保持一致）
//	@return	序列化后的行格式时间字符串
//
// 返回示例: 2006/01/02 15:04:05.000
func appendTime(dst []byte, t time.Time, format string) []byte {
	// 追加时间值
	return appendTimeValue(false, dst, t, format)
}

// 追加JSON格式的时间
//
//	@param	dst		目标切片
//	@param	key		时间的JSON键名
//	@param	t		实际时间
//	@param	format	序列化格式（与time.Format保持一致）
//	@return	序列化后的JSON格式时间字符串
//
// 返回示例: "time": "2006/01/02 15:04:05.000"` || `"time": 123456789000
func appendTimeJSON(dst []byte, key string, t time.Time, format string) []byte {
	// 拼接键名
	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, `": `...)

	// 追加时间值
	return appendTimeValue(true, dst, t, format)
}
