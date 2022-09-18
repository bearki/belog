/**
 * @Title 整型及整型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc int64|int32|int|int16|int8|byte
 * @Desc *int64|*int32|*int|*int16|*int8|*byte
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"strconv"
	"unsafe"

	"github.com/bearki/belog/v2/internal/convert"
)

// Int64 格式化int64类型字段信息
//
// 拼接格式  "index": 20
func Int64(name string, value int64) Field {
	f := Field{
		KeyBytes:       convert.StringToBytes(name),
		ValPrefixBytes: normalValPrefix[:],
		ValBytes:       strconv.AppendInt(eightCapBytesPool.Get(), value, 10),
		ValSuffixBytes: normalValSuffix[:],
		valBytesPut:    eightCapBytesPool.Put,
	}
	return f
}

// Int32 格式化int32类型字段信息
//
// 拼接格式  "index": 20
func Int32(name string, value int32) Field {
	return Int64(name, int64(value))
}

// Int 格式化int类型字段信息
//
// 拼接格式  "index": 20
func Int(name string, value int) Field {
	return Int64(name, int64(value))
}

// Int16 格式化int16类型字段信息
//
// 拼接格式  "index": 20
func Int16(name string, value int16) Field {
	return Int64(name, int64(value))
}

// Int8 格式化int8类型字段信息
//
// 拼接格式  "index": 20
func Int8(name string, value int8) Field {
	return Int64(name, int64(value))
}

// Int64p 格式化*int64类型字段信息
//
// 拼接格式  "index": 20
func Int64p(name string, valuep *int64) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Int64(name, int64(*valuep))
}

// Int32p 格式化*int32类型字段信息
//
// 拼接格式  "index": 20
func Int32p(name string, valuep *int32) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Int64(name, int64(*valuep))
}

// Intp 格式化*int类型字段信息
//
// 拼接格式  "index": 20
func Intp(name string, valuep *int) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Int64(name, int64(*valuep))
}

// Int16p 格式化*int16类型字段信息
//
// 拼接格式  "index": 20
func Int16p(name string, valuep *int16) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Int64(name, int64(*valuep))
}

// Int8p 格式化*int8类型字段信息
//
// 拼接格式  "index": 20
func Int8p(name string, valuep *int8) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Int64(name, int64(*valuep))
}
