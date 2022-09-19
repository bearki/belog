/**
 * @Title 复数型及复数型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc complex64|complex128
 * @Desc *complex64|*complex128
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"strconv"
	"unsafe"

	"github.com/bearki/belog/v2/internal/convert"
)

// Complex128 格式化complex128类型字段信息
//
// 拼接格式  "name": (1.1+1i)
func Complex128(name string, value complex128) Field {
	f := Field{
		KeyBytes:       convert.StringToBytes(name),
		ValPrefixBytes: normalValPrefix[:],
		ValBytes:       convert.StringToBytes(strconv.FormatComplex(value, 'f', -1, 128)),
		ValSuffixBytes: normalValSuffix[:],
		valBytesPut:    nil,
	}
	return f
}

// Complex64 格式化complex64类型字段信息
//
// 拼接格式  "name": (1.1+1i)
func Complex64(name string, value complex64) Field {
	f := Field{
		KeyBytes:       convert.StringToBytes(name),
		ValPrefixBytes: normalValPrefix[:],
		ValBytes:       convert.StringToBytes(strconv.FormatComplex(complex128(value), 'f', -1, 64)),
		ValSuffixBytes: normalValSuffix[:],
		valBytesPut:    nil,
	}
	return f
}

// Complex128p 格式化*complex128类型字段信息
//
// 拼接格式  "name": (1.1+1i)
func Complex128p(name string, valuep *complex128) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Complex128(name, *valuep)
}

// Complex64 格式化*complex64类型字段信息
//
// 拼接格式  "name": (1.1+1i)
func Complex64p(name string, valuep *complex64) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Complex64(name, *valuep)
}
