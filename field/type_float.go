/**
 * @Title 浮点型及浮点型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc float32|float64
 * @Desc *float32|*float64
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"strconv"
	"unsafe"

	"github.com/bearki/belog/v2/internal/convert"
)

// Float64 格式化float64类型字段信息
//
// 拼接格式  "name": 128.4566764
func Float64(name string, value float64) Field {
	f := Field{
		KeyBytes:       convert.StringToBytes(name),
		ValPrefixBytes: normalValPrefix[:],
		ValBytes:       strconv.AppendFloat(eightCapBytesPool.Get(), value, 'E', -1, 64),
		ValSuffixBytes: normalValSuffix[:],
		valBytesPut:    eightCapBytesPool.Put,
	}
	return f
}

// Float32 格式化float32类型字段信息
//
// 拼接格式  "name": 128.4566764
func Float32(name string, value float32) Field {
	f := Field{
		KeyBytes:       convert.StringToBytes(name),
		ValPrefixBytes: normalValPrefix[:],
		ValBytes:       strconv.AppendFloat(eightCapBytesPool.Get(), float64(value), 'E', -1, 32),
		ValSuffixBytes: normalValSuffix[:],
		valBytesPut:    eightCapBytesPool.Put,
	}
	return f
}

// Float64p 格式化*float64类型字段信息
//
// 注意: 空指针将输出null
//
// 拼接格式  "name": 128.4566764
func Float64p(name string, valuep *float64) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Float64(name, *valuep)
}

// Float32p 格式化*float32类型字段信息
//
// 注意: 空指针将输出null
//
// 拼接格式  "name": 128.4566764
func Float32p(name string, valuep *float32) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Float32(name, *valuep)
}
