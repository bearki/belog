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

//------------------------------ 值类型转换 ------------------------------//

// Bool 格式化bool类型字段信息
func Bool(name string, value bool) Field {
	return Field{
		Key:     name,
		ValType: TypeBool,
		Integer: int64(convert.BoolToInt(value)),
	}
}

//------------------------------ 指针类型转换 ------------------------------//

// Boolp 格式化*bool类型字段信息
func Boolp(name string, valuep *bool) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Bool(name, *valuep)
}
