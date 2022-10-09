/**
 * @Title 整型及整型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc uint64|uint32|uint|uint16|uint8|byte
 * @Desc *uint64|*uint32|*uint|*uint16|*uint8|*byte
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

//------------------------------ 值类型转换 ------------------------------//

// Byte 格式化byte类型字段信息
func Byte(name string, value byte) Field {
	return intn(name, int64(value), TypeByte)
}

// Uint8 格式化uint8类型字段信息
func Uint8(name string, value uint8) Field {
	return intn(name, int64(value), TypeUint8)
}

// Uint16 格式化uint16类型字段信息
func Uint16(name string, value uint16) Field {
	return intn(name, int64(value), TypeUint16)
}

// Uint 格式化uint类型字段信息
func Uint(name string, value uint) Field {
	return intn(name, int64(value), TypeUint)
}

// Uintptr 格式化uintptr类型字段信息
func Uintptr(name string, value uintptr) Field {
	return intn(name, int64(value), TypeUintptr)
}

// Uint32 格式化uint32类型字段信息
func Uint32(name string, value uint32) Field {
	return intn(name, int64(value), TypeUint32)
}

// Uint64 格式化uint64类型字段信息
func Uint64(name string, value uint64) Field {
	return intn(name, int64(value), TypeUint64)
}

//------------------------------ 指针类型转换 ------------------------------//

// Bytep 格式化*byte类型字段信息
func Bytep(name string, valuep *byte) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Byte(name, *valuep)
}

// Uint8p 格式化*uint8类型字段信息
func Uint8p(name string, valuep *uint8) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Uint8(name, *valuep)
}

// Uint16p 格式化*uint16类型字段信息
func Uint16p(name string, valuep *uint16) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Uint16(name, *valuep)
}

// Uintp 格式化*uint类型字段信息
func Uintp(name string, valuep *uint) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Uint(name, *valuep)
}

// Uintptrp 格式化*uintptr类型字段信息
func Uintptrp(name string, valuep *uintptr) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Uintptr(name, *valuep)
}

// Uint32p 格式化*uint32类型字段信息
func Uint32p(name string, valuep *uint32) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Uint32(name, *valuep)
}

// Uint64p 格式化*uint64类型字段信息
func Uint64p(name string, valuep *uint64) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Uint64(name, *valuep)
}

//------------------------------ 切片类型转换 ------------------------------//

// Bytes 格式化byte类型字段信息
func Bytes(name string, values []byte) Field {
	return Field{Key: name, ValType: TypeBytes, Interface: values}
}

// Uint8s 格式化uint8类型字段信息
func Uint8s(name string, values []uint8) Field {
	return Field{Key: name, ValType: TypeUint8s, Interface: values}
}

// Uint16s 格式化uint16类型字段信息
func Uint16s(name string, values []uint16) Field {
	return Field{Key: name, ValType: TypeUint16s, Interface: values}
}

// Uints 格式化uint类型字段信息
func Uints(name string, values []uint) Field {
	return Field{Key: name, ValType: TypeUints, Interface: values}
}

// Uintptrs 格式化uintptr类型字段信息
func Uintptrs(name string, values []uintptr) Field {
	return Field{Key: name, ValType: TypeUintptrs, Interface: values}
}

// Uint32s 格式化uint32类型字段信息
func Uint32s(name string, values []uint32) Field {
	return Field{Key: name, ValType: TypeUint32s, Interface: values}
}

// Uint64s 格式化uint64类型字段信息
func Uint64s(name string, values []uint64) Field {
	return Field{Key: name, ValType: TypeUint64s, Interface: values}
}
