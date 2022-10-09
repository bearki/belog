package field

import (
	"time"
	"unsafe"
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

//------------------------------ 指针类型转换 ------------------------------//

// Timep 格式化*time.Time类型字段信息
func Timep(name string, valuep *time.Time, format ...string) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Time(name, *valuep, format...)
}
