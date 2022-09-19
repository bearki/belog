/**
 * @Title 字符串及字符串指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc string|*string
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"unsafe"

	"github.com/bearki/belog/v2/internal/convert"
)

// Bool 格式化string类型字段信息
//
// 拼接格式  "name": "value"
func String(name string, value string) Field {
	f := Field{
		KeyBytes:       convert.StringToBytes(name),
		ValPrefixBytes: stringValPrefix[:],
		ValBytes:       convert.StringToBytes(value),
		ValSuffixBytes: stringValSuffix[:],
		valBytesPut:    nil,
	}
	return f
}

// Boolp 格式化*string类型字段信息
//
// 拼接格式  "name": "valuep"
func Stringp(name string, valuep *string) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return String(name, *valuep)
}
