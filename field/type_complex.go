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
func Complex64(name string, value complex64) Field {
	return Field{
		Key:       name,
		ValType:   TypeComplex64,
		Interface: value,
	}
}

// Complex128 格式化complex128类型字段信息
func Complex128(name string, value complex128) Field {
	return Field{
		Key:       name,
		ValType:   TypeComplex128,
		Interface: value,
	}
}

//------------------------------ 指针类型转换 ------------------------------//

// Complex64 格式化*complex64类型字段信息
func Complex64p(name string, valuep *complex64) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Complex64(name, *valuep)
}

// Complex128p 格式化*complex128类型字段信息
func Complex128p(name string, valuep *complex128) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Complex128(name, *valuep)
}
