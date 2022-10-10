/**
 * @Title 整型及整型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc int64|int32|int|int16|int8
 * @Desc *int64|*int32|*int|*int16|*int8
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

// 组装~int类型的字段结构
func intn(key string, val int64, vt Type) Field {
	return Field{Key: key, Type: vt, Integer: val}
}

//------------------------------ 值类型转换 ------------------------------//

// Int8 格式化int8类型字段信息
func Int8(key string, val int8) Field {
	return intn(key, int64(val), TypeInt8)
}

// Int16 格式化int16类型字段信息
func Int16(key string, val int16) Field {
	return intn(key, int64(val), TypeInt16)
}

// Int 格式化int类型字段信息
func Int(key string, val int) Field {
	return intn(key, int64(val), TypeInt)
}

// Int32 格式化int32类型字段信息
func Int32(key string, val int32) Field {
	return intn(key, int64(val), TypeInt32)
}

// Int64 格式化int64类型字段信息
func Int64(key string, val int64) Field {
	return intn(key, val, TypeInt64)
}

//------------------------------ 指针类型转换 ------------------------------//

// Int8p 格式化*int8类型字段信息
func Int8p(key string, valp *int8) Field {
	if valp == nil {
		return nullField(key)
	}
	return Int8(key, *valp)
}

// Int16p 格式化*int16类型字段信息
func Int16p(key string, valp *int16) Field {
	if valp == nil {
		return nullField(key)
	}
	return Int16(key, *valp)
}

// Intp 格式化*int类型字段信息
func Intp(key string, valp *int) Field {
	if valp == nil {
		return nullField(key)
	}
	return Int(key, *valp)
}

// Int32p 格式化*int32类型字段信息
func Int32p(key string, valp *int32) Field {
	if valp == nil {
		return nullField(key)
	}
	return Int32(key, *valp)
}

// Int64p 格式化*int64类型字段信息
func Int64p(key string, valp *int64) Field {
	if valp == nil {
		return nullField(key)
	}
	return Int64(key, *valp)
}

//------------------------------ 切片类型转换 ------------------------------//

// Int8s 格式化[]int8类型字段信息
func Int8s(key string, vals []int8) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeInt8s, Interface: vals}
}

// Int16s 格式化[]int16类型字段信息
func Int16s(key string, vals []int16) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeInt16s, Interface: vals}
}

// Ints 格式化[]int类型字段信息
func Ints(key string, vals []int) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeInts, Interface: vals}
}

// Int32s 格式化[]int32类型字段信息
func Int32s(key string, vals []int32) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeInt32s, Interface: vals}
}

// Int64s 格式化[]int64类型字段信息
func Int64s(key string, vals []int64) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeInt64s, Interface: vals}
}
