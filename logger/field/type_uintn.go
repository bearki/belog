/**
 * @Title 整型及整型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc uint64|uint32|uint|uint16|uint8|byte
 * @Desc *uint64|*uint32|*uint|*uint16|*uint8|*byte
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"strconv"
	"unsafe"

	"github.com/bearki/belog/v2/internal/convert"
)

// Uint64 格式化uint64类型字段信息
//
// 拼接格式  "index": 20
func Uint64(name string, value uint64) Field {
	f := Field{
		KeyBytes:       convert.StringToBytes(name),
		ValPrefixBytes: normalValPrefix[:],
		ValBytes:       strconv.AppendUint(numberBytesPool.Get(), value, 10),
		ValSuffixBytes: normalValSuffix[:],
		valBytesPut:    numberBytesPool.Put,
	}
	return f
}

// Uint32 格式化uint32类型字段信息
//
// 拼接格式  "index": 20
func Uint32(name string, value uint32) Field {
	return Uint64(name, uint64(value))
}

// Uint 格式化uint类型字段信息
//
// 拼接格式  "index": 20
func Uint(name string, value uint) Field {
	return Uint64(name, uint64(value))
}

// Uint16 格式化uint16类型字段信息
//
// 拼接格式  "index": 20
func Uint16(name string, value uint16) Field {
	return Uint64(name, uint64(value))
}

// Uint8 格式化uint8类型字段信息
//
// 拼接格式  "index": 20
func Uint8(name string, value uint8) Field {
	return Uint64(name, uint64(value))
}

// Byte 格式化byte类型字段信息
//
// 拼接格式  "index": 255
func Byte(name string, value byte) Field {
	return Uint64(name, uint64(value))
}

// Uint64p 格式化*uint64类型字段信息
//
// 拼接格式  "index": 20
func Uint64p(name string, valuep *uint64) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uint32p 格式化*uint32类型字段信息
//
// 拼接格式  "index": 20
func Uint32p(name string, valuep *uint32) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uintp 格式化*uint类型字段信息
//
// 拼接格式  "index": 20
func Uintp(name string, valuep *uint) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uint16p 格式化*uint16类型字段信息
//
// 拼接格式  "index": 20
func Uint16p(name string, valuep *uint16) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uint8p 格式化*uint8类型字段信息
//
// 拼接格式  "index": 20
func Uint8p(name string, valuep *uint8) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Bytep 格式化*byte类型字段信息
//
// 拼接格式  "index": 255
func Bytep(name string, valuep *byte) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}
