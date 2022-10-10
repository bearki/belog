/**
 * @Title 复数型及复数型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc complex64|complex128
 * @Desc *complex64|*complex128
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

//------------------------------ 值类型转换 ------------------------------//

// Complex64 格式化complex64类型字段信息
func Complex64(key string, val complex64) Field {
	return Field{Key: key, Type: TypeComplex64, Interface: val}
}

// Complex128 格式化complex128类型字段信息
func Complex128(key string, val complex128) Field {
	return Field{Key: key, Type: TypeComplex128, Interface: val}
}

//------------------------------ 指针类型转换 ------------------------------//

// Complex64 格式化*complex64类型字段信息
func Complex64p(key string, valp *complex64) Field {
	if valp == nil {
		return nullField(key)
	}
	return Complex64(key, *valp)
}

// Complex128p 格式化*complex128类型字段信息
func Complex128p(key string, valp *complex128) Field {
	if valp == nil {
		return nullField(key)
	}
	return Complex128(key, *valp)
}

//------------------------------ 切片类型转换 ------------------------------//

// Complex64s 格式化[]complex64类型字段信息
func Complex64s(key string, vals []complex64) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeComplex64s, Interface: vals}
}

// Complex128s 格式化[]complex128类型字段信息
func Complex128s(key string, vals []complex128) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeComplex128s, Interface: vals}
}
