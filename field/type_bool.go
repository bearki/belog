/**
 * @Title 布尔型及布尔型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc bool|*bool
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"github.com/bearki/belog/v2/internal/convert"
)

//------------------------------ 值类型转换 ------------------------------//

// Bool 格式化bool类型字段信息
func Bool(key string, val bool) Field {
	return Field{Key: key, Type: TypeBool, Integer: int64(convert.BoolToInt(val))}
}

//------------------------------ 指针类型转换 ------------------------------//

// Boolp 格式化*bool类型字段信息
func Boolp(key string, valp *bool) Field {
	if valp == nil {
		return nullField(key)
	}
	return Bool(key, *valp)
}

//------------------------------ 切片类型转换 ------------------------------//

// Bools 格式化[]bool类型字段信息
func Bools(key string, vals []bool) Field {
	return Field{Key: key, Type: TypeBools, Interface: vals}
}
