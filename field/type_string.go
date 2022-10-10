/**
 * @Title 字符串及字符串指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc string|*string
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

//------------------------------ 值类型转换 ------------------------------//

// String 格式化string类型字段信息
func String(key string, val string) Field {
	return Field{Key: key, Type: TypeString, String: val}
}

//------------------------------ 指针类型转换 ------------------------------//

// Boolp 格式化*string类型字段信息
func Stringp(key string, valp *string) Field {
	if valp == nil {
		return nullField(key)
	}
	return String(key, *valp)
}

//------------------------------ 切片类型转换 ------------------------------//

// Strings 格式化[]string类型字段信息
func Strings(key string, vals []string) Field {
	return Field{Key: key, Type: TypeStrings, Interface: vals}
}
