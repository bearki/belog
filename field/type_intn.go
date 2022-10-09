/**
 * @Title 整型及整型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc int64|int32|int|int16|int8|byte
 * @Desc *int64|*int32|*int|*int16|*int8|*byte
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

// 组装~int类型的字段结构
func intn(name string, value int64, vt Type) Field {
	return Field{Key: name, ValType: vt, Integer: value}
}

//------------------------------ 值类型转换 ------------------------------//

// Int8 格式化int8类型字段信息
func Int8(name string, value int8) Field {
	return intn(name, int64(value), TypeInt8)
}

// Int16 格式化int16类型字段信息
func Int16(name string, value int16) Field {
	return intn(name, int64(value), TypeInt16)
}

// Int 格式化int类型字段信息
func Int(name string, value int) Field {
	return intn(name, int64(value), TypeInt)
}

// Int32 格式化int32类型字段信息
func Int32(name string, value int32) Field {
	return intn(name, int64(value), TypeInt32)
}

// Int64 格式化int64类型字段信息
func Int64(name string, value int64) Field {
	return intn(name, value, TypeInt64)
}

//------------------------------ 指针类型转换 ------------------------------//

// Int8p 格式化*int8类型字段信息
func Int8p(name string, valuep *int8) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Int8(name, *valuep)
}

// Int16p 格式化*int16类型字段信息
func Int16p(name string, valuep *int16) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Int16(name, *valuep)
}

// Intp 格式化*int类型字段信息
func Intp(name string, valuep *int) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Int(name, *valuep)
}

// Int32p 格式化*int32类型字段信息
func Int32p(name string, valuep *int32) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Int32(name, *valuep)
}

// Int64p 格式化*int64类型字段信息
func Int64p(name string, valuep *int64) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Int64(name, *valuep)
}

//------------------------------ 切片类型转换 ------------------------------//

// Int8s 格式化[]int8类型字段信息
func Int8s(name string, values []int8) Field {
	return Field{Key: name, ValType: TypeInt8s, Interface: values}
}

// Int16s 格式化[]int16类型字段信息
func Int16s(name string, values []int16) Field {
	return Field{Key: name, ValType: TypeInt16s, Interface: values}
}

// Ints 格式化[]int类型字段信息
func Ints(name string, values []int) Field {
	return Field{Key: name, ValType: TypeInts, Interface: values}
}

// Int32s 格式化[]int32类型字段信息
func Int32s(name string, values []int32) Field {
	return Field{Key: name, ValType: TypeInt32s, Interface: values}
}

// Int64s 格式化[]int64类型字段信息
func Int64s(name string, values []int64) Field {
	return Field{Key: name, ValType: TypeInt64s, Interface: values}
}
