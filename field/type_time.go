package field

import (
	"time"
)

//------------------------------ 值类型转换 ------------------------------//

// Time 格式化time.Time类型字段信息
func Time(name string, value time.Time, format ...string) Field {
	var timeFmt string
	if len(format) > 0 {
		timeFmt = format[0]
	}
	return Field{
		Key:     name,
		ValType: TypeTime,
		Integer: value.UnixMicro(), // 微秒时间戳
		String:  timeFmt,
	}
}

// Duration 格式化time.Duration类型字段信息
func Duration(name string, value time.Duration) Field {
	return Field{
		Key:     name,
		ValType: TypeDuration,
		Integer: int64(value),
	}
}

//------------------------------ 指针类型转换 ------------------------------//

// Timep 格式化*time.Time类型字段信息
func Timep(name string, valuep *time.Time, format ...string) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Time(name, *valuep, format...)
}

// Durationp 格式化*time.Duration类型字段信息
func Durationp(name string, valuep *time.Duration) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Duration(name, *valuep)
}
