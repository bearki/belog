/**
 * @Title 布尔型及布尔型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc bool|*bool
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"unsafe"

	"github.com/bearki/belog/v2/internal/convert"
)

// 布尔类型常量数组，用于切片的底层引用
var (
	trueConstArray  = [4]byte{'t', 'r', 'u', 'e'}
	falseConstArray = [5]byte{'f', 'a', 'l', 's', 'e'}
)

func boolBytes(value bool) []byte {
	if value {
		return trueConstArray[:]
	}
	return falseConstArray[:]
}

// Bool 格式化bool类型字段信息
//
// 拼接格式  "name": true
func Bool(name string, value bool) Field {
	f := Field{
		KeyBytes:       convert.StringToBytes(name),
		ValPrefixBytes: normalValPrefix[:],
		ValBytes:       boolBytes(value),
		ValSuffixBytes: normalValSuffix[:],
		valBytesPut:    nil,
	}
	return f
}

// Boolp 格式化*bool类型字段信息
//
// 注意: 空指针将输出null
//
// 拼接格式  "name": true
func Boolp(name string, valuep *bool) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Bool(name, *valuep)
}
