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
)

//------------------------------ 值类型转换 ------------------------------//

// String 格式化string类型字段信息
func String(name string, value string) Field {
	f := Field{
		Key:     name,
		ValType: TypeString,
		String:  value,
	}
	return f
}

//------------------------------ 指针类型转换 ------------------------------//

// Boolp 格式化*string类型字段信息
func Stringp(name string, valuep *string) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return String(name, *valuep)
}
